// 本文件以 UTF-8 无 BOM 格式编码.
#ifndef MESSAGE_TO_BYTES_H
#define MESSAGE_TO_BYTES_H

#include "txdata.pb.h"
#include <QSharedPointer>

using GPMSGPTR = ::QSharedPointer<::google::protobuf::Message>;

class m2b
{
public:
    static bool msg2pkg(::txdata::MsgType msgType, const ::google::protobuf::Message& msgIn, ::QByteArray& pkgOut)
    {
        //令(int i = 1)则assert(reinterpret_cast<char*>(&i)[0] == 1)
        pkgOut.clear();
        pkgOut.append(reinterpret_cast<char*>(&msgType), 2);
        std::string tmpData;
        if (msgIn.SerializeToString(&tmpData) == false)
            return false;
        pkgOut.append(tmpData.data(), tmpData.size());
        return true;
    }

    static QByteArray msg2pkg(::txdata::MsgType msgType, const ::google::protobuf::Message& msgIn)
    {
        QByteArray pkgOut;
        bool retVal = msg2pkg(msgType, msgIn, pkgOut);
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

    static bool pkg2msg(const ::QByteArray &pkgIn, ::txdata::MsgType& typeOut, GPMSGPTR& msgOut)
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
