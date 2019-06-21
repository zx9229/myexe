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
            text: qsTr("状态变迁")
        }
        TextField {
            id: roState
            readOnly: true
            Layout.fillWidth: true
            text: qsTr("QtRemoteObjects状态")
        }
        TextField {
            id: serviceState
            readOnly: true
            Layout.fillWidth: true
            text: qsTr("服务状态")
        }
        Button {
            Layout.fillWidth: true
            text: qsTr("启动服务")
            onClicked: dataExch.startService()
        }
        Button {
            Layout.fillWidth: true
            text: qsTr("查询QtRemoteObjects状态")
            onClicked: {
                var localDT = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss")
                roState.text = "%1  [%2]".arg(dataExch.remoteObjectState()).arg(localDT)
            }
        }
        Button {
            Layout.fillWidth: true
            text: qsTr("查询service状态")
            onClicked: {
                var localDT = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss")
                serviceState.text = "%1  [%2]".arg(dataExch.serviceState()).arg(localDT)
            }
        }
        Rectangle {
            Layout.fillHeight: true
        }
    }
    Connections {
        target: dataExch
        onSigStateChanged: {
            var localDT = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss")
            roStateTransition.text = "%1(%2) => %3(%4)  [%5]".arg(sOldState).arg(iOldState).arg(sState).arg(iState).arg(localDT)
        }
    }
}
