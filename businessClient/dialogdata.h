#ifndef DIALOGDATA_H
#define DIALOGDATA_H

#include <QDialog>
#include "sqlstruct.h"

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
    void getData(QByteArray& dataOut, int32_t& typeOut);
    void setData(const QCommonNtosReq& data);
    void setData(const QCommonNtosRsp& data);

private:
    void initUI();
    void switchMode(bool isInputNotOutput);

private slots:
    void slotCurrentIndexChanged(const QString &text);
    void slotAccept();

private:
    Ui::DialogData* ui;
    QByteArray      m_currData;
    int32_t         m_currType;
};

#endif // DIALOGDATA_H
