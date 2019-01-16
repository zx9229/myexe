//"CompleteThis" => "Ctrl+Space" 我喜欢改成 "Alt+/" 或 "Ctrl+Shift+Space"
//[工具]>[选项]>[环境]>[键盘]
//F2     在光标选中对象的声明和定义之间切换.
//F4     头文件和源文件之间切换.
//Ctrl+I (Auto-indent Selection)自动缩进选中的部分.
//////////////////////////////////////////////////////////////////////////
#include "mainwindow.h"
#include <QApplication>
#include "logindialog.h"
#include "dataexchanger.h"

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);

    DataExchanger m_dataExchanger;

    LoginDialog loginDlg(&m_dataExchanger);
    MainWindow w(&m_dataExchanger);
    if (loginDlg.exec() == QDialog::Accepted)
    {
        w.show();
        return a.exec();
    }
    else
    {
        return 0;
    }
}
