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
#include <QDebug>

namespace {
    //使用非const的引用,可以强校验入参的类型,基本可以杜绝隐式类型转换.
    inline bool Valid(bool& /*data*/) { return true; }
    inline bool Valid(int32_t& data) { return INT32_MAX != data; }
    inline bool Valid(int64_t& data) { return INT64_MAX != data; }
    inline bool Valid(float& data) { return FLT_MAX != data; }
    inline bool Valid(double& data) { return DBL_MAX <= data; }
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
    KeyValue() {}
    KeyValue(const QString& k, const QString& v) :Key(k), Value(v) {}
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
    bool refresh_data(QSqlQuery& query)
    {
        bool isOk = false;
        QString sqlStr = QObject::tr("SELECT Value FROM %1 WHERE Key=:Key").arg(static_table_name());
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        query.bindValue(":Key", this->Key);
        isOk = query.exec();
        Q_ASSERT(isOk);
        isOk = query.next();
        this->Value = isOk ? query.value("Value").toString() : QString("");//QString()时insert失败.
        Q_ASSERT(query.next() == false);
        return isOk;
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
    static QString select_data(QSqlQuery& query, const QString& Key, bool* isExists = nullptr)
    {
        QString Value;
        QString sqlStr = QObject::tr("SELECT Value FROM %1 WHERE Key=:Key").arg(static_table_name());
        bool isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        query.bindValue(":Key", Key);
        isOk = query.exec();
        Q_ASSERT(isOk);
        if (query.next())
        {
            if (isExists) { *isExists = true; }
            Value = query.value("Value").toString();
        }
        while (query.next()) { Q_ASSERT(false); }//肯定走不到这里.
        return Value;
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

class CommonData
{
public:
    bool       SendOK;    //Req发送成功.
    int32_t    RspCnt;    //与Req对应的Rsp的Cnt.
    int32_t    MsgType;   //Common2Req,Common2Rsp,Common1Req,Common1Rsp
    QString    MsgTypeTxt;//MsgType的文本.
    QString    PeerID;//对端(参考python3的[help(socket.socket.getpeername)]).
    QString    UserID;//本端.
    int64_t    MsgNo;
    int32_t    SeqNo;
    int64_t    BatchNo;
    int64_t    RefNum;
    QString    RefText;
    QString    SenderID;
    QString    RecverID;
    bool       ToRoot;
    bool       IsLog;
    bool       IsSafe;
    bool       IsPush;
    bool       UpCache;
    int32_t    TxType;//通信的对象的类型.
    QString    TxTypeTxt;
    QByteArray TxData;//通信的对象经pb序列化后的二进制数据.
    QString    TxDataTxt;//通信的对象转换成json字符串.
    QDateTime  TxTime;
    QDateTime  InsertTime;//插入时刻(插入之后,不再修改它).
    bool       IsLast;
public:
    CommonData()
    {
        this->SendOK = false;
        this->RspCnt = INT32_MAX;
        this->MsgType = INT32_MAX;
        this->MsgNo = INT32_MAX;
        this->SeqNo = INT32_MAX;
        this->BatchNo = INT64_MAX;
        this->RefNum = INT64_MAX;
        this->ToRoot = false;
        this->IsLog = false;
        this->IsSafe = false;
        this->IsPush = false;
        this->UpCache = false;
        this->TxType = INT32_MAX;
        this->IsLast = false;
    }
public:
    static QString static_table_name()
    {
        return "CommonData";
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            SendOK     INTEGER     NULL ,\
            RspCnt     INTEGER     NULL ,\
            MsgType    INTEGER     NULL ,\
            MsgTypeTxt TEXT        NULL ,\
            PeerID     TEXT    NOT NULL ,\
            UserID     TEXT    NOT NULL ,\
            MsgNo      INTEGER NOT NULL ,\
            SeqNo      INTEGER NOT NULL ,\
            BatchNo    INTEGER     NULL ,\
            RefNum     INTEGER     NULL ,\
            RefText    TEXT        NULL ,\
            SenderID   TEXT        NULL ,\
            RecverID   TEXT        NULL ,\
            ToRoot     INTEGER     NULL ,\
            IsLog      INTEGER     NULL ,\
            IsSafe     INTEGER     NULL ,\
            IsPush     INTEGER     NULL ,\
            UpCache    INTEGER     NULL ,\
            TxType     INTEGER     NULL ,\
            TxTypeTxt  TEXT        NULL ,\
            TxData     BLOB        NULL ,\
            TxDataTxt  TEXT        NULL ,\
            TxTime     TEXT        NULL ,\
            InsertTime TEXT        NULL ,\
            IsLast     INTEGER     NULL )"
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
        if (Valid(this->SendOK)) { cols.append("SendOK"); }
        if (Valid(this->RspCnt)) { cols.append("RspCnt"); }
        if (Valid(this->MsgType)) { cols.append("MsgType"); }
        if (Valid(this->MsgTypeTxt)) { cols.append("MsgTypeTxt"); }
        if (Valid(this->PeerID)) { cols.append("PeerID"); }
        if (Valid(this->UserID)) { cols.append("UserID"); }
        if (Valid(this->MsgNo)) { cols.append("MsgNo"); }
        if (Valid(this->SeqNo)) { cols.append("SeqNo"); }
        if (Valid(this->BatchNo)) { cols.append("BatchNo"); }
        if (Valid(this->RefNum)) { cols.append("RefNum"); }
        if (Valid(this->RefText)) { cols.append("RefText"); }
        if (Valid(this->SenderID)) { cols.append("SenderID"); }
        if (Valid(this->RecverID)) { cols.append("RecverID"); }
        if (Valid(this->ToRoot)) { cols.append("ToRoot"); }
        if (Valid(this->IsLog)) { cols.append("IsLog"); }
        if (Valid(this->IsSafe)) { cols.append("IsSafe"); }
        if (Valid(this->IsPush)) { cols.append("IsPush"); }
        if (Valid(this->UpCache)) { cols.append("UpCache"); }
        if (Valid(this->TxType)) { cols.append("TxType"); }
        if (Valid(this->TxTypeTxt)) { cols.append("TxTypeTxt"); }
        if (Valid(this->TxData)) { cols.append("TxData"); }
        if (Valid(this->TxDataTxt)) { cols.append("TxDataTxt"); }
        if (Valid(this->TxTime)) { cols.append("TxTime"); }
        if (Valid(this->InsertTime)) { cols.append("InsertTime"); }
        if (Valid(this->IsLast)) { cols.append("IsLast"); }
        //
        QString sqlStr = QObject::tr("INSERT %1 INTO %2 (%3) VALUES (%4)").arg(insertNotReplace ? "" : "OR REPLACE").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->SendOK)) { query.bindValue(":SendOK", this->SendOK); }
        if (Valid(this->RspCnt)) { query.bindValue(":RspCnt", this->RspCnt); }
        if (Valid(this->MsgType)) { query.bindValue(":MsgType", this->MsgType); }
        if (Valid(this->MsgTypeTxt)) { query.bindValue(":MsgTypeTxt", this->MsgTypeTxt); }
        if (Valid(this->PeerID)) { query.bindValue(":PeerID", this->PeerID); }
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->MsgNo)) { query.bindValue(":MsgNo", this->MsgNo); }
        if (Valid(this->SeqNo)) { query.bindValue(":SeqNo", this->SeqNo); }
        if (Valid(this->BatchNo)) { query.bindValue(":BatchNo", this->BatchNo); }
        if (Valid(this->RefNum)) { query.bindValue(":RefNum", this->RefNum); }
        if (Valid(this->RefText)) { query.bindValue(":RefText", this->RefText); }
        if (Valid(this->SenderID)) { query.bindValue(":SenderID", this->SenderID); }
        if (Valid(this->RecverID)) { query.bindValue(":RecverID", this->RecverID); }
        if (Valid(this->ToRoot)) { query.bindValue(":ToRoot", this->ToRoot); }
        if (Valid(this->IsLog)) { query.bindValue(":IsLog", this->IsLog); }
        if (Valid(this->IsSafe)) { query.bindValue(":IsSafe", this->IsSafe); }
        if (Valid(this->IsPush)) { query.bindValue(":IsPush", this->IsPush); }
        if (Valid(this->UpCache)) { query.bindValue(":UpCache", this->UpCache); }
        if (Valid(this->TxType)) { query.bindValue(":TxType", this->TxType); }
        if (Valid(this->TxTypeTxt)) { query.bindValue(":TxTypeTxt", this->TxTypeTxt); }
        if (Valid(this->TxData)) { query.bindValue(":TxData", this->TxData); }
        if (Valid(this->TxDataTxt)) { query.bindValue(":TxDataTxt", this->TxDataTxt); }
        if (Valid(this->TxTime)) { query.bindValue(":TxTime", this->TxTime); }
        if (Valid(this->InsertTime)) { query.bindValue(":InsertTime", this->InsertTime); }
        if (Valid(this->IsLast)) { query.bindValue(":IsLast", this->IsLast); }
        //
        isOk = query.exec();
        if (isOk && lastInsertId) { *lastInsertId = query.lastInsertId().toLongLong(); }
        return isOk;
    }
    bool updateRequestData(QSqlQuery& query)
    {
        if (Valid(this->SeqNo) && (this->SeqNo == 0)) { return false; }
        QString sqlStr = QString("UPDATE %1 SET RspCnt=RspCnt+1, IsLast=MAX(IsLast,:IsLast) WHERE UserID=:UserID AND MsgNo=:MsgNo AND SeqNo=0").arg(static_table_name());
        bool isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        query.bindValue(":IsLast", this->IsLast);
        query.bindValue(":UserID", this->UserID);
        query.bindValue(":MsgNo", this->MsgNo);
        isOk = query.exec();
        return isOk;
    }
    static bool delete_data(QSqlQuery& query, const QString& whereCond)
    {
        QString sqlStr = QObject::tr("DELETE FROM %1").arg(static_table_name());
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        return query.exec(sqlStr);
    }
};
Q_DECLARE_METATYPE(CommonData);

