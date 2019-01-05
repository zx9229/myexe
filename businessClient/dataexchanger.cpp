#include "dataexchanger.h"

DataExchanger::DataExchanger()
{

}

DataExchanger::~DataExchanger()
{

}

QString DataExchanger::Login(const QString &url, const QString &username, const QString &password)
{
    QString message;

    m_ws.start(url);

    return message;
}

MyWebsock& DataExchanger::ws()
{
    return m_ws;
}
