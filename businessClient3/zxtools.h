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
        dst->MsgTypeTxt=QString::fromStdString(::txdata::MsgType_Name(static_cast<txdata::MsgType>(dst->MsgType)));
        dst->PeerID = QString::fromStdString(src->recverid());
        dst->UserID = QString::fromStdString(src->senderid());
        dst->MsgNo = src->msgno();
        dst->SeqNo = src->seqno();
        dst->SenderID = QString::fromStdString(src->senderid());
        dst->RecverID = QString::fromStdString(src->recverid());
        dst->ToRoot = src->toroot();
        dst->IsLog = src->islog();
        dst->IsSafe = false;
        dst->IsPush = src->ispush();
        dst->UpCache = false;
        dst->TxType = static_cast<int32_t>(src->reqtype());
        dst->TxTypeTxt=QString::fromStdString(::txdata::MsgType_Name(src->reqtype()));
        dst->TxData = QString::fromStdString(src->reqdata());
        dst->TxDataTxt.clear();//TODO:
        gpt2qdt(dst->TxTime, src->reqtime());
        dst->InsertTime = QDateTime::currentDateTime();
        dst->IsLast = false;
    }
    static void CommonData2Common1Req(txdata::Common1Req* dst, const CommonData* src)
    {
        dst->set_msgno(src->MsgNo);
        dst->set_seqno(src->SeqNo);
        dst->set_senderid(src->SenderID.toStdString());
        dst->set_recverid(src->RecverID.toStdString());
        dst->set_toroot(src->ToRoot);
        dst->set_islog(src->IsLog);
        dst->set_ispush(src->IsPush);
        dst->set_reqtype(static_cast<txdata::MsgType>(src->TxType));
        dst->set_reqdata(src->TxData.toStdString());
        qdt2gpt(*dst->mutable_reqtime(), src->TxTime);
    }
    static void Common1Rsp2CommonData(CommonData* dst, const txdata::Common1Rsp* src)
    {
        dst->RspCnt = 0;
        dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
        dst->PeerID = QString::fromStdString(src->senderid());
        dst->UserID = QString::fromStdString(src->recverid());
        dst->MsgNo = src->msgno();
        dst->SeqNo = src->seqno();
        dst->SenderID = QString::fromStdString(src->senderid());
        dst->RecverID = QString::fromStdString(src->recverid());
        dst->ToRoot = src->toroot();
        dst->IsLog = src->islog();
        dst->IsSafe = false;
        dst->IsPush = src->ispush();
        dst->UpCache = false;
        dst->TxType = static_cast<int32_t>(src->rsptype());
        dst->TxTypeTxt= QString::fromStdString(::txdata::MsgType_Name(src->rsptype()));
        dst->TxData = QString::fromStdString(src->rspdata());
        dst->TxDataTxt.clear();
        gpt2qdt(dst->TxTime, src->rsptime());
        dst->InsertTime = QDateTime::currentDateTime();
        dst->IsLast = src->islast();
    }
    static void CommonData2Common1Rsp(txdata::Common1Rsp* dst, const CommonData* src)
    {
        dst->set_msgno(src->MsgNo);
        dst->set_seqno(src->SeqNo);
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
};

#endif // ZXTOOLS_H
