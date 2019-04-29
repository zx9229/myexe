#ifndef TEMPUTILS_H
#define TEMPUTILS_H

#include<QString>
#include"protobuf/m2b.h"

class QAtomicKey
{
public:
    QString ZoneName;
    QString NodeName;
    QString ExecType;
    QString ExecName;
};

class QConnInfoEx
{
public:
    QAtomicKey  UserKey;
    QString     UserID;
    QAtomicKey  BelongKey;
    QString     BelongID;
    QString     Version;
    QString     LinkMode;
    int         ExePid;
    QString     ExePath;
    QString     Remark;
    QStringList Pathway;
public:
};
Q_DECLARE_METATYPE(QConnInfoEx);

#endif // TEMPUTILS_H
