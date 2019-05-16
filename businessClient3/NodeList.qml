import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.0
import MySqlTableModel 0.1

Item {
    signal sigShowNodeRequest(string UserID)

    ColumnLayout {
        anchors.fill: parent

        Button {
            text: qsTr("刷新NodeList")
            onClicked: mstm.select()
        }

        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true

            model: MySqlTableModel {
                id: mstm
                selectStatement: "SELECT * FROM ConnInfoEx"
            }

            delegate: Rectangle {
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
                    }
                    onPressAndHold: {
                        sigShowNodeRequest(UserID)
                    }
                }
            }

            ScrollBar.vertical: ScrollBar {}
        }
    }
}

/*##^## Designer {
    D{i:0;autoSize:true;height:480;width:640}
}
 ##^##*/
