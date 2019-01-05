#include "mainwindow.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    connect(ui->pushButton_link, &QPushButton::clicked, this, &MainWindow::slotClickedLink);
    connect(ui->pushButton_send, &QPushButton::clicked, this, &MainWindow::slotClickedSend);
    ui->lineEdit_url->setText("ws://localhost:10083/websocket");
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::slotClickedLink()
{
    m_ws.stop();
    m_ws.start(ui->lineEdit_url->text());
}

#include "m2b.h"
void MainWindow::slotClickedSend()
{
    txdata::ConnectedData tmpData = {};
    tmpData.mutable_info()->set_exetype(txdata::ConnectionInfo_AppType::ConnectionInfo_AppType_CLIENT);

    QByteArray data;
    {
        m2b::msg2slice(txdata::ID_ConnectedData, tmpData, data);
    }
    m_ws.sendBinaryMessage(data);
}
