#include "dataexchanger.h"
#include <QApplication>
#include <QSqlError>
#include "m2b.h"
#include "sqlstruct.h"

std::string atomicKey2str(const txdata::AtomicKey& src)
{
    return QString("/%1/%2/%3/%4")
        .arg(src.zonename().data())
        .arg(src.nodename().data())
        .arg(src.exectype())
        .arg(src.execname().data())
        .toStdString();
}

DataExchanger::DataExchanger(QObject *parent) :
    QObject(parent),
    m_ws(parent)
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

MyWebsock& DataExchanger::ws()
{
    return m_ws;
}

bool DataExchanger::start()
{
    m_ws.stop(true);
    return m_ws.start(m_url);
}

void DataExchanger::setURL(const QString &url)
{
    m_url = url;
}

void DataExchanger::setUserKey(const QString &zoneName, const QString &nodeName, txdata::ProgramType execType, const QString &execName)
{
    m_ownInfo.mutable_userkey()->set_zonename(zoneName.toStdString());
    m_ownInfo.mutable_userkey()->set_nodename(nodeName.toStdString());
    m_ownInfo.mutable_userkey()->set_exectype(execType);
    m_ownInfo.mutable_userkey()->set_execname(execName.toStdString());

    m_ownInfo.set_userid(atomicKey2str(m_ownInfo.userkey()));
}

void DataExchanger::setBelongKey(const QString &zoneName, const QString &nodeName, txdata::ProgramType execType, const QString &execName)
{
    m_ownInfo.mutable_belongkey()->set_zonename(zoneName.toStdString());
    m_ownInfo.mutable_belongkey()->set_nodename(nodeName.toStdString());
    m_ownInfo.mutable_belongkey()->set_exectype(execType);
    m_ownInfo.mutable_belongkey()->set_execname(execName.toStdString());

    m_ownInfo.set_belongid(atomicKey2str(m_ownInfo.belongkey()));
}

bool DataExchanger::sendCommonNtosReq(QCommonNtosReq& reqData, bool needResp, bool needSave)
{
    int64_t lastInsertId = 0;
    reqData.UserID = QString::fromStdString(m_ownInfo.userid());

    bool opFinish = false;
    int64_t iRequestID = 0;
    int64_t iSeqNo = 0;
    if (needResp)
    {
        iRequestID = m_RequestID.Value.toLongLong();
        reqData.RequestID = ++iRequestID;
    }
    if (needSave)
    {
        iSeqNo = m_SeqNo.Value.toLongLong();
        reqData.SeqNo = ++iSeqNo;
    }
    if (true)
    {
        bool isOk = false;
        QSqlQuery sqlQuery;
        do
        {
            opFinish = m_db.transaction();//临时用(操作结束的标志)当做(事务开启成功的标志)
            if (!opFinish) { break; }
            if (iRequestID)
            {
                m_RequestID.Value.number(iRequestID);
                isOk = m_RequestID.update_data(sqlQuery);
                if (!isOk) { break; }
            }
            if (iSeqNo)
            {
                m_SeqNo.Value.number(iSeqNo);
                isOk = m_SeqNo.update_data(sqlQuery);
                if (!isOk) { break; }
            }
            reqData.ReqTime = QDateTime::currentDateTime();
            isOk = reqData.insert_data(sqlQuery, &lastInsertId);
            if (!isOk) { break; }
            isOk = m_db.commit();
            if (!isOk) { break; }
        } while (false);
        if (!isOk)//如果操作数据库失败.
        {
            reqData.ReqTime = QDateTime();
            qDebug() << sqlQuery.lastError();
            if (opFinish)//如果开启了事务,就需要回滚.
            {
                isOk = m_db.rollback();
                Q_ASSERT(isOk);//如果回滚失败,那我也没有办法了.
                opFinish = false;
            }
        }
    }
    if (opFinish)
    {
        reqData.RefNum = lastInsertId;
        txdata::CommonNtosReq data4send;
        toCommonNtosReq(reqData, data4send);
        m_ws.sendBinaryMessage(m2b::msg2pkg(data4send));
    }
    return opFinish;
}

