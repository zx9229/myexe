#ifndef DATAEXCHANGER_H
#define DATAEXCHANGER_H

#include "mywebsock.h"

class DataExchanger
{
public:
    DataExchanger();
    ~DataExchanger();

public:
    QString Login(const QString& url, const QString& username, const QString& password);
    MyWebsock& ws();

private:
    MyWebsock m_ws;
};

#endif // DATAEXCHANGER_H
