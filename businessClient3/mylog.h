#ifndef MY_LOG_H
#define MY_LOG_H
#include <QDebug>
#include <QDateTime>

class mylog
{
public:
    static void initialize(const char* filename)
    {
        FILE*& log_fp = get_instance();
        if (nullptr != log_fp)
            return;
        log_fp = fopen(filename, "wt");//write_text.
        Q_ASSERT(nullptr != log_fp);
        {
            atexit(mylog::terminate);
            qInstallMessageHandler(mylog::my_message_handler);
        }
        qDebug() << "[log initialize]";
    }
private:
    static FILE*& get_instance()
    {
        static FILE* static_log_fp = nullptr;
        return static_log_fp;
    }
    static void terminate()
    {
        FILE*& log_fp = get_instance();
        qDebug() << "[log terminate]";
        fflush(log_fp);
        fclose(log_fp);
        log_fp = nullptr;
    }
    static void my_message_handler(QtMsgType msgType, const QMessageLogContext& context, const QString& message)
    {
        FILE*& log_fp = get_instance();
        QByteArray curTimeUtf8 = QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss.zzz").toUtf8();
        QByteArray messageUtf8 = message.toUtf8();

        fprintf(log_fp, "%s [%s] %s (%s:%u, %s)\n",
            curTimeUtf8.constData(), to_string(msgType), messageUtf8.constData(),
            context.file, context.line, context.function);

        static int cnt = 0;
        cnt += 1; cnt = cnt % 1;
        if (0 == cnt)
        {
            fflush(log_fp);
        }

        if (msgType == QtMsgType::QtFatalMsg)
        {
            abort();
        }
    }
    static const char* to_string(QtMsgType msgType)
    {
        switch (msgType) {
        case QtMsgType::QtDebugMsg:
            return "Debug";
        case QtMsgType::QtWarningMsg:
            return "Warn";
        case QtMsgType::QtCriticalMsg:
            return "Critical";
        case QtMsgType::QtFatalMsg:
            return "Fatal";
        case QtMsgType::QtInfoMsg:
            return "Info";
        }
        return "Unknown";
    }
};

#endif//MY_LOG_H
