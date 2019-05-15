import QtQuick 2.11
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4
import QtQuick.Controls 1.4 as Controls1
import QtQuick.Controls.Styles 1.4
import MySqlTableModel 0.1

Pane {
    signal sigShowReqRsp(string UserID, var MsgNo)//The var type is a generic property type that can refer to any data type.

    property alias statement: mstm.selectStatement

    ColumnLayout {
        anchors.fill: parent

        Button {
            text: qsTr("刷新NodeRequest")
            onClicked: mstm.select()
        }

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id:mstm
            }
            delegate: Column {
                id: idColumn
                width: listView.width
                height: txt1.implicitHeight + messageText.implicitHeight + 10
                Label {
                    id: txt1
                    text: "%1 %2 %3 %4".arg(InsertTime).arg(TxTypeTxt).arg(RspCnt).arg(IsLast)
                }
                Rectangle {
                    id: idRect
                    radius: 5
                    height: messageText.implicitHeight + 24
                    width: Math.min(messageText.implicitWidth + 24,listView.width)
                    TextEdit {
                        id: messageText
                        text: TxDataTxt
                        anchors.fill: parent
                        anchors.margins: 12
                        wrapMode: Label.Wrap
                        readOnly: true
                        selectByMouse: true
                    }
                    MouseArea {
                        anchors.fill: parent
                        onClicked: {
                            idColumn.ListView.view.currentIndex = index
                        }
                        onPressAndHold: {
                            sigShowReqRsp(UserID, MsgNo)
                        }
                    }
                }
                states: State {
                    when: idColumn.ListView.isCurrentItem
                    PropertyChanges {
                        target: idRect
                        color: "green"
                    }
                }
            }
            ScrollBar.vertical: ScrollBar{}
        }

        Button {
            Layout.fillWidth: true
            text: qsTr("显示发送面板")
            onClicked: {
                paneSend.visible = !paneSend.visible
                text = paneSend.visible ? qsTr("隐藏发送面板") : qsTr("显示发送面板")
            }
        }

        Pane {
            id: paneSend
            visible: false
            Layout.fillWidth: true
            ColumnLayout {
                GroupBox {
                    Row {
                        Controls1.CheckBox {
                            id: cbIsLog
                            text: qsTr("IsLog")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbIsSafe
                            text: qsTr("IsSafe")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbIsPush
                            text: qsTr("IsPush")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbUpCache
                            text: qsTr("UpCache")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbToDB
                            text: qsTr("ToDB")
                            checked: false
                        }
                    }
                }
                GroupBox {
                    RowLayout {
                        Controls1.RadioButton {
                            id:rbC1REQ
                            text: qsTr("C1Req")
                            checked: true
                        }
                        Controls1.RadioButton {
                            text: qsTr("C2Req")
                        }
                    }
                }
                Controls1.ComboBox {
                    id:idComboBox
                    Layout.fillWidth: true
                    model: dataExch.getTxMsgTypeNameList()
                }
                Button {
                    text: qsTr("填充示例JSON")
                    onClicked: {
                        idTextArea.text = dataExch.jsonExample(idComboBox.currentText)
                    }
                }
                Button {
                    text: qsTr("发送")
                    onClicked: {
                        var message = dataExch.demoFun(idComboBox.currentText,idTextArea.text,"",cbIsLog.checked,cbIsSafe.checked,cbIsPush.checked,cbUpCache.checked,rbC1REQ.checked)
                        ToolTip.show("SUCCESS:"+message, 5000)
                    }
                }
                TextArea {
                    id: idTextArea
                    Layout.fillWidth: true
                    Layout.fillHeight: true
                    background: Rectangle {
                        border.width: 2
                        border.color: "blue"
                    }
                }
            }
        }
    }
}
