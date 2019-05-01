#include "dataexchanger.h"
#include <QCoreApplication>
#include <QSqlError>
#include "m2b.h"
#include "google/protobuf/util/json_util.h"
// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Timestamp
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

QString DataExchanger::jsonByMsgObje(const google::protobuf::Message &msgObj, bool *isOk)
{
    if (isOk) { *isOk = true; }
    google::protobuf::util::JsonOptions jsonOpt;
    if (true) {
        jsonOpt.add_whitespace = true;
        jsonOpt.always_print_primitive_fields = true;
        jsonOpt.preserve_proto_field_names = true;
    }
    std::string jsonStr;
    if (google::protobuf::util::MessageToJsonString(msgObj, &jsonStr, jsonOpt) != google::protobuf::util::Status::OK)
    {
        jsonStr.clear();
        if (isOk) { *isOk = false; }
    }
    return QString::fromStdString(jsonStr);
}

QString DataExchanger::nameByMsgType(txdata::MsgType msgType, int flag, bool* isOk)
{
    //flag=0  结果类似(txdata.ConnectionInfo)
    //flag=1  结果类似(ConnectionInfo)
    if (isOk) { *isOk = false; }
    QString retName;
    do
    {
        if (txdata::MsgType_IsValid(msgType) == false)
            break;
        std::string name = txdata::MsgType_Name(msgType);
        const std::string HEAD("ID_");
        std::size_t pos = name.find(HEAD);
        if (std::string::npos == pos)
            break;
        name = name.substr(HEAD.size() + pos);
        if (flag == 0)
            retName = "txdata." + QString::fromStdString(name);
        else if (flag == 1)
            retName = QString::fromStdString(name);
        else
            break;
        if (isOk) { *isOk = true; }
    } while (false);
    return retName;
}

QString DataExchanger::jsonByMsgType(txdata::MsgType msgType, const QByteArray &serializedData, bool *isOk)
{
    QString jsonStr;

    std::string name = nameByMsgType(msgType, 0, isOk).toStdString();
    if (name.empty())
        return jsonStr;

    const google::protobuf::Descriptor* curDesc = google::protobuf::DescriptorPool::generated_pool()->FindMessageTypeByName(name);
    if (nullptr == curDesc)
        return jsonStr;
    google::protobuf::Message* curMesg = google::protobuf::MessageFactory::generated_factory()->GetPrototype(curDesc)->New();
    if (nullptr == curMesg)
        return jsonStr;

    QSharedPointer<google::protobuf::Message> guard(curMesg);

    if (curMesg->ParseFromArray(serializedData.constData(), serializedData.size()) == false)
        return jsonStr;

    jsonStr = jsonByMsgObje(*curMesg, isOk);

    return jsonStr;
}

bool DataExchanger::calcObjByName(const QString &typeName, QSharedPointer<google::protobuf::Message> &objOut)
{
    objOut.clear();
    // https://blog.csdn.net/riopho/article/details/80372510
    const google::protobuf::Descriptor* desc = google::protobuf::DescriptorPool::generated_pool()->FindMessageTypeByName(typeName.toStdString());
    if (nullptr == desc) { return false; }
    // desc->index();
    google::protobuf::Message* mesg = google::protobuf::MessageFactory::generated_factory()->GetPrototype(desc)->New();
    objOut.reset(mesg);
    return (mesg ? true : false);
}

QString DataExchanger::jsonToObjAndS(const QString &typeName, const QString &jsonStr, txdata::MsgType &msgType, QByteArray &serializedData)
{
    QString message;
    //
    msgType = txdata::MsgType::Zero1;
    serializedData.clear();
    do
    {
        QSharedPointer<google::protobuf::Message> curObj;
        if (calcObjByName(typeName, curObj) == false)
        {
            message = "calc object by name fail";
            break;
        }
        if (google::protobuf::util::JsonStringToMessage(jsonStr.toStdString(), curObj.data()) != google::protobuf::util::Status::OK)
        {
            message = "fill object by json fail";
            break;
        }
        std::string binData;
        if (curObj->SerializeToString(&binData) == false)
        {
            message = "serialize object fail";
            break;
        }
        serializedData.append(binData.data(), binData.size());
        msgType = static_cast<txdata::MsgType>(curObj->GetDescriptor()->index());
    } while (false);
    //
    return message;
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

bool DataExchanger::start()
{
    m_ws.stop(true);
    return m_ws.start(m_url);
}

void DataExchanger::initDB()
{
    m_db = QSqlDatabase::addDatabase("QSQLITE");
    m_db.setDatabaseName(false ? (":memory:") : ("_zx_test.db"));
}

void DataExchanger::initOwnInfo()
{
    m_ownInfo.set_version("Ver20190108");
    m_ownInfo.set_linkmode(txdata::ConnectionInfo_LinkType_CONNECT);
    m_ownInfo.set_exepid(static_cast<int>(QCoreApplication::applicationPid()));
    m_ownInfo.set_exepath(QCoreApplication::applicationFilePath().toStdString());
    m_ownInfo.set_remark("");
}

void DataExchanger::handle_MessageAck(QSharedPointer<txdata::MessageAck> data)
{
    Q_ASSERT(data.data() != nullptr);
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << QString::fromStdString(data->GetTypeName());
    //TODO:
}

void DataExchanger::handle_ConnectReq(QSharedPointer<txdata::ConnectReq> data)
{
    Q_ASSERT(data.data() != nullptr);
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << QString::fromStdString(data->GetTypeName());

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
    }
    else
    {
        m_ws.interrupt();
    }
}

void DataExchanger::handle_ConnectRsp(QSharedPointer<txdata::ConnectRsp> data)
{
    Q_ASSERT(data.data() != nullptr);
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << QString::fromStdString(data->GetTypeName());
    //TODO:
}

void DataExchanger::slotOnConnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnConnected";

    {
        txdata::ConnectReq data4send = {};
        data4send.mutable_inforeq()->CopyFrom(m_ownInfo);
        data4send.add_pathway(data4send.inforeq().userid());

        m_ws.sendBinaryMessage(m2b::msg2package(data4send));
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
    if (m2b::package2msg(message, theType, theMsg) == false)
    {
        qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnMessage, slice2msg, failure";
        return;
    }
    switch (theType) {
    case txdata::MsgType::ID_MessageAck:
        handle_MessageAck(qSharedPointerDynamicCast<txdata::MessageAck>(theMsg));
        break;
    case txdata::MsgType::ID_CommonReq:
        break;
    case txdata::MsgType::ID_CommonRsp:
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

    sigWebsocketError(error);
}
