#ifndef MYANDROIDCLS_H
#define MYANDROIDCLS_H
#include <QObject>

class MyAndroidCls
{
public:
    static int androidSdkVersion();
    static QString getExternalStorageDirectory(bool* isOk = nullptr);
    static QString calcExternalStorageDirectory();
};

class android_tool
{
public:
    static void logVerbose(const QString& tag, const QString& msg);
    static void toastShow(const QString& message);
    static void startTheService();
    static void ttsInit();
    static bool ttsSpeak(const QString& text);
};

#endif // MYANDROIDCLS_H
