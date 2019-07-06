#include "datarosvr.h"
#include <QObject>

DataROSvr::DataROSvr(QObject *parent) :DataROSource(parent)
{
    QObject::connect(&m_dataExch, &DataExchanger::sigReady, this, &DataROSource::sigReady);
    QObject::connect(&m_dataExch, &DataExchanger::sigStatusError, this, &DataROSource::sigStatusError);
    QObject::connect(&m_dataExch, &DataExchanger::sigTableChanged, this, &DataROSource::sigTableChanged);
}

void DataROSvr::doRun()
{
    bool setHostUrlRet = m_node.setHostUrl(QUrl(QStringLiteral(LOCAL_RO_URL)));
    Q_ASSERT(true == setHostUrlRet);
    bool enableRemotingRet = m_node.enableRemoting(this);
    Q_ASSERT(true == enableRemotingRet);
}

QString DataROSvr::dbLoadValue(const QString & key)
{
    return m_dataExch.dbLoadValue(key);
}
bool    DataROSvr::dbSaveValue(const QString & key, const QString & value)
{
    return m_dataExch.dbSaveValue(key, value);
}
QString DataROSvr::memGetData(const QString & varName)
{
    return m_dataExch.memGetData(varName);
}
bool    DataROSvr::memSetData(const QString & varName, const QString & value)
{
    return m_dataExch.memSetData(varName, value);
}
QString DataROSvr::memGetInfo(const QString & varName, const QStringList & paths)
{
    return m_dataExch.memGetInfo(varName, paths);
}
bool    DataROSvr::memSetInfo(const QString & varName, const QStringList & paths, const QString & value)
{
    return m_dataExch.memSetInfo(varName, paths, value);
}
QString DataROSvr::serviceState()
{
    return m_dataExch.serviceState();
}
bool DataROSvr::start()
{
    return m_dataExch.start();
}
QString DataROSvr::sendReq(const QString & typeName, const QString & jsonText, const QString & rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB)
{
    return m_dataExch.sendReq(typeName, jsonText, rID, isLog, isSafe, isPush, isUpCache, isC1NotC2, fillMsgNo, forceToDB);
}
QStringList DataROSvr::getTxMsgTypeNameList()
{
    return m_dataExch.getTxMsgTypeNameList();
}
QString DataROSvr::jsonExample(const QString & typeName)
{
    return m_dataExch.jsonExample(typeName);
}
bool DataROSvr::deleteCommonData1(const QString& userid, qint64 msgno)
{
    return m_dataExch.deleteCommonData1(userid, msgno);
}
bool DataROSvr::deleteCommonData2(const QString& userid, qint64 msgno, int seqno)
{
    return m_dataExch.deleteCommonData2(userid, msgno, seqno);
}
bool DataROSvr::deletePushWrap(const QString& userid, const QString& peerid, qint64 msgno)
{
    return m_dataExch.deletePushWrap(userid, peerid, msgno);
}
QString DataROSvr::serviceInfo()
{
    return m_dataExch.serviceInfo();
}
QString DataROSvr::sendCommonReq(const QStringList& kvs, bool isC1NotC2)
{
    return m_dataExch.sendCommonReq(kvs, isC1NotC2);
}
