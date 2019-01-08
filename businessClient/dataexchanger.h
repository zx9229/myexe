#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H

#include "mywebsock.h"
#include <QObject>
#include "txdata.pb.h"

class DataExchanger : public QObject
{
    Q_OBJECT

public:
    explicit DataExchanger(QObject *parent = 0);
    ~DataExchanger();

public:
    MyWebsock& ws();
    bool start();

    void setURL(const QString& url);
    void setUserKey(const QString& zoneName, const QString& nodeName, txdata::ProgramType execType, const QString& execName);
    void setBelongKey(const QString& zoneName, const QString& nodeName, txdata::ProgramType execType, const QString& execName);

signals:
    void sigLoginProgress(int curPos, int errCode, const QString& errMsg);
    void sigReady();

private:
    void initOwnInfo();
    void handle_ConnectedData(QSharedPointer<txdata::ConnectedData> data);
    void handle_CommonNtosRsp(QSharedPointer<txdata::CommonNtosRsp> data);

private slots:
    void slotOnConnected();
    void slotOnDisconnected();
    void slotOnMessage(const QByteArray &message);

public slots:
    void slotReqServerCache();

private:
    MyWebsock m_ws;
    QString   m_url;

    txdata::ConnectionInfo m_ownInfo;
    txdata::ConnectionInfo m_parentInfo;
};

#endif // DATAEXCHANGER_H
