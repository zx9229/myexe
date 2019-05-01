#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include "dataexchanger.h"

int main(int argc, char *argv[])
{
#if defined(Q_OS_WIN)
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
#endif

    QGuiApplication app(argc, argv);

    DataExchanger dataExch;

    QQmlApplicationEngine engine;
    {
        //注意: 此处的"dataExch"必须小写字母开头，QML才能访问C++对象的函数与属性.
        engine.rootContext()->setContextProperty("dataExch", &dataExch);
    }
    engine.load(QUrl(QStringLiteral("qrc:/main.qml")));
    if (engine.rootObjects().isEmpty())
        return -1;

    return app.exec();
}
