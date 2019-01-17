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
    static bool msg2pkg(const ::google::protobuf::Message& msgIn, QByteArray& pkgOut)
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

    static QByteArray msg2pkg(const ::google::protobuf::Message& msgIn)
    {
        QByteArray pkgOut;
        bool retVal = msg2pkg(msgIn, pkgOut);
        Q_ASSERT(retVal);
        return pkgOut;
    }

    static std::string msg2bin(const ::google::protobuf::Message& src)
    {
        std::string dst;
        bool retVal = src.SerializeToString(&dst);
        Q_ASSERT(retVal);
        return dst;
    }

    static bool pkg2msg(const QByteArray &pkgIn, ::txdata::MsgType& typeOut, GPMSGPTR& msgOut)
    {
        msgOut.clear();

        const char* pkgData = pkgIn.constData();
        char* type4 = reinterpret_cast<char*>(&typeOut);
        type4[0] = pkgData[0]; type4[1] = pkgData[1]; type4[2] = 0; type4[3] = 0;

        switch (typeOut)
        {
        case ::txdata::MsgType::ID_ConnectedData:
            msgOut = QSharedPointer<txdata::ConnectedData>(new txdata::ConnectedData);
            break;
        case ::txdata::MsgType::ID_DisconnectedData:
            msgOut = QSharedPointer<txdata::DisconnectedData>(new txdata::DisconnectedData);
            break;
        case ::txdata::MsgType::ID_ExecuteCommandReq:
            msgOut = QSharedPointer<txdata::ExecuteCommandReq>(new txdata::ExecuteCommandReq);
            break;
        case ::txdata::MsgType::ID_ExecuteCommandRsp:
            msgOut = QSharedPointer<txdata::ExecuteCommandRsp>(new txdata::ExecuteCommandRsp);
            break;
        case ::txdata::MsgType::ID_CommonNtosReq:
            msgOut = QSharedPointer<txdata::CommonNtosReq>(new txdata::CommonNtosReq);
            break;
        case ::txdata::MsgType::ID_CommonNtosRsp:
            msgOut = QSharedPointer<txdata::CommonNtosRsp>(new txdata::CommonNtosRsp);
            break;
        case ::txdata::MsgType::ID_ParentDataReq:
            msgOut = QSharedPointer<txdata::ParentDataReq>(new txdata::ParentDataReq);
            break;
        case ::txdata::MsgType::ID_ParentDataRsp:
            msgOut = QSharedPointer<txdata::ParentDataRsp>(new txdata::ParentDataRsp);
            break;
        default:
            break;
        }

        if (msgOut && msgOut->ParseFromArray(pkgData + 2, pkgIn.size() - 2) == true)
        {
            return true;
        }

        return false;
    }
};

#endif//MESSAGE_TO_BYTES_H
