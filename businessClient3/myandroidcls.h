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

#endif // MYANDROIDCLS_H
