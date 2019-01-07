#include "mainwindow.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(DataExchanger* p, QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow),
    m_dataExch(p)
{
    ui->setupUi(this);
    connect(ui->pushButton_send, &QPushButton::clicked, this, &MainWindow::slotClickedSend);
}

MainWindow::~MainWindow()
{
    delete ui;
}

#include "m2b.h"
void MainWindow::slotClickedSend()
{
    txdata::ConnectedData tmpData = {};
    tmpData.mutable_info()->mutable_userkey()->set_exectype(txdata::ProgramType::CLIENT);

    QByteArray data;
    {
        m2b::msg2slice(txdata::ID_ConnectedData, tmpData, data);
    }
    m_dataExch->ws().sendBinaryMessage(data);
}
