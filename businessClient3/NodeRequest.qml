import QtQuick 2.11
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4
import QtQuick.Controls 1.4 as Controls1
import QtQuick.Controls.Styles 1.4
import MySqlTableModel 0.1

Item {
    signal sigPickMsgNo(string MsgNo)
    //property alias statement: mstm.selectStatement
    property string peerid

    ColumnLayout {
        anchors.fill: parent

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id:mstm
                selectStatement: "SELECT * FROM CommonData WHERE MsgType IN(1,3) AND PeerID='%1'".arg(peerid)
            }
            delegate: Column {
                id: idColumn
                width: listView.width
                height: txt1.implicitHeight + messageText.implicitHeight + 10
                Label {
                    id: txt1
                    text: "%1 %2 (%3) %4".arg(InsertTime).arg(TxTypeTxt).arg(IsLast).arg(RspCnt)
                }
                Rectangle {
                    id: idRect
                    radius: 5
                    color: "lightgray"
                    height: messageText.implicitHeight + 24
                    width: Math.min(messageText.implicitWidth + 24,listView.width)
                    TextEdit {
                        id: messageText
                        text: TxDataTxt
                        anchors.fill: parent
                        anchors.margins: 12
                        wrapMode: Label.Wrap
                        readOnly: true
                        selectByMouse: false
                    }
                    MouseArea {
                        anchors.fill: parent
                        onClicked: {
                            idColumn.ListView.view.currentIndex = index
                            sigPickMsgNo(MsgNo.toString())
                        }
                        onDoubleClicked: {
                            messageText.selectByMouse = !messageText.selectByMouse
                        }
                        onPressAndHold: {
                            menu.userid = UserID
                            menu.msgno = MsgNo
                            menu.jsonText = TxDataTxt
                            menu.popup()
                        }
                    }
                }
                states: State {
                    when: idColumn.ListView.isCurrentItem
                    PropertyChanges {
                        target: idRect
                        color: "tan"
                    }
                }
            }
            ScrollBar.vertical: ScrollBar {
                id: verScrollBar
            }
            header: RefreshView {
                id: rv_refresh
                tips: qsTr("刷新中...")
                onRefeash: {
                    timer.start()
                }
                Timer {
                    id: timer
                    interval:300; running: false; repeat: false
                    onTriggered: {
                        mstm.select()
                        rv_refresh.hideView()
                    }
                }
            }
            footer: RefreshView {
                id: rv_load
                tips: qsTr("加载更多")
                onRefeash: {
                    loadMoreTimer.start()
                }
                Timer {
                    id: loadMoreTimer
                    interval:300; running: false; repeat: false
                    onTriggered: {
                        mstm.select()
                        rv_load.hideView()
                    }
                }
            }
        }

        Button {
            Layout.fillWidth: true
            text: qsTr("显示发送面板")
            onClicked: {
                paneSend.visible = !paneSend.visible
                text = paneSend.visible ? "隐藏发送面板" : "显示发送面板"
            }
        }

        Pane {
            id: paneSend
            visible: false
            Layout.fillWidth: true
            background: Rectangle {
                color: "silver"
            }

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
                            visible: false
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
                            visible: false
                        }
                        Controls1.CheckBox {
                            id: cbFillMsgNo
                            text: qsTr("FillMsgNo")
                            checked: true
                            visible: true
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
                            Controls1.ExclusiveGroup { id: common1Req_common2Req }
                            Controls1.RadioButton {
                                exclusiveGroup: common1Req_common2Req
                                id: rbC1Req
                                text: qsTr("C1Req")
                                checked: true
                                onClicked: {
                                    cbIsSafe.checked = false
                                    cbIsSafe.visible = !rbC1Req.checked
                                    cbIsUpCache.checked = false
                                    cbIsUpCache.visible = !rbC1Req.checked
                                    cbFillMsgNo.checked = true
                                    cbFillMsgNo.visible = rbC1Req.checked
                                }
                            }
                            Controls1.RadioButton {
                                exclusiveGroup: common1Req_common2Req
                                id: rbC2Req
                                text: qsTr("C2Req")
                                onClicked: {
                                    cbIsSafe.checked = false
                                    cbIsSafe.visible = !rbC1Req.checked
                                    cbIsUpCache.checked = false
                                    cbIsUpCache.visible = !rbC1Req.checked
                                    cbFillMsgNo.checked = true
                                    cbFillMsgNo.visible = rbC1Req.checked
                                }
                            }
                        }
                    }
                    Button {
                        text: qsTr("填充示例JSON")
                        onClicked: idTextArea.text = dataExch.jsonExample(idComboBox.currentText)
                    }
                    Button {
                        text: qsTr("发送")
                        onClicked: {
                            var message = dataExch.sendReq(idComboBox.currentText,idTextArea.text,peerid,cbIsLog.checked,cbIsSafe.checked,cbIsPush.checked,cbIsUpCache.checked,rbC1Req.checked,cbFillMsgNo.checked,cbForceToDB.checked)
                            ToolTip.show("send: "+message, 5000)
                        }
                    }
                }

                Controls1.ComboBox {
                    id: idComboBox
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

    Connections {
        target: dataExch
        onSigTableChanged: {
            if (tableName === "CommonData") {
                mstm.select()
                verScrollBar.setPosition(1.0)
            }
        }
    }

    Menu {
        id: menu
        property var userid: undefined
        property var msgno : undefined
        property var jsonText : undefined
        MenuItem {
            text: "复制"
            onTriggered: dataExch.copyText(jsonText)
            visible: true
        }
        MenuItem {
            text: "删除整个MsgNo"
            onTriggered: dataExch.deleteCommonData1(menu.userid,menu.msgno)
            visible: true
        }
        MenuItem {
            text: "TTS朗读json"
            onTriggered: dataExch.ttsSpeak(menu.jsonText)
            visible: true
        }
    }
}
