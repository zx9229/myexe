#include "logindialog.h"
#include "ui_logindialog.h"

LoginDialog::LoginDialog(DataExchanger* p, QWidget *parent) :
    QDialog(parent),
    ui(new Ui::LoginDialog),
    m_dataExch(p)
{
    ui->setupUi(this);
    connect(ui->pushButton_cancel, &QPushButton::clicked, this, &LoginDialog::reject);
    connect(ui->pushButton_login, &QPushButton::clicked, this, &LoginDialog::slotClickedLogin);

    ui->lineEdit_url->setText("ws://localhost:10083/websocket");
}

LoginDialog::~LoginDialog()
{
    delete ui;
}

void LoginDialog::slotClickedLogin()
{
    if (ui->lineEdit_url->text().trimmed().isEmpty())
    {
        QString url = QString("ws://%1:%2%3")
            .arg(ui->lineEdit_host->text().trimmed())
            .arg(ui->lineEdit_port->text().trimmed())
            .arg(ui->lineEdit_path->text().trimmed());
        ui->lineEdit_url->setText(url);
    }
    m_dataExch->Login(ui->lineEdit_url->text().trimmed(), "", "");
    //TODO:登录成功
    this->accept();
}
