#ifndef SAFESYNCHCACHE_H
#define SAFESYNCHCACHE_H
#include <QString>
#include "m2b.h"

class UniSym
{
public:
    QString UserID;
    int64_t MsgNo;
    int32_t SeqNo;
public:
    UniSym():MsgNo(0),SeqNo(0){}
    void fromUniKey(const txdata::UniKey& uniKey)
    {
        this->UserID=QString::fromStdString( uniKey.userid());
        this->MsgNo=uniKey.msgno();
        this->SeqNo=uniKey.seqno();
    }
    bool operator<(const UniSym& other)const
    {
        if(this->UserID<other.UserID)
        {
            return true;
        }
        else if(this->UserID==other.UserID)
        {
            if(this->MsgNo<other.MsgNo)
            {
                return true;
            }
            else if(this->MsgNo==other.MsgNo)
            {
                if(this->SeqNo<other.SeqNo)
                {
                    return true;
                }
                else
                {
                    return false;
                }
            }
            else
            {
                return false;
            }
        }else
        {
            return false;
        }
    }
};

class Node4Sync
{
public:
    UniSym      Key;
    bool        TxToRoot;
    std::string RecverID;
    GPMSGPTR    Data;
public:
    Node4Sync():TxToRoot(false){}
};
using Node4SyncPtr = std::shared_ptr<Node4Sync>;

class SafeSynchCache
{
private:
    std::mutex                    m_mutex;
    std::map<UniSym,Node4SyncPtr> m_MAP;
public:
    bool insertData(const txdata::UniKey& uniKey,bool toR,const std::string& rID,GPMSGPTR pm)
    {
        Node4SyncPtr node = Node4SyncPtr(new Node4Sync);
        node->Key.fromUniKey(uniKey);
        node->TxToRoot=toR;
        node->RecverID=rID;
        node->Data=pm;
        std::lock_guard<std::mutex>lg(m_mutex);
        auto it = m_MAP.find(node->Key);
        bool doInsert = (m_MAP.end()==it);
        if(doInsert)
        {
            m_MAP[node->Key]=node;
        }
        return doInsert;
    }
    Node4SyncPtr deleteData(const txdata::UniKey& uniKey)
    {
        UniSym uniSym;
        uniSym.fromUniKey(uniKey);
        std::lock_guard<std::mutex>lg(m_mutex);
        auto it = m_MAP.find(uniSym);
        return (m_MAP.end()==it)?nullptr:it->second;
    }
};

#endif // SAFESYNCHCACHE_H
