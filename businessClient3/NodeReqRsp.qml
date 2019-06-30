import QtQuick 2.0
import QtQuick.Controls 2.4
import QtQuick.Layouts 1.3
import MySqlTableModel 0.1

Item {
    property string userid
    property string msgno

    ColumnLayout {
        anchors.fill: parent

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM CommonData WHERE UserID='%1' AND MsgNo='%2'".arg(userid).arg(msgno)
                //selectStatement: "SELECT * FROM CommonData WHERE UserID='%1' AND MsgNo='%2' AND json_extract(TxDataTxt,'$.SecGap')>0".arg(userid).arg(msgno)
            }
            delegate: Column {
                id: idColumn
                property bool isReq: (1==MsgType)||(3==MsgType)
                spacing: 6
                anchors.right: isReq ? parent.right : undefined
                Label {
                    id: txt1
                    text: "%1 (%2) %3".arg(InsertTime).arg(IsLast).arg(RspCnt)
                    color: "gray"
                    anchors.right: isReq ? parent.right : undefined
                }
                Row {
                    id: messageRow
                    spacing: 6
                    anchors.right: isReq ? parent.right : undefined
                    Rectangle {
                        id: avatarLeft
                        width: labelAL.implicitWidth + 4
                        height: width
                        visible: !isReq
                        border.color: "gray"
                        border.width: 1
                        Label {//为了缩减字母,选用(Q&A)代表请求和响应.
                            id: labelAL
                            anchors.centerIn: parent
                            text: MsgType==1?"C1Q":(MsgType==2?"C1A":(MsgType==3?"C2Q":MsgType==4?"C2A":"NIL"))
                        }
                    }
                    Rectangle {
                        id: idRect
                        radius: 5
                        color: "lightgray"
                        height: messageText.implicitHeight + 24
                        width: Math.min(messageText.implicitWidth + 24, listView.width - (isReq ? avatarRight.width : avatarLeft.width) - messageRow.spacing)
                        TextEdit {
                            id: messageText
                            text: TxDataTxt
                            anchors.fill: parent
                            anchors.margins: 12
                            wrapMode: TextEdit.Wrap
                            readOnly: true
                            selectByMouse: false
                        }
                        MouseArea {
                            anchors.fill: parent
                            onClicked: {
                                idColumn.ListView.view.currentIndex = index
                            }
                            onDoubleClicked: {
                                messageText.selectByMouse = !messageText.selectByMouse
                            }
                            onPressAndHold: {
                                menu.userid = UserID
                                menu.msgno = MsgNo
                                menu.seqno = SeqNo
                                menu.jsonText = TxDataTxt
                                menu.popup()
                            }
                        }
                    }
                    Rectangle {
                        id: avatarRight
                        width: labelAR.implicitWidth + 4
                        height: width
                        visible: isReq
                        border.color: "gray"
                        border.width: 1
                        Label {//为了缩减字母,选用(Q&A)代表请求和响应.
                            id: labelAR
                            anchors.centerIn: parent
                            text: MsgType==1?"C1Q":(MsgType==2?"C1A":(MsgType==3?"C2Q":MsgType==4?"C2A":"NIL"))
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
        property var seqno: undefined
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
            text: "删除单个SeqNo"
            onTriggered: dataExch.deleteCommonData2(menu.userid,menu.msgno,menu.seqno)
            visible: true
        }
        MenuItem {
            text: "TTS朗读json"
            onTriggered: dataExch.ttsSpeak(menu.jsonText)
            visible: true
        }
        MenuItem {
            text: "TTS朗读json.Subject"
            onTriggered: {
                var jsonObj = JSON.parse(menu.jsonText);
                dataExch.ttsSpeak(jsonObj.Subject)
            }
        }
        MenuItem {
            text: "TTS朗读json.Content"
            onTriggered: {
                var jsonObj = JSON.parse(menu.jsonText);
                dataExch.ttsSpeak(jsonObj.Content)
            }
        }
    }
}
