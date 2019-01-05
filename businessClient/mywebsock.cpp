#include "mywebsock.h"
#include <QDateTime>
#include <QDebug>

MyWebsock::MyWebsock(QObject *parent /* = Q_NULLPTR */) :
    QObject(parent),
    m_interval(5),
    m_ws(QString(), QWebSocketProtocol::VersionLatest, parent),
    m_timer(parent)
{
    QObject::connect(&m_ws, &QWebSocket::connected, this, &MyWebsock::connected);
    QObject::connect(&m_ws, &QWebSocket::disconnected, this, &MyWebsock::disconnected);
    QObject::connect(&m_ws, &QWebSocket::binaryMessageReceived, this, &MyWebsock::binaryMessageReceived);
    QObject::connect(&m_ws, static_cast<void(QWebSocket::*)(QAbstractSocket::SocketError)>(&QWebSocket::error), this, &MyWebsock::error);
    QObject::connect(&m_ws, &QWebSocket::pong, this, &MyWebsock::pong);

    m_timer.setSingleShot(true);
    QObject::connect(&m_timer, &QTimer::timeout, this, &MyWebsock::reconnect);
}

bool MyWebsock::start(const QUrl& url)
{
    if (!m_url.isEmpty() || url.isEmpty())
        return false;

    m_url = url;
    m_ws.open(m_url);

    return true;
}

void MyWebsock::stop()
{
    m_timer.stop();
    m_ws.abort();
    m_url.clear();
}

qint64 MyWebsock::sendBinaryMessage(const QByteArray &data)
{
    return m_ws.sendBinaryMessage(data);
}

void MyWebsock::reconnect()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::reconnect";
    m_ws.abort();
    m_ws.open(m_url);
}

void MyWebsock::connected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::connected";
    m_timer.stop();
}

void MyWebsock::disconnected()
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::disconnected";

    m_timer.start(m_interval * 1000);
}

void MyWebsock::error(QAbstractSocket::SocketError error)
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::error";
}

void MyWebsock::pong(quint64 elapsedTime, const QByteArray &payload)
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::pong";
}

#include "m2b.h"
void MyWebsock::binaryMessageReceived(const QByteArray &message)
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::binaryMessageReceived";
    onMessage(message);

    {
        txdata::MsgType theType = {};
        GPMSGPTR theMsg;
        m2b::slice2msg(message, theType, theMsg);
        QSharedPointer<txdata::ConnectedData> txData = qSharedPointerDynamicCast<txdata::ConnectedData>(theMsg);
        printf("%s\n", txData->info().uniqueid().c_str());
    }
}
