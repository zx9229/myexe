//"CompleteThis"=>"Ctrl+Space" 我改成了 "Alt+/"
//F2 在光标选中对象的声明和定义之间切换.
//F4 头文件和源文件之间切换.
//////////////////////////////////////////////////////////////////////////
#include "mainwindow.h"
#include <QApplication>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    MainWindow w;
    w.show();

    return a.exec();
}
