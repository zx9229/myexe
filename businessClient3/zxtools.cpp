// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Timestamp
#include <QDebug>
#include "google/protobuf/util/json_util.h"
#include "zxtools.h"

void zxtools::qdt2gpt(::google::protobuf::Timestamp& gptDst, const QDateTime& qdtSrc)
{
    //因为 iVal = 5; 即[dst = src], 估计因为这个原因, 标准库里的一些函数也采用了这个布局, 比如strcpy函数:
    //char * strcpy(char * _Dest,const char * _Source);
    //boost库也大多使用该布局, 比如split:
    //inline SequenceSequenceT& split(SequenceSequenceT& Result, RangeT& Input, PredicateT Pred, token_compress_mode_type eCompress=token_compress_off);
    //所以这里也(尽量)采用这种方式.
    gptDst.set_seconds(qdtSrc.toSecsSinceEpoch());
    gptDst.set_nanos(qdtSrc.time().msec() * 1000 * 1000);
}

void zxtools::gpt2qdt(QDateTime& qdtDst, const ::google::protobuf::Timestamp& gptSrc)
{
    qdtDst = QDateTime::fromMSecsSinceEpoch(gptSrc.seconds() * 1000 + gptSrc.nanos() / 1000 / 1000);
}

GPMSGPTR zxtools::name2object(const std::string& name)
{
    //(::txdata::Common1Req)的名字是(txdata.Common1Req)
    // https://blog.csdn.net/riopho/article/details/80372510
    const google::protobuf::Descriptor* desc = google::protobuf::DescriptorPool::generated_pool()->FindMessageTypeByName(name);
    if (nullptr == desc) { return nullptr; } // desc->index();
    google::protobuf::Message* object = google::protobuf::MessageFactory::generated_factory()->GetPrototype(desc)->New();
    return (object ? GPMSGPTR(object) : nullptr);
}

QString zxtools::MsgTypeName2MsgClassName(const QString& msgTypeName)
{
    const static QString ID_("ID_");
    return msgTypeName.startsWith(ID_) ? ("txdata." + msgTypeName.mid(ID_.size())) : msgTypeName;
}

QString zxtools::object2json(const google::protobuf::Message &msgObj, bool *isOk)
{
    if (isOk) { *isOk = true; }
    google::protobuf::util::JsonOptions jsonOpt;
    if (true) {
        jsonOpt.add_whitespace = true;
        jsonOpt.always_print_primitive_fields = true;
        jsonOpt.preserve_proto_field_names = true;
    }
    std::string jsonStr;
    if (google::protobuf::util::MessageToJsonString(msgObj, &jsonStr, jsonOpt) != google::protobuf::util::Status::OK)
    {
        jsonStr.clear();
        if (isOk) { *isOk = false; }
    }
    return QString::fromStdString(jsonStr).trimmed();
}

QString zxtools::binary2json(txdata::MsgType msgType, const std::string& binData, bool *isOk)
{
    GPMSGPTR msgObj;
    if (m2b::slice2msg(binData, msgType, msgObj))
    {
        return object2json(*msgObj, isOk);
    }
    else
    {
        if (isOk) { *isOk = false; }
        return "";
    }
}

GPMSGPTR zxtools::json2object(const QString& msgTypeName, const QString& jsonText, txdata::MsgType* msgType)
{
    GPMSGPTR curObject;
    txdata::MsgType curMsgType = txdata::MsgType::Zero1;
    do
    {
        if (txdata::MsgType_Parse(msgTypeName.toStdString(), &curMsgType) == false)
            break;
        QString curMsgClassName = MsgTypeName2MsgClassName(msgTypeName);
        curObject = name2object(curMsgClassName.toStdString());
        if (nullptr == curObject)
            break;
        if (google::protobuf::util::JsonStringToMessage(jsonText.toStdString(), curObject.data()) != google::protobuf::util::Status::OK)
        {
            curObject.clear();
            break;
        }
    } while (false);
    if (msgType) { *msgType = curMsgType; }
    return curObject;
}

bool zxtools::json2binary(const QString& msgTypeName, const QString& jsonText, txdata::MsgType& msgType, std::string& binData)
{
    bool isOk = false;
    do
    {
        GPMSGPTR msgObj = json2object(msgTypeName, jsonText, &msgType);
        if (nullptr == msgObj)
            break;
        if (msgObj->SerializeToString(&binData) == false)
            break;
        isOk = true;
    } while (false);
    return isOk;
}

