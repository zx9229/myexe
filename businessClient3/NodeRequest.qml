import QtQuick 2.11
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4
import QtQuick.Controls 1.4 as Controls1
import QtQuick.Controls.Styles 1.4
import MySqlTableModel 0.1

Pane {
    signal sigShowNodeList()
    signal sigShowNodeReqRsp(string UserID, var MsgNo)//The var type is a generic property type that can refer to any data type.
    //property alias statement: mstm.selectStatement
    property string peerid

    ColumnLayout {
        anchors.fill: parent

        RowLayout {
            Button {
                text: qsTr("<[返回]")
                onClicked: sigShowNodeList()
            }
            Button {
                text: qsTr("刷新NodeRequest")
                onClicked: mstm.select()
            }
        }

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id:mstm
                selectStatement: "SELECT * FROM CommonData WHERE MsgType IN(3,5) AND PeerID='%1'".arg(peerid)
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
                            sigShowNodeReqRsp(UserID, MsgNo)
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
                            text: qsTr("Log")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbIsSafe
                            text: qsTr("Safe")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbIsPush
                            text: qsTr("Push")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbIsUpCache
                            text: qsTr("UpCache")
                            checked: false
                        }
                        Controls1.CheckBox {
                            id: cbForceToDB
                            text: qsTr("DB")
                            checked: false
                        }
                    }
                }
                RowLayout {
                    GroupBox {
                        RowLayout {
                            Controls1.RadioButton {
                                id:rbC1REQ
                                text: qsTr("C1Q")
                                checked: true
                            }
                            Controls1.RadioButton {
                                text: qsTr("C2Q")
                            }
                        }
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
                            var message = dataExch.demoFun(idComboBox.currentText,idTextArea.text,peerid,cbIsLog.checked,cbIsSafe.checked,cbIsPush.checked,cbIsUpCache.checked,rbC1REQ.checked,cbForceToDB.checked)
                            ToolTip.show("SUCCESS:"+message, 5000)
                        }
                    }
                }

                Controls1.ComboBox {
                    id:idComboBox
                    Layout.fillWidth: true
                    model: dataExch.getTxMsgTypeNameList()
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
