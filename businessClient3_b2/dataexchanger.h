#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H
//Exposing Attributes of C++ Types to QML
//https://doc.qt.io/qt-5/qtqml-cppintegration-exposecppattributes.html
#include <QObject>
#include <QSharedPointer>
#include <QSqlQuery>
#include"protobuf/m2b.h"
#include "mywebsock.h"
#include "sqlstruct.h"
#include "safesynchcache.h"

class DataExchanger : public QObject
{
    Q_OBJECT

public:
    explicit DataExchanger(QObject *parent = nullptr);
    ~DataExchanger();

public:
    MyWebsock& ws();

    static QString jsonByMsgObje(const google::protobuf::Message &msgObj, bool *isOk = nullptr);
    static QString nameByMsgType(txdata::MsgType msgType, int flag = 0, bool *isOk = nullptr);
    static QString jsonByMsgType(txdata::MsgType msgType, const QByteArray& serializedData, bool *isOk = nullptr);
    static bool    calcObjByName(const QString& typeName, QSharedPointer<google::protobuf::Message>& objOut);
    static QString jsonToObjAndS(const QString& typeName, const QString& jsonStr, txdata::MsgType& msgType, QByteArray& serializedData);
    static void qdt2gpt(::google::protobuf::Timestamp& gptDst, const QDateTime& qdtSrc);

    Q_INVOKABLE void setURL(const QString& url);
    Q_INVOKABLE void setOwnInfo(const QString& userID, const QString& belongID);
    Q_INVOKABLE bool start();
    Q_INVOKABLE QString demoFun(const QString& typeName, const QString& jsonText, const QString& rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2);
    Q_INVOKABLE QString QryConnInfoReq(const QString& userId);

public slots:

signals:
    void sigReady();
    void sigWebsocketError(QAbstractSocket::SocketError error);
    void sigStatusError(const QString& errMessage, int errType);

private:
    void initDB();
    void initOwnInfo();
    QString toC1C2(const QString& typeName, const QString& jsonText, const QString& rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, GPMSGPTR& msgOut);
    QString sendCommon1Req(QSharedPointer<txdata::Common1Req> data);
    QString sendCommon2Req(QSharedPointer<txdata::Common2Req> data);
    void handle_Common2Ack(QSharedPointer<txdata::Common2Ack> data);
    void handle_Common2Rsp(QSharedPointer<txdata::Common2Rsp> data);
    void handle_Common1Rsp(QSharedPointer<txdata::Common1Rsp> data);
    void handle_ConnectReq(QSharedPointer<txdata::ConnectReq> data);
    void handle_ConnectRsp(QSharedPointer<txdata::ConnectRsp> data);
    void deal_QryConnInfoRsp(QSharedPointer<txdata::QryConnInfoRsp> data);

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

    QSqlDatabase m_db;

    SafeSynchCache m_cacheSync;
};

#endif // DATAEXCHANGER_H
