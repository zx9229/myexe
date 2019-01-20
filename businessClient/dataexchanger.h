#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H

#include "mywebsock.h"
#include <QObject>
#include <QSharedPointer>
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

    static QString jsonByMsgObje(const google::protobuf::Message &msgObj, bool *isOk = nullptr);
    static QString nameByMsgType(txdata::MsgType msgType, int flag = 0, bool *isOk = nullptr);
    static QString jsonByMsgType(txdata::MsgType msgType, const QByteArray& serializedData, bool *isOk = nullptr);
    static bool    calcObjByName(const QString& typeName, QSharedPointer<google::protobuf::Message>& objOut);
    static QString jsonToObjAndS(const QString& typeName, const QString& jsonStr, txdata::MsgType& msgType, QByteArray& serializedData);

    static void CommonNtosReqQ2TX(const QCommonNtosReq& src, txdata::CommonNtosReq& dst);
    static void CommonNtosRspQ2TX(const QCommonNtosRsp& src, txdata::CommonNtosRsp& dst);
    static void CommonNtosRspTX2Q(const txdata::CommonNtosRsp& src, QCommonNtosRsp& dst);

    void setURL(const QString& url);
    void setUserKey(const QString& zoneName, const QString& nodeName, txdata::ProgramType execType, const QString& execName);
    void setBelongKey(const QString& zoneName, const QString& nodeName, txdata::ProgramType execType, const QString& execName);
    bool sendCommonNtosReq(QCommonNtosReq& reqData, bool needResp, bool needSave);
    void sendCommonNtosReq4resend(QCommonNtosReq& reqData);

public slots:
    void slotParentDataReq();

signals:
    void sigReady();
    void sigParentData(const QMap<QString, QConnInfoEx>& data);
    void sigWebsocketError(QAbstractSocket::SocketError error);
    void sigCommonNtosRsp(qint64 RefNum);

private:
    void initDB();
    void initOwnInfo();
    void handle_ConnectedData(QSharedPointer<txdata::ConnectedData> data);
    void handle_CommonNtosRsp(QSharedPointer<txdata::CommonNtosRsp> data);
    void handle_ParentDataRsp(QSharedPointer<txdata::ParentDataRsp> data);

private slots:
    void slotOnConnected();
    void slotOnDisconnected();
    void slotOnMessage(const QByteArray &message);
    void slotOnError(QAbstractSocket::SocketError error);

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
