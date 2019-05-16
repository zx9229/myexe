import QtQuick 2.4
import QtQuick.Layouts 1.3
import QtQuick.Controls 1.6
import QtQuick.Controls.Styles 1.4

Item {
    id: itemitem
    width: 400
    height: 400
    property alias labelDT: labelDT
    property alias textAreaMessage: textAreaMessage

    GridLayout {
        id: gridLayout
        rows: 7
        columns: 2
        anchors.fill: parent

        Label {
            id: labelDT
            text: qsTr("DATE_TIME")
            horizontalAlignment: Text.AlignHCenter
            verticalAlignment: Text.AlignVCenter
            Layout.fillWidth: true
            Layout.columnSpan: 2
        }

        Label {
            id: labelURL
            text: qsTr("URL")
        }

        TextField {
            id: textFieldURL
            Layout.fillWidth: true
            readOnly: true
            style: TextFieldStyle {
                background: Rectangle {
                    color: "lightgray"
                }
            }
        }

        Label {
            id: labelHost
            text: qsTr("Host")
        }

        TextField {
            id: textFieldHost
            Layout.fillWidth: true
            placeholderText: qsTr("主机")
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
            placeholderText: qsTr("端口")
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
            placeholderText: qsTr("父节点的名字")
        }

        Label {
            id: labelUserID
            text: qsTr("UserID")
        }

        TextField {
            id: textFieldUserID
            Layout.fillWidth: true
            placeholderText: qsTr("本节点的名字")
        }

        RowLayout {
            id: rowLayout
            Layout.fillWidth: true
            Layout.columnSpan: 2

            Button {
                id: buttonSaveConf
                text: qsTr("保存配置")
                Layout.fillWidth: true
            }

            Button {
                id: buttonLoadConf
                text: qsTr("载入配置")
                Layout.fillWidth: true
            }
        }

        Button {
            id: buttonSignIn
            text: qsTr("登录")
            Layout.fillWidth: true
            Layout.columnSpan: 2
        }

        TextArea {
            id: textAreaMessage
            text: qsTr("Text Area")
            Layout.columnSpan: 2
            Layout.fillHeight: true
            Layout.fillWidth: true
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

    Connections {
        target: textFieldHost
        onEditingFinished: {
            textFieldURL.text = "ws://%1:%2/websocket".arg(
                        textFieldHost.text).arg(textFieldPort.text)
        }
    }

    Connections {
        target: textFieldPort
        onEditingFinished: {
            textFieldURL.text = "ws://%1:%2/websocket".arg(
                        textFieldHost.text).arg(textFieldPort.text)
        }
    }

    Connections {
        target: buttonSaveConf
        onClicked: {
            dataExch.saveValue("Host", textFieldHost.text)
            dataExch.saveValue("Port", textFieldPort.text)
            dataExch.saveValue("BelongID", textFieldBelongID.text)
            dataExch.saveValue("UserID", textFieldUserID.text)
        }
    }

    Connections {
        target: buttonLoadConf
        onClicked: {
            textFieldHost.text = dataExch.loadValue("Host")
            textFieldPort.text = dataExch.loadValue("Port")
            textFieldBelongID.text = dataExch.loadValue("BelongID")
            textFieldUserID.text = dataExch.loadValue("UserID")
            if (textFieldHost.text == "" && textFieldPort.text == "") {
                textFieldHost.text = qsTr("192.168.3.157")
                textFieldPort.text = qsTr("40078")
                textFieldBelongID.text = qsTr("n4")
                textFieldUserID.text = qsTr("ZXCVB")
            }
            textFieldURL.text = "ws://%1:%2/websocket".arg(
                        textFieldHost.text).arg(textFieldPort.text)
        }
    }
}


/*##^## Designer {
    D{i:1;anchors_height:100;anchors_width:100;anchors_x:107;anchors_y:81}
}
 ##^##*/
