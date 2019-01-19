#include "dialogdata.h"
#include "ui_dialogdata.h"
#include <QMessageBox>
#include "google/protobuf/util/json_util.h"
#include "txdata.pb.h"
#include "dataexchanger.h"

//如果获取了(google::protobuf::Message)那么需要自行析构.
bool CalcObjectByName(const std::string name, const google::protobuf::Descriptor** desc, google::protobuf::Message** mesg)
{
    if (nullptr != desc) { *desc = nullptr; }
    if (nullptr != mesg) { *mesg = nullptr; }
    // https://blog.csdn.net/riopho/article/details/80372510
    const google::protobuf::Descriptor* curDesc = google::protobuf::DescriptorPool::generated_pool()->FindMessageTypeByName(name);
    if (nullptr == curDesc)
        return false;
    //desc->index();
    google::protobuf::Message* curMesg = google::protobuf::MessageFactory::generated_factory()->GetPrototype(curDesc)->New();
    if (nullptr == curMesg)
        return false;
    if (nullptr != desc)
        *desc = curDesc;
    if (nullptr != mesg)
        *mesg = curMesg;
    else
        delete mesg;
    return ((nullptr != curDesc) && (nullptr != curMesg));
}

class MessageGuard
{
public:
    MessageGuard(google::protobuf::Message* message) :m_message(message) {}
    ~MessageGuard() { if (m_message) { delete m_message; } }
private:
    google::protobuf::Message* m_message;
};

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
    ui->comboBox_msgType->setEnabled(isInputNotOutput);
    ui->checkBox_needResp->setEnabled(isInputNotOutput);
    ui->checkBox_needSave->setEnabled(isInputNotOutput);
    ui->lineEdit_MsgType->setEnabled(!isInputNotOutput);
    ui->pushButton_outer->setEnabled(!isInputNotOutput);
    ui->pushButton_inner->setEnabled(!isInputNotOutput);
}

void DialogData::slotCurrentIndexChanged(const QString &text)
{
    google::protobuf::Message* mesg = nullptr;
    if (CalcObjectByName(text.toStdString(), nullptr, &mesg) == false)
        return;
    MessageGuard guard(mesg);
    ui->plainTextEdit->setPlainText(DataExchanger::jsonByMsgObje(*mesg));
}

void DialogData::slotAccept()
{
    m_currType = 0;
    m_currData.clear();
    const google::protobuf::Descriptor* desc = nullptr;
    google::protobuf::Message* mesg = nullptr;
    if (CalcObjectByName(ui->comboBox_msgType->currentText().toStdString(), &desc, &mesg) == false)
    {
        QMessageBox::information(this, tr("name->object"), tr("名字->对象, 失败."));
        return;
    }
    MessageGuard guard(mesg);

    std::string jsonStr = ui->plainTextEdit->toPlainText().toStdString();
    if (google::protobuf::util::JsonStringToMessage(jsonStr, mesg) != google::protobuf::util::Status::OK)
    {
        QMessageBox::information(this, tr("json->object"), tr("JSON->对象, 失败."));
        return;
    }

    std::string binData;
    if (mesg->SerializeToString(&binData) == false)
    {
        QMessageBox::information(this, tr("object->binary"), tr("对象->二进制, 失败."));
        return;
    }

    m_currData.append(binData.data(), binData.size());
    m_currType = desc->index();

    this->accept();
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
    dataOut = m_currData;
    typeOut = m_currType;
}

void DialogData::setData(const QCommonNtosReq &data)
{
    switchMode(false);

    txdata::MsgType innerType = static_cast<txdata::MsgType>(data.ReqType);
    ui->pushButton_inner->setProperty("MsgType", DataExchanger::nameByMsgType(innerType));
    ui->pushButton_inner->setProperty("MsgJson", DataExchanger::jsonByMsgType(innerType, data.ReqData));

    txdata::CommonNtosReq dataTx;
    DataExchanger::toCommonNtosReq(data, dataTx);
    ui->pushButton_outer->setProperty("MsgType", QString::fromStdString(dataTx.GetDescriptor()->name()));
    ui->pushButton_outer->setProperty("MsgJson", DataExchanger::jsonByMsgObje(dataTx));
}

void DialogData::setData(const QCommonNtosRsp &data)
{
    switchMode(false);
}
