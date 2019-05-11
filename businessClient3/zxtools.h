#ifndef ZXTOOLS_H
#define ZXTOOLS_H
#include  "sqlstruct.h"
#include "protobuf/m2b.h"

class zxtools
{
public:
    static void qdt2gpt(::google::protobuf::Timestamp& gptDst, const QDateTime& qdtSrc)
    {
        if (qdtSrc.isValid())
        {
            gptDst.set_seconds(qdtSrc.offsetFromUtc());
            gptDst.set_nanos(qdtSrc.time().msec() * 1000 * 1000);
        }
    }
    static void gpt2qdt(QDateTime& qdtDst,const ::google::protobuf::Timestamp& gptSrc)
    {
        qdtDst.fromTime_t(static_cast<uint>(gptSrc.seconds()));
    }
    static void Common1Req2CommonData(CommonData* dst, const txdata::Common1Req* src)
    {
        dst->RspCnt=0;
        dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
        dst->TargetID = QString::fromStdString(src->recverid());
        dst->UserID = QString::fromStdString(src->senderid());
        dst->MsgNo= src->msgno();
        dst->SeqNo = src->seqno();
        dst->SenderID=QString::fromStdString(src->senderid());
        dst->RecverID=QString::fromStdString(src->recverid());
        dst->TxToRoot=src->txtoroot();
        dst->IsLog=src->islog();
        dst->IsSafe=false;
        dst->IsPush=src->ispush();
        dst->UpCache=false;
        dst->InnerType=static_cast<int32_t>(src->reqtype()) ;
        dst->InnerData=QString::fromStdString(src->reqdata());
        gpt2qdt(dst->InnerTime,src->reqtime());
        dst->InsertTime=QDateTime::currentDateTime();
        dst->IsLast=false;
    }
    static void Common1Rsp2CommonData(CommonData* dst, const txdata::Common1Rsp* src)
    {
        dst->RspCnt=0;
        dst->MsgType = static_cast<int32_t>(m2b::CalcMsgType(*src));
        dst->TargetID=QString::fromStdString(src->senderid());
        dst->UserID=QString::fromStdString(src->recverid());
        dst->MsgNo=src->msgno();
        dst->SeqNo=src->seqno();
        dst->SenderID=QString::fromStdString( src->senderid());
        dst->RecverID=QString::fromStdString( src->recverid());
        dst->TxToRoot=src->txtoroot();
        dst->IsLog=src->islog();
        dst->IsSafe=false;
        dst->IsPush=src->ispush();
        dst->UpCache=false;
        dst->InnerType=static_cast<int32_t>(src->rsptype());
        dst->InnerData=QString::fromStdString(src->rspdata());
        gpt2qdt(dst->InnerTime,src->rsptime());
        dst->InsertTime=QDateTime::currentDateTime();
        dst->IsLast=src->islast();
    }
};

#endif // ZXTOOLS_H
