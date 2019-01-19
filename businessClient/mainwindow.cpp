#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QMessageBox>
#include "dataexchanger.h"
#include "dialogdata.h"
#include "dialogreqrsp.h"

MainWindow::MainWindow(DataExchanger* p, QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow),
    m_dataExch(p)
{
    ui->setupUi(this);
    initUI();
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::initUI()
{
    if (true) {
        ui->tableWidget->setSelectionBehavior(QAbstractItemView::SelectRows); //它俩组合在一起以设置整行选中.
        ui->tableWidget->setSelectionMode(QAbstractItemView::SingleSelection);//它俩组合在一起以设置整行选中.
        ui->tableWidget->setAlternatingRowColors(true);//设置隔一行变一颜色,即:一灰一白.
        ui->tableWidget->setEditTriggers(QAbstractItemView::NoEditTriggers);
        connect(ui->tableWidget, &QTableWidget::cellDoubleClicked, this, &MainWindow::slotCellDoubleClicked);
    }
    connect(m_dataExch, &DataExchanger::sigParentData, this, &MainWindow::slotParentData);
    connect(m_dataExch, &DataExchanger::sigCommonNtosRsp, this, &MainWindow::slotCommonNtosRsp);
    connect(ui->pushButton_ParentDataReq, &QPushButton::clicked, this, &MainWindow::slotClickedParentDataReq);
    connect(ui->pushButton_show, &QPushButton::clicked, this, &MainWindow::slotClickedShow);
    connect(ui->pushButton_send, &QPushButton::clicked, this, &MainWindow::slotClickedSend);
    connect(ui->pushButton_RefNumShow, &QPushButton::clicked, this, &MainWindow::slotClickedRefNumShow);
    if (true) {
        connect(ui->lineEdit_UserID, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_u_ZoneName, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_u_NodeName, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_u_ExecType, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_u_ExecName, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_b_ZoneName, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_b_NodeName, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_b_ExecType, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_b_ExecName, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_UserID, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_BelongID, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_Version, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_LinkMode, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_ExePid, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_ExePath, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
        connect(ui->lineEdit_cie_Pathway, &QLineEdit::cursorPositionChanged, this, &MainWindow::slotCursorPositionChanged);
    }
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

void MainWindow::slotCommonNtosRsp(qint64 RefNum)
{
    if ((0 < RefNum && ui->lineEdit_RefNum->text().toLongLong() == RefNum) == false)
        return;
    QSqlQuery sqlQuery;
    int32_t cntTime = 0;
    QString minTime, maxTime;
    if (QCommonNtosRsp::select_stat_info(sqlQuery, RefNum, cntTime, minTime, maxTime) == false)
        return;
    ui->lineEdit_RspTimeCnt->setText(QString::number(cntTime));
    ui->lineEdit_RspTimeMin->setText(minTime);
    ui->lineEdit_RspTimeMax->setText(maxTime);
}

void MainWindow::slotClickedShow()
{
    ui->widget_ConnInfoEx->setVisible(!ui->widget_ConnInfoEx->isVisible());
}

void MainWindow::slotClickedSend()
{
    DialogData dlgData;
    if (dlgData.exec() != QDialog::Accepted)
        return;
    QCommonNtosReq reqData;
    dlgData.getData(reqData.ReqData, reqData.ReqType);
    if (m_dataExch->sendCommonNtosReq(reqData, dlgData.needResp(), dlgData.needSave()) == false)
    {
        QMessageBox::information(this, tr("发送请求"), tr("发送请求失败"));
        return;
    }
    ui->lineEdit_RefNum->setText(QString::number(reqData.RefNum));
    ui->lineEdit_ReqTime->setText(reqData.ReqTime.toString("yyyy-MM-dd HH:mm:ss"));
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

void MainWindow::slotCursorPositionChanged(int iOld, int iNew)
{
    //在手机上数据显示不全,想要一种方式可以看到全部数据,准备双击然后弹窗然后发现不太好写,
    //遂绑定此信号以代替双击信息,同时尽量避免误操作时的弹窗,遂有此函数.
    QLineEdit* curLineEdit = qobject_cast<QLineEdit*>(sender());
    if (curLineEdit == nullptr)
        return;
    if (iOld == 0 || iOld == iNew || iNew == 0)
        return;
    int curTextSize = curLineEdit->text().size();
    if (iOld == curTextSize || iNew == curTextSize)
        return;
    QMessageBox::information(this, "QLineEdit", curLineEdit->text());
}

void MainWindow::slotClickedRefNumShow()
{
    DialogReqRsp dlgReqRsp(m_dataExch, this);
    int64_t refNum = ui->lineEdit_RefNum->text().toLongLong();
    dlgReqRsp.setRefNum(refNum);
    dlgReqRsp.exec();
}
