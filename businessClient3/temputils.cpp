#include <QCoreApplication>
#include <QProcess>
#include <QStorageInfo>
#include <QDateTime>
#include <QDebug>
#include <QBuffer>
#include <QSysInfo>
#include "temputils.h"

#if defined(Q_OS_ANDROID)
#include "myandroidcls.h"
#endif


QString TempUtils::calcAppDataDir(bool doMkpath)
{
    QString value;
    do {
        QString pathValue;
#if defined(Q_OS_WIN)
        if (QSysInfo::kernelType().startsWith("win", Qt::CaseInsensitive) == false)
            break;
        pathValue = ".";
#elif defined(Q_OS_ANDROID)
        bool isOk = false;
        pathValue = MyAndroidCls::getExternalStorageDirectory(&isOk);
        if (!isOk || pathValue.isEmpty())
            break;
#else
/**/#error please_set_application_data_dir
#endif
        pathValue = QString("%1/_%2").arg(pathValue).arg(qApp->applicationName());
        pathValue = QDir(pathValue).absolutePath();
        if (doMkpath)
        {
            if (QDir().mkdir(pathValue) == false)
                break;
        }
        value = pathValue;
    } while (false);
    return value;
}

QString TempUtils::calcAppLogDir(bool doMkpath)
{
    QString value;
    do {
        QString pathValue = calcAppDataDir();
        if (pathValue.isEmpty())
            break;
        pathValue = QString("%1/%2").arg(pathValue).arg("app_log");
        pathValue = QDir(pathValue).absolutePath();
        if (doMkpath)
        {
            if (QDir().mkpath(pathValue) == false)
                break;
        }
        value = pathValue;
    } while (false);
    return value;
}

QString TempUtils::calcAppLogName(bool* mkpathOk)
{
    if (mkpathOk) { *mkpathOk = false; }
    QString value;
    do {
        QString pathValue = TempUtils::calcAppLogDir(true);
        if (pathValue.isEmpty())
            break;
        QString curDtStr = QDateTime::currentDateTime().toString("yyyy_MM_dd-hh_mm_ss");
        pathValue = QString("%1/zx_app_%2.log").arg(pathValue).arg(curDtStr);
        value = QDir(pathValue).absolutePath();
        if (mkpathOk) { *mkpathOk = true; }
    } while (false);
    return value;
}

QString TempUtils::calcServiceLogDir(bool doMkpath)
{
    QString value;
    do {
        QString pathValue = calcAppDataDir();
        if (pathValue.isEmpty())
            break;
        pathValue = QString("%1/%2").arg(pathValue).arg("service_log");
        pathValue = QDir(pathValue).absolutePath();
        if (doMkpath)
        {
            if (QDir().mkpath(pathValue) == false)
                break;
        }
        value = pathValue;
    } while (false);
    return value;
}

QString TempUtils::calcServiceLogName(bool* mkpathOk)
{
    if (mkpathOk) { *mkpathOk = false; }
    QString value;
    do {
        QString pathValue = TempUtils::calcServiceLogDir(true);
        if (pathValue.isEmpty())
            break;
        QString curDtStr = QDateTime::currentDateTime().toString("yyyy_MM_dd-hh_mm_ss");
        pathValue = QString("%1/zx_srv_%2.log").arg(pathValue).arg(curDtStr);
        value = QDir(pathValue).absolutePath();
        if (mkpathOk) { *mkpathOk = true; }
    } while (false);
    return value;
}

QString TempUtils::calcEnvironmentVariable(const char* name, bool* isOk /*= nullptr*/)
{
    QString value;
    if (isOk)
    {
        *isOk = false;
    }
    if (nullptr == name || 0x0 == *name)
    {
        return value;
    }
    QString key; key.sprintf("%s=", name);
    QStringList kvList = QProcess::systemEnvironment();
    for (auto& kv : kvList)
    {
        if (kv.startsWith(key, Qt::CaseSensitive))
        {
            value = kv.mid(key.size());
            if (isOk)
            {
                *isOk = true;
            }
            break;
        }
    }
    return value;
}
