#ifndef SQL_STRUCT_H
#define SQL_STRUCT_H
#include <cstdint>
#include <cfloat>
#include <QObject>
#include <QVariant>
#include <QString>
#include <QByteArray>
#include <QDateTime>
#include <QSqlQuery>

namespace {
    //使用非const的引用,可以强校验入参的类型,基本可以杜绝隐式类型转换.
    inline bool Valid(int32_t& data) { return INT32_MAX != data; }
    inline bool Valid(int64_t& data) { return INT64_MAX != data; }
    inline bool Valid(float& data) { return FLT_MAX != data; }
    inline bool Valid(double& data) { return DBL_MAX != data; }
    inline bool Valid(QString& data) { return !data.isNull(); }
    inline bool Valid(QByteArray& data) { return !data.isNull(); }
    inline bool Valid(QDateTime& data) { return data.isValid(); }
    inline void fromQVariant(int32_t& dst, const QVariant& src) { dst = src.toInt(nullptr); }
    inline void fromQVariant(int64_t& dst, const QVariant& src) { dst = src.toLongLong(nullptr); }
    inline void fromQVariant(QString& dst, const QVariant& src) { dst = src.toString(); }
    inline void fromQVariant(QByteArray& dst, const QVariant& src) { dst = src.toByteArray(); }
    inline void fromQVariant(QDateTime& dst, const QVariant& src) { dst = QDateTime::fromString(src.toString(), Qt::ISODateWithMs); }
}

class KeyValue
{
public:
    QString Key;
    QString Value;
public:
    static QString static_table_name()
    {
        return "KeyValue";
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            Key   TEXT NOT NULL PRIMARY KEY ,\
            Value TEXT NOT NULL )"
        ).QString::arg(static_table_name());
        return sql;
    }
    bool insert_data(QSqlQuery& query, bool insertNotReplace)
    {
        bool isOk = false;
        QString sqlStr = QObject::tr("INSERT %1 INTO %2 (Key, Value) VALUES (:Key, :Value)")
            .arg(insertNotReplace ? "" : "OR REPLACE").arg(static_table_name());
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        query.bindValue(":Key", this->Key);
        query.bindValue(":Value", this->Value);
        return query.exec();
    }
    bool update_data(QSqlQuery& query)
    {
        bool isOk = false;
        QString sqlStr = QObject::tr("UPDATE %1 SET Value=:Value WHERE Key=:Key").arg(static_table_name());
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        query.bindValue(":Key", this->Key);
        query.bindValue(":Value", this->Value);
        return query.exec();
    }
    static void select_data(QSqlQuery& query, QList<KeyValue>& results)
    {
        results.clear();
        bool isOk = false;
        QString sqlStr = QObject::tr("SELECT Key,Value FROM %1").arg(static_table_name());
        isOk = query.exec(sqlStr);
        Q_ASSERT(isOk);
        while (query.next()) {
            KeyValue result;
            result.Key = query.value("Key").toString();
            result.Value = query.value("Value").toString();
            results.append(result);
        }
    }
};

