#ifndef ZXTOOLS_H
#define ZXTOOLS_H
#include  "sqlstruct.h"
#include "protobuf/m2b.h"

class zxtools
{
public:
    static void qdt2gpt(::google::protobuf::Timestamp& gptDst, const QDateTime& qdtSrc)
    {
        //因为 iVal = 5; 即[dst = src], 估计因为这个原因, 标准库里的一些函数也采用了这个布局, 比如strcpy函数:
        //char * strcpy(char * _Dest,const char * _Source);
        //boost库也大多使用该布局, 比如split:
        //inline SequenceSequenceT& split(SequenceSequenceT& Result, RangeT& Input, PredicateT Pred, token_compress_mode_type eCompress=token_compress_off);
        //所以这里也(尽量)采用这种方式.
        gptDst.set_seconds(qdtSrc.offsetFromUtc());
        gptDst.set_nanos(qdtSrc.time().msec() * 1000 * 1000);
    }
    static void gpt2qdt(QDateTime& qdtDst, const ::google::protobuf::Timestamp& gptSrc)
    {
        qdtDst.fromTime_t(static_cast<uint>(gptSrc.seconds()));
    }
    static void Common1Req2CommonData(CommonData* dst, const txdata::Common1Req* src)
    {
        dst->RspCnt = 0;
        dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
        dst->TargetID = QString::fromStdString(src->recverid());
        dst->UserID = QString::fromStdString(src->senderid());
        dst->MsgNo = src->msgno();
        dst->SeqNo = src->seqno();
        dst->SenderID = QString::fromStdString(src->senderid());
        dst->RecverID = QString::fromStdString(src->recverid());
        dst->TxToRoot = src->txtoroot();
        dst->IsLog = src->islog();
        dst->IsSafe = false;
        dst->IsPush = src->ispush();
        dst->UpCache = false;
        dst->InnerType = static_cast<int32_t>(src->reqtype());
        dst->InnerData = QString::fromStdString(src->reqdata());
        gpt2qdt(dst->InnerTime, src->reqtime());
        dst->InsertTime = QDateTime::currentDateTime();
        dst->IsLast = false;
    }
    static void CommonData2Common1Req(txdata::Common1Req* dst, const CommonData* src)
    {
        dst->set_msgno(src->MsgNo);
        dst->set_seqno(src->SeqNo);
        dst->set_senderid(src->SenderID.toStdString());
        dst->set_recverid(src->RecverID.toStdString());
        dst->set_txtoroot(src->TxToRoot);
        dst->set_islog(src->IsLog);
        dst->set_ispush(src->IsPush);
        dst->set_reqtype(static_cast<txdata::MsgType>(src->InnerType));
        dst->set_reqdata(src->InnerData.toStdString());
        qdt2gpt(*dst->mutable_reqtime(), src->InnerTime);
    }
    static void Common1Rsp2CommonData(CommonData* dst, const txdata::Common1Rsp* src)
    {
        dst->RspCnt = 0;
        dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
        dst->TargetID = QString::fromStdString(src->senderid());
        dst->UserID = QString::fromStdString(src->recverid());
        dst->MsgNo = src->msgno();
        dst->SeqNo = src->seqno();
        dst->SenderID = QString::fromStdString(src->senderid());
        dst->RecverID = QString::fromStdString(src->recverid());
        dst->TxToRoot = src->txtoroot();
        dst->IsLog = src->islog();
        dst->IsSafe = false;
        dst->IsPush = src->ispush();
        dst->UpCache = false;
        dst->InnerType = static_cast<int32_t>(src->rsptype());
        dst->InnerData = QString::fromStdString(src->rspdata());
        gpt2qdt(dst->InnerTime, src->rsptime());
        dst->InsertTime = QDateTime::currentDateTime();
        dst->IsLast = src->islast();
    }
    static void CommonData2Common1Rsp(txdata::Common1Rsp* dst, const CommonData* src)
    {
        dst->set_msgno(src->MsgNo);
        dst->set_seqno(src->SeqNo);
        dst->set_senderid(src->SenderID.toStdString());
        dst->set_recverid(src->RecverID.toStdString());
        dst->set_txtoroot(src->TxToRoot);
        dst->set_islog(src->IsLog);
        dst->set_ispush(src->IsPush);
        dst->set_rsptype(static_cast<txdata::MsgType>(src->InnerType));
        dst->set_rspdata(src->InnerData.toStdString());
        qdt2gpt(*dst->mutable_rsptime(), src->InnerTime);
        dst->set_islast(src->IsLast);
    }
};

#endif // ZXTOOLS_H
