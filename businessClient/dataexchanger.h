#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H

#include <QObject>
#include "mywebsock.h"
#include "txdata.pb.h"

class DataExchanger : public QObject
{
    Q_OBJECT

public:
    explicit DataExchanger(QObject *parent = 0);
    ~DataExchanger();

public:
    QString login();
    MyWebsock& ws();

    void setURL(const QString& url);
    void setUserKey(const QString& zoneName, const QString& nodeName, txdata::ProgramType execType, const QString& execName);
    void setBelongKey(const QString& zoneName, const QString& nodeName, txdata::ProgramType execType, const QString& execName);

signals:
    void sigLoginProgress(int curPos, int errCode, const QString& errMsg);

private slots:
    void slotOnConnected();
    void slotOnDisconnected();
    void slotOnMessage(const QByteArray &message);

private:
    MyWebsock m_ws;
    QString   m_url;
    QString   m_username;
    QString   m_password;
    const int m_totalPos;//总进度.

    txdata::AtomicKey m_userKey;
    txdata::AtomicKey m_belongKey;
};

#endif // DATAEXCHANGER_H
