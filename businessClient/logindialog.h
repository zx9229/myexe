#ifndef LOGINDIALOG_H
#define LOGINDIALOG_H

#include <QDialog>
#include "dataexchanger.h"

namespace Ui {
    class LoginDialog;
}

class LoginDialog : public QDialog
{
    Q_OBJECT

public:
    explicit LoginDialog(DataExchanger* p, QWidget *parent = 0);
    ~LoginDialog();

private slots:
    void slotClickedLogin();

private:
    Ui::LoginDialog* ui;
    DataExchanger*   m_dataExch;
};

#endif // LOGINDIALOG_H
