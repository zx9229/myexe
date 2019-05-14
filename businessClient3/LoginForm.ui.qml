import QtQuick 2.4
import QtQuick.Controls 1.6
import QtQuick.Layouts 1.3

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
            horizontalAlignment: Text.AlignHCenter
            verticalAlignment: Text.AlignVCenter
            Layout.fillWidth: true
            Layout.columnSpan: 3
        }

        Label {
            id: labelURL
            text: qsTr("URL")
        }

        TextField {
            id: textFieldURL
            Layout.fillWidth: true
            placeholderText: qsTr("Text Field")
        }

        Button {
            id: buttonQuickFill
            text: qsTr("速填")
        }

        Label {
            id: labelHost
            text: qsTr("Host")
        }

        TextField {
            id: textFieldHost
            Layout.fillWidth: true
            Layout.columnSpan: 2
            placeholderText: qsTr("Text Field")
            validator: RegExpValidator {
                regExp: /[^`~!@#$%^&*()_=+\[\]{}\\|;:'",<>/?]+/
            }
        }

        Label {
            id: labelPort
            text: qsTr("Port")
        }

        TextField {
            id: textFieldPort
            Layout.fillWidth: true
            Layout.columnSpan: 2
            placeholderText: qsTr("Text Field")
            validator: IntValidator {
                bottom: 0
                top: 65535
            }
        }

        Label {
            id: labelBelongID
            text: qsTr("BelongID")
        }

        TextField {
            id: textFieldBelongID
            Layout.fillWidth: true
            Layout.columnSpan: 2
            placeholderText: qsTr("Text Field")
        }

        Label {
            id: labelUserID
            text: qsTr("UserID")
        }

        TextField {
            id: textFieldUserID
            Layout.fillWidth: true
            Layout.columnSpan: 2
            placeholderText: qsTr("Text Field")
        }

        Button {
            id: buttonSignIn
            text: qsTr("登录")
            Layout.fillWidth: true
            Layout.columnSpan: 3
        }

        TextArea {
            id: textAreaMessage
            text: qsTr("Text Area")
            Layout.columnSpan: 3
            Layout.fillHeight: true
            Layout.fillWidth: true
        }
    }

    Connections {
        target: buttonQuickFill
        onClicked: {
            textFieldURL.text = qsTr("ws://localhost:65535/websocket")
            textFieldHost.text = qsTr("192.168.3.157")
            textFieldPort.text = qsTr("40078")
            textFieldURL.text = "ws://%1:%2/websocket".arg(
                        textFieldHost.text).arg(textFieldPort.text)
            textFieldUserID.text = qsTr("ZXCVB")
            textFieldBelongID.text = qsTr("n4")
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
    D{i:1;anchors_height:100;anchors_width:100;anchors_x:107;anchors_y:81}
}
 ##^##*/
