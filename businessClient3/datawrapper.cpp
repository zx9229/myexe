#include "datawrapper.h"

DataWrapper::DataWrapper(bool isRO, QObject *parent) :QObject(parent)
{
    if (isRO)
    {
        m_client = QSharedPointer<DataROReplica>(new DataROReplica);
    }
    else
    {
        m_daExch = QSharedPointer<DataExchanger>(new DataExchanger);
    }
}

QString DataWrapper::dbLoadValue(const QString & key)
{
    if (m_daExch)
    {
        return m_daExch->dbLoadValue(key);
    }
    else
    {
        auto reply = m_client->dbLoadValue(key);
        reply.waitForFinished();
        return reply.returnValue();
    }
}

bool    DataWrapper::dbSaveValue(const QString & key, const QString & value)
{
    if (m_daExch)
    {
        return m_daExch->dbSaveValue(key, value);
    }
    else
    {
        auto reply = m_client->dbSaveValue(key, value);
        reply.waitForFinished();
        return reply.returnValue();
    }
}

QString DataWrapper::memGetData(const QString & varName)
{
    if (m_daExch)
    {
        return m_daExch->memGetData(varName);
    }
    else
    {
        auto reply = m_client->memGetData(varName);
        reply.waitForFinished();
        return reply.returnValue();
    }
}

bool    DataWrapper::memSetData(const QString & varName, const QString & value)
{
    if (m_daExch)
    {
        return m_daExch->memSetData(varName, value);
    }
    else
    {
        auto reply = m_client->memSetData(varName, value);
        reply.waitForFinished();
        return reply.returnValue();
    }
}

QString DataWrapper::memGetInfo(const QString & varName, const QStringList & paths)
{
    if (m_daExch)
    {
        return m_daExch->memGetInfo(varName, paths);
    }
    else
    {
        auto reply = m_client->memGetInfo(varName, paths);
        reply.waitForFinished();
        return reply.returnValue();
    }
}

bool    DataWrapper::memSetInfo(const QString & varName, const QStringList & paths, const QString & value)
{
    if (m_daExch)
    {
        return m_daExch->memSetInfo(varName, paths, value);
    }
    else
    {
        auto reply = m_client->memSetInfo(varName, paths, value);
        reply.waitForFinished();
        return reply.returnValue();
    }
}

bool DataWrapper::start()
{
    if (m_daExch)
    {
        return m_daExch->start();
    }
    else
    {
        auto reply = m_client->start();
        reply.waitForFinished();
        return reply.returnValue();
    }
}

QString DataWrapper::sendReq(const QString & typeName, const QString & jsonText, const QString & rID, bool isLog, bool isSafe, bool isPush, bool isUpCache, bool isC1NotC2, bool fillMsgNo, bool forceToDB)
{
    if (m_daExch)
    {
        return m_daExch->sendReq(typeName, jsonText, rID, isLog, isSafe, isPush, isUpCache, isC1NotC2, fillMsgNo, forceToDB);
    }
    else
    {
        auto reply = m_client->sendReq(typeName, jsonText, rID, isLog, isSafe, isPush, isUpCache, isC1NotC2, fillMsgNo, forceToDB);
        if (reply.waitForFinished())
            return reply.returnValue();
        else
            return "unknown_error";
    }
}

QStringList DataWrapper::getTxMsgTypeNameList()
{
    if (m_daExch)
    {
        return m_daExch->getTxMsgTypeNameList();
    }
    else
    {
        auto reply = m_client->getTxMsgTypeNameList();
        reply.waitForFinished();
        return reply.returnValue();
    }
}

QString DataWrapper::jsonExample(const QString & typeName)
{
    if (m_daExch)
    {
        return m_daExch->jsonExample(typeName);
    }
    else
    {
        auto reply = m_client->jsonExample(typeName);
        reply.waitForFinished();
        return reply.returnValue();
    }
}
