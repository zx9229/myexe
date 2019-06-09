#include <QVector>
#include <QDebug>
#include <QStorageInfo>
#include <QCoreApplication>
#ifdef Q_OS_ANDROID
#include <QAndroidJniEnvironment>
#include <QAndroidJniObject>
#include <QtAndroid>
#endif
#include "myandroidcls.h"

int MyAndroidCls::androidSdkVersion()
{
    int value = -1;
#ifdef Q_OS_ANDROID
    value = QtAndroid::androidSdkVersion();
#endif
    return value;
}

QString MyAndroidCls::getExternalStorageDirectory(bool* isOk)
{
    QString value;
    if (isOk)
    {
        *isOk = false;
    }
#ifdef Q_OS_ANDROID
    //https://developer.android.com/reference/android/os/Environment.html
    QString externalStorageState = QAndroidJniObject::callStaticObjectMethod("android/os/Environment", "getExternalStorageState", "()Ljava/lang/String;").toString();
    QString MEDIA_MOUNTED = QAndroidJniObject::getStaticObjectField("android/os/Environment", "MEDIA_MOUNTED", "Ljava/lang/String;").toString();
    if (externalStorageState == MEDIA_MOUNTED)
    {
        QAndroidJniObject storageDirectory = QAndroidJniObject::callStaticObjectMethod("android/os/Environment", "getExternalStorageDirectory", "()Ljava/io/File;");
        value = storageDirectory.callObjectMethod("toString", "()Ljava/lang/String;").toString();
        if (isOk && !value.isEmpty())
        {
            *isOk = true;
        }
    }
#endif
    return value;
}

QString MyAndroidCls::calcExternalStorageDirectory()
{
    QString value;

    const char* _storage = "/storage";
    QString internalStorageDir = MyAndroidCls::getExternalStorageDirectory();
    if (internalStorageDir.isEmpty() || !internalStorageDir.startsWith(_storage, Qt::CaseSensitive))
        return value;

    for (auto& info : QStorageInfo::mountedVolumes())
    {
        QString rootPath = info.rootPath();
        if (rootPath.startsWith(_storage, Qt::CaseSensitive) == false)
            continue;
        if (internalStorageDir.startsWith(rootPath, Qt::CaseSensitive) == false)
        {
            value = rootPath;
            break;
        }
    }

    return value;
}
