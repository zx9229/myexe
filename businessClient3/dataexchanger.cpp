#include "dataexchanger.h"
#include <QCoreApplication>
#include <QSqlQuery>
#include <QSqlError>
#include <QtCore/QMetaEnum>
#include "m2b.h"
#include "google/protobuf/util/json_util.h"
#include "zxtools.h"

enum StatusErrorType
{
    WebsockError = 1,
    WebsockDisconnected = 2,
    ConnectReqError = 3,
    ConnectRspError = 4,
};

DataExchanger::DataExchanger(QObject *parent) :
    QObject(parent),
    m_ws(parent),
    m_MsgNo("MsgNo", QString())
{
    connect(&m_ws, &MyWebsock::sigConnected, this, &DataExchanger::slotOnConnected);
    connect(&m_ws, &MyWebsock::sigDisconnected, this, &DataExchanger::slotOnDisconnected);
    connect(&m_ws, &MyWebsock::sigMessage, this, &DataExchanger::slotOnMessage);
    connect(&m_ws, &MyWebsock::sigError, this, &DataExchanger::slotOnError);

    initOwnInfo();
    initDB();
}

DataExchanger::~DataExchanger()
{

}

QString DataExchanger::dbLoadValue(const QString& key)
{
    QSqlQuery sqlQuery;
    KeyValue kv(key, QString());
    return kv.refresh_data(sqlQuery) ? kv.Value : QString();
}

bool DataExchanger::dbSaveValue(const QString& key, const QString& value)
{
    QSqlQuery sqlQuery;
    KeyValue kv(key, value);
    return kv.insert_data(sqlQuery, false);
}

void DataExchanger::setURL(const QString &url)
{
    m_url = url;  // 例如【ws://localhost:65535/websocket】.
}

void DataExchanger::setOwnInfo(const QString &userID, const QString &belongID)
{
    m_ownInfo.set_userid(userID.toStdString());
    m_ownInfo.set_belongid(belongID.toStdString());
}

QString DataExchanger::getUserID()
{
    return QString::fromStdString(m_ownInfo.userid());
}

bool DataExchanger::start()
{
    m_ws.stop(true);
    return m_ws.start(m_url);
}

QString DataExchanger::sendReq(const QString &typeName, const QString &jsonText, const QString &rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB)
{
    QString message;
    int64_t msgNo = 0;
    if (fillMsgNo || !isC1NotC2)
    {
        msgNo = m_MsgNo.Value.toLongLong() + 1;
    }
    GPMSGPTR msgData;
    message = toC1C2(msgNo, typeName, jsonText, rID, isLog, isSafe, isPush, isUpCache, isC1NotC2, msgData);
    if (!message.isEmpty())
        return message;

    CommonData tmpCommonData;

    if (isC1NotC2)
    {
        QSharedPointer<txdata::Common1Req> c1data = qSharedPointerDynamicCast<txdata::Common1Req>(msgData);
        message = sendCommon1Req(c1data);
        if (message.isEmpty() || forceToDB)
        {
            zxtools::Common1Req2CommonData(&tmpCommonData, c1data.get());
            tmpCommonData.SendOK = message.isEmpty();
        }
    }
    else
    {
        QSharedPointer<txdata::Common2Req> c2data = qSharedPointerDynamicCast<txdata::Common2Req>(msgData);
        message = sendCommon2Req(c2data);
        if (message.isEmpty() || forceToDB)
        {
            zxtools::Common2Req2CommonData(&tmpCommonData, c2data.get());
            tmpCommonData.SendOK = message.isEmpty();
        }
    }

    if ((0 != tmpCommonData.MsgNo && message.isEmpty()) || forceToDB)
    {
        QSqlQuery sqlQuery;
        tmpCommonData.insert_data(sqlQuery, true, nullptr);
        m_MsgNo.Value.setNum(msgNo);
        m_MsgNo.update_data(sqlQuery);
        emit sigTableChanged(tmpCommonData.static_table_name());
    }

    return message;
}

QString DataExchanger::QryConnInfoReq(const QString &userId)
{
    txdata::QryConnInfoReq tmpData;
    QString typeName = m2b::CalcMsgTypeName(tmpData);
    QString jsonText = zxtools::object2json(tmpData, nullptr);
    return sendReq(typeName, jsonText, userId, false, false, false, false, true, false, false);
}

