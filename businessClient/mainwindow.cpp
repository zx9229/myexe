#include "mainwindow.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(DataExchanger* p, QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow),
    m_dataExch(p)
{
    ui->setupUi(this);
    connect(ui->pushButton_send, &QPushButton::clicked, this, &MainWindow::slotClickedSend);
    connect(ui->pushButton_reqServerCache, &QPushButton::clicked, this, &MainWindow::slotClickedReqServerCache);
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::slotClickedSend()
{
}

void MainWindow::slotClickedReqServerCache()
{
    m_dataExch->slotReqServerCache();
}
