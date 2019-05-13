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
        source: "qrc:/Login.qml"
        onLoaded: {
            if (false) {
            } else if (source == "qrc:/Login.qml") {
                item.sigSignIn.connect(function(){
                    pageLoader.source = "qrc:/NodeList.qml"
                })
            } else if (source == "qrc:/NodeList.qml") {
                item.sigShowNode.connect(function(PeerID){
                    var sqlStatement = "SELECT * FROM CommonData WHERE MsgType IN(3,5) AND PeerID='%1'".arg(PeerID)
                    pageLoader.setSource("qrc:/NodeRequest.qml", {"statement":sqlStatement})
                })
            } else if (source == "qrc:/NodeRequest.qml") {
                item.sigShowReqRsp.connect(function(UserID, MsgNo){
                    var sqlStatement = "SELECT * FROM CommonData WHERE UserID='%1' AND MsgNo='%2'".arg(UserID).arg(MsgNo.toString())
                    pageLoader.setSource("qrc:/NodeReqRsp.qml", {"statement":sqlStatement})
                })
            }
        }
    }
}
