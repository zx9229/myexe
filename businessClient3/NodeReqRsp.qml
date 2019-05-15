import QtQuick 2.0
import QtQuick.Controls 2.4 //(Page QML Type)
import QtQuick.Layouts 1.3
import MySqlTableModel 0.1

//Item {
Page {
    signal sigShowNodeRequest(string UserID)
    //property alias statement: mstm.selectStatement
    property string peerid
    property string userid
    property string msgno

    ColumnLayout {
        anchors.fill: parent

        RowLayout {
            Button {
                text: qsTr("<[返回]")
                onClicked: sigShowNodeRequest(peerid)
            }

            Button {
                text: qsTr("刷新NodeReqRsp")
            }
        }

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM CommonData WHERE UserID='%1' AND MsgNo='%2'".arg(userid).arg(msgno)
            }
            delegate: Column {
                property bool isReq: (3==MsgType)||(5==MsgType)
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
                        height: 32
                        width: height
                        visible: !isReq
                        border.color: "gray"
                        border.width: 1
                        Label {//为了缩减字母,选用(Q&A)代表请求和响应.
                            anchors.centerIn: parent
                            text: MsgType==3?"C2Q":(MsgType==4?"C2A":(MsgType==5?"C1Q":MsgType==6?"C1A":"NIL"))
                        }
                    }
                    Rectangle {
                        height: messageText.implicitHeight + 24
                        width: Math.min(messageText.implicitWidth + 24, listView.width - (isReq ? avatarRight.width : avatarLeft.width) - messageRow.spacing)
                        color: isReq ? "lightgrey":"steelblue"
                        Label {
                            id: messageText
                            text: TxDataTxt
                            color: isReq ? "black":"white"
                            anchors.fill: parent
                            anchors.margins: 12
                            wrapMode: Label.Wrap
                        }
                    }
                    Rectangle {
                        id: avatarRight
                        height: 32
                        width: height
                        visible: isReq
                        border.color: "gray"
                        border.width: 1
                        Label {//为了缩减字母,选用(Q&A)代表请求和响应.
                            anchors.centerIn: parent
                            text: MsgType==3?"C2Q":(MsgType==4?"C2A":(MsgType==5?"C1Q":MsgType==6?"C1A":"NIL"))
                        }
                    }
                }
            }
            ScrollBar.vertical: ScrollBar {
                id: verScrollBar
            }
        }
    }
}
