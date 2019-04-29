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
public:
    void from_txdata_AtomicKey(const ::txdata::AtomicKey& src)
    {
        this->ZoneName = QString::fromStdString(src.zonename());
        this->NodeName = QString::fromStdString(src.nodename());
        this->ExecType = QString::fromStdString(txdata::ProgramType_Name(src.exectype()));
        this->ExecName = QString::fromStdString(src.execname());
    }
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
    void from_txdata_ConnectedData(const ::txdata::ConnectedData& src)
    {
        this->UserKey.from_txdata_AtomicKey(src.info().userkey());
        this->UserID = QString::fromStdString(src.info().userid());
        this->BelongKey.from_txdata_AtomicKey(src.info().belongkey());
        this->BelongID = QString::fromStdString(src.info().belongid());
        this->Version = QString::fromStdString(src.info().version());
        this->LinkMode = QString::fromStdString(txdata::ConnectionInfo::LinkType_Name(src.info().linkmode()));
        this->ExePid = src.info().exepid();
        this->ExePath = QString::fromStdString(src.info().exepath());
        this->Remark = QString::fromStdString(src.info().remark());
        for (int i = 0; i < src.pathway_size(); ++i) { this->Pathway.append(QString::fromStdString(src.pathway(i))); }
    }
};
Q_DECLARE_METATYPE(QConnInfoEx);

#endif // TEMPUTILS_H
