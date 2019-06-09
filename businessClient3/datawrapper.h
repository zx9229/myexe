#ifndef DATA_WRAPPER_H
#define DATA_WRAPPER_H
#include <QScopedPointer>
#include "rep_dataro_replica.h"
#include "dataexchanger.h"
#include "datarosvr.h"

class DataWrapper :public QObject
{
    Q_OBJECT
public:
    DataWrapper(bool useRO,bool isServer, QObject *parent = nullptr);
public Q_SLOTS:
    Q_INVOKABLE QString dbLoadValue(const QString & key);
    Q_INVOKABLE bool    dbSaveValue(const QString & key, const QString & value);
    Q_INVOKABLE QString memGetData(const QString & varName);
    Q_INVOKABLE bool    memSetData(const QString & varName, const QString & value);
    Q_INVOKABLE QString memGetInfo(const QString & varName, const QStringList & paths);
    Q_INVOKABLE bool    memSetInfo(const QString & varName, const QStringList & paths, const QString & value);
    Q_INVOKABLE bool start();
    Q_INVOKABLE QString sendReq(const QString & typeName, const QString & jsonText, const QString & rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB);
    Q_INVOKABLE QStringList getTxMsgTypeNameList();
    Q_INVOKABLE QString jsonExample(const QString & typeName);
signals:
    void sigReady();
    void sigStatusError(const QString & errMessage, int errType);
    void sigTableChanged(const QString & tableName);
private:
    QSharedPointer<DataExchanger> m_daExch;
    QSharedPointer<DataROReplica> m_client;
    QSharedPointer<QRemoteObjectHost> m_node;
    QSharedPointer<DataROSvr> m_server;
    QSqlDatabase m_db;
};

#endif // DATA_WRAPPER_H
