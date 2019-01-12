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
        QTableWidget* curTableWidget = qobject_cast<QTableWidget*>(sender());
        QVariant qVariant = curTableWidget->item(row, 0)->data(Qt::UserRole);
        Q_ASSERT(qVariant.canConvert<QConnInfoEx>());
        QConnInfoEx cie = qVariant.value<QConnInfoEx>();

        ui->lineEdit_UserID->setText(cie.UserID);
        ui->lineEdit_cie_u_ZoneName->setText(cie.UserKey.ZoneName);
        ui->lineEdit_cie_u_NodeName->setText(cie.UserKey.NodeName);
        ui->lineEdit_cie_u_ExecType->setText(cie.UserKey.ExecType);
        ui->lineEdit_cie_u_ExecName->setText(cie.UserKey.ExecName);
        ui->lineEdit_cie_b_ZoneName->setText(cie.BelongKey.ZoneName);
        ui->lineEdit_cie_b_NodeName->setText(cie.BelongKey.NodeName);
        ui->lineEdit_cie_b_ExecType->setText(cie.BelongKey.ExecType);
        ui->lineEdit_cie_b_ExecName->setText(cie.BelongKey.ExecName);
        ui->lineEdit_cie_UserID->setText(cie.UserID);
        ui->lineEdit_cie_BelongID->setText(cie.BelongID);
        ui->lineEdit_cie_Version->setText(cie.Version);
        ui->lineEdit_cie_LinkMode->setText(cie.LinkMode);
        ui->lineEdit_cie_ExePid->setText(QString::number(cie.ExePid));
        ui->lineEdit_cie_ExePath->setText(cie.ExePath);
        ui->lineEdit_cie_Pathway->setText(cie.Pathway.join("✖"));
        //按住Alt不放,在小键盘区输入8251,松开Alt键,可得※符号.
        //Alt并08251(10进制的08251=16进制的203B)※.
        //Alt并09745(10进制的09745=16进制的2611)☑.
        //Alt并10004(10进制的10004=16进制的2714)✔.
        //Alt并10006(10进制的10006=16进制的2716)✖.
    }
    ui->tabWidget->setCurrentIndex(ui->tabWidget->count() - 1);
}
