import QtQuick 2.4
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3

Item {
    width: 400
    height: 400

    GridLayout {
        id: gridLayout
        rows: 7
        columns: 2
        anchors.fill: parent

        Label {
            id: label
            text: qsTr("Label")
        }

        TextField {
            id: textField
            text: qsTr("Text Field")
            Layout.fillWidth: true
        }

        Label {
            id: label1
            text: qsTr("Label")
        }

        TextField {
            id: textField1
            text: qsTr("Text Field")
            Layout.fillWidth: true
        }

        Label {
            id: label2
            text: qsTr("Label")
        }

        TextField {
            id: textField2
            text: qsTr("Text Field")
            Layout.fillWidth: true
        }

        Button {
            id: button
            text: qsTr("Button")
            Layout.columnSpan: 2
            Layout.fillWidth: true
        }

        TextArea {
            id: textArea
            text: qsTr("Text Area")
            Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.columnSpan: 2
            background: Rectangle {
                radius: 2
                border.color: "blue"
                border.width: 2
                //implicitHeight: 50
                //implicitWidth: 100
            }
        }
    }
}
