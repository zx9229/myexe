import QtQuick 2.0
import QtQuick.Controls 1.6
import QtQuick.Layouts 1.3
import MySqlTableModel 0.1

Item {
    ColumnLayout {
        anchors.fill: parent
        Button{
            onClicked: {
                mstm.select()
                function refreshColumn(objTableView, lstFieldName) {
                    for(var i = objTableView.columnCount-1; i >= 0; i--)
                    {
                        objTableView.removeColumn(i)
                    }
                    for(i = 0; i < lstFieldName.length; i++)
                    {
                        var qmlStr = "import QtQuick 2.4; import QtQuick.Controls 1.4; TableViewColumn { width: 100; role: \""+lstFieldName[i]+"\"; title: \""+lstFieldName[i]+"\" }"
                        objTableView.addColumn(Qt.createQmlObject(qmlStr, objTableView, "dynamicSnippet1"))
                    }
                }
                refreshColumn(tView,mstm.nameList())
            }
        }

        TableView {
            id: tView
            Layout.fillWidth: true
            Layout.fillHeight: true
            TableViewColumn {
                role: "UserID"
                title: "UserID"
                width: 100
            }
            TableViewColumn {
                role: "BelongID"
                title: "BelongID"
                width: 100
            }
            TableViewColumn {
                role: "Version"
                title: "Version"
                width: 100
            }
            TableViewColumn {
                role: "ExePid"
                title: "ExePid"
                width: 100
            }
            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM ConnInfoEx"
            }
        }
    }
}
