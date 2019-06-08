import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4

Item {
    signal sigShowNodeList()
    signal sigShowPathwayInfo()
    signal sigShowNodePushWrap(string PeerID)
    signal sigShowNodeRequest(string PeerID)
    signal sigShowNodeReqRsp(string UserID, string MsgNo)
    signal sigShowSystemSettings()
    //The var type is a generic property type that can refer to any data type.
    property var userID: dataExch.getUserID()
    property var peerID: undefined
    property var msgNo : undefined

    ColumnLayout {
        anchors.fill: parent
        GridLayout {
            columns: 2
            rows: 3
            HomePageComponent {
                txtText: qsTr("\n")
                btnText: qsTr("节点列表")
                onSigClicked: sigShowNodeList()
            }
            HomePageComponent {
                txtText: '本端节点:[%1]\n对端节点:[%2]'.arg(userID).arg(peerID)
                btnText: qsTr("节点推送消息")
                onSigClicked: {
                    if (typeof(peerID) === "undefined") {
                        ToolTip.show("请先指定对端节点", 5000)
                    } else {
                        sigShowNodePushWrap(peerID)
                    }
                }
            }
            HomePageComponent {
                txtText: '本端节点:[%1]\n对端节点:[%2]\n'.arg(userID).arg(peerID)
                btnText: qsTr("节点请求")
                onSigClicked: {
                    if (typeof(peerID) === "undefined") {
                        ToolTip.show("请先指定对端节点", 5000)
                    } else {
                        sigShowNodeRequest(peerID)
                    }
                }
            }
            HomePageComponent {
                txtText: '本端节点:[%1]\n对端节点:[%2]\n消息编号:[%3]'.arg(userID).arg(peerID).arg(msgNo)
                btnText: qsTr("节点请求响应")
                onSigClicked: {
                    if (false) {
                    } else if (typeof(peerID) === "undefined") {
                        ToolTip.show("请先指定对端节点", 5000)
                    } else if (typeof(msgNo)  === "undefined") {
                        ToolTip.show("请先指定消息编号", 5000)
                    } else {
                        sigShowNodeReqRsp(userID, msgNo)
                    }
                }
            }
            HomePageComponent{
                txtText: qsTr("")
                btnText: qsTr("系统设置")
            }
            HomePageComponent{
                txtText: qsTr("")
                btnText: qsTr("路径信息")
                onSigClicked: sigShowPathwayInfo()
            }
        }
        Rectangle {
            Layout.fillHeight: true
        }
    }
}

/*##^## Designer {
    D{i:0;autoSize:true;height:480;width:640}
}
 ##^##*/
