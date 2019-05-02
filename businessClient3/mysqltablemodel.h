#ifndef MYSQLTABLEMODEL_H
#define MYSQLTABLEMODEL_H

#include <QObject>
#include <QSqlTableModel>
#include <QQmlApplicationEngine>
#include <QDebug>
#include <QSqlQuery>

class MySqlTableModel : public QSqlTableModel
{
    Q_OBJECT
    Q_PROPERTY(QString selectStatement READ getSelectStatement WRITE setSelectStatement NOTIFY selectStatementChanged)

public:
    explicit MySqlTableModel(QObject *parent = nullptr, QSqlDatabase db = QSqlDatabase()) : QSqlTableModel (parent, db)
    {
    }

    static void doQmlRegisterType()
    {
        qmlRegisterType<MySqlTableModel>("MySqlTableModel", 0, 1, "MySqlTableModel");
    }

public:
    QString getSelectStatement() const
    {
        return m_selectStatement;
    }

    void setSelectStatement(const QString& statement)
    {
        m_selectStatement = statement;
        emit selectStatementChanged();
    }

signals:
    void selectStatementChanged();

public slots:
    virtual bool select() override
    {
        bool bRet = QSqlTableModel::select();
        if(bRet) {
            while (canFetchMore()) {
                fetchMore();
            }
        }
        return bRet;
    }

protected:
    QString selectStatement() const override
    {
        //TODO:暂不考虑filter()等.
        return m_selectStatement;
    }

private:
    QString m_selectStatement;
};

#endif // MYSQLTABLEMODEL_H
