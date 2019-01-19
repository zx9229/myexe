#ifndef DIALOGREQRSP_H
#define DIALOGREQRSP_H

#include <QDialog>

namespace Ui {
    class DialogReqRsp;
}

class DataExchanger;
class DialogReqRsp : public QDialog
{
    Q_OBJECT

public:
    explicit DialogReqRsp(DataExchanger* p, QWidget *parent = 0);
    ~DialogReqRsp();

public:
    void initUI();
    void setRefNum(int64_t refNum);

private slots:
    void slotReload();
    void slotCellDoubleClicked(int row, int column);

private:
    Ui::DialogReqRsp* ui;
    DataExchanger*    m_dataExch;
};

#endif // DIALOGREQRSP_H