class CommonNtosDataNode
{
public:
    int64_t    ID;//PK
    int64_t    RequestID;
    int64_t    SeqNo;//UNIQUE(For the purposes of UNIQUE constraints, NULL values are considered distinct from all other values, including other NULLs.)
    QString    UserID;
    QString    ReqType;
    QByteArray ReqData;
    QDateTime  ReqTime;
    QDateTime  CreateTime;//insert_time
    int32_t    State;
    int32_t    ErrNo;
    QString    ErrMsg;
    QString    RspType;
    QByteArray RspData;
public:
    CommonNtosDataNode()
    {
        this->RequestID = INT64_MAX;
        this->ID = INT64_MAX;
        this->SeqNo = INT64_MAX;
        this->State = INT32_MAX;
        this->ErrNo = INT32_MAX;
    }
public:
    static QString static_table_name()
    {
        return "CommonNtosDataNode";//可能会有(object_table_name)函数.
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            ID         INTEGER NOT NULL PRIMARY KEY ,\
            RequestID  INTEGER     NULL ,\
            SeqNo      INTEGER     NULL UNIQUE ,\
            UserID     TEXT    NOT NULL ,\
            ReqType    TEXT    NOT NULL ,\
            ReqData    BLOB        NULL ,\
            ReqTime    TEXT        NULL ,\
            CreateTime TEXT        NULL ,\
            State      INTEGER NOT NULL ,\
            ErrNo      INTEGER NOT NULL ,\
            ErrMsg     TEXT    NOT NULL ,\
            RspType    TEXT    NOT NULL ,\
            RspData    BLOB        NULL   )"
        ).QString::arg(static_table_name());
        return sql;
    }
    bool insert_data(QSqlQuery& query, int64_t* lastInsertId = nullptr)
    {
        //请外部保证在同一个(上下文/先后顺序/总之就是加锁的意思).
        bool isOk = false;
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){cols.append\("$1"\);}】.
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){query.bindValue\(":$1",this->$1\);}】.
        //注意(NOT NULL)要特殊处理.
        this->CreateTime = QDateTime::currentDateTime();
        QStringList cols;
        if (Valid(this->ID)) { cols.append("ID"); }
        if (Valid(this->RequestID)) { cols.append("RequestID"); }
        if (Valid(this->SeqNo)) { cols.append("SeqNo"); }
        if (Valid(this->UserID)) { cols.append("UserID"); }
        if (Valid(this->ReqType)) { cols.append("ReqType"); }
        if (Valid(this->ReqData)) { cols.append("ReqData"); }
        if (Valid(this->ReqTime)) { cols.append("ReqTime"); }
        if (Valid(this->CreateTime)) { cols.append("CreateTime"); }
        if (Valid(this->State)) { cols.append("State"); }
        if (Valid(this->ErrNo)) { cols.append("ErrNo"); }
        if (Valid(this->ErrMsg)) { cols.append("ErrMsg"); }
        if (Valid(this->RspType)) { cols.append("RspType"); }
        if (Valid(this->RspData)) { cols.append("RspData"); }
        //
        QString sqlStr = QObject::tr("INSERT INTO %1 (%2) VALUES (%3)").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->ID)) { query.bindValue(":ID", this->ID); }
        if (Valid(this->RequestID)) { query.bindValue(":RequestID", this->RequestID); }
        if (Valid(this->SeqNo)) { query.bindValue(":SeqNo", this->SeqNo); }
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->ReqType)) { query.bindValue(":ReqType", this->ReqType); }
        if (Valid(this->ReqData)) { query.bindValue(":ReqData", this->ReqData); }
        if (Valid(this->ReqTime)) { query.bindValue(":ReqTime", this->ReqTime); }
        if (Valid(this->CreateTime)) { query.bindValue(":CreateTime", this->CreateTime); }
        if (Valid(this->State)) { query.bindValue(":State", this->State); }
        if (Valid(this->ErrNo)) { query.bindValue(":ErrNo", this->ErrNo); }
        if (Valid(this->ErrMsg)) { query.bindValue(":ErrMsg", this->ErrMsg); }
        if (Valid(this->RspType)) { query.bindValue(":RspType", this->RspType); }
        if (Valid(this->RspData)) { query.bindValue(":RspData", this->RspData); }
        //
        isOk = query.exec();
        if (isOk && lastInsertId) { *lastInsertId = query.lastInsertId().toLongLong(); }
        return isOk;
    }
    bool update_data(QSqlQuery& query, QString& whereCond)
    {
        bool isOk = false;
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){cols.append\("$1=:$1"\);}】.
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){query.bindValue\(":$1",this->$1\);}】.
        QStringList cols;
        if (Valid(this->ID)) { cols.append("ID=:ID"); }
        if (Valid(this->RequestID)) { cols.append("RequestID=:RequestID"); }
        if (Valid(this->SeqNo)) { cols.append("SeqNo=:SeqNo"); }
        if (Valid(this->UserID)) { cols.append("UserID=:UserID"); }
        if (Valid(this->ReqType)) { cols.append("ReqType=:ReqType"); }
        if (Valid(this->ReqData)) { cols.append("ReqData=:ReqData"); }
        if (Valid(this->ReqTime)) { cols.append("ReqTime=:ReqTime"); }
        if (Valid(this->CreateTime)) { cols.append("CreateTime=:CreateTime"); }
        if (Valid(this->State)) { cols.append("State=:State"); }
        if (Valid(this->ErrNo)) { cols.append("ErrNo=:ErrNo"); }
        if (Valid(this->ErrMsg)) { cols.append("ErrMsg=:ErrMsg"); }
        if (Valid(this->RspType)) { cols.append("RspType=:RspType"); }
        if (Valid(this->RspData)) { cols.append("RspData=:RspData"); }
        //
        QString sqlStr = QObject::tr("UPDATE %1 SET %2").arg(static_table_name()).arg(cols.join(" AND "));
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->ID)) { query.bindValue(":ID", this->ID); }
        if (Valid(this->RequestID)) { query.bindValue(":RequestID", this->RequestID); }
        if (Valid(this->SeqNo)) { query.bindValue(":SeqNo", this->SeqNo); }
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->ReqType)) { query.bindValue(":ReqType", this->ReqType); }
        if (Valid(this->ReqData)) { query.bindValue(":ReqData", this->ReqData); }
        if (Valid(this->ReqTime)) { query.bindValue(":ReqTime", this->ReqTime); }
        if (Valid(this->CreateTime)) { query.bindValue(":CreateTime", this->CreateTime); }
        if (Valid(this->State)) { query.bindValue(":State", this->State); }
        if (Valid(this->ErrNo)) { query.bindValue(":ErrNo", this->ErrNo); }
        if (Valid(this->ErrMsg)) { query.bindValue(":ErrMsg", this->ErrMsg); }
        if (Valid(this->RspType)) { query.bindValue(":RspType", this->RspType); }
        if (Valid(this->RspData)) { query.bindValue(":RspData", this->RspData); }
        //
        return query.exec();
    }
};

