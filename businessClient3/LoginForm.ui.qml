import QtQuick 2.4
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3

Item {
    id: itemitem
    width: 400
    height: 400
    property alias labelDT: labelDT

    GridLayout {
        id: gridLayout
        rows: 7
        columns: 3
        anchors.fill: parent

        Label {
            id: labelDT
            text: qsTr("Label")
            verticalAlignment: Text.AlignVCenter
            horizontalAlignment: Text.AlignHCenter
            Layout.fillWidth: true
            Layout.columnSpan: 3
        }

        Label {
            id: labelURL
            text: qsTr("URL")
        }

        TextField {
            id: textFieldURL
            text: qsTr("Text Field")
            Layout.fillWidth: true
        }

        Button {
            id: buttonQuickFill
            text: qsTr("速填")
            Layout.fillWidth: false
        }

        Label {
            id: labelBelongID
            text: qsTr("BelongID")
        }

        TextField {
            id: textFieldBelongID
            text: qsTr("Text Field")
            Layout.columnSpan: 2
            Layout.fillWidth: true
        }

        Label {
            id: labelUserID
            text: qsTr("UserID")
        }

        TextField {
            id: textFieldUserID
            text: qsTr("Text Field")
            Layout.columnSpan: 2
            Layout.fillWidth: true
        }

        Button {
            id: buttonSignIn
            text: qsTr("登录")
            Layout.columnSpan: 3
            Layout.fillWidth: true
        }

        TextArea {
            id: textAreaMessage
            text: qsTr("Text Area")
            Layout.columnSpan: 3
            Layout.fillHeight: true
            Layout.fillWidth: true
            background: Rectangle {
                radius: 2
                border.color: "blue"
                border.width: 2
                //implicitHeight: 50
                //implicitWidth: 100
            }
        }
    }

    Connections {
        target: buttonQuickFill
        onClicked: {
            textFieldURL.text = qsTr("ws://localhost:65535/websocket")
        }
    }

    Connections {
        target: buttonSignIn
        onClicked: {
            dataExch.setURL(textFieldURL.text)
            dataExch.setOwnInfo(textFieldUserID.text, textFieldBelongID.text)
            dataExch.start()
        }
    }
}
