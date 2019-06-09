//"CompleteThis" => "Ctrl+Space" 我喜欢改成 "Alt+/" 或 "Ctrl+Shift+Space"
//[工具]>[选项]>[环境]>[键盘]
//F2     在光标选中对象的声明和定义之间切换.
//F4     头文件和源文件之间切换.
//Ctrl+I (Auto-indent Selection)自动缩进选中的部分.
//////////////////////////////////////////////////////////////////////////
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include "mylog.h"
#include "mysqltablemodel.h"
#include "datawrapper.h"
#include "temputils.h"

int main(int argc, char *argv[])
{
    bool useRO = true;
#if defined(Q_OS_WIN)
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    //useRO=false;
#endif

    QGuiApplication app(argc, argv);

    bool isService = ((2 == argc) && (qstrcmp(argv[1], "-service") == 0));
    isService = useRO && isService;

    {
        QByteArray logName;
        logName = isService ? TempUtils::calcServiceLogName().toUtf8() : TempUtils::calcAppLogName().toUtf8();
        Q_ASSERT(logName.isEmpty() == false);
        mylog::initialize(logName.constData());
    }

    MySqlTableModel::doQmlRegisterType();

    DataWrapper dataWrap(useRO, isService);//DataExchanger dataExch;

    QSharedPointer<QQmlApplicationEngine> engine;
    if (!useRO || !isService)
    {
        engine = QSharedPointer<QQmlApplicationEngine>(new QQmlApplicationEngine);
    }
    if (engine)
    {
        {
            //注意: 此处的"dataExch"必须小写字母开头，QML才能访问C++对象的函数与属性.
            engine->rootContext()->setContextProperty("dataExch", &dataWrap);
        }
        engine->load(QUrl(QStringLiteral("qrc:/main.qml")));
        if (engine->rootObjects().isEmpty())
            return -1;
    }
    return app.exec();
}
