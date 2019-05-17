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
        for (int i = 0; i < fieldNameList.size(); i++)
        {
            roleNameHash[Qt::UserRole + i] = fieldNameList[i].toLatin1();
        }
        return roleNameHash;
    }

    Q_INVOKABLE QStringList nameList() const
    {
        //function refreshColumn(objTableView, lstFieldName) {
        //    for(var i = objTableView.columnCount-1; i >= 0; i--)
        //    {
        //        objTableView.removeColumn(i)
        //    }
        //    for(i = 0; i < lstFieldName.length; i++)
        //    {
        //        var qmlStr = "import QtQuick 2.4; import QtQuick.Controls 1.4; TableViewColumn { width: 100; role: \""+lstFieldName[i]+"\"; title: \""+lstFieldName[i]+"\" }"
        //        objTableView.addColumn(Qt.createQmlObject(qmlStr, objTableView, "dynamicSnippet1"))
        //    }
        //}
        QStringList fieldNameList;
        QSqlRecord rec = this->record();//having only the field names.
        for (int i = 0; i < rec.count(); i++)
        {
            fieldNameList.append(rec.field(i).name());
        }
        return fieldNameList;
    }

    Q_INVOKABLE QVariant qryData(int rowIdx, int colIdx) const
    {
        return QSqlTableModel::data(index(rowIdx, colIdx), Qt::DisplayRole);
    }

    Q_INVOKABLE QStringList tmpQmlList(QObject* objTableView)
    {
        QStringList qmlList;
        QStringList fieldNameList = this->nameList();
        for (int i = 0; i < fieldNameList.size(); i++) {
            QString sss;
            sss += QString(R"(import QtQuick 2.4;)");
            sss += QString(R"(import QtQuick.Controls 1.4;)");
            sss += QString(R"(TableViewColumn { width : 100 ; role : "%1" ; title : "%1" } )").arg(fieldNameList[i]);
            qmlList.append(sss);
        }
        return qmlList;
        do
        {
            bool isOk = false;
            QVariant columnCount = objTableView->property("columnCount");
            if (!columnCount.isValid()) { break; }
            for (int i = columnCount.toInt() - 1; i >= 0; i--)
            {
                QVariant returnedValue;
                QVariant idx = i;
                isOk = QMetaObject::invokeMethod(objTableView, "removeColumn", Q_RETURN_ARG(QVariant, returnedValue), Q_ARG(QVariant, idx));
                if (!isOk) { break; }
            }
            QStringList fieldNameList = this->nameList();
            QVariant vObjTableView;
            vObjTableView.setValue(objTableView);
            QVariant dynamicSnippet1 = "dynamicSnippet1";
            for (int i = 0; i < fieldNameList.size(); i++) {
                QString sss;
                sss += QString(R"(import QtQuick 2.4;)");
                sss += QString(R"(import QtQuick.Controls 1.4;)");
                sss += QString(R"(TableViewColumn { width : 100 ; role : "%1" ; title : "%1" } )").arg(fieldNameList[i]);
                QVariant vQmlObj;
                //Qt.createQmlObject(sss,tView,"dynamicSnippet1")
                isOk = QMetaObject::invokeMethod(nullptr, "createQmlObject", Q_RETURN_ARG(QVariant, vQmlObj), Q_ARG(QVariant, QVariant(sss)), Q_ARG(QVariant, vObjTableView), Q_ARG(QVariant, dynamicSnippet1));
                if (!isOk) { break; }
                if (!vQmlObj.isValid()) { break; }
                QVariant returnedValue;
                isOk = QMetaObject::invokeMethod(objTableView, "addColumn", Q_RETURN_ARG(QVariant, returnedValue), Q_ARG(QVariant, vQmlObj));
                if (!isOk) { break; }
            }
        } while (false);
        return qmlList;
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
        else {
            qDebug() << "MySqlTableModel" << bRet << this->selectStatement();
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