void zxtools::Common1Req2CommonData(CommonData* dst, const txdata::Common1Req* src)
{
    dst->RspCnt = 0;
    dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
    dst->MsgTypeTxt = QString::fromStdString(::txdata::MsgType_Name(static_cast<txdata::MsgType>(dst->MsgType)));
    dst->PeerID = QString::fromStdString(src->recverid());
    dst->UserID = QString::fromStdString(src->senderid());
    dst->MsgNo = src->msgno();
    dst->SeqNo = src->seqno();
    dst->BatchNo = src->batchno();
    dst->RefNum = src->refnum();
    dst->RefText = QString::fromStdString(src->reftext());
    dst->SenderID = QString::fromStdString(src->senderid());
    dst->RecverID = QString::fromStdString(src->recverid());
    dst->ToRoot = src->toroot();
    dst->IsLog = src->islog();
    dst->IsSafe = false;
    dst->IsPush = src->ispush();
    dst->UpCache = false;
    dst->TxType = static_cast<int32_t>(src->reqtype());
    dst->TxTypeTxt = QString::fromStdString(::txdata::MsgType_Name(src->reqtype()));
    dst->TxData = QByteArray::fromStdString(src->reqdata());
    dst->TxDataTxt = binary2json(src->reqtype(), dst->TxData.toStdString()).trimmed();
    gpt2qdt(dst->TxTime, src->reqtime());
    dst->InsertTime = QDateTime::currentDateTime();
    dst->IsLast = false;
}

void zxtools::CommonData2Common1Req(txdata::Common1Req* dst, const CommonData* src)
{
    dst->set_msgno(src->MsgNo);
    dst->set_seqno(src->SeqNo);
    dst->set_batchno(src->BatchNo);
    dst->set_refnum(src->RefNum);
    dst->set_reftext(src->RefText.toStdString());
    dst->set_senderid(src->SenderID.toStdString());
    dst->set_recverid(src->RecverID.toStdString());
    dst->set_toroot(src->ToRoot);
    dst->set_islog(src->IsLog);
    dst->set_ispush(src->IsPush);
    dst->set_reqtype(static_cast<txdata::MsgType>(src->TxType));
    dst->set_reqdata(src->TxData.toStdString());
    qdt2gpt(*dst->mutable_reqtime(), src->TxTime);
}

void zxtools::Common1Rsp2CommonData(CommonData* dst, const txdata::Common1Rsp* src)
{
    dst->RspCnt = 0;
    dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
    dst->MsgTypeTxt = QString::fromStdString(::txdata::MsgType_Name(static_cast<txdata::MsgType>(dst->MsgType)));
    dst->PeerID = QString::fromStdString(src->senderid());
    dst->UserID = QString::fromStdString(src->recverid());
    dst->MsgNo = src->msgno();
    dst->SeqNo = src->seqno();
    dst->BatchNo = src->batchno();
    dst->RefNum = src->refnum();
    dst->RefText = QString::fromStdString(src->reftext());
    dst->SenderID = QString::fromStdString(src->senderid());
    dst->RecverID = QString::fromStdString(src->recverid());
    dst->ToRoot = src->toroot();
    dst->IsLog = src->islog();
    dst->IsSafe = false;
    dst->IsPush = src->ispush();
    dst->UpCache = false;
    dst->TxType = static_cast<int32_t>(src->rsptype());
    dst->TxTypeTxt = QString::fromStdString(::txdata::MsgType_Name(src->rsptype()));
    dst->TxData = QByteArray::fromStdString(src->rspdata());
    dst->TxDataTxt = binary2json(src->rsptype(), dst->TxData.toStdString()).trimmed();
    gpt2qdt(dst->TxTime, src->rsptime());
    dst->InsertTime = QDateTime::currentDateTime();
    dst->IsLast = src->islast();
}

