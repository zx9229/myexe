#include "dataexchanger.h"
#include <QApplication>
#include "m2b.h"

DataExchanger::DataExchanger(QObject *parent) :
    QObject(parent),
    m_ws(parent),
    m_totalPos(3)
{
    connect(&m_ws, &MyWebsock::sigConnected, this, &DataExchanger::slotOnConnected);
    connect(&m_ws, &MyWebsock::sigDisconnected, this, &DataExchanger::slotOnDisconnected);
    connect(&m_ws, &MyWebsock::sigMessage, this, &DataExchanger::slotOnMessage);
}

DataExchanger::~DataExchanger()
{

}

QString DataExchanger::Login(const QString &url, const QString &username, const QString &password)
{
    QString message;

    m_ws.start(url);

    return message;
}

MyWebsock& DataExchanger::ws()
{
    return m_ws;
}

void DataExchanger::slotOnConnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnConnected";

    {
        txdata::ConnectedData tmpData = {};
        {
            tmpData.mutable_info()->set_userid("");
            tmpData.mutable_info()->set_belongid("");
            tmpData.mutable_info()->set_version("20190106");
            tmpData.mutable_info()->set_exetype(txdata::ProgramType::CLIENT);
            tmpData.mutable_info()->set_isleaf(false);//(当AppType/ExeType为NODE时,表示NODE是否为LeafNode)
            tmpData.mutable_info()->set_linkmode(txdata::ConnectionInfo_LinkType_CONNECT);
            tmpData.mutable_info()->set_exepid(static_cast<int>(QCoreApplication::applicationPid()));
            tmpData.mutable_info()->set_exepath(QCoreApplication::applicationFilePath().toStdString());
        }
        QByteArray data;
        m2b::msg2slice(txdata::ID_ConnectedData, tmpData, data);
        //
        m_ws.sendBinaryMessage(data);
    }
}

void DataExchanger::slotOnDisconnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnDisconnected";
}

void DataExchanger::slotOnMessage(const QByteArray &message)
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "slotOnMessage";

    if (true) {
        txdata::MsgType theType = {};
        GPMSGPTR theMsg;
        m2b::slice2msg(message, theType, theMsg);
        QSharedPointer<txdata::ConnectedData> txData = qSharedPointerDynamicCast<txdata::ConnectedData>(theMsg);
        printf("%s\n", txData->info().userid().c_str());
    }
}