class QCommonNtosReq
{
public:
    int64_t    RefNum;
    int64_t    RequestID;
    QString    UserID;
    int64_t    SeqNo;
    int32_t    ReqType;
    QByteArray ReqData;
    QDateTime  ReqTime;
public:
    QCommonNtosReq()
    {
        this->RefNum = INT64_MAX;
        this->RequestID = INT64_MAX;
        this->UserID.clear();
        this->SeqNo = INT64_MAX;
        this->ReqType = INT32_MAX;
        this->ReqData.clear();
        this->ReqTime = QDateTime();
    }
public:
    static QString static_table_name()
    {
        return "QCommonNtosReq";//可能会有(object_table_name)函数.
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            RefNum    INTEGER NOT NULL PRIMARY KEY ,\
            RequestID INTEGER     NULL ,\
            UserID    TEXT    NOT NULL ,\
            SeqNo     INTEGER     NULL UNIQUE ,\
            ReqType   INTEGER NOT NULL ,\
            ReqData   BLOB        NULL ,\
            ReqTime   TEXT        NULL )"
        ).QString::arg(static_table_name());
        return sql;
    }
    bool insert_data(QSqlQuery& query, int64_t* lastInsertId = nullptr)
    {
        //请外部保证在同一个(上下文/先后顺序/总之就是加锁的意思).
        bool isOk = false;
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){cols.append\("$1"\);}】.
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){query.bindValue\(":$1",this->$1\);}】.
        QStringList cols;
        if (Valid(this->RefNum)) { cols.append("RefNum"); }
        if (Valid(this->RequestID)) { cols.append("RequestID"); }
        if (Valid(this->UserID)) { cols.append("UserID"); }
        if (Valid(this->SeqNo)) { cols.append("SeqNo"); }
        if (Valid(this->ReqType)) { cols.append("ReqType"); }
        if (Valid(this->ReqData)) { cols.append("ReqData"); }
        if (Valid(this->ReqTime)) { cols.append("ReqTime"); }
        //
        QString sqlStr = QObject::tr("INSERT INTO %1 (%2) VALUES (%3)").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->RefNum)) { query.bindValue(":RefNum", this->RefNum); }
        if (Valid(this->RequestID)) { query.bindValue(":RequestID", this->RequestID); }
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->SeqNo)) { query.bindValue(":SeqNo", this->SeqNo); }
        if (Valid(this->ReqType)) { query.bindValue(":ReqType", this->ReqType); }
        if (Valid(this->ReqData)) { query.bindValue(":ReqData", this->ReqData); }
        if (Valid(this->ReqTime)) { query.bindValue(":ReqTime", this->ReqTime); }
        //
        isOk = query.exec();
        if (isOk && lastInsertId) { *lastInsertId = query.lastInsertId().toLongLong(); }
        return isOk;
    }
    static bool select_data(QSqlQuery& query, int64_t RefNum, QList<QCommonNtosReq>& dataOut)
    {
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【fromQVariant\(curData.$1,query.value\("$1"\)\);】.
        QString sql = QObject::tr("SELECT * FROM %1 WHERE RefNum=%2").QString::arg(static_table_name()).arg(RefNum);
        if (query.exec(sql) == false)
            return false;
        while (query.next()) {
            QCommonNtosReq curData;
            fromQVariant(curData.RefNum, query.value("RefNum"));
            fromQVariant(curData.RequestID, query.value("RequestID"));
            fromQVariant(curData.UserID, query.value("UserID"));
            fromQVariant(curData.SeqNo, query.value("SeqNo"));
            fromQVariant(curData.ReqType, query.value("ReqType"));
            fromQVariant(curData.ReqData, query.value("ReqData"));
            fromQVariant(curData.ReqTime, query.value("ReqTime"));
            dataOut.append(curData);
        }
        return true;
    }
};
Q_DECLARE_METATYPE(QCommonNtosReq);

