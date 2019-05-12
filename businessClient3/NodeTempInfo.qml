import QtQuick 2.0
import QtQuick.Controls 1.6
import QtQuick.Controls.Styles 1.4
import QtQuick.Layouts 1.3
import MySqlTableModel 0.1

Item {
    ColumnLayout {
        anchors.fill: parent
        Button{
            text: qsTr("NodeTempInfo刷新")
            Layout.fillWidth: true
            onClicked: {
                mstm.select()
                function refreshColumn(objTableView, lstFieldName) {
                    for(var i = objTableView.columnCount-1; i >= 0; i--)
                    {
                        objTableView.removeColumn(i)
                    }
                    for(i = 0; i < lstFieldName.length; i++)
                    {
                        var qmlStr = "import QtQuick 2.0; import QtQuick.Controls 1.6; TableViewColumn { width: 100; role: \""+lstFieldName[i]+"\"; title: \""+lstFieldName[i]+"\" }"
                        objTableView.addColumn(Qt.createQmlObject(qmlStr, objTableView, "dynamicSnippet1"))
                    }
                }
                refreshColumn(tView, mstm.nameList())
            }
        }

        TableView {
            id: tView
            Layout.fillWidth: true   //QML Item: Binding loop detected for property "width"
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
            onClicked: console.log("onClicked", row)
            onDoubleClicked: console.log("onDoubleClicked", row)
            onPressAndHold: console.log("onPressAndHold", row, mstm.qryData(row, 0))
        }
    }
}