QStringList DataExchanger::getTxMsgTypeNameList()
{
    QStringList typeNameList;
    typeNameList << QString::fromStdString(::txdata::MsgType_Name(txdata::ID_EchoItem));
    typeNameList << QString::fromStdString(::txdata::MsgType_Name(txdata::ID_EmailItem));
    typeNameList << QString::fromStdString(::txdata::MsgType_Name(txdata::ID_ExecCmdReq));
    typeNameList << QString::fromStdString(::txdata::MsgType_Name(txdata::ID_QueryRecordReq));
    typeNameList << QString::fromStdString(::txdata::MsgType_Name(txdata::ID_PushItem));
    typeNameList << QString::fromStdString(::txdata::MsgType_Name(txdata::ID_SubscribeReq));
    return  typeNameList;
}

QString DataExchanger::jsonExample(const QString& typeName)
{
    txdata::MsgType msgType = txdata::MsgType::Zero1;
    txdata::MsgType_Parse(typeName.toStdString(), &msgType);
    if (msgType == txdata::MsgType::ID_EchoItem)
    {
        txdata::EchoItem tmpObj;
        tmpObj.set_data("DATA");
        tmpObj.set_rspcnt(2);
        tmpObj.set_secgap(0);
        return zxtools::object2json(tmpObj);
    }
    if (msgType == txdata::MsgType::ID_PushItem)
    {
        txdata::PushItem tmpObj;
        tmpObj.set_subject("Subject");
        tmpObj.set_content("Content");
        return zxtools::object2json(tmpObj);
    }
    if (msgType == txdata::MsgType::ID_SubscribeReq)
    {
        txdata::SubscribeReq tmpObj;
        return zxtools::object2json(tmpObj);
    }
    return "";
}

