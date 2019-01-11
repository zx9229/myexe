#include "dataexchanger.h"
#include <QApplication>
#include "m2b.h"

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

    initOwnInfo();
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

void DataExchanger::slotOnConnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnConnected";

    {
        txdata::ConnectedData tmpData = {};
        tmpData.mutable_info()->CopyFrom(m_ownInfo);
        tmpData.add_pathway(tmpData.info().userid());

        m_ws.sendBinaryMessage(m2b::msg2pkg(txdata::ID_ConnectedData, tmpData));
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
    m_ws.sendBinaryMessage(m2b::msg2pkg(txdata::MsgType::ID_ParentDataReq, reqData));
}
