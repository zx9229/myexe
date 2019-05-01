import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12

Item {
    id: itemitem
    width: 400
    height: 400
    property alias labelDT: labelDT
    property alias textAreaMessage: textAreaMessage

    GridLayout {
        id: gridLayout
        rows: 7
        columns: 3
        anchors.fill: parent

        Label {
            id: labelDT
            text: qsTr("DATE_TIME")
            verticalAlignment: Text.AlignVCenter
            horizontalAlignment: Text.AlignHCenter
            Layout.columnSpan: 3
            Layout.fillWidth: true
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




/*##^## Designer {
    D{i:1;anchors_height:100;anchors_width:100;anchors_x:91;anchors_y:77}
}
 ##^##*/
