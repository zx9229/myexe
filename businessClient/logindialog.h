#ifndef LOGINDIALOG_H
#define LOGINDIALOG_H

#include <QDialog>
#include <QComboBox>
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

private:
    void initUI();
    void setComboBox4ProgramType(QComboBox* obj);

private slots:
    void slotClickedLogin();
    void slotClickedClear();
    void slotClickedQuickFill();

private:
    Ui::LoginDialog* ui;
    DataExchanger*   m_dataExch;
};

#endif // LOGINDIALOG_H
