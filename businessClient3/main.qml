import QtQuick 2.11
import QtQuick.Window 2.11

Window {
    visible: true
    width: 324
    height: 702
    title: qsTr("MyHelloWorld")
    color: "lightgray"

    Loader {
        id: pageLoader
        anchors.fill: parent
        source: "qrc:/Login.qml"
        onLoaded: {
            if (false) {
            } else if (source == "qrc:/Login.qml") {
                item.sigShowNodeList.connect(function(){
                    pageLoader.source = "qrc:/NodeList.qml"
                })
            } else if (source == "qrc:/NodeList.qml") {
                item.sigShowNodeRequest.connect(function(PeerID){
                    pageLoader.setSource("qrc:/NodeRequest.qml", {"peerid":PeerID})
                })
            } else if (source == "qrc:/NodeRequest.qml") {
                item.sigShowNodeList.connect(function(){
                    pageLoader.source = "qrc:/NodeList.qml"
                })
                item.sigShowNodeReqRsp.connect(function(UserID, MsgNo){
                    pageLoader.setSource("qrc:/NodeReqRsp.qml", {"peerid":item.peerid, "userid":UserID, "msgno":MsgNo.toString()})
                })
            } else if (source == "qrc:/NodeReqRsp.qml") {
                item.sigShowNodeRequest.connect(function(PeerID){
                    pageLoader.setSource("qrc:/NodeRequest.qml", {"peerid":PeerID})
                })
            }
        }
    }
}