void zxtools::CommonData2Common1Rsp(txdata::Common1Rsp* dst, const CommonData* src)
{
    dst->set_msgno(src->MsgNo);
    dst->set_seqno(src->SeqNo);
    dst->set_batchno(src->BatchNo);
    dst->set_refnum(src->RefNum);
    dst->set_reftext(src->RefText.toStdString());
    dst->set_senderid(src->SenderID.toStdString());
    dst->set_recverid(src->RecverID.toStdString());
    dst->set_toroot(src->ToRoot);
    dst->set_islog(src->IsLog);
    dst->set_ispush(src->IsPush);
    dst->set_rsptype(static_cast<txdata::MsgType>(src->TxType));
    dst->set_rspdata(src->TxData.toStdString());
    qdt2gpt(*dst->mutable_rsptime(), src->TxTime);
    dst->set_islast(src->IsLast);
}

void zxtools::Common2Req2CommonData(CommonData* dst, const txdata::Common2Req* src)
{
    dst->RspCnt = 0;
    dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
    dst->MsgTypeTxt = QString::fromStdString(::txdata::MsgType_Name(static_cast<txdata::MsgType>(dst->MsgType)));
    dst->PeerID = QString::fromStdString(src->recverid());
    dst->UserID = QString::fromStdString(src->key().userid());
    dst->MsgNo = src->key().msgno();
    dst->SeqNo = src->key().seqno();
    dst->BatchNo = src->batchno();
    dst->RefNum = src->refnum();
    dst->RefText = QString::fromStdString(src->reftext());
    dst->SenderID = QString::fromStdString(src->senderid());
    dst->RecverID = QString::fromStdString(src->recverid());
    dst->ToRoot = src->toroot();
    dst->IsLog = src->islog();
    dst->IsSafe = src->issafe();
    dst->IsPush = src->ispush();
    dst->UpCache = src->upcache();
    dst->TxType = static_cast<int32_t>(src->reqtype());
    dst->TxTypeTxt = QString::fromStdString(::txdata::MsgType_Name(src->reqtype()));
    dst->TxData = QByteArray::fromStdString(src->reqdata());
    dst->TxDataTxt = binary2json(src->reqtype(), dst->TxData.toStdString()).trimmed();
    gpt2qdt(dst->TxTime, src->reqtime());
    dst->InsertTime = QDateTime::currentDateTime();
    dst->IsLast = false;
}

void zxtools::Common2Rsp2CommonData(CommonData* dst, const txdata::Common2Rsp* src)
{
    dst->RspCnt = 0;
    dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
    dst->MsgTypeTxt = QString::fromStdString(::txdata::MsgType_Name(static_cast<txdata::MsgType>(dst->MsgType)));
    dst->PeerID = QString::fromStdString(src->senderid());
    dst->UserID = QString::fromStdString(src->key().userid());
    dst->MsgNo = src->key().msgno();
    dst->SeqNo = src->key().seqno();
    dst->BatchNo = src->batchno();
    dst->RefNum = src->refnum();
    dst->RefText = QString::fromStdString(src->reftext());
    dst->SenderID = QString::fromStdString(src->senderid());
    dst->RecverID = QString::fromStdString(src->recverid());
    dst->ToRoot = src->toroot();
    dst->IsLog = src->islog();
    dst->IsSafe = src->issafe();
    dst->IsPush = src->ispush();
    dst->UpCache = src->upcache();
    dst->TxType = static_cast<int32_t>(src->rsptype());
    dst->TxTypeTxt = QString::fromStdString(::txdata::MsgType_Name(src->rsptype()));
    dst->TxData = QByteArray::fromStdString(src->rspdata());
    dst->TxDataTxt = binary2json(src->rsptype(), dst->TxData.toStdString()).trimmed();
    gpt2qdt(dst->TxTime, src->rsptime());
    dst->InsertTime = QDateTime::currentDateTime();
    dst->IsLast = src->islast();
}

void zxtools::PushWrap2PushWrap(PushWrap* dst, const txdata::PushWrap* src, const std::string& peerid)
{
    dst->PeerID = QString::fromStdString(peerid);
    dst->MsgNo = src->msgno();
    dst->UserID = QString::fromStdString(src->userid());
    gpt2qdt(dst->PshTime, src->msgtime());
    dst->PshType = static_cast<int32_t>(src->msgtype());
    dst->PshTypeTxt = QString::fromStdString(::txdata::MsgType_Name(src->msgtype()));
    dst->PshData = QByteArray::fromStdString(src->msgdata());
    dst->PshDataTxt = binary2json(src->msgtype(), dst->PshData.toStdString()).trimmed();
}

