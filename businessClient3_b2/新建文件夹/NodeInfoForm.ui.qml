import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import MySqlTableModel 0.1
import QtQml.Models 2.12

Item {
    width: 400
    height: 400

    ColumnLayout {
        id: columnLayout
        anchors.fill: parent
        Layout.fillHeight: true
        Layout.fillWidth: true

        RowLayout {
            id: rowLayout

            //width: 100
            //height: 100
            Button {
                id: button
                text: qsTr("Button")
            }

            Button {
                id: button1
                text: qsTr("Button")
            }
        }
        TableView {
            clip: true
            rowSpacing: 1
            columnSpacing: 1
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM ConnInfoEx"
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

    Connections {
        target: button
        onClicked: mstm.select()
    }
}




/*##^## Designer {
    D{i:1;anchors_height:100;anchors_width:100;anchors_x:120;anchors_y:81}
}
 ##^##*/
