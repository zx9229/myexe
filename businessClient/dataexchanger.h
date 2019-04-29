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

    void setURL(const QString& url);

public slots:
    void slotParentDataReq();

signals:
    void sigReady();
    void sigWebsocketError(QAbstractSocket::SocketError error);

private:
    void initDB();
    void initOwnInfo();
    void handle_MessageAck(QSharedPointer<txdata::MessageAck> data);

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
};

#endif // DATAEXCHANGER_H
