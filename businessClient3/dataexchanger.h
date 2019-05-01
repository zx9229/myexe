#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H
//Exposing Attributes of C++ Types to QML
//https://doc.qt.io/qt-5/qtqml-cppintegration-exposecppattributes.html
#include "mywebsock.h"
#include <QObject>
#include <QSharedPointer>
#include <QSqlQuery>
#include"protobuf/m2b.h"

class DataExchanger : public QObject
{
    Q_OBJECT

public:
    explicit DataExchanger(QObject *parent = 0);
    ~DataExchanger();

public:
    MyWebsock& ws();

    static QString jsonByMsgObje(const google::protobuf::Message &msgObj, bool *isOk = nullptr);
    static QString nameByMsgType(txdata::MsgType msgType, int flag = 0, bool *isOk = nullptr);
    static QString jsonByMsgType(txdata::MsgType msgType, const QByteArray& serializedData, bool *isOk = nullptr);
    static bool    calcObjByName(const QString& typeName, QSharedPointer<google::protobuf::Message>& objOut);
    static QString jsonToObjAndS(const QString& typeName, const QString& jsonStr, txdata::MsgType& msgType, QByteArray& serializedData);

    Q_INVOKABLE void setURL(const QString& url);
    Q_INVOKABLE void setOwnInfo(const QString& userID, const QString& belongID);
    Q_INVOKABLE bool start();

public slots:

signals:
    void sigReady();
    void sigWebsocketError(QAbstractSocket::SocketError error);
    void sigStatusError(const QString& errMessage, int errType);

private:
    void initDB();
    void initOwnInfo();
    void handle_MessageAck(QSharedPointer<txdata::MessageAck> data);
    void handle_ConnectReq(QSharedPointer<txdata::ConnectReq> data);
    void handle_ConnectRsp(QSharedPointer<txdata::ConnectRsp> data);

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
