#include "mywebsock.h"
#include <QDateTime>
#include <QDebug>
#include <QThread>

MyWebsock::MyWebsock(QObject *parent /* = Q_NULLPTR */) :
    QObject(parent),
    m_interval(5),
    m_timer(parent),
    m_url(QUrl()),
    m_ws(QString(), QWebSocketProtocol::VersionLatest, parent),
    m_alive(false)
{
    QObject::connect(&m_ws, &QWebSocket::connected, this, &MyWebsock::slotConnected);
    QObject::connect(&m_ws, &QWebSocket::disconnected, this, &MyWebsock::slotDisconnected);
    QObject::connect(&m_ws, &QWebSocket::binaryMessageReceived, this, &MyWebsock::slotBinaryMessageReceived);
    QObject::connect(&m_ws, static_cast<void(QWebSocket::*)(QAbstractSocket::SocketError)>(&QWebSocket::error), this, &MyWebsock::slotError);
    QObject::connect(&m_ws, &QWebSocket::pong, this, &MyWebsock::slotPong);

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

void MyWebsock::interrupt()
{
    m_ws.abort(); // 类似于"拔网线"的操作, 破坏连接, 但是不破坏重连机制.
}

void MyWebsock::stop(bool sync /* = false */)
{
    m_timer.stop();
    m_ws.abort();
    m_url.clear();
    // 让回调函数重置m_alive的值.
    while (sync && m_alive)
    {
        QThread::msleep(1);
    }
}

qint64 MyWebsock::sendBinaryMessage(const QByteArray &data)
{
    if (!m_alive) { return 0; }
    qint64 sendBytes = m_ws.sendBinaryMessage(data);
    Q_ASSERT(data.size() == sendBytes);
    return sendBytes;
}

void MyWebsock::reconnect()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::reconnect";

    m_ws.abort();
    m_ws.open(m_url);
}

void MyWebsock::slotConnected()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::slotConnected";

    m_timer.stop();
    Q_ASSERT(m_alive == false);
    m_alive = true;

    emit sigConnected();
}

void MyWebsock::slotDisconnected()
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::slotDisconnected";

    if (m_alive)
    {
        m_alive = false;
        emit sigDisconnected();
    }

    Q_ASSERT(m_alive == false);

    m_timer.start(m_interval * 1000);
}

void MyWebsock::slotBinaryMessageReceived(const QByteArray &message)
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::slotBinaryMessageReceived";

    emit sigMessage(message);
}

void MyWebsock::slotError(QAbstractSocket::SocketError error)
{
    //qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::slotError "<< error;

    sigError(error);
}

void MyWebsock::slotPong(quint64 elapsedTime, const QByteArray &payload)
{
    qDebug() << QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss") << "MyWebsock::slotPong";
}
