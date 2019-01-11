#ifndef LOGINDIALOG_H
#define LOGINDIALOG_H

#include <QDialog>
#include <QComboBox>

namespace Ui {
    class LoginDialog;
}

class DataExchanger;
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
    void slotReady();
    void slotClickedLogin();
    void slotClickedClear();
    void slotClickedQuickFill();

private:
    Ui::LoginDialog* ui;
    DataExchanger*   m_dataExch;
};

#endif // LOGINDIALOG_H