bool zxtools::needTTS(const PushWrap* src, QString& text)
{
    bool isTTS = false;
    if (static_cast<txdata::MsgType>(src->PshType) == txdata::MsgType::ID_PushItem)
    {
        GPMSGPTR theMsg;
        if (m2b::slice2msg(src->PshData.toStdString(), txdata::MsgType::ID_PushItem, theMsg))
        {
            auto pushItem = qSharedPointerDynamicCast<txdata::PushItem>(theMsg);
            for (int i = 0; i < pushItem->modes_size(); i++)
            {
                if (pushItem->modes(i) == "tts")
                {
                    isTTS = true;
                    break;
                }
            }
            if (isTTS)
            {
                text = QString::fromStdString(pushItem->subject());
            }
        }
    }
    return isTTS;
}

QMap<QString, QString> zxtools::fromOddEven(const QStringList& kvs)
{
    QMap<QString, QString> dstMap;
    do
    {
        if (kvs.size() % 2)
        {
            dstMap.clear();
            break;
        }
        for (int i = 0; i < kvs.size() / 2; i++)
        {
            if (dstMap.count(kvs.at(2 * i)))
            {
                dstMap.clear();
                break;
            }
            dstMap.insert(kvs.at(2 * i), kvs.at(2 * i + 1));
        }
    } while (false);
    return dstMap;
}

QString zxtools::map2Common1Req(txdata::Common1Req* dst, const QMap<QString, QString>& src)
{
    QString message;
    do
    {
        ::txdata::MsgType tmpMsgType = ::txdata::MsgType::Zero1;
        std::string tmpMsgData;
        if (!zxtools::json2binary(src.value("ReqType"), src.value("ReqData"), tmpMsgType, tmpMsgData))
        {
            message = "json2binary fail";
            break;
        }
        dst->set_msgno(src.value("MsgNo").toLongLong());
        dst->set_seqno(src.value("SeqNo").toInt());
        dst->set_batchno(src.value("BatchNo").toLong());
        dst->set_refnum(src.value("RefNum").toLong());
        dst->set_reftext(src.value("RefText").toStdString());
        dst->set_senderid(src.value("SenderID").toStdString());
        dst->set_recverid(src.value("RecverID").toStdString());
        dst->set_toroot(src.value("ToRoot").toInt() != 0);
        dst->set_islog(src.value("IsLog").toInt() != 0);
        dst->set_ispush(src.value("IsPush").toInt() != 0);
        dst->set_reqtype(tmpMsgType);
        dst->set_reqdata(tmpMsgData);
        zxtools::qdt2gpt(*dst->mutable_reqtime(), QDateTime::fromString(src.value("ReqTime"), "yyyy-MM-dd hh:mm:ss"));
    } while (false);
    return message;
}

QString zxtools::map2Common2Req(txdata::Common2Req* dst, const QMap<QString, QString>& src)
{
    QString message;
    do
    {
        ::txdata::MsgType tmpMsgType = ::txdata::MsgType::Zero1;
        std::string tmpMsgData;
        if (!zxtools::json2binary(src.value("ReqType"), src.value("ReqData"), tmpMsgType, tmpMsgData))
        {
            message = "json2binary fail";
            break;
        }
        dst->mutable_key()->set_userid(src.value("UserID").toStdString());
        dst->mutable_key()->set_msgno(src.value("MsgNo").toLongLong());
        dst->mutable_key()->set_seqno(src.value("SeqNo").toInt());
        dst->set_batchno(src.value("BatchNo").toLong());
        dst->set_refnum(src.value("RefNum").toLong());
        dst->set_reftext(src.value("RefText").toStdString());
        dst->set_senderid(src.value("SenderID").toStdString());
        dst->set_recverid(src.value("RecverID").toStdString());
        dst->set_toroot(src.value("ToRoot").toInt() != 0);
        dst->set_islog(src.value("IsLog").toInt() != 0);
        dst->set_issafe(src.value("IsSafe").toInt() != 0);
        dst->set_ispush(src.value("IsPush").toInt() != 0);
        dst->set_upcache(src.value("UpCache").toInt() != 0);
        dst->set_reqtype(tmpMsgType);
        dst->set_reqdata(tmpMsgData);
        zxtools::qdt2gpt(*dst->mutable_reqtime(), QDateTime::fromString(src.value("ReqTime"), "yyyy-MM-dd hh:mm:ss"));
    } while (false);
    return message;
}
