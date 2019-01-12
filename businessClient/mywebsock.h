#ifndef MY_WEBSOCK_H
#define MY_WEBSOCK_H

#include <QObject>
#include <QTimer>          //(QT += core)
#include <QAbstractSocket> //(QT += network)
#include <QWebSocket>      //(QT += websockets)

class MyWebsock :public QObject
{
    Q_OBJECT

public:
    explicit MyWebsock(QObject *parent = Q_NULLPTR);

public:
    bool start(const QUrl& url);
    void interrupt();
    void stop(bool sync = false);
    qint64 sendBinaryMessage(const QByteArray &data);

signals:
    void sigConnected();
    void sigDisconnected();
    void sigMessage(const QByteArray &message);
    void sigError(QAbstractSocket::SocketError error);

private:
    void reconnect();

private slots:
    void slotConnected();
    void slotDisconnected();
    void slotBinaryMessageReceived(const QByteArray &message);
    void slotError(QAbstractSocket::SocketError error);
    void slotPong(quint64 elapsedTime, const QByteArray &payload);

private:
    const int  m_interval; //重连间隔n秒.
    QTimer     m_timer;
    QUrl       m_url;
    QWebSocket m_ws;
    bool       m_alive;    //socket是否正常.
};

#endif // MY_WEBSOCK_H
