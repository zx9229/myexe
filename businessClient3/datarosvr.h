#ifndef DATA_RO_SVR_H
#define DATA_RO_SVR_H

#include "rep_dataro_source.h"
#include "dataexchanger.h"

class DataROSvr :public DataROSource
{
    Q_OBJECT
public:
    DataROSvr(QObject *parent = nullptr);
    void doRun();
public Q_SLOTS:
    virtual QString dbLoadValue(const QString & key);
    virtual bool    dbSaveValue(const QString & key, const QString & value);
    virtual QString memGetData(const QString & varName);
    virtual bool    memSetData(const QString & varName, const QString & value);
    virtual QString memGetInfo(const QString & varName, const QStringList & paths);
    virtual bool    memSetInfo(const QString & varName, const QStringList & paths, const QString & value);
    virtual QString serviceState();
    virtual bool start();
    virtual QString sendReq(const QString & typeName, const QString & jsonText, const QString & rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB);
    virtual QStringList getTxMsgTypeNameList();
    virtual QString jsonExample(const QString & typeName);
    virtual bool deleteCommonData(const QString & userid, int64_t msgno);
    virtual bool deleteCommonData(const QString & userid, int64_t msgno, int seqno);
    virtual bool deletePushWrap(const QString & userid, const QString & peerid, int64_t msgno);
    virtual QString serverStartTime();
public:
    DataExchanger m_dataExch;
    QRemoteObjectHost m_node;
};

#endif // DATA_RO_SVR_H
