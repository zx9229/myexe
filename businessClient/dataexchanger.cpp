#include "dataexchanger.h"
#include "m2b.h"

DataExchanger::DataExchanger(QObject *parent) :
    QObject(parent),
    m_ws(parent)
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
        printf("%s\n", txData->info().uniqueid().c_str());
    }
}
