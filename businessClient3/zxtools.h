#ifndef ZXTOOLS_H
#define ZXTOOLS_H
#include  "sqlstruct.h"
#include "protobuf/m2b.h"

class zxtools
{
public:
    static void qdt2gpt(::google::protobuf::Timestamp& gptDst, const QDateTime& qdtSrc);
    static void gpt2qdt(QDateTime& qdtDst, const ::google::protobuf::Timestamp& gptSrc);
    static QString MsgTypeName2MsgClassName(const QString& msgTypeName);
    static GPMSGPTR name2object(const std::string& name);
    static QString object2json(const google::protobuf::Message &msgObj, bool *isOk = nullptr);
    static QString binary2json(txdata::MsgType msgType, const std::string& binData, bool *isOk = nullptr);
    static GPMSGPTR json2object(const QString& msgTypeName, const QString& jsonText, txdata::MsgType* msgType = nullptr);
    static bool json2binary(const QString& msgTypeName, const QString& jsonText, txdata::MsgType& msgType, std::string& binData);
    static void Common1Req2CommonData(CommonData* dst, const txdata::Common1Req* src);
    static void CommonData2Common1Req(txdata::Common1Req* dst, const CommonData* src);
    static void Common1Rsp2CommonData(CommonData* dst, const txdata::Common1Rsp* src);
    static void CommonData2Common1Rsp(txdata::Common1Rsp* dst, const CommonData* src);
    static void Common2Req2CommonData(CommonData* dst, const txdata::Common2Req* src);
    static void Common2Rsp2CommonData(CommonData* dst, const txdata::Common2Rsp* src);
};

#endif // ZXTOOLS_H
