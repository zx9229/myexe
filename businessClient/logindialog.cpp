#include "logindialog.h"
#include "ui_logindialog.h"
#include <QMetaEnum>
#include "dataexchanger.h"

LoginDialog::LoginDialog(DataExchanger* p, QWidget *parent) :
    QDialog(parent),
    ui(new Ui::LoginDialog),
    m_dataExch(p)
{
    ui->setupUi(this);
    this->setWindowFlags(this->windowFlags() | Qt::WindowMinimizeButtonHint | Qt::WindowMaximizeButtonHint);
    this->showFullScreen();
    initUI();
    QObject::connect(ui->pushButton_cancel, &QPushButton::clicked, this, &LoginDialog::reject);
    QObject::connect(ui->pushButton_login, &QPushButton::clicked, this, &LoginDialog::slotClickedLogin);
    QObject::connect(ui->pushButton_clear, &QPushButton::clicked, this, &LoginDialog::slotClickedClear);
    QObject::connect(ui->pushButton_quickFill, &QPushButton::clicked, this, &LoginDialog::slotClickedQuickFill);
    QObject::connect(m_dataExch, &DataExchanger::sigReady, this, &LoginDialog::slotReady);
    QObject::connect(m_dataExch, &DataExchanger::sigWebsocketError, this, &LoginDialog::slotWebsocketError);
}

LoginDialog::~LoginDialog()
{
    delete ui;
}

void LoginDialog::initUI()
{
    ui->lineEdit_port->setValidator(new QIntValidator());
    if (true) {
        ui->lineEdit_UserZoneName->setValidator(new QRegExpValidator(QRegExp("^[^/]+$")));
        ui->lineEdit_UserNodeName->setValidator(new QRegExpValidator(QRegExp("^[^/]+$")));
        ui->lineEdit_UserExecName->setValidator(new QRegExpValidator(QRegExp("^[^/]+$")));
        setComboBox4ProgramType(ui->comboBox_UserExecType);
    }
    if (true) {
        ui->lineEdit_BelongZoneName->setValidator(new QRegExpValidator(QRegExp("^[^/]+$")));
        ui->lineEdit_BelongNodeName->setValidator(new QRegExpValidator(QRegExp("^[^/]+$")));
        ui->lineEdit_BelongExecName->setValidator(new QRegExpValidator(QRegExp("^[^/]+$")));
        setComboBox4ProgramType(ui->comboBox_BelongExecType);
    }
}

void LoginDialog::setComboBox4ProgramType(QComboBox *comboBox)
{
    comboBox->clear();
    for (int i = static_cast<int>(txdata::ProgramType_MIN); i <= static_cast<int>(txdata::ProgramType_MAX); ++i)
    {
        txdata::ProgramType curType = static_cast<txdata::ProgramType>(i);
        std::string curDesc = txdata::ProgramType_Name(curType);
        comboBox->addItem(QString::fromStdString(curDesc), curType);
    }
}

void LoginDialog::slotWebsocketError(QAbstractSocket::SocketError error)
{
    QString text = QString("%1 - %2")
        .arg(QDateTime::currentDateTime().toString("yyyy-MM-dd HH:mm:ss"))
        .arg(QMetaEnum::fromType<QAbstractSocket::SocketError>().valueToKey(error));

    ui->plainTextEdit_message->setPlainText(text);
}

void LoginDialog::slotReady()
{
    QObject::disconnect(m_dataExch, &DataExchanger::sigWebsocketError, this, &LoginDialog::slotWebsocketError);
    QObject::disconnect(m_dataExch, &DataExchanger::sigReady, this, &LoginDialog::slotReady);
    this->accept();
}

void LoginDialog::slotClickedLogin()
{
    if (!ui->lineEdit_host->text().trimmed().isEmpty() &&
        !ui->lineEdit_port->text().trimmed().isEmpty() &&
        !ui->lineEdit_path->text().trimmed().isEmpty())
    {
        QString url = QString("ws://%1:%2%3")
            .arg(ui->lineEdit_host->text().trimmed())
            .arg(ui->lineEdit_port->text().trimmed())
            .arg(ui->lineEdit_path->text().trimmed());
        ui->lineEdit_url->setText(url);
    }
    m_dataExch->setURL(ui->lineEdit_url->text().trimmed());

    m_dataExch->setUserKey(
        ui->lineEdit_UserZoneName->text().trimmed(),
        ui->lineEdit_UserNodeName->text().trimmed(),
        static_cast<txdata::ProgramType>(ui->comboBox_UserExecType->currentData().toInt()),
        ui->lineEdit_UserExecName->text().trimmed());

    m_dataExch->setBelongKey(
        ui->lineEdit_BelongZoneName->text().trimmed(),
        ui->lineEdit_BelongNodeName->text().trimmed(),
        static_cast<txdata::ProgramType>(ui->comboBox_BelongExecType->currentData().toInt()),
        ui->lineEdit_BelongExecName->text().trimmed());

    m_dataExch->start();
}

void LoginDialog::slotClickedClear()
{
    ui->lineEdit_host->clear();
    ui->lineEdit_port->clear();
    ui->lineEdit_path->clear();
    ui->lineEdit_url->clear();
    if (true) {
        ui->lineEdit_UserZoneName->clear();
        ui->lineEdit_UserNodeName->clear();
        ui->lineEdit_UserExecName->clear();
        ui->comboBox_UserExecType->setCurrentIndex(0);
    }
    if (true) {
        ui->lineEdit_BelongZoneName->clear();
        ui->lineEdit_BelongNodeName->clear();
        ui->lineEdit_BelongExecName->clear();
        ui->comboBox_BelongExecType->setCurrentIndex(0);
    }
}

void LoginDialog::slotClickedQuickFill()
{
    this->slotClickedClear();
    ui->lineEdit_host->setText("localhost");
    ui->lineEdit_port->setText("10083");
    ui->lineEdit_path->setText("/websocket");
    ui->lineEdit_url->setText("ws://localhost:10083/websocket");
}
