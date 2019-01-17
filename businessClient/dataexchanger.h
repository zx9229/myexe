#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H

#include "mywebsock.h"
#include <QObject>
#include "temputils.h"
#include "sqlstruct.h"

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
    bool sendCommonNtosReq(QCommonNtosReq& reqData, bool needResp, bool needSave);

signals:
    void sigReady();
    void sigParentData(const QMap<QString, QConnInfoEx>& data);
    void sigWebsocketError(QAbstractSocket::SocketError error);

private:
    void initDB();
    void initOwnInfo();
    void handle_ConnectedData(QSharedPointer<txdata::ConnectedData> data);
    void handle_CommonNtosRsp(QSharedPointer<txdata::CommonNtosRsp> data);
    void handle_ParentDataRsp(QSharedPointer<txdata::ParentDataRsp> data);
    void toCommonNtosReq(QCommonNtosReq& src, txdata::CommonNtosReq& dst);

private slots:
    void slotOnConnected();
    void slotOnDisconnected();
    void slotOnMessage(const QByteArray &message);
    void slotOnError(QAbstractSocket::SocketError error);

public slots:
    void slotParentDataReq();

private:
    MyWebsock m_ws;
    QString   m_url;

    txdata::ConnectionInfo m_ownInfo;
    txdata::ConnectionInfo m_parentInfo;

    QMap<QString, QConnInfoEx> m_parentData;
    KeyValue m_RequestID;
    KeyValue m_SeqNo;

    QSqlDatabase m_db;
};

#endif // DATAEXCHANGER_H
