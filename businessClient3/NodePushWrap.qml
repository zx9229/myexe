import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4
import QtQuick.Controls.Styles 1.4
import MySqlTableModel 0.1

Item {
    property string peerid

    ColumnLayout {
        anchors.fill: parent

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
                        onPressAndHold: {}
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
            if (tableName === "PushWrap") {
                mstm.select()
                verScrollBar.setPosition(1.0)
            }
        }
    }
}
