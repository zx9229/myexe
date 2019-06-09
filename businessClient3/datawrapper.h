#ifndef DATA_WRAPPER_H
#define DATA_WRAPPER_H
#include <QScopedPointer>
#include "rep_dataro_replica.h"
#include "dataexchanger.h"

class DataWrapper :public QObject
{
public:
    DataWrapper(bool isRO, QObject *parent = nullptr);
public Q_SLOTS:
    QString dbLoadValue(const QString & key);
    bool    dbSaveValue(const QString & key, const QString & value);
    QString memGetData(const QString & varName);
    bool    memSetData(const QString & varName, const QString & value);
    QString memGetInfo(const QString & varName, const QStringList & paths);
    bool    memSetInfo(const QString & varName, const QStringList & paths, const QString & value);
    bool start();
    QString sendReq(const QString & typeName, const QString & jsonText, const QString & rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB);
    QStringList getTxMsgTypeNameList();
    QString jsonExample(const QString & typeName);
public:
    QSharedPointer<DataExchanger> m_daExch;
    QSharedPointer<DataROReplica> m_client;
};

#endif // DATA_WRAPPER_H
