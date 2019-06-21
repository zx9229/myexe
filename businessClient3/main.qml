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
            source: "qrc:/HomePage.qml"
            onLoaded: {
                item.sigShowServiceStatus.connect(function(){
                    var urlStr = "qrc:/ServiceStatus.qml"
                    if (loader1.source == urlStr) { swipeView.setCurrentIndex(1) } else {
                        loader1.source = urlStr
                    }
                })
                item.sigShowLogin.connect(function(){
                    var urlStr = "qrc:/Login.qml"
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

        Loader {
            id: loader1
            source: "qrc:/ServiceStatus.qml"
            onLoaded: {
                swipeView.setCurrentIndex(1)
                if (item.hasOwnProperty('sigPickUserID')) {
                    item.sigPickUserID.connect(function(userid){
                        loader0.item.userID = userid
                    })
                }
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
                if (item.hasOwnProperty('sigShowHomePage')) {
                    item.sigShowHomePage.connect(function(){
                        swipeView.setCurrentIndex(0)
                    })
                }
            }
        }
    }
}