class PushWrap
{
public:
    QString    PeerID;
    int64_t    MsgNo;
    QString    UserID;
    QDateTime  PshTime;
    int32_t    PshType;
    QString    PshTypeTxt;
    QByteArray PshData;   //通信的对象经pb序列化后的二进制数据.
    QString    PshDataTxt;//通信的对象转换成json字符串.
public:
    PushWrap()
    {
        this->MsgNo = INT64_MAX;
        this->PshType = INT32_MAX;
    }
public:
    static QString static_table_name()
    {
        return "PushWrap";
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            PeerID     TEXT    NOT NULL,\
            MsgNo      INTEGER NOT NULL,\
            UserID     TEXT        NULL,\
            PshTime    TEXT        NULL,\
            PshType    INTEGER     NULL,\
            PshTypeTxt TEXT        NULL,\
            PshData    BLOB        NULL,\
            PshDataTxt TEXT        NULL,\
              PRIMARY KEY(PeerID,MsgNo) )"
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
        if (Valid(this->PeerID)) { cols.append("PeerID"); }
        if (Valid(this->MsgNo)) { cols.append("MsgNo"); }
        if (Valid(this->UserID)) { cols.append("UserID"); }
        if (Valid(this->PshTime)) { cols.append("PshTime"); }
        if (Valid(this->PshType)) { cols.append("PshType"); }
        if (Valid(this->PshTypeTxt)) { cols.append("PshTypeTxt"); }
        if (Valid(this->PshData)) { cols.append("PshData"); }
        if (Valid(this->PshDataTxt)) { cols.append("PshDataTxt"); }
        //
        QString sqlStr = QObject::tr("INSERT %1 INTO %2 (%3) VALUES (%4)").arg(insertNotReplace ? "" : "OR REPLACE").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->PeerID)) { query.bindValue(":PeerID", this->PeerID); }
        if (Valid(this->MsgNo)) { query.bindValue(":MsgNo", this->MsgNo); }
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->PshTime)) { query.bindValue(":PshTime", this->PshTime); }
        if (Valid(this->PshType)) { query.bindValue(":PshType", this->PshType); }
        if (Valid(this->PshTypeTxt)) { query.bindValue(":PshTypeTxt", this->PshTypeTxt); }
        if (Valid(this->PshData)) { query.bindValue(":PshData", this->PshData); }
        if (Valid(this->PshDataTxt)) { query.bindValue(":PshDataTxt", this->PshDataTxt); }
        //
        isOk = query.exec();
        if (isOk && lastInsertId) { *lastInsertId = query.lastInsertId().toLongLong(); }
        return isOk;
    }
    static bool select_data(QSqlQuery& query, const QString& whereCond, QList<PushWrap>& dataOut)
    {
        //查找【^.+? ([a-zA-Z0-9_]+);.*$】替换【fromQVariant\(curData.$1,query.value\("$1"\)\);】.
        QString sqlStr = QObject::tr("SELECT * FROM %1").QString::arg(static_table_name());
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        if (query.exec(sqlStr) == false)
            return false;
        while (query.next()) {
            PushWrap curData;
            fromQVariant(curData.PeerID, query.value("PeerID"));
            fromQVariant(curData.MsgNo, query.value("MsgNo"));
            fromQVariant(curData.UserID, query.value("UserID"));
            fromQVariant(curData.PshTime, query.value("PshTime"));
            fromQVariant(curData.PshType, query.value("PshType"));
            fromQVariant(curData.PshTypeTxt, query.value("PshTypeTxt"));
            fromQVariant(curData.PshData, query.value("PshData"));
            fromQVariant(curData.PshDataTxt, query.value("PshDataTxt"));
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
Q_DECLARE_METATYPE(PushWrap);

class PathwayInfo
{
public:
    QString    UserID;
    QString    Pathway;
public:
    static QString static_table_name()
    {
        return "PathwayInfo";
    }
    static QString static_create_sql()
    {
        QString sql = QObject::tr(
            "CREATE TABLE IF NOT EXISTS %1 (\
            UserID  TEXT NOT NULL PRIMARY KEY ,\
            Pathway TEXT     NULL )"
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
        if (Valid(this->Pathway)) { cols.append("Pathway"); }
        //
        QString sqlStr = QObject::tr("INSERT %1 INTO %2 (%3) VALUES (%4)").arg(insertNotReplace ? "" : "OR REPLACE").arg(static_table_name()).arg(cols.join(',')).arg(":" + cols.join(", :"));
        isOk = query.prepare(sqlStr);
        Q_ASSERT(isOk);
        //
        if (Valid(this->UserID)) { query.bindValue(":UserID", this->UserID); }
        if (Valid(this->Pathway)) { query.bindValue(":Pathway", this->Pathway); }
        //
        isOk = query.exec();
        if (isOk && lastInsertId) { *lastInsertId = query.lastInsertId().toLongLong(); }
        return isOk;
    }
    static bool delete_data(QSqlQuery& query, const QString& whereCond)
    {
        QString sqlStr = QObject::tr("DELETE FROM %1").arg(static_table_name());
        if (!whereCond.isEmpty()) { sqlStr += " WHERE " + whereCond; }
        return query.exec(sqlStr);
    }
};
Q_DECLARE_METATYPE(PathwayInfo);

#endif // SQL_STRUCT_H
