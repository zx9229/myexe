import QtQuick 2.0
import QtQuick.Controls 2.4
import QtQuick.Layouts 1.3

Item {
    ColumnLayout {
        anchors.fill: parent
        TextField {
            id: roStateTransition
            readOnly: true
            Layout.fillWidth: true
            text: "状态变迁"
        }
        TextField {
            id: roState
            readOnly: true
            Layout.fillWidth: true
            text: "当前状态"
        }
        Button {
            Layout.fillWidth: true
            text: "查询当前状态"
            onClicked: roState.text = dataExch.remoteObjectState()
        }
        Button {
            Layout.fillWidth: true
            text: "启动服务"
            onClicked: dataExch.startService()
        }
        Rectangle {
            Layout.fillHeight: true
        }
    }
    Connections {
        target: dataExch
        onSigStateChanged: {
            roStateTransition.text = "%1(%2) => %3(%4)".arg(sOldState).arg(iOldState).arg(sState).arg(iState)
        }
    }
}
