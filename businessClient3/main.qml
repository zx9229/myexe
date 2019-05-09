import QtQuick 2.11
import QtQuick.Window 2.11

Window {
    visible: true
    width: 640
    height: 480
    title: qsTr("MyHelloWorld")
    color: "silver"

    Loader {
        id: pageLoader
        anchors.fill: parent
        source: "Login.qml"
    }

    Connections {
        target: dataExch
        onSigReady: pageLoader.source = "NodeTempInfo.qml"
    }
}
