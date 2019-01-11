#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QMessageBox>
#include "dataexchanger.h"

MainWindow::MainWindow(DataExchanger* p, QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow),
    m_dataExch(p)
{
    ui->setupUi(this);
    connect(ui->pushButton_ParentDataReq, &QPushButton::clicked, this, &MainWindow::slotClickedParentDataReq);
    connect(m_dataExch, &DataExchanger::sigParentData, this, &MainWindow::slotParentData);
    ui->tableWidget->setSelectionBehavior(QAbstractItemView::SelectRows); //它俩组合在一起以设置整行选中.
    ui->tableWidget->setSelectionMode(QAbstractItemView::SingleSelection);//它俩组合在一起以设置整行选中.
    ui->tableWidget->setAlternatingRowColors(true);//设置隔一行变一颜色,即:一灰一白.
    ui->tableWidget->setEditTriggers(QAbstractItemView::NoEditTriggers);
    connect(ui->tableWidget, &QTableWidget::cellDoubleClicked, this, &MainWindow::slotCellDoubleClicked);
    connect(ui->pushButton_show, &QPushButton::clicked, this, &MainWindow::slotClickedShow);
    connect(ui->pushButton_send, &QPushButton::clicked, this, &MainWindow::slotClickedSend);
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::slotClickedParentDataReq()
{
    m_dataExch->slotParentDataReq();
}

void MainWindow::slotParentData(const QMap<QString, QConnInfoEx>& data)
{
    const int colCnt = ui->tableWidget->columnCount();
    ui->tableWidget->clearContents();
    int rowIdx = -1;
    for (auto&p : data)
    {
        ++rowIdx;
        if (ui->tableWidget->item(rowIdx, 0) == nullptr)
        {
            ui->tableWidget->insertRow(rowIdx);
            for (int i = 0; i < colCnt; ++i) { ui->tableWidget->setItem(rowIdx, i, new QTableWidgetItem()); }
        }
        ui->tableWidget->item(rowIdx, 0)->setText(p.UserID);
        ui->tableWidget->item(rowIdx, 1)->setText(p.UserKey.ExecType);
        ui->tableWidget->item(rowIdx, 2)->setText(p.BelongID);
        ui->tableWidget->item(rowIdx, 3)->setText(p.Version);
        Q_ASSERT(colCnt == 4);
        QVariant qVariant; qVariant.setValue(p);
        ui->tableWidget->item(rowIdx, 0)->setData(Qt::UserRole, qVariant);
    }

    ui->tableWidget->setRowCount(rowIdx + 1);
}

void MainWindow::slotClickedShow()
{
    ui->widget_ConnInfoEx->setVisible(!ui->widget_ConnInfoEx->isVisible());
}

void MainWindow::slotClickedSend()
{
    QMessageBox::information(this, "SEND", "Not Implemented");
}

void MainWindow::slotCellDoubleClicked(int row, int column)
{
    if (true) {
        QTableWidget* curObj = qobject_cast<QTableWidget*>(sender());
        QVariant qVariant = curObj->item(row, 0)->data(Qt::UserRole);
        Q_ASSERT(qVariant.canConvert<QConnInfoEx>());
        QConnInfoEx qci = qVariant.value<QConnInfoEx>();

        ui->lineEdit_UserID->setText(qci.UserID);
        ui->lineEdit_cie_UserID->setText(qci.UserID);
        ui->lineEdit_cie_BelongID->setText(qci.BelongID);
        ui->lineEdit_cie_Pathway->setText(qci.Pathway.join("=>"));
    }
    ui->tabWidget->setCurrentIndex(ui->tabWidget->count() - 1);
}
