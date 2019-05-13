import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.0
import MySqlTableModel 0.1

Item {
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
    }
}
