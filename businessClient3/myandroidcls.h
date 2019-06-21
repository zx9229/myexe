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


#ifdef Q_OS_ANDROID
#include <QAndroidJniEnvironment>
#include <QAndroidJniObject>
#include <QtAndroid>
#endif
class android_tool
{
public:
    static void logVerbose(const QString& tag, const QString& msg)
    {
#ifdef Q_OS_ANDROID
        QAndroidJniObject AJO_tag = QAndroidJniObject::fromString(tag);
        QAndroidJniObject AJO_msg = QAndroidJniObject::fromString(msg);
        QAndroidJniObject::callStaticMethod<void>("zx/qtproject/example/AndroidTool", "logVerbose", "(Ljava/lang/String;Ljava/lang/String;)V", AJO_tag.object(), AJO_msg.object());
#endif
    }
    static void toastShow(const QString& message)
    {
#ifdef Q_OS_ANDROID
        QAndroidJniObject AJO_message = QAndroidJniObject::fromString(message);
        QAndroidJniObject::callStaticMethod<void>("zx/qtproject/example/AndroidTool", "toastShow", "(Landroid/content/Context;Ljava/lang/String;)V", QtAndroid::androidActivity().object(), AJO_message.object());
#endif
    }
    static QString funTest()
    {
        QString msg;
#ifdef Q_OS_ANDROID
        msg += "ZxActivity,beg,";
        QAndroidJniObject::callStaticMethod<void>("zx/qtproject/example/ZxActivity", "funTest", "(Landroid/content/Context;)V", QtAndroid::androidActivity().object());
        msg += "end,";
#endif
        return msg;
    }
};

#endif // MYANDROIDCLS_H
