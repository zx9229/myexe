#include "dialogdata.h"
#include "ui_dialogdata.h"
#include <QMessageBox>
#include "google/protobuf/util/json_util.h"
#include "txdata.pb.h"
#include "dataexchanger.h"

DialogData::DialogData(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::DialogData)
{
    ui->setupUi(this);
    this->setWindowFlags(this->windowFlags() | Qt::WindowMinimizeButtonHint | Qt::WindowMaximizeButtonHint);
#ifdef Q_OS_ANDROID
    this->showFullScreen();
#endif
    initUI();
    QObject::connect(ui->pushButton_reject, &QPushButton::clicked, this, &DialogData::reject);
    QObject::connect(ui->pushButton_accept, &QPushButton::clicked, this, &DialogData::slotAccept);
    QObject::connect(ui->comboBox_msgType, static_cast<void(QComboBox::*)(const QString &)>(&QComboBox::currentIndexChanged), this, &DialogData::slotCurrentIndexChanged);
    QObject::connect(ui->pushButton_inner, &QPushButton::clicked, this, &DialogData::slotClickedShowInnerOuter);
    QObject::connect(ui->pushButton_outer, &QPushButton::clicked, this, &DialogData::slotClickedShowInnerOuter);
}

DialogData::~DialogData()
{
    delete ui;
}

void DialogData::initUI()
{
    QStringList nameList = { "" };
    for (int i = 1; i <= static_cast<int>(txdata::MsgType_MAX); ++i)
    {
        bool isOk = false;
        nameList << DataExchanger::nameByMsgType(static_cast<txdata::MsgType>(i), 0, &isOk);
        Q_ASSERT(isOk);
    }
    for (auto it = nameList.begin(); it != nameList.end(); ++it)
    {
        ui->comboBox_msgType->addItem(*it);
    }
    switchMode(true);
}

void DialogData::switchMode(bool isInputNotOutput)
{
    ui->pushButton_accept->setEnabled(isInputNotOutput);
    ui->comboBox_msgType->setEnabled(isInputNotOutput);
    ui->checkBox_needResp->setEnabled(isInputNotOutput);
    ui->checkBox_needSave->setEnabled(isInputNotOutput);
    ui->lineEdit_MsgType->setEnabled(!isInputNotOutput);
    ui->pushButton_outer->setEnabled(!isInputNotOutput);
    ui->pushButton_inner->setEnabled(!isInputNotOutput);
}

void DialogData::slotCurrentIndexChanged(const QString &text)
{
    QString jsonStr;
    QSharedPointer<google::protobuf::Message> curObj;
    if (DataExchanger::calcObjByName(text, curObj))
    {
        jsonStr = DataExchanger::jsonByMsgObje(*curObj);
    }
    ui->plainTextEdit->setPlainText(jsonStr);
}

void DialogData::slotAccept()
{
    m_curType = 0;
    m_curData.clear();
    txdata::MsgType currType = txdata::MsgType::Zero1;

    QString typeStr = ui->comboBox_msgType->currentText();
    QString jsonStr = ui->plainTextEdit->toPlainText();

    QString message = DataExchanger::jsonToObjAndS(typeStr, jsonStr, currType, m_curData);
    m_curType = static_cast<int32_t>(currType);
    if (!message.isEmpty())
    {
        QMessageBox::information(this, tr("json->obj->bin"), message);
        return;
    }

    this->accept();
}

void DialogData::slotClickedShowInnerOuter()
{
    QPushButton* curPushButton = qobject_cast<QPushButton*>(sender());
    QString nameStr = curPushButton->property("MsgType").toString();
    QString jsonStr = curPushButton->property("MsgJson").toString();
    ui->lineEdit_MsgType->setText(nameStr);
    ui->plainTextEdit->setPlainText(jsonStr);
}

bool DialogData::needResp()
{
    return (ui->checkBox_needResp->checkState() == Qt::Checked);
}

bool DialogData::needSave()
{
    return (ui->checkBox_needSave->checkState() == Qt::Checked);
}

void DialogData::getData(QByteArray &dataOut, int32_t& typeOut)
{
    dataOut = m_curData;
    typeOut = m_curType;
}

void DialogData::setData(const QCommonNtosReq &data)
{
    switchMode(false);

    txdata::MsgType innerType = static_cast<txdata::MsgType>(data.ReqType);
    ui->pushButton_inner->setProperty("MsgType", DataExchanger::nameByMsgType(innerType));
    ui->pushButton_inner->setProperty("MsgJson", DataExchanger::jsonByMsgType(innerType, data.ReqData));

    //txdata::CommonNtosReq dataTx;
    //DataExchanger::CommonNtosReqQ2TX(data, dataTx);
    //ui->pushButton_outer->setProperty("MsgType", QString::fromStdString(dataTx.GetDescriptor()->name()));
    //ui->pushButton_outer->setProperty("MsgJson", DataExchanger::jsonByMsgObje(dataTx));
}

void DialogData::setData(const QCommonNtosRsp &data)
{
    switchMode(false);

    txdata::MsgType innerType = static_cast<txdata::MsgType>(data.RspType);
    ui->pushButton_inner->setProperty("MsgType", DataExchanger::nameByMsgType(innerType));
    if (true) {
        QString jsonStr = DataExchanger::jsonByMsgType(innerType, data.RspData);
        ui->pushButton_inner->setProperty("MsgJson", jsonStr.isEmpty() ? QString(data.RspData) : jsonStr);
    }
    //txdata::CommonNtosRsp dataTx;
    //DataExchanger::CommonNtosRspQ2TX(data, dataTx);
    //ui->pushButton_outer->setProperty("MsgType", QString::fromStdString(dataTx.GetDescriptor()->name()));
    //ui->pushButton_outer->setProperty("MsgJson", DataExchanger::jsonByMsgObje(dataTx));
}
