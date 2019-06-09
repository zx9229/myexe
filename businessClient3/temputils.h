//"CompleteThis"=>"Ctrl+Space" 我改成了 "Alt+/"
#ifndef TEMPUTILS_H
#define TEMPUTILS_H
#include <QString>

class TempUtils
{
public:
    static QString calcAppDataDir(bool doMkpath = false);
    static QString calcAppLogDir(bool doMkpath = false);
    static QString calcAppLogName(bool* mkpathOk = nullptr);
    static QString calcServiceLogDir(bool doMkpath = false);
    static QString calcServiceLogName(bool* mkpathOk = nullptr);
    static QString calcEnvironmentVariable(const char* name, bool* isOk = nullptr);
};

#endif // TEMPUTILS_H
