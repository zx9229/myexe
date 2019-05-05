#ifndef MYSQLTABLEMODEL_H
#define MYSQLTABLEMODEL_H

#include <QObject>
#include <QSqlTableModel>
#include <QQmlApplicationEngine>
#include <QDebug>
#include <QSqlQuery>
#include <QSqlRecord>
#include <QSqlField>

class MySqlTableModel : public QSqlTableModel
{
    Q_OBJECT
    Q_PROPERTY(QString selectStatement READ getSelectStatement WRITE setSelectStatement NOTIFY selectStatementChanged)

public:
    explicit MySqlTableModel(QObject *parent = nullptr, QSqlDatabase db = QSqlDatabase()) : QSqlTableModel(parent, db)
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
        this->select();
        emit selectStatementChanged();
    }

    int rowCount(const QModelIndex &parent = QModelIndex()) const override
    {
        return QSqlTableModel::rowCount(parent);
    }

    int columnCount(const QModelIndex &parent = QModelIndex()) const override
    {
        return QSqlTableModel::columnCount(parent);
    }

    QVariant data(const QModelIndex &idx, int role = Qt::DisplayRole) const override
    {
        if (Qt::UserRole <= role)
        {
            int col = role - Qt::UserRole;
            QModelIndex newIdx = idx.siblingAtColumn(col);
            return QSqlTableModel::data(newIdx, Qt::DisplayRole);
        }
        else
        {
            return QSqlTableModel::data(idx, role);
        }
    }

    QHash<int, QByteArray> roleNames() const override
    {
        //return QSqlTableModel::roleNames();
        QHash<int, QByteArray> roleNameHash;
        QStringList fieldNameList = this->nameList();
        for (int i = 0; i<fieldNameList.size(); i++)
        {
            roleNameHash[Qt::UserRole + i] = fieldNameList[i].toLatin1();
        }
        return roleNameHash;
    }

    Q_INVOKABLE QStringList nameList() const
    {
        QStringList fieldNameList;
        QSqlRecord rec = this->record();//having only the field names.
        for (int i = 0; i<rec.count(); i++)
        {
            fieldNameList.append(rec.field(i).name());
        }
        return fieldNameList;
    }

signals:
    void selectStatementChanged();

public slots:
    virtual bool select() override
    {
        bool bRet = QSqlTableModel::select();
        if (bRet) {
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

/*
void JustForTest()
{
    const char* QT_SQL_DEFAULT_CONNECTION = "qt_sql_default_connection";
    QSqlDatabase db;
    if (db.contains(QT_SQL_DEFAULT_CONNECTION))
    {
        db = QSqlDatabase::database(QT_SQL_DEFAULT_CONNECTION);
        qDebug() << "database" << QT_SQL_DEFAULT_CONNECTION;
    }
    else
    {
        db = QSqlDatabase::addDatabase("QSQLITE");
        db.setDatabaseName("_just4test.sqlite.db");
        qDebug() << "addDatabase" << "setDatabaseName";
    }
    Q_ASSERT(db.open());
    QSqlQuery sqlQuery;
    sqlQuery.exec(QObject::tr("CREATE TABLE IF NOT EXISTS student (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)"));
    sqlQuery.exec(QObject::tr("INSERT INTO student VALUES (3,'张三',23)"));
    sqlQuery.exec(QObject::tr("INSERT INTO student VALUES (4,'李四',24)"));
    sqlQuery.exec(QObject::tr("INSERT INTO student VALUES (5,'王五',25)"));
}
*/

/*
//文件JustForTest.qml的内容:
import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import MySqlTableModel 0.1
Item {
    ColumnLayout {
        anchors.fill: parent
        Button {
            Layout.fillWidth: true
            text: qsTr("刷新")
            onClicked: mstm.select()
        }
        TableView {
            Layout.fillHeight: true
            Layout.fillWidth: true
            columnSpacing: 1
            rowSpacing: 1
            clip: true
            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM student"
            }
            delegate: Rectangle {
                implicitHeight: 50
                implicitWidth: 100
                Text {
                    text: display
                }
            }
        }
    }
}
*/
