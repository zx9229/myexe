import QtQuick 2.11
import QtQuick.Window 2.11
import QtQuick.Controls 2.4

Window {
    visible: true
    width: 324
    height: 702
    title: qsTr("MyHelloWorld")
    color: "silver"

    SwipeView {
        id: swipeView
        currentIndex: 0
        anchors.fill: parent

        Loader {
            id: loader0
            source: "qrc:/Login.qml"
            onLoaded: {
                if (false) {
                } else if (source == "qrc:/Login.qml") {
                    item.sigShowHomePage.connect(function(){
                        loader0.source = "qrc:/HomePage.qml"
                    })
                } else if (source == "qrc:/HomePage.qml") {
                    item.sigShowNodeList.connect(function(){
                        var urlStr = "qrc:/NodeList.qml"
                        if (loader1.source == urlStr) { swipeView.setCurrentIndex(1) } else {
                            loader1.source = urlStr
                        }
                    })
                    item.sigShowPathwayInfo.connect(function(){
                        var urlStr = "qrc:/PathwayInfo.qml"
                        if (loader1.source == urlStr) { swipeView.setCurrentIndex(1) } else {
                            loader1.source = urlStr
                        }
                    })
                    item.sigShowNodePushWrap.connect(function(PeerID){
                        var urlStr = "qrc:/NodePushWrap.qml"
                        if (loader1.source == urlStr) { swipeView.setCurrentIndex(1) } else {
                            var attrMap = {"peerid":PeerID}
                            loader1.setSource(urlStr, attrMap)
                        }
                    })
                    item.sigShowNodeRequest.connect(function(PeerID){
                        var urlStr = "qrc:/NodeRequest.qml"
                        if (loader1.source == urlStr) { swipeView.setCurrentIndex(1) } else {
                            var attrMap = {"peerid":PeerID}
                            loader1.setSource(urlStr, attrMap)
                        }
                    })
                    item.sigShowNodeReqRsp.connect(function(UserID, MsgNo){
                        var urlStr = "qrc:/NodeReqRsp.qml"
                        if (loader1.source == urlStr) { swipeView.setCurrentIndex(1) } else {
                            var attrMap = {"userid":UserID, "msgno":MsgNo}
                            loader1.setSource(urlStr, attrMap)
                        }
                    })
                }
            }
        }

        Loader {
            id: loader1
            onLoaded: {
                swipeView.setCurrentIndex(1)
                if (item.hasOwnProperty('sigPickPeerID')) {
                    item.sigPickPeerID.connect(function(peerid){
                        //当指定peerID的时候,说明msgNo已经失效了.
                        loader0.item.peerID = peerid
                        loader0.item.msgNo  = undefined
                    })
                }
                if (item.hasOwnProperty('sigPickMsgNo')) {
                    item.sigPickMsgNo.connect(function(msgno){
                        //当指定msgNo的时候,必定已经指定了peerID.
                        loader0.item.msgNo = msgno
                    })
                }
            }
        }
    }
}
