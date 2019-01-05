#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H

#include <QObject>
#include "mywebsock.h"

class DataExchanger : public QObject
{
    Q_OBJECT

public:
    explicit DataExchanger(QObject *parent = 0);
    ~DataExchanger();

public:
    QString Login(const QString& url, const QString& username, const QString& password);
    MyWebsock& ws();

private slots:
    void slotOnConnected();
    void slotOnDisconnected();
    void slotOnMessage(const QByteArray &message);

private:
    MyWebsock m_ws;
    QString   m_url;
    QString   m_username;
    QString   m_password;
};

#endif // DATAEXCHANGER_H
