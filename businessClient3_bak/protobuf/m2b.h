// 本文件以 UTF-8 无 BOM 格式编码.
#ifndef MESSAGE_TO_BYTES_H
#define MESSAGE_TO_BYTES_H
/*
如果编译出错了, 请搜索第一个【errno() const】并修改成类似下面的内容, 仅修改第一个就可以了:
#ifdef errno
#undef errno
::google::protobuf::int32 errno() const;
#endif
*/
#include "txdata.pb.h"
#include <QSharedPointer>

using GPMSGPTR = QSharedPointer<::google::protobuf::Message>;

class m2b
{
public:
    static bool msg2package(const ::google::protobuf::Message& msgIn, QByteArray& pkgOut)
    {
        pkgOut.clear();
        //在定义(.proto)文件的时候,就必须将其对应正确.
        ::txdata::MsgType msgType = static_cast<::txdata::MsgType>(msgIn.GetDescriptor()->index());
        //令(int i = 1)则assert(reinterpret_cast<char*>(&i)[0] == 1)
        pkgOut.append(reinterpret_cast<char*>(&msgType), 2);
        std::string tmpData;
        if (msgIn.SerializeToString(&tmpData) == false)
            return false;
        pkgOut.append(tmpData.data(), tmpData.size());
        return true;
    }

    static QByteArray msg2package(const ::google::protobuf::Message& msgIn)
    {
        QByteArray pkgOut;
        bool retVal = msg2package(msgIn, pkgOut);
        Q_ASSERT(retVal);
        return pkgOut;
    }

    static std::string msg2slice(const ::google::protobuf::Message& src)
    {
        std::string dst;
        bool retVal = src.SerializeToString(&dst);
        Q_ASSERT(retVal);
        return dst;
    }

    static bool slice2msg(const char* data, int size, ::txdata::MsgType msgType, GPMSGPTR& msgOut)
    {
        // 需要在shell下,先创建ff函数,再执行ff函数.
        // ff(){ sed -n '/^enum MsgType/,/}/p' "$1" | sed 's/[ \t]*\?\(ID_\)\([^ \t]\+\).*/case ::txdata::MsgType::\1\2: \n msgOut = QSharedPointer<txdata::\2>(new txdata::\2); \n break;/g' ; }
        // ff  txdata.proto
        switch (msgType)
        {
        case ::txdata::MsgType::ID_MessageAck:
            msgOut = QSharedPointer<txdata::MessageAck>(new txdata::MessageAck);
            break;
        case ::txdata::MsgType::ID_CommonErr:
            msgOut = QSharedPointer<txdata::CommonErr>(new txdata::CommonErr);
            break;
        case ::txdata::MsgType::ID_CommonReq:
            msgOut = QSharedPointer<txdata::CommonReq>(new txdata::CommonReq);
            break;
        case ::txdata::MsgType::ID_CommonRsp:
            msgOut = QSharedPointer<txdata::CommonRsp>(new txdata::CommonRsp);
            break;
        case ::txdata::MsgType::ID_ConnectionInfo:
            msgOut = QSharedPointer<txdata::ConnectionInfo>(new txdata::ConnectionInfo);
            break;
        case ::txdata::MsgType::ID_DisconnectedData:
            msgOut = QSharedPointer<txdata::DisconnectedData>(new txdata::DisconnectedData);
            break;
        case ::txdata::MsgType::ID_ConnectReq:
            msgOut = QSharedPointer<txdata::ConnectReq>(new txdata::ConnectReq);
            break;
        case ::txdata::MsgType::ID_ConnectRsp:
            msgOut = QSharedPointer<txdata::ConnectRsp>(new txdata::ConnectRsp);
            break;
        case ::txdata::MsgType::ID_OnlineNotice:
            msgOut = QSharedPointer<txdata::OnlineNotice>(new txdata::OnlineNotice);
            break;
        case ::txdata::MsgType::ID_SystemReport:
            msgOut = QSharedPointer<txdata::SystemReport>(new txdata::SystemReport);
            break;
        case ::txdata::MsgType::ID_QueryRecordReq:
            msgOut = QSharedPointer<txdata::QueryRecordReq>(new txdata::QueryRecordReq);
            break;
        case ::txdata::MsgType::ID_QueryRecordRsp:
            msgOut = QSharedPointer<txdata::QueryRecordRsp>(new txdata::QueryRecordRsp);
            break;
        case ::txdata::MsgType::ID_ExecCmdReq:
            msgOut = QSharedPointer<txdata::ExecCmdReq>(new txdata::ExecCmdReq);
            break;
        case ::txdata::MsgType::ID_ExecCmdRsp:
            msgOut = QSharedPointer<txdata::ExecCmdRsp>(new txdata::ExecCmdRsp);
            break;
        case ::txdata::MsgType::ID_EchoItem:
            msgOut = QSharedPointer<txdata::EchoItem>(new txdata::EchoItem);
            break;
        case ::txdata::MsgType::ID_ReportDataItem:
            msgOut = QSharedPointer<txdata::ReportDataItem>(new txdata::ReportDataItem);
            break;
        case ::txdata::MsgType::ID_SendMailItem:
            msgOut = QSharedPointer<txdata::SendMailItem>(new txdata::SendMailItem);
            break;
        default:
            break;
        }
        if (msgOut && msgOut->ParseFromArray(data, size) == true)
        {
            return true;
        }
        return false;
    }

    static bool package2msg(const QByteArray &pkgIn, ::txdata::MsgType& typeOut, GPMSGPTR& msgOut)
    {
        msgOut.clear();

        const char* pkgData = pkgIn.constData();
        char* type4 = reinterpret_cast<char*>(&typeOut);
        type4[0] = pkgData[0]; type4[1] = pkgData[1]; type4[2] = 0; type4[3] = 0;

        return slice2msg(pkgData + 2, pkgIn.size() - 2, typeOut, msgOut);
    }
};

#endif//MESSAGE_TO_BYTES_H
