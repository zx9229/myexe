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
    bool start(const QUrl& url);
    void stop();
    qint64 sendBinaryMessage(const QByteArray &data);

signals:
    void onMessage(const QByteArray &message);

private:
    void reconnect();

private slots:
    void connected();
    void disconnected();
    void binaryMessageReceived(const QByteArray &message);
    void error(QAbstractSocket::SocketError error);
    void pong(quint64 elapsedTime, const QByteArray &payload);

private:
    const int  m_interval;
    QUrl       m_url;
    QWebSocket m_ws;
    QTimer     m_timer;
};

#endif // MY_WEBSOCK_H