class QCommonNtosRsp
{
public:
    int64_t    ID;
    QDateTime  InsertTime;
    int64_t    RequestID;
    QString    Pathway;
    int64_t    SeqNo;
    int32_t    RspType;
    QByteArray RspData;
    QDateTime  RspTime;
    int32_t    FromServer;
    int32_t    ErrNo;
    QString    ErrMsg;
    int64_t    RefNum;
public:
    QCommonNtosRsp()
    {
        this->ID = INT64_MAX;
        this->InsertTime = QDateTime();
        this->RequestID = INT64_MAX;
        this->Pathway.clear();
        this->SeqNo = INT64_MAX;
        this->RspType = INT32_MAX;
        this->RspData.clear();
        this->RspTime = QDateTime();
        this->FromServer = INT32_MAX;
        this->ErrNo = INT32_MAX;
        this->ErrMsg.clear();
        this->RefNum = INT64_MAX;
    }
public:
    static QString static_table_name()
    {
        return "QCommonNtosRsp";//可能会有(object_table_name)函数.
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            ID         INTEGER NOT NULL PRIMARY KEY ,\
            InsertTime TEXT    NOT NULL ,\
            RequestID  INTEGER     NULL ,\
            Pathway    TEXT        NULL ,\
            SeqNo      INTEGER     NULL ,\
            RspType    INTEGER     NULL ,  /*not_null*/ \
            RspData    BLOB        NULL ,\
            RspTime    TEXT        NULL ,  /*not_null*/ \
            FromServer INTEGER     NULL ,  /*not_null*/ \
            ErrNo      INTEGER     NULL ,  /*not_null*/ \
            ErrMsg     TEXT        NULL ,\
            RefNum     INTEGER     NULL    /*not_null*/ )"
        ).QString::arg(static_table_name());
        return sql;
    }
    bool insert_data(QSqlQuery& query)
    {
        this->InsertTime = QDateTime::currentDateTime();
        //请外部保证在同一个(上下文/先后顺序/总之就是加锁的意思).
        bool isOk = false;
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){cols.append\("$1"\);}】.
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){query.bindValue\(":$1",this->$1\);}】.
        QStringList cols;
        if (Valid(this->ID)) { cols.append("ID"); }
        if (Valid(this->InsertTime)) { cols.append("InsertTime"); }
        if (Valid(this->RequestID)) { cols.append("RequestID"); }
        if (Valid(this->Pathway)) { cols.append("Pathway"); }
        if (Valid(this->SeqNo)) { cols.append("SeqNo"); }
        if (Valid(this->RspType)) { cols.append("RspType"); }
        if (Valid(this->RspData)) { cols.append("RspData"); }
        if (Valid(this->RspTime)) { cols.append("RspTime"); }
        if (Valid(this->FromServer)) { cols.append("FromServer"); }
        if (Valid(this->ErrNo)) { cols.append("ErrNo"); }
        if (Valid(this->ErrMsg)) { cols.append("ErrMsg"); }
        if (Valid(this->RefNum)) { cols.append("RefNum"); }
        //
        QString sqlStr = QObject::tr("INSERT INTO %1 (%2) VALUES (%3)").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->ID)) { query.bindValue(":ID", this->ID); }
        if (Valid(this->InsertTime)) { query.bindValue(":InsertTime", this->InsertTime); }
        if (Valid(this->RequestID)) { query.bindValue(":RequestID", this->RequestID); }
        if (Valid(this->Pathway)) { query.bindValue(":Pathway", this->Pathway); }
        if (Valid(this->SeqNo)) { query.bindValue(":SeqNo", this->SeqNo); }
        if (Valid(this->RspType)) { query.bindValue(":RspType", this->RspType); }
        if (Valid(this->RspData)) { query.bindValue(":RspData", this->RspData); }
        if (Valid(this->RspTime)) { query.bindValue(":RspTime", this->RspTime); }
        if (Valid(this->FromServer)) { query.bindValue(":FromServer", this->FromServer); }
        if (Valid(this->ErrNo)) { query.bindValue(":ErrNo", this->ErrNo); }
        if (Valid(this->ErrMsg)) { query.bindValue(":ErrMsg", this->ErrMsg); }
        if (Valid(this->RefNum)) { query.bindValue(":RefNum", this->RefNum); }
        //
        return query.exec();
    }
    static bool select_stat_info(QSqlQuery& query, int64_t RefNum, int&cntTime, QString&minTime, QString& maxTime)
    {
        bool isOK = false;
        QString sqlStr = QString("SELECT COUNT(InsertTime) AS CntTime, MIN(InsertTime) AS MinTime, MAX(InsertTime) AS MaxTime FROM QCommonNtosRsp WHERE RefNum=%1").arg(RefNum);
        isOK = query.exec(sqlStr);
        while (query.next()) {
            fromQVariant(cntTime, query.value("CntTime"));
            fromQVariant(minTime, query.value("MinTime"));
            fromQVariant(maxTime, query.value("MaxTime"));
        }
        return isOK;
    }
    static bool select_data(QSqlQuery& query, int64_t RefNum, QList<QCommonNtosRsp>& dataOut)
    {
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【fromQVariant\(curData.$1,query.value\("$1"\)\);】.
        QString sql = QObject::tr("SELECT * FROM %1 WHERE RefNum=%2 ORDER BY InsertTime ASC").QString::arg(static_table_name()).arg(RefNum);
        if (query.exec(sql) == false)
            return false;
        while (query.next()) {
            QCommonNtosRsp curData;
            fromQVariant(curData.ID, query.value("ID"));
            fromQVariant(curData.InsertTime, query.value("InsertTime"));
            fromQVariant(curData.RequestID, query.value("RequestID"));
            fromQVariant(curData.Pathway, query.value("Pathway"));
            fromQVariant(curData.SeqNo, query.value("SeqNo"));
            fromQVariant(curData.RspType, query.value("RspType"));
            fromQVariant(curData.RspData, query.value("RspData"));
            fromQVariant(curData.RspTime, query.value("RspTime"));
            fromQVariant(curData.FromServer, query.value("FromServer"));
            fromQVariant(curData.ErrNo, query.value("ErrNo"));
            fromQVariant(curData.ErrMsg, query.value("ErrMsg"));
            fromQVariant(curData.RefNum, query.value("RefNum"));
            dataOut.append(curData);
        }
        return true;
    }
};
Q_DECLARE_METATYPE(QCommonNtosRsp);

#endif // SQL_STRUCT_H
