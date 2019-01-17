#ifndef DIALOGDATA_H
#define DIALOGDATA_H

#include <QDialog>

namespace Ui {
    class DialogData;
}

class DialogData : public QDialog
{
    Q_OBJECT

public:
    explicit DialogData(QWidget *parent = 0);
    ~DialogData();

public:
    bool needResp();
    bool needSave();
    void getData(QByteArray& dataOut, int& typeOut);

private:
    void initUI();

private slots:
    void slotCurrentIndexChanged(const QString &text);
    void slotAccept();

private:
    Ui::DialogData* ui;
    QByteArray      m_currData;
    int             m_currType;
};

#endif // DIALOGDATA_H