void DataExchanger::initDB()
{
    m_db = QSqlDatabase::addDatabase("QSQLITE");
    m_db.setDatabaseName(QString().isEmpty() ? ("_zx_test.db") : (":memory:"));
    bool isOk = false;
    isOk = m_db.open();
    Q_ASSERT(isOk);
    QSqlQuery sqlQuery;
    if (true) {
        isOk = m_db.transaction();
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(KeyValue::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(ConnInfoEx::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(CommonData::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(PushWrap::static_create_sql());
        Q_ASSERT(isOk);
        isOk = m_db.commit();
        Q_ASSERT(isOk);
    }
    if (true) {
        m_MsgNo.refresh_data(sqlQuery);
        isOk = m_MsgNo.insert_data(sqlQuery, false);
        Q_ASSERT(isOk);
    }
}

void DataExchanger::initOwnInfo()
{
    m_ownInfo.set_version("Ver20190108");
    m_ownInfo.set_linkmode(txdata::ConnectionInfo_LinkType_CONNECT);
    m_ownInfo.set_exepid(static_cast<int>(QCoreApplication::applicationPid()));
    m_ownInfo.set_exepath(QCoreApplication::applicationFilePath().toStdString());
    m_ownInfo.set_remark("");
}

QString DataExchanger::toC1C2(int64_t msgNo, const QString &typeName, const QString &jsonText, const QString &rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, GPMSGPTR &msgOut)
{
    QString message;

    std::string binData;
    txdata::MsgType curType = txdata::MsgType::Zero1;
    if (zxtools::json2binary(typeName, jsonText, curType, binData) == false)
        message = "serialized data failed";

    if (!message.isEmpty())
        return message;

    if (isC1NotC2)
    {
        QSharedPointer<txdata::Common1Req> c1req = QSharedPointer<txdata::Common1Req>(new txdata::Common1Req);
        c1req->set_msgno(msgNo);
        c1req->set_seqno(0);
        c1req->set_senderid(this->m_ownInfo.userid());
        c1req->set_recverid(rID.toStdString());
        c1req->set_toroot(true);
        c1req->set_islog(isLog);
        c1req->set_ispush(isPush);
        c1req->set_reqtype(curType);
        c1req->set_reqdata(binData);
        zxtools::qdt2gpt(*(c1req->mutable_reqtime()), QDateTime::currentDateTime());
        msgOut = c1req;
    }
    else
    {
        QSharedPointer<txdata::Common2Req> c2req = QSharedPointer<txdata::Common2Req>(new txdata::Common2Req);
        c2req->mutable_key()->set_userid(m_ownInfo.userid());
        c2req->mutable_key()->set_msgno(msgNo);
        c2req->mutable_key()->set_seqno(0);
        c2req->set_senderid(m_ownInfo.userid());
        c2req->set_recverid(rID.toStdString());
        c2req->set_toroot(true);
        c2req->set_islog(isLog);
        c2req->set_issafe(isSafe);
        c2req->set_ispush(isPush);
        c2req->set_upcache(isUpCache);
        c2req->set_reqtype(curType);
        c2req->set_reqdata(binData);
        zxtools::qdt2gpt(*(c2req->mutable_reqtime()), QDateTime::currentDateTime());
        msgOut = c2req;
    }

    return message;
}

QString DataExchanger::sendCommon1Req(QSharedPointer<txdata::Common1Req> data)
{
    qint64 sendBytes = m_ws.sendBinaryMessage(m2b::msg2package(*data));
    //一个字节都没发出去,肯定就发送失败了.
    return (0 == sendBytes) ? "send failed" : "";
}

QString DataExchanger::sendCommon2Req(QSharedPointer<txdata::Common2Req> data)
{
    if (data->issafe())
    {
        bool isOk = m_cacheSync.insertData(data->key(), data->toroot(), data->recverid(), data);
        Q_ASSERT(isOk);
    }
    qint64 sendBytes = m_ws.sendBinaryMessage(m2b::msg2package(*data));
    //一个字节都没发出去,肯定就发送失败了.
    return (0 == sendBytes) ? "send failed" : "";
}

void DataExchanger::handle_Common2Ack(QSharedPointer<txdata::Common2Ack> data)
{
    Q_ASSERT(data.data() != nullptr);
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << QString::fromStdString(data->GetTypeName());
    //TODO:
}

void DataExchanger::handle_Common2Rsp(QSharedPointer<txdata::Common2Rsp> msgData)
{
    Q_ASSERT(msgData.data() != nullptr);
    if (msgData->islog())
    {
        qDebug() << QString::fromStdString(msgData->DebugString());
    }
    GPMSGPTR curData;
    if (!m2b::slice2msg(msgData->rspdata(), msgData->rsptype(), curData))
    {
        qDebug() << "handle_Common2Rsp," << "slice2msg fail," << msgData->rsptype();
        return;
    }
    Q_ASSERT(msgData->rsptype() == m2b::CalcMsgType(*curData));

    if (msgData->mutable_key()->msgno() != 0)
    {
        CommonData tmpData;
        zxtools::Common2Rsp2CommonData(&tmpData, msgData.get());
        QSqlQuery sqlQuery;
        bool isOk = tmpData.insert_data(sqlQuery, true, nullptr);
        Q_ASSERT(isOk);
        isOk = tmpData.updateRequestData(sqlQuery);
        Q_ASSERT(isOk);
        emit sigTableChanged(tmpData.static_table_name());
    }
}

void DataExchanger::handle_Common1Req(QSharedPointer<txdata::Common1Req> msgData)
{
    Q_ASSERT(msgData.data() != nullptr);
    if (msgData->islog())
    {
        qDebug() << QString::fromStdString(msgData->DebugString());
    }
    GPMSGPTR curData;
    if (!m2b::slice2msg(msgData->reqdata(), msgData->reqtype(), curData))
    {
        qDebug() << "handle_Common1Rsp," << "slice2msg fail," << msgData->reqtype();
        return;
    }
    Q_ASSERT(msgData->reqtype() == m2b::CalcMsgType(*curData));

    if (msgData->msgno() != 0)
    {
        CommonData tmpData;
        zxtools::Common1Req2CommonData(&tmpData, msgData.get());
        QSqlQuery sqlQuery;
        bool isOk = tmpData.insert_data(sqlQuery, true, nullptr);
        Q_ASSERT(isOk);
        emit sigTableChanged(tmpData.static_table_name());
    }

    if (msgData->reqtype() == txdata::MsgType::ID_PushWrap)
    {
        GPMSGPTR tmpData;
        if (m2b::slice2msg(msgData->reqdata(), msgData->reqtype(), tmpData))
        {
            QSharedPointer<txdata::PushWrap> tmp2Data = qSharedPointerDynamicCast<txdata::PushWrap>(tmpData);
            PushWrap pshData;
            zxtools::PushWrap2PushWrap(&pshData, tmp2Data.get(), msgData->senderid());
            QSqlQuery sqlQuery;
            bool isOk = pshData.insert_data(sqlQuery, false, nullptr);
            Q_ASSERT(isOk);
            emit sigTableChanged(pshData.static_table_name());
        }
    }

    switch (msgData->reqtype()) {
    default:
        break;
    }
}

void DataExchanger::handle_Common1Rsp(QSharedPointer<txdata::Common1Rsp> msgData)
{
    Q_ASSERT(msgData.data() != nullptr);
    if (msgData->islog())
    {
        qDebug() << QString::fromStdString(msgData->DebugString());
    }
    GPMSGPTR curData;
    if (!m2b::slice2msg(msgData->rspdata(), msgData->rsptype(), curData))
    {
        qDebug() << "handle_Common1Rsp," << "slice2msg fail," << msgData->rsptype();
        return;
    }
    Q_ASSERT(msgData->rsptype() == m2b::CalcMsgType(*curData));

    if (msgData->msgno() != 0)
    {
        CommonData tmpData;
        zxtools::Common1Rsp2CommonData(&tmpData, msgData.get());
        QSqlQuery sqlQuery;
        bool isOk = tmpData.insert_data(sqlQuery, true, nullptr);
        Q_ASSERT(isOk);
        isOk = tmpData.updateRequestData(sqlQuery);
        Q_ASSERT(isOk);
        emit sigTableChanged(tmpData.static_table_name());
    }

    switch (msgData->rsptype()) {
    case txdata::ID_QryConnInfoRsp:
        deal_QryConnInfoRsp(qSharedPointerDynamicCast<txdata::QryConnInfoRsp>(curData));
        break;
    default:
        break;
    }
}

void DataExchanger::handle_ConnectReq(QSharedPointer<txdata::ConnectReq> data)
{
    Q_ASSERT(data.data() != nullptr);
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << QString::fromStdString(data->GetTypeName());
    txdata::ConnectRsp data4send = {};
    {
        data4send.mutable_inforeq()->CopyFrom(data->inforeq());
        data4send.mutable_inforsp()->CopyFrom(m_ownInfo);
        data4send.set_errno(0);
    }
    bool checkOk = false;
    do
    {
        if (data->inforeq().userid() != m_ownInfo.belongid())
        {
            data4send.set_errmsg("(req.UserID != rsp.BelongID)");
            break;
        }
        if (data->pathway_size() != 1)
        {
            data4send.set_errmsg("len(req.Pathway) != 1");
            break;
        }
        if (data->pathway(0) != data->inforeq().userid())
        {
            data4send.set_errmsg("req.UserID != req.Pathway[0]");
            break;
        }
        checkOk = true;
    } while (false);
    m_ws.sendBinaryMessage(m2b::msg2package(data4send));
    if (checkOk)
    {
        m_parentInfo.CopyFrom(data->inforeq());
        emit sigReady();
        QString mesg = this->QryConnInfoReq("");
    }
    else
    {
        {
            QStringList strList;
            strList.append("ConnectReq");
            strList.append(QString::fromStdString(data4send.errmsg()));
            strList.append(m_ws.url().toString());
            emit sigStatusError(strList.join('\n'), StatusErrorType::ConnectReqError);
        }
        m_ws.interrupt();
    }
}

void DataExchanger::handle_ConnectRsp(QSharedPointer<txdata::ConnectRsp> data)
{
    Q_ASSERT(data.data() != nullptr);
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << QString::fromStdString(data->GetTypeName());
    if (data->errno() != 0)
    {
        QStringList strList;
        strList.append("ConnectRsp");
        strList.append(QString::fromStdString(data->errmsg()));
        strList.append(m_ws.url().toString());
        emit sigStatusError(strList.join('\n'), StatusErrorType::ConnectRspError);
    }
}

void DataExchanger::deal_QryConnInfoRsp(QSharedPointer<txdata::QryConnInfoRsp> msgData)
{
    Q_ASSERT(msgData.data() != nullptr);
    QSqlQuery sqlQuery;
    ConnInfoEx::delete_data(sqlQuery, "");
    for (auto&p : msgData->cache())
    {
        const txdata::ConnectReq& curData = p.second;
        ConnInfoEx cie;
        cie.UserID = QString::fromStdString(curData.inforeq().userid());
        cie.BelongID = QString::fromStdString(curData.inforeq().belongid());
        cie.Version = QString::fromStdString(curData.inforeq().version());
        cie.ExePid = curData.inforeq().exepid();
        cie.ExePath = QString::fromStdString(curData.inforeq().exepath());
        cie.Remark = QString::fromStdString(curData.inforeq().remark());
        QStringList pathway;
        for (int i = 0; i < curData.pathway_size(); i++) { pathway.append(QString::fromStdString(curData.pathway().Get(i))); }
        cie.Pathway = pathway.join("->");
        cie.insert_data(sqlQuery, false);
    }
    emit sigTableChanged(ConnInfoEx::static_table_name());
}

void DataExchanger::slotOnConnected()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnConnected";
    {
        txdata::ConnectReq data4send = {};
        data4send.mutable_inforeq()->CopyFrom(m_ownInfo);
        data4send.add_pathway(data4send.inforeq().userid());

        m_ws.sendBinaryMessage(m2b::msg2package(data4send));
    }
}

void DataExchanger::slotOnDisconnected()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnDisconnected";
    {
        QStringList strList;
        strList.append("Disconnected");
        strList.append(m_ws.url().toString());
        emit sigStatusError(strList.join('\n'), StatusErrorType::WebsockDisconnected);
    }
}

void DataExchanger::slotOnMessage(const QByteArray &message)
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnMessage";
    txdata::MsgType theType = {};
    GPMSGPTR theMsg;
    if (m2b::package2msg(message, theType, theMsg) == false)
    {
        qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnMessage, slice2msg, failure";
        return;
    }
    switch (theType) {
    case txdata::MsgType::ID_Common2Ack:
        handle_Common2Ack(qSharedPointerDynamicCast<txdata::Common2Ack>(theMsg));
        break;
    case txdata::MsgType::ID_Common2Req:
        break;
    case txdata::MsgType::ID_Common2Rsp:
        handle_Common2Rsp(qSharedPointerDynamicCast<txdata::Common2Rsp>(theMsg));
        break;
    case txdata::MsgType::ID_Common1Req:
        handle_Common1Req(qSharedPointerDynamicCast<txdata::Common1Req>(theMsg));
        break;
    case txdata::MsgType::ID_Common1Rsp:
        handle_Common1Rsp(qSharedPointerDynamicCast<txdata::Common1Rsp>(theMsg));
        break;
    case txdata::MsgType::ID_DisconnectedData:
        break;
    case txdata::MsgType::ID_ConnectReq:
        handle_ConnectReq(qSharedPointerDynamicCast<txdata::ConnectReq>(theMsg));
        break;
    case txdata::MsgType::ID_ConnectRsp:
        handle_ConnectRsp(qSharedPointerDynamicCast<txdata::ConnectRsp>(theMsg));
        break;
    case txdata::MsgType::ID_OnlineNotice:
        break;
    case txdata::MsgType::ID_SystemReport:
        break;
    default:
        break;
    }
}

void DataExchanger::slotOnError(QAbstractSocket::SocketError error)
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "DataExchanger::slotOnError " << error;
    {
        QMetaEnum metaEnum = QMetaEnum::fromType<QAbstractSocket::SocketError>();
        const char* errorValue = metaEnum.valueToKey(error);

        QStringList strList;
        strList.append(errorValue);
        strList.append(m_ws.url().toString());
        emit sigStatusError(strList.join('\n'), StatusErrorType::WebsockError);
    }
}
