#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "dataexchanger.h"

namespace Ui {
    class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(DataExchanger* p, QWidget *parent = 0);
    ~MainWindow();

private slots:
    void slotClickedSend();
    void slotClickedReqServerCache();

private:
    Ui::MainWindow* ui;
    DataExchanger*  m_dataExch;
};

#endif // MAINWINDOW_H
