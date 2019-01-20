#include "dialogreqrsp.h"
#include "ui_dialogreqrsp.h"
#include <QMessageBox>
#include "dataexchanger.h"
#include "dialogdata.h"

DialogReqRsp::DialogReqRsp(DataExchanger* p, QWidget *parent) :
    QDialog(parent),
    ui(new Ui::DialogReqRsp),
    m_dataExch(p)
{
    ui->setupUi(this);
    this->setWindowFlags(this->windowFlags() | Qt::WindowMinimizeButtonHint | Qt::WindowMaximizeButtonHint);
#ifdef Q_OS_ANDROID
    this->showFullScreen();
#endif
    initUI();
    QObject::connect(ui->pushButton_reject, &QPushButton::clicked, this, &DialogReqRsp::reject);
    QObject::connect(ui->pushButton_resend, &QPushButton::clicked, this, &DialogReqRsp::slotClickedResend);
    QObject::connect(ui->pushButton_reload, &QPushButton::clicked, this, &DialogReqRsp::slotClickedReload);
    if (true) {
        ui->tableWidget->setSelectionBehavior(QAbstractItemView::SelectRows); //它俩组合在一起以设置整行选中.
        ui->tableWidget->setSelectionMode(QAbstractItemView::SingleSelection);//它俩组合在一起以设置整行选中.
        ui->tableWidget->setAlternatingRowColors(true);//设置隔一行变一颜色,即:一灰一白.
        ui->tableWidget->setEditTriggers(QAbstractItemView::NoEditTriggers);
        connect(ui->tableWidget, &QTableWidget::cellDoubleClicked, this, &DialogReqRsp::slotCellDoubleClicked);
    }
}

DialogReqRsp::~DialogReqRsp()
{
    delete ui;
}

void DialogReqRsp::initUI()
{
    ui->lineEdit_RefNum->setValidator(new QIntValidator());
}

void DialogReqRsp::setRefNum(int64_t refNum)
{
    refNum = (0 <= refNum) && (refNum <= INT64_MAX) ? refNum : 0;
    ui->lineEdit_RefNum->setText(QString::number(refNum));
    slotClickedReload();
}

void DialogReqRsp::slotClickedReload()
{
    int64_t refNum = ui->lineEdit_RefNum->text().toLongLong();

    QSqlQuery sqlQuery;

    QList<QCommonNtosReq> listReq;
    QCommonNtosReq::select_data(sqlQuery, refNum, listReq);
    QList<QCommonNtosRsp> listRsp;
    QCommonNtosRsp::select_data(sqlQuery, refNum, listRsp);

    const int colCnt = ui->tableWidget->columnCount();
    ui->tableWidget->clearContents();
    int rowIdx = -1;
    for (auto&node : listReq)
    {
        ++rowIdx;
        if (ui->tableWidget->item(rowIdx, 0) == nullptr)
        {
            ui->tableWidget->insertRow(rowIdx);
            for (int i = 0; i < colCnt; ++i) { ui->tableWidget->setItem(rowIdx, i, new QTableWidgetItem()); }
        }
        QVariant qVariant; qVariant.setValue(node);
        ui->tableWidget->item(rowIdx, 0)->setData(Qt::UserRole, qVariant);
        ui->tableWidget->item(rowIdx, 0)->setText("Req");
        ui->tableWidget->item(rowIdx, 1)->setText(node.ReqTime.toString("yyyy-MM-dd HH:mm:ss"));
        ui->tableWidget->item(rowIdx, 2)->setText(DataExchanger::nameByMsgType(static_cast<txdata::MsgType>(node.ReqType)));
        Q_ASSERT(colCnt == 3);
    }
    for (auto&node : listRsp)
    {
        ++rowIdx;
        if (ui->tableWidget->item(rowIdx, 0) == nullptr)
        {
            ui->tableWidget->insertRow(rowIdx);
            for (int i = 0; i < colCnt; ++i) { ui->tableWidget->setItem(rowIdx, i, new QTableWidgetItem()); }
        }
        QVariant qVariant; qVariant.setValue(node);
        ui->tableWidget->item(rowIdx, 0)->setData(Qt::UserRole, qVariant);
        ui->tableWidget->item(rowIdx, 0)->setText("Rsp");
        ui->tableWidget->item(rowIdx, 1)->setText(node.InsertTime.toString("yyyy-MM-dd HH:mm:ss"));
        ui->tableWidget->item(rowIdx, 2)->setText(DataExchanger::nameByMsgType(static_cast<txdata::MsgType>(node.RspType)));
        Q_ASSERT(colCnt == 3);
    }
    ui->tableWidget->setRowCount(rowIdx + 1);
}

void DialogReqRsp::slotClickedResend()
{
    int curRowIdx = ui->tableWidget->currentRow();
    if (curRowIdx < 0)
    {
        QMessageBox::information(this, tr("resend"), tr("请先选中一行数据"));
        return;
    }
    QVariant qVariant = ui->tableWidget->item(curRowIdx, 0)->data(Qt::UserRole);
    if (qVariant.canConvert<QCommonNtosReq>() == false)
    {
        QMessageBox::information(this, tr("resend"), tr("当前行不是(QCommonNtosReq)数据, 无法发送"));
        return;
    }
    QCommonNtosReq node = qVariant.value<QCommonNtosReq>();
    m_dataExch->sendCommonNtosReq4resend(node);
}

void DialogReqRsp::slotCellDoubleClicked(int row, int column)
{
    QTableWidget* curTableWidget = qobject_cast<QTableWidget*>(sender());
    QVariant qVariant = curTableWidget->item(row, 0)->data(Qt::UserRole);
    if (curTableWidget->item(row, 0)->text() == "Req")
    {
        Q_ASSERT(qVariant.canConvert<QCommonNtosReq>());
        QCommonNtosReq node = qVariant.value<QCommonNtosReq>();
        DialogData dlgData;
        dlgData.setData(node);
        dlgData.exec();
    }
    else if (curTableWidget->item(row, 0)->text() == "Rsp")
    {
        Q_ASSERT(qVariant.canConvert<QCommonNtosRsp>());
        QCommonNtosRsp node = qVariant.value<QCommonNtosRsp>();
        DialogData dlgData;
        dlgData.setData(node);
        dlgData.exec();
    }
    else
    {
        Q_ASSERT(false);
    }
}
