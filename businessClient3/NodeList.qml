import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.0
import MySqlTableModel 0.1

Item {
    signal sigPickPeerID(string peerid)

    ColumnLayout {
        anchors.fill: parent

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true

            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM ConnInfoEx"
            }

            delegate: Rectangle {
                property string testid: UserID
                id: idRect
                radius: 5
                color: "lightgray"
                border.color: "gray"
                border.width: 1
                width: listView.width
                height: txt1.implicitHeight + txt2.implicitHeight + 20
                Column {
                    anchors.verticalCenter: parent.verticalCenter
                    Label {
                        id: txt1
                        text: "UserID:" + UserID + ", BelongID:" + BelongID
                    }
                    Label {
                        id: txt2
                        text: "Pathway:" + Pathway
                    }
                }

                states: State {
                    when: idRect.ListView.isCurrentItem
                    PropertyChanges {
                        target: idRect
                        color: "tan"
                    }
                }

                MouseArea {
                    anchors.fill: parent
                    onClicked: {
                        idRect.ListView.view.currentIndex = index //https://blog.csdn.net/x356982611/article/details/53008236
                        sigPickPeerID(UserID)
                    }
                    onPressAndHold: {}
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
            if (tableName === "ConnInfoEx") {
                mstm.select()
            }
        }
    }
}

/*##^## Designer {
    D{i:0;autoSize:true;height:480;width:640}
}
 ##^##*/