void DataExchanger::initDB()
{
    m_db = QSqlDatabase::addDatabase("QSQLITE");
    m_db.setDatabaseName(false ? (":memory:") : ("_zx_test.db"));
    bool isOk = false;
    isOk = m_db.open();
    Q_ASSERT(isOk);
    QSqlQuery sqlQuery;
    if (true) {
        isOk = m_db.transaction();
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(KeyValue::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(QCommonNtosReq::static_create_sql());
        Q_ASSERT(isOk);
        isOk = sqlQuery.exec(QCommonNtosRsp::static_create_sql());
        Q_ASSERT(isOk);
        isOk = m_db.commit();
        Q_ASSERT(isOk);
    }
    if (true) {
        QList<KeyValue> kvList;
        KeyValue::select_data(sqlQuery, kvList);
        for (auto&node : kvList)
        {
            if (node.Key == "RequestID") { m_RequestID = node; }
            if (node.Key == "SeqNo") { m_SeqNo = node; }
        }
        if (m_RequestID.Key.isEmpty())
        {
            m_RequestID = { "RequestID","0" };
            m_RequestID.insert_data(sqlQuery, true);
        }
        if (m_SeqNo.Key.isEmpty())
        {
            m_SeqNo = { "SeqNo","0" };
            m_SeqNo.insert_data(sqlQuery, true);
        }
    }
}

void DataExchanger::initOwnInfo()
{
    m_ownInfo.mutable_userkey();
    m_ownInfo.set_userid(atomicKey2str(*m_ownInfo.mutable_userkey()));
    m_ownInfo.mutable_belongkey();
    m_ownInfo.set_belongid(atomicKey2str(*m_ownInfo.mutable_belongkey()));
    m_ownInfo.set_version("Ver20190108");
    m_ownInfo.set_linkmode(txdata::ConnectionInfo_LinkType_CONNECT);
    m_ownInfo.set_exepid(static_cast<int>(QCoreApplication::applicationPid()));
    m_ownInfo.set_exepath(QCoreApplication::applicationFilePath().toStdString());
    m_ownInfo.set_remark("");
}

void DataExchanger::handle_ConnectedData(QSharedPointer<txdata::ConnectedData> data)
{
    Q_ASSERT(data.data() != nullptr);

    bool checkOK = false;

    do
    {
        if (data->info().userid() != m_ownInfo.belongid())
            break;
        if (atomicKey2str(data->info().userkey()) != data->info().userid())
            break;
        if (atomicKey2str(data->info().belongkey()) != data->info().belongid())
            break;
        if (data->pathway_size() != 1)
            break;
        if (data->pathway(0) != data->info().userid())
            break;
        checkOK = true;
    } while (false);

    if (checkOK)
    {
        m_parentInfo.CopyFrom(data->info());
        emit sigReady();
    }
    else
    {
        m_ws.interrupt();
    }
}

void DataExchanger::handle_CommonNtosRsp(QSharedPointer<txdata::CommonNtosRsp> data)
{
    Q_ASSERT(data.data() != nullptr);
    Q_ASSERT(data->pathway_size() == 1);
    Q_ASSERT(data->pathway(0) == m_ownInfo.userid());
    QCommonNtosRsp rspData;
    toCommonNtosRsp(*data, rspData);
    rspData.insert_data(QSqlQuery());
}

void DataExchanger::handle_ParentDataRsp(QSharedPointer<txdata::ParentDataRsp> data)
{
    m_parentData.clear();
    for (int i = 0; i < data->data_size(); ++i)
    {
        QConnInfoEx curData;
        curData.from_txdata_ConnectedData(data->data(i));
        m_parentData.insert(curData.UserID, curData);
    }
    sigParentData(m_parentData);
}

void DataExchanger::toCommonNtosReq(QCommonNtosReq &src, txdata::CommonNtosReq &dst)
{
    dst.set_requestid(src.RequestID);
    dst.set_userid(src.UserID.toStdString());
    dst.set_seqno(src.SeqNo);
    dst.set_reqtype(static_cast<txdata::MsgType>(src.ReqType));
    dst.set_reqdata(src.ReqData.data(), src.ReqData.size());
    if (src.ReqTime.isValid())
    {
        dst.mutable_reqtime()->set_seconds(src.ReqTime.offsetFromUtc());
        dst.mutable_reqtime()->set_nanos(src.ReqTime.time().msec() * 1000 * 1000);
    }
    dst.set_refnum(src.RefNum);
}

void DataExchanger::toCommonNtosRsp(txdata::CommonNtosRsp &src, QCommonNtosRsp &dst)
{
    //dst.ID = INT64_MAX;
    //dst.InsertTime = QDateTime::currentDateTime();
    dst.RequestID = src.requestid();
    for (int i = 0; i < src.pathway_size(); ++i) dst.Pathway.append(src.pathway(i).c_str());
    dst.SeqNo = src.seqno();
    dst.RspType = src.rsptype();
    dst.RspData.append(src.rspdata().data(), src.rspdata().size());
    dst.RspTime.fromTime_t(src.rsptime().seconds());
    dst.FromServer = src.fromserver() ? 1 : 0;
    dst.ErrNo = src.errno();
    dst.ErrMsg = QString::fromStdString(src.errmsg());
    dst.RefNum = src.refnum();
}

void DataExchanger::slotOnConnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnConnected";

    {
        txdata::ConnectedData data4send = {};
        data4send.mutable_info()->CopyFrom(m_ownInfo);
        data4send.add_pathway(data4send.info().userid());

        m_ws.sendBinaryMessage(m2b::msg2pkg(data4send));
    }
}

void DataExchanger::slotOnDisconnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnDisconnected";
}

void DataExchanger::slotOnMessage(const QByteArray &message)
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnMessage";

    txdata::MsgType theType = {};
    GPMSGPTR theMsg;
    if (m2b::pkg2msg(message, theType, theMsg) == false)
    {
        qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnMessage, slice2msg, failure";
        return;
    }

    if (theType == txdata::MsgType::ID_ConnectedData)
    {
        handle_ConnectedData(qSharedPointerDynamicCast<txdata::ConnectedData>(theMsg));
    }
    else if (theType == txdata::MsgType::ID_CommonNtosRsp)
    {
        handle_CommonNtosRsp(qSharedPointerDynamicCast<txdata::CommonNtosRsp>(theMsg));
    }
    else if (theType == txdata::MsgType::ID_ParentDataRsp)
    {
        handle_ParentDataRsp(qSharedPointerDynamicCast<txdata::ParentDataRsp>(theMsg));
    }
}

void DataExchanger::slotParentDataReq()
{
    txdata::ParentDataReq reqData;
    reqData.set_requestid(0);
    // https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Timestamp
    reqData.mutable_reqtime()->set_seconds(time(NULL));
    reqData.mutable_reqtime()->set_nanos(0);
    m_ws.sendBinaryMessage(m2b::msg2pkg(reqData));
}

void DataExchanger::slotOnError(QAbstractSocket::SocketError error)
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "DataExchanger::slotOnError " << error;

    sigWebsocketError(error);
}
