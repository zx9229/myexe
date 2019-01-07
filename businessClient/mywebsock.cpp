#include "mywebsock.h"
#include <QDateTime>
#include <QDebug>

MyWebsock::MyWebsock(QObject *parent /* = Q_NULLPTR */) :
    QObject(parent),
    m_interval(5),
    m_timer(parent),
    m_url(QUrl()),
    m_ws(QString(), QWebSocketProtocol::VersionLatest, parent),
    m_alive(false)
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
    if (!m_url.isEmpty() || url.isEmpty() || m_alive)
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
    //让回调函数重置m_alive的值.
}

qint64 MyWebsock::sendBinaryMessage(const QByteArray &data)
{
    qint64 sendBytes = m_ws.sendBinaryMessage(data);
    Q_ASSERT(data.size()== sendBytes);
    return sendBytes;
}

void MyWebsock::reconnect()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::reconnect";

    m_ws.abort();
    m_ws.open(m_url);
}

void MyWebsock::connected()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::connected";

    m_timer.stop();
    Q_ASSERT(m_alive == false);
    m_alive = true;

    sigConnected();
}

void MyWebsock::disconnected()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::disconnected";

    if (m_alive)
    {
        m_alive = false;
        sigDisconnected();
    }

    Q_ASSERT(m_alive == false);

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

void MyWebsock::binaryMessageReceived(const QByteArray &message)
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::binaryMessageReceived";

    sigMessage(message);
}
