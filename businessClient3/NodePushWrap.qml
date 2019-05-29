import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4
import QtQuick.Controls.Styles 1.4
import MySqlTableModel 0.1

Item {
    signal sigShowNodeList()
    property string peerid
    ColumnLayout {
        anchors.fill: parent

        RowLayout {
            Button {
                text: qsTr("<[返回]")
                onClicked: sigShowNodeList()
            }
            Button {
                text: qsTr("刷新NodePushWrap")
            }
        }

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM PushWrap WHERE PeerID='%1'".arg(peerid)
            }
            delegate: Column {
                id: idColumn
                width: listView.width
                Label {
                    id: txt1
                    text: "%1 %2 %3".arg(UserID).arg(PshTime).arg(PshTypeTxt)
                }
                Rectangle {
                    id: idRect
                    radius: 5
                    color: "lightgray"
                    height: messageText.implicitHeight + 24
                    width: Math.min(messageText.implicitWidth + 24,listView.width)
                    TextEdit {
                        id: messageText
                        text: PshDataTxt
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
                        }
                        onDoubleClicked: {
                            messageText.selectByMouse = !messageText.selectByMouse
                        }
                        onPressAndHold: {
                            //sigShowNodeReqRsp(UserID, MsgNo)
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
        }
    }
}
