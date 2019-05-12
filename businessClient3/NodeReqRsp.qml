import QtQuick 2.0
import QtQuick.Controls 2.4 //(Page QML Type)
import QtQuick.Layouts 1.3
import MySqlTableModel 0.1

//Item {
Page {
    property alias statement: mstm.selectStatement

    ColumnLayout {
        anchors.fill: parent

        Button {
            text: qsTr("刷新NodeReqRsp")
        }

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id: mstm
                selectStatement: ""
            }
            delegate: Column {
                property bool isReq: (3==MsgType)||(5==MsgType)
                spacing: 6
                anchors.right: isReq ? parent.right : undefined
                Label {
                    id: txt1
                    text: "%1 %2 %3".arg(InsertTime).arg(IsLast).arg(RspCnt)
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
                        Label {
                            text: MsgType==3?"C2REQ":(MsgType==4?"C2RSP":(MsgType==5?"C1REQ":MsgType==6?"C1RSP":"NULL"))
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
                        Label {
                            text: MsgType==3?"C2REQ":(MsgType==4?"C2RSP":(MsgType==5?"C1REQ":MsgType==6?"C1RSP":"NULL"))
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
