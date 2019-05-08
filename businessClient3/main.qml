import QtQuick 2.4
import QtQuick.Window 2.10

Window {
    visible: true
    width: 800
    height: 400
    title: qsTr("MyHelloWorld")
    color: "silver"
    Loader{
        id:pageLoader
        anchors.fill: parent
        source: "Login.qml"
    }
    Connections {
        target: dataExch
        onSigReady: pageLoader.source="NodeTempInfo.qml"
    }
}
