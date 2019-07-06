#include "dataexchanger.h"
#include <fstream>
#include <QCoreApplication>
#include <QtCore/QMetaEnum>
#include <QDebug>
#include <QJsonDocument>
#include <QJsonObject>
#include <QSqlError>
#include <QSqlQuery>
#include "google/protobuf/util/json_util.h"
#include "zxtools.h"
#include "mytts.h"
#include "m2b.h"

QString qjoGetSet(QJsonObject& clObj, const QStringList& paths, const QString* value)
{
    Q_ASSERT(!paths.empty());
    if (paths.size() == 1)
    {
        if (nullptr == value)
        {
            return clObj.value(paths.at(0)).toString();
        }
        else
        {
            clObj.insert(paths.at(0), *value);
            return *value;
        }
    }
    else
    {
        QJsonObject childObj = clObj.value(paths.at(0)).toObject();
        QString retVal = qjoGetSet(childObj, paths.mid(1), value);
        if (nullptr != value)
        {
            clObj.insert(paths.at(0), childObj);
        }
        return retVal;
    }
}

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
    m_MsgNo("MsgNo", QString()),
    m_lastFind(false),
    m_subUser("ROOT")
{
    connect(&m_ws, &MyWebsock::sigConnected, this, &DataExchanger::slotOnConnected);
    connect(&m_ws, &MyWebsock::sigDisconnected, this, &DataExchanger::slotOnDisconnected);
    connect(&m_ws, &MyWebsock::sigMessage, this, &DataExchanger::slotOnMessage);
    connect(&m_ws, &MyWebsock::sigError, this, &DataExchanger::slotOnError);

    initOwnInfo();
    initDB();
    initTypeToJson();
    m_startDateTime = QDateTime::currentDateTime();
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

QString DataExchanger::memGetData(const QString& varName)
{
    if (varName == "MsgNo")
        return m_MsgNo.Value;
    if (varName == "PathwayInfo")
        return m_PathwayInfo ? zxtools::object2json(*m_PathwayInfo, nullptr) : "{}";
    return varName.isEmpty() ? "" : "";
}

bool DataExchanger::memSetData(const QString& varName, const QString& value)
{
    if (varName == "url")
    {
        m_url = value;  // 例如【ws://localhost:65535/websocket】.
        return true;
    }
    else
    {
        return false;
    }
}

QString DataExchanger::memGetInfo(const QString& varName, const QStringList& paths)
{
    txdata::ConnectionInfo* pCI = nullptr;
    if (varName == "myself")
    {
        pCI = &m_ownInfo;
    }
    else if (varName == "parent")
    {
        pCI = &m_parentInfo;
    }
    else
    {
        return "";
    }
    bool isOk = false;
    QString jsonText = zxtools::object2json(*pCI, &isOk);
    Q_ASSERT(isOk);
    QByteArray jsonTextBA = jsonText.toLocal8Bit();
    //
    QJsonParseError jsonParseError;
    QJsonDocument jsonDoc = QJsonDocument::fromJson(jsonTextBA, &jsonParseError);
    Q_ASSERT(jsonParseError.error == QJsonParseError::NoError);
    QJsonObject jsonRoot = jsonDoc.object();
    return  qjoGetSet(jsonRoot, paths, nullptr);
}

bool DataExchanger::memSetInfo(const QString& varName, const QStringList& paths, const QString& value)
{
    //auto filename2content = [](const std::string& filename, std::string& content)
    //{
    //    std::copy(std::istreambuf_iterator<char>(std::ifstream(filename).rdbuf()), std::istreambuf_iterator<char>(), std::back_inserter(content));
    //};
    txdata::ConnectionInfo* pCI = nullptr;
    if (varName == "myself")
    {
        pCI = &m_ownInfo;
    }
    else if (varName == "parent")
    {
        pCI = &m_parentInfo;
    }
    else
    {
        return false;
    }
    bool isOk = false;
    QString jsonText = zxtools::object2json(*pCI, &isOk);
    Q_ASSERT(isOk);
    QByteArray jsonTextBA = jsonText.toLocal8Bit();
    //
    QJsonParseError jsonParseError;
    QJsonDocument jsonDoc = QJsonDocument::fromJson(jsonTextBA, &jsonParseError);
    Q_ASSERT(jsonParseError.error == QJsonParseError::NoError);
    QJsonObject jsonRoot = jsonDoc.object();
    qjoGetSet(jsonRoot, paths, &value);
    jsonDoc.setObject(jsonRoot);
    //
    jsonTextBA = jsonDoc.toJson();
    jsonText = QString::fromLocal8Bit(jsonTextBA);
    //
    QString typeName = m2b::CalcMsgTypeName(m_ownInfo);
    GPMSGPTR gpmsgptr = zxtools::json2object(typeName, jsonText);
    QSharedPointer<txdata::ConnectionInfo> spCI = qSharedPointerDynamicCast<txdata::ConnectionInfo>(gpmsgptr);
    *pCI = *spCI;
    return true;
}

QString DataExchanger::serviceState()
{
    QMetaEnum metaEnum = QMetaEnum::fromType<QAbstractSocket::SocketState>();
    return  metaEnum.valueToKey(m_ws.state());
}

bool DataExchanger::start()
{
    m_ws.stop(true);
    return m_ws.start(m_url);
}

QString DataExchanger::sendReq(const QString &typeName, const QString &jsonText, const QString &rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB)
{
    QMap<QString,QString> kvMap;
    {
        kvMap["ReqType"]=typeName;
        kvMap["ReqData"]=jsonText;
        kvMap["RecverID"]=rID;
        kvMap["IsLog"]=isLog;
        kvMap["IsSafe"]=isSafe;
        kvMap["IsPush"]=isPush;
        kvMap["UpCache"]=isUpCache;
    }
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

QStringList DataExchanger::getTxMsgTypeNameList()
{
    QStringList typeNameList;
    for (auto it = m_typeTOjson.begin(); it != m_typeTOjson.end(); ++it)
    {
        typeNameList << QString::fromStdString(::txdata::MsgType_Name(it.key()));
    }
    return  typeNameList;
}

QString DataExchanger::jsonExample(const QString& typeName)
{
    txdata::MsgType msgType = txdata::MsgType::Zero1;
    txdata::MsgType_Parse(typeName.toStdString(), &msgType);
    auto it = m_typeTOjson.find(msgType);
    return (m_typeTOjson.end() == it) ? "" : it.value();
}

bool DataExchanger::deleteCommonData1(const QString& userid, qint64 msgno)
{
    QSqlQuery sqlQuery;
    QString whereCond = QString("UserID='%1' AND MsgNo='%2'").arg(userid).arg(QString::number(msgno));
    bool isOk = CommonData::delete_data(sqlQuery, whereCond);
    if (isOk) { emit sigTableChanged(CommonData::static_table_name()); }
    return isOk;
}

bool DataExchanger::deleteCommonData2(const QString& userid, qint64 msgno, int seqno)
{
    QSqlQuery sqlQuery;
    QString whereCond = QString("UserID='%1' AND MsgNo='%2' AND SeqNo='%3'").arg(userid).arg(QString::number(msgno)).arg(QString::number(seqno));
    qDebug() << whereCond;
    bool isOk = CommonData::delete_data(sqlQuery, whereCond);
    if (isOk) { emit sigTableChanged(CommonData::static_table_name()); }
    return isOk;
}

bool DataExchanger::deletePushWrap(const QString& userid, const QString& peerid, qint64 msgno)
{
    QSqlQuery sqlQuery;
    QString whereCond = QString("UserID='%1' AND PeerID='%2' AND MsgNo='%3'").arg(userid).arg(peerid).arg(QString::number(msgno));
    bool isOk = PushWrap::delete_data(sqlQuery, whereCond);
    if (isOk) { emit sigTableChanged(PushWrap::static_table_name()); }
    return isOk;
}

QString DataExchanger::serviceInfo()
{
    return m_startDateTime.toString("yyyy-MM-dd hh:mm:ss");
}

QString DataExchanger::sendCommonReq(const QStringList& kvs, bool isC1NotC2)
{
    QString message;
    do
    {
        QMap<QString, QString> kvMap = zxtools::fromOddEven(kvs);
        if (kvMap.isEmpty())
        {
            message = "param illegal";
            break;
        }
        CommonData tmpCommonData;
        if (isC1NotC2)
        {
            QSharedPointer<txdata::Common1Req> c1req = QSharedPointer<txdata::Common1Req>(new txdata::Common1Req);
            message = zxtools::map2Common1Req(c1req.get(), kvMap);
            if (!message.isEmpty())
                break;
            message = sendCommon1Req(c1req);
            if (!message.isEmpty())
                break;
            zxtools::Common1Req2CommonData(&tmpCommonData, c1req.get());
        }
        else
        {
            QSharedPointer<txdata::Common2Req> c2req = QSharedPointer<txdata::Common2Req>(new txdata::Common2Req);
            message = zxtools::map2Common2Req(c2req.get(), kvMap);
            if (!message.isEmpty())
                break;
            message = sendCommon2Req(c2req);
            if (!message.isEmpty())
                break;
            zxtools::Common2Req2CommonData(&tmpCommonData, c2req.get());
        }
        tmpCommonData.SendOK = message.isEmpty();
        if (tmpCommonData.SendOK && tmpCommonData.MsgNo != 0)
        {
            QSqlQuery sqlQuery;
            tmpCommonData.insert_data(sqlQuery, true, nullptr);
            if (m_MsgNo.Value.toLongLong() < tmpCommonData.MsgNo)
            {
                m_MsgNo.Value.setNum(tmpCommonData.MsgNo);
                m_MsgNo.update_data(sqlQuery);
            }
            emit sigTableChanged(tmpCommonData.static_table_name());
        }
    } while (false);
    return message;
}

void DataExchanger::initDB()
{
    m_db = QSqlDatabase::addDatabase("QSQLITE");
    m_db.setDatabaseName(QString().isEmpty() ? (SQLITE_DB_NAME) : (":memory:"));
    bool isOk = false;
    isOk = m_db.open();
    Q_ASSERT(isOk);
    QSqlQuery sqlQuery;
    if (true) {
        isOk = m_db.transaction();
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(KeyValue::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(CommonData::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(PushWrap::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(PathwayInfo::static_create_sql());
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

void DataExchanger::initTypeToJson()
{
    if (true) {
        txdata::EchoItem tmpObj;
        tmpObj.set_data("DATA");
        tmpObj.set_rspcnt(2);
        tmpObj.set_secgap(0);
        m_typeTOjson[m2b::CalcMsgType(tmpObj)] = zxtools::object2json(tmpObj);
    }
    if (true) {
        txdata::PushItem tmpObj;
        tmpObj.set_subject("Subject");
        tmpObj.set_content("Content");
        tmpObj.add_modes("tts");
        m_typeTOjson[m2b::CalcMsgType(tmpObj)] = zxtools::object2json(tmpObj);
    }
    if (true) {
        txdata::SubscribeReq tmpObj;
        m_typeTOjson[m2b::CalcMsgType(tmpObj)] = zxtools::object2json(tmpObj);
    }
    if (true) {
        txdata::QrySubscribeReq tmpObj;
        m_typeTOjson[m2b::CalcMsgType(tmpObj)] = zxtools::object2json(tmpObj);
    }
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
    if (msgData->issafe())
    {
        QSharedPointer<txdata::Common2Ack> ackOut;
        genAck4Common2Rsp(msgData, ackOut);
        m_ws.sendBinaryMessage(m2b::msg2package(*ackOut));
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
            QSharedPointer<txdata::PushWrap> txdataPushWrap = qSharedPointerDynamicCast<txdata::PushWrap>(tmpData);
            PushWrap pshData;
            zxtools::PushWrap2PushWrap(&pshData, txdataPushWrap.get(), msgData->senderid());
            QSqlQuery sqlQuery;
            bool isOk = pshData.insert_data(sqlQuery, false, nullptr);
            Q_ASSERT(isOk);
            emit sigTableChanged(pshData.static_table_name());
            {
                QString text;
                if (zxtools::needTTS(&pshData, text))
                {
                    MyTTS::staticSpeak(text);
                }
            }
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

    GPMSGPTR innerRspData;
    m2b::slice2msg(msgData->rspdata(), msgData->rsptype(), innerRspData);
    switch (msgData->rsptype()) {
    case txdata::ID_QryConnInfoRsp:
        break;
    case txdata::ID_SubscribeRsp:
        handle_SubscribeRsp(qSharedPointerDynamicCast<txdata::SubscribeRsp>(innerRspData));
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
        if (data->inforeq().userid().empty())
        {
            data4send.set_errmsg("(req.UserID == EMPTYSTR)");
            break;
        }
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

void DataExchanger::handle_PathwayInfo(QSharedPointer<txdata::PathwayInfo> data)
{
    Q_ASSERT(data.data() != nullptr);
    m_PathwayInfo = data;
    QSqlQuery sqlQuery;
    bool curFind = (data->info().find(m_subUser.toStdString()) != data->info().end());
    if (m_lastFind != curFind)
    {
        if (curFind)
        {
            QString whereCond = QString("PeerID='%1' ORDER BY MsgNo DESC LIMIT 1").arg(m_subUser);
            QList<PushWrap> dataList;
            bool isOk = PushWrap::select_data(sqlQuery, whereCond, dataList);
            Q_ASSERT(isOk);
            txdata::SubscribeReq tmpData;
            tmpData.set_frommsgno(dataList.empty() ? 0 : dataList[0].MsgNo);
            QString jsonText = zxtools::object2json(tmpData, &isOk);
            Q_ASSERT(isOk);
            QString typeName = m2b::CalcMsgTypeName(tmpData);
            sendReq(typeName, jsonText, m_subUser, false, false, false, false, true, false, false);
            qDebug() << m2b::CalcMsgTypeName(tmpData) << QString::fromStdString(tmpData.DebugString());
        }
        m_lastFind = curFind;
    }
    PathwayInfo::delete_data(sqlQuery, "");
    for (auto it = data->info().begin(); it != data->info().end(); ++it)
    {
        PathwayInfo tmpData;
        QStringList pathway;
        for (int i = 0; i < it->second.data_size(); i++) { pathway.append(QString::fromStdString(it->second.data(i))); }
        tmpData.UserID = QString::fromStdString(it->first);
        tmpData.Pathway = pathway.join("->");
        tmpData.insert_data(sqlQuery, false);
    }
    emit sigTableChanged(PathwayInfo::static_table_name());
}

void DataExchanger::handle_SubscribeRsp(QSharedPointer<txdata::SubscribeRsp> data)
{
    Q_ASSERT(data.data() != nullptr);
    qDebug() << m2b::CalcMsgTypeName(*data) << QString::fromStdString(data->DebugString());
}

void DataExchanger::genAck4Common2Rsp(QSharedPointer<txdata::Common2Rsp> rspIn, QSharedPointer<txdata::Common2Ack>& ackOut)
{
    Q_ASSERT(rspIn->issafe() == true);
    Q_ASSERT(rspIn->ispush() == false);
    ackOut = QSharedPointer<txdata::Common2Ack>(new txdata::Common2Ack);
    ackOut->mutable_key()->set_userid(rspIn->key().userid());
    ackOut->mutable_key()->set_msgno(rspIn->key().msgno());
    ackOut->mutable_key()->set_seqno(rspIn->key().seqno());
    ackOut->set_senderid(this->m_ownInfo.userid());
    ackOut->set_recverid(rspIn->senderid());
    ackOut->set_toroot(!rspIn->toroot());
    ackOut->set_islog(rspIn->islog());
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
    case txdata::MsgType::ID_PathwayInfo:
        handle_PathwayInfo(qSharedPointerDynamicCast<txdata::PathwayInfo>(theMsg));
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
