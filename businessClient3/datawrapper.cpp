#include <QSqlDatabase>
#include "datawrapper.h"
#include "myandroidcls.h"

DataWrapper::DataWrapper(bool useRO, bool isServer, QObject *parent) :QObject(parent)
{
    if (useRO)
    {
        if (isServer)
        {
            m_server = QSharedPointer<DataROSvr>(new DataROSvr);
            m_server->doRun();
            qDebug() << "server init ok";
        }
        else
        {
            m_node.reset(new QRemoteObjectHost);
            bool connectToNodeRet = m_node->connectToNode(QUrl(QStringLiteral(LOCAL_RO_URL)));
            qDebug() << "connectToNodeRet=" << connectToNodeRet;
            //Q_ASSERT(connectToNodeRet);
            m_client.reset(m_node->acquire<DataROReplica>());
            bool waitForSourceRet = m_client->waitForSource(3000);
            qDebug() << "waitForSourceRet=" << waitForSourceRet;
            //Q_ASSERT(waitForSourceRet);
            QObject::connect(m_client.get(), &DataROReplica::sigReady, this, &DataWrapper::sigReady);
            QObject::connect(m_client.get(), &DataROReplica::sigStatusError, this, &DataWrapper::sigStatusError);
            QObject::connect(m_client.get(), &DataROReplica::sigTableChanged, this, &DataWrapper::sigTableChanged);
            QObject::connect(m_client.get(), &DataROReplica::stateChanged, this, &DataWrapper::onStateChanged);
            {
                m_db = QSqlDatabase::addDatabase("QSQLITE");
                m_db.setDatabaseName(QString().isEmpty() ? (SQLITE_DB_NAME) : (":memory:"));
                bool isOk = m_db.open();
                Q_ASSERT(isOk);
            }
            qDebug() << "client init ok";
        }
    }
    else
    {
        m_daExch = QSharedPointer<DataExchanger>(new DataExchanger);
        QObject::connect(m_daExch.get(), &DataExchanger::sigReady, this, &DataWrapper::sigReady);
        QObject::connect(m_daExch.get(), &DataExchanger::sigStatusError, this, &DataWrapper::sigStatusError);
        QObject::connect(m_daExch.get(), &DataExchanger::sigTableChanged, this, &DataWrapper::sigTableChanged);
        qDebug() << "DataExchanger init ok";
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

QString  DataWrapper::serviceState()
{
    if (m_daExch)
    {
        return m_daExch->serviceState();
    }
    else
    {
        auto reply = m_client->serviceState();
        reply.waitForFinished(1000);
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
        if (reply.waitForFinished(1000))
        {
            return reply.returnValue();
        }
        else
        {
            emit sigStatusError("start timeout", 0);
            return false;
        }
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

Q_INVOKABLE QString DataWrapper::remoteObjectState()
{
    QRemoteObjectReplica::State curState = QRemoteObjectReplica::Uninitialized;
    QMetaEnum metaEnum = QMetaEnum::fromType<QRemoteObjectReplica::State>();
    if (m_client)
    {
        curState = m_client->state();
    }
    return metaEnum.valueToKey(curState);
}

Q_INVOKABLE void DataWrapper::startTheService()
{
    android_tool::startTheService();
}

Q_INVOKABLE void DataWrapper::ttsSpeak(const QString & text)
{
    android_tool::ttsInit();
    android_tool::ttsSpeak(text);
}

void DataWrapper::onStateChanged(QRemoteObjectReplica::State state, QRemoteObjectReplica::State oldState)
{
    QMetaEnum metaEnum = QMetaEnum::fromType<QRemoteObjectReplica::State>();
    QString sState = metaEnum.valueToKey(state);
    QString sOldState = metaEnum.valueToKey(oldState);
    emit sigStateChanged(state, sState, oldState, sOldState);
}
