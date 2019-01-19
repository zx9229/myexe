#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "temputils.h"

namespace Ui {
    class MainWindow;
}

class DataExchanger;
class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(DataExchanger* p, QWidget *parent = 0);
    ~MainWindow();

private:
    void initUI();

private slots:
    void slotClickedParentDataReq();
    void slotParentData(const QMap<QString, QConnInfoEx>& data);
    void slotCommonNtosRsp(qint64 RefNum);
    void slotClickedShow();
    void slotClickedSend();
    void slotCellDoubleClicked(int row, int column);
    void slotCursorPositionChanged(int iOld, int iNew);
    void slotClickedRefNumShow();

private:
    Ui::MainWindow* ui;
    DataExchanger*  m_dataExch;
};

#endif // MAINWINDOW_H
