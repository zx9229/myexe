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

#define LOCAL_RO_URL    "local:data_ro_svr"
#define SQLITE_DB_NAME  "_zx_test.db"

class DataExchanger : public QObject
{
    Q_OBJECT

public:
    explicit DataExchanger(QObject *parent = nullptr);
    ~DataExchanger();

public:
    Q_INVOKABLE QString dbLoadValue(const QString& key);
    Q_INVOKABLE bool    dbSaveValue(const QString& key, const QString& value);
    Q_INVOKABLE QString memGetData(const QString& varName);
    Q_INVOKABLE bool    memSetData(const QString& varName, const QString& value);
    Q_INVOKABLE QString memGetInfo(const QString& varName, const QStringList& paths);
    Q_INVOKABLE bool    memSetInfo(const QString& varName, const QStringList& paths, const QString& value);
    Q_INVOKABLE QString serviceState();
    Q_INVOKABLE bool start();
    Q_INVOKABLE QString sendReq(const QString& typeName, const QString& jsonText, const QString& rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB);
    Q_INVOKABLE QStringList getTxMsgTypeNameList();
    Q_INVOKABLE QString jsonExample(const QString& typeName);

signals:
    void sigReady();
    //void sigWebsocketError(QAbstractSocket::SocketError error);
    void sigStatusError(const QString& errMessage, int errType);
    void sigTableChanged(const QString &tableName);

private:
    void initDB();
    void initOwnInfo();
    QString toC1C2(int64_t msgNo, const QString& typeName, const QString& jsonText, const QString& rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, GPMSGPTR& msgOut);
    QString sendCommon1Req(QSharedPointer<txdata::Common1Req> data);
    QString sendCommon2Req(QSharedPointer<txdata::Common2Req> data);
    void handle_Common2Ack(QSharedPointer<txdata::Common2Ack> data);
    void handle_Common2Rsp(QSharedPointer<txdata::Common2Rsp> data);
    void handle_Common1Req(QSharedPointer<txdata::Common1Req> data);
    void handle_Common1Rsp(QSharedPointer<txdata::Common1Rsp> data);
    void handle_ConnectReq(QSharedPointer<txdata::ConnectReq> data);
    void handle_ConnectRsp(QSharedPointer<txdata::ConnectRsp> data);
    void handle_PathwayInfo(QSharedPointer<txdata::PathwayInfo> data);
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
    KeyValue     m_MsgNo;

    bool m_lastFind;
    QString m_subUser;
    SafeSynchCache m_cacheSync;
};

#endif // DATAEXCHANGER_H
