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
    inline bool Valid(QStringList& data) { return true; }
    inline bool Valid(QByteArray& data) { return !data.isNull(); }
    inline bool Valid(QDateTime& data) { return data.isValid(); }
    inline void fromQVariant(int32_t& dst, const QVariant& src) { dst = src.toInt(nullptr); }
    inline void fromQVariant(int64_t& dst, const QVariant& src) { dst = src.toLongLong(nullptr); }
    inline void fromQVariant(QString& dst, const QVariant& src) { dst = src.toString(); }
    inline void fromQVariant(QStringList& dst, const QVariant& src) { dst = src.toStringList(); }
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

class ConnInfoEx
{
public:
    QString UserID;//PK
    QString BelongID;
    QString Version;
    int32_t LinkMode;
    int32_t ExePid;
    QString ExePath;
    QString Remark;
    QString Pathway;
public:
    ConnInfoEx()
    {
        this->LinkMode = INT32_MAX;
        this->ExePid = INT32_MAX;
    }
public:
    static QString static_table_name()
    {
        return "ConnInfoEx";
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            UserID   TEXT    NOT NULL PRIMARY KEY ,\
            BelongID TEXT    NOT NULL ,\
            Version  TEXT        NULL ,\
            LinkMode INTEGER     NULL ,\
            ExePid   INTEGER     NULL ,\
            ExePath  TEXT        NULL ,\
            Remark   TEXT        NULL ,\
            Pathway  TEXT        NULL )"
        ).QString::arg(static_table_name());
        return  sql;
    }
    bool insert_data(QSqlQuery& query, bool insertNotReplace, int64_t* lastInsertId = nullptr)
    {
        //请外部保证在同一个(上下文/先后顺序/总之就是加锁的意思).
        bool isOk = false;
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){cols.append\("$1"\);}】.
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){query.bindValue\(":$1",this->$1\);}】.
        //注意(NOT NULL)要特殊处理.
        QStringList cols;
        if (Valid(this->UserID)) { cols.append("UserID"); }
        if (Valid(this->BelongID)) { cols.append("BelongID"); }
        if (Valid(this->Version)) { cols.append("Version"); }
        if (Valid(this->LinkMode)) { cols.append("LinkMode"); }
        if (Valid(this->ExePid)) { cols.append("ExePid"); }
        if (Valid(this->ExePath)) { cols.append("ExePath"); }
        if (Valid(this->Remark)) { cols.append("Remark"); }
        if (Valid(this->Pathway)) { cols.append("Pathway"); }
        //
        QString sqlStr = QObject::tr("INSERT %1 INTO %2 (%3) VALUES (%4)").arg(insertNotReplace ? "" : "OR REPLACE").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->BelongID)) { query.bindValue(":BelongID", this->BelongID); }
        if (Valid(this->Version)) { query.bindValue(":Version", this->Version); }
        if (Valid(this->LinkMode)) { query.bindValue(":LinkMode", this->LinkMode); }
        if (Valid(this->ExePid)) { query.bindValue(":ExePid", this->ExePid); }
        if (Valid(this->ExePath)) { query.bindValue(":ExePath", this->ExePath); }
        if (Valid(this->Remark)) { query.bindValue(":Remark", this->Remark); }
        if (Valid(this->Pathway)) { query.bindValue(":Pathway", this->Pathway); }
        //
        isOk = query.exec();
        if (isOk && lastInsertId) { *lastInsertId = query.lastInsertId().toLongLong(); }
        return isOk;
    }
    bool update_data(QSqlQuery& query, const QString& whereCond)
    {
        bool isOk = false;
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){cols.append\("$1=:$1"\);}】.
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【if\(Valid\(this->$1\)\){query.bindValue\(":$1",this->$1\);}】.
        QStringList cols;
        if (Valid(this->UserID)) { cols.append("UserID=:UserID"); }
        if (Valid(this->BelongID)) { cols.append("BelongID=:BelongID"); }
        if (Valid(this->Version)) { cols.append("Version=:Version"); }
        if (Valid(this->LinkMode)) { cols.append("LinkMode=:LinkMode"); }
        if (Valid(this->ExePid)) { cols.append("ExePid=:ExePid"); }
        if (Valid(this->ExePath)) { cols.append("ExePath=:ExePath"); }
        if (Valid(this->Remark)) { cols.append("Remark=:Remark"); }
        if (Valid(this->Pathway)) { cols.append("Pathway=:Pathway"); }
        //
        QString sqlStr = QObject::tr("UPDATE %1 SET %2").arg(static_table_name()).arg(cols.join(" AND "));
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->BelongID)) { query.bindValue(":BelongID", this->BelongID); }
        if (Valid(this->Version)) { query.bindValue(":Version", this->Version); }
        if (Valid(this->LinkMode)) { query.bindValue(":LinkMode", this->LinkMode); }
        if (Valid(this->ExePid)) { query.bindValue(":ExePid", this->ExePid); }
        if (Valid(this->ExePath)) { query.bindValue(":ExePath", this->ExePath); }
        if (Valid(this->Remark)) { query.bindValue(":Remark", this->Remark); }
        if (Valid(this->Pathway)) { query.bindValue(":Pathway", this->Pathway); }
        //
        return query.exec();
    }
    static bool select_data(QSqlQuery& query, const QString& whereCond, QList<ConnInfoEx>& dataOut)
    {
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【fromQVariant\(curData.$1,query.value\("$1"\)\);】.
        QString sqlStr = QObject::tr("SELECT * FROM %1").QString::arg(static_table_name());
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        if (query.exec(sqlStr) == false)
            return false;
        while (query.next()) {
            ConnInfoEx curData;
            fromQVariant(curData.UserID, query.value("UserID"));
            fromQVariant(curData.BelongID, query.value("BelongID"));
            fromQVariant(curData.Version, query.value("Version"));
            fromQVariant(curData.LinkMode, query.value("LinkMode"));
            fromQVariant(curData.ExePid, query.value("ExePid"));
            fromQVariant(curData.ExePath, query.value("ExePath"));
            fromQVariant(curData.Remark, query.value("Remark"));
            fromQVariant(curData.Pathway, query.value("Pathway"));
            dataOut.append(curData);
        }
        return true;
    }
    static bool delete_data(QSqlQuery& query, const QString& whereCond)
    {
        QString sqlStr = QObject::tr("DELETE FROM %1").arg(static_table_name());
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        return query.exec(sqlStr);
    }
};
Q_DECLARE_METATYPE(ConnInfoEx);

#endif // SQL_STRUCT_H
