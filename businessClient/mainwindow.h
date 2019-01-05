#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
//////////////////////////////////////////////////////////////////////////
#include "mywebsock.h"

namespace Ui {
    class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();

private slots:
    void slotClickedLink();
    void slotClickedSend();

private:
    Ui::MainWindow *ui;
    MyWebsock m_ws;
};

#endif // MAINWINDOW_H
