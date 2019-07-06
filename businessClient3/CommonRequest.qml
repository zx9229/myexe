import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 1.4 as Controls1
import QtQuick.Controls 2.4

Item {
    function funLoaded() {
        tfUserID.autoFillData()
        tfMsgNo.autoFillData()
        tfSenderID.autoFillData()
        cbRecverID.autoFillData()
        cbReqType.autoFillData()
    }

    Controls1.ScrollView {
        anchors.fill: parent
        GridLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            id: gridLayout
            rows: 17
            columns: 3
            Label {
                text: qsTr("C1C2")
            }
            GroupBox {
                Layout.fillWidth: true
                Layout.columnSpan: 2
                RowLayout {
                    Controls1.ExclusiveGroup { id: common1Req_common2Req }
                    Controls1.RadioButton {
                        exclusiveGroup: common1Req_common2Req
                        id: rbC1Req
                        text: qsTr("Common11Req")
                        checked: true
                        onClicked: {}
                    }
                    Controls1.RadioButton {
                        exclusiveGroup: common1Req_common2Req
                        id: rbC2Req
                        text: qsTr("Common12Req")
                        onClicked: {}
                    }
                }
            }
            Label {
                text: "UserID"
            }
            TextField {
                id: tfUserID
                text: ""
                Layout.fillWidth: true
                readOnly: tfUserID_cb.checked
                function autoFillData() {
                    tfUserID.text = dataExch.memGetInfo("myself", ["UserID"])
                }
            }
            Controls1.CheckBox {
                id: tfUserID_cb
                text: "AutoFill"
                checked: true
                onClicked: {
                    tfUserID.readOnly = tfUserID_cb.checked
                    tfUserID.autoFillData()
                }
            }
            Label {
                text: "MsgNo"
            }
            TextField {
                id: tfMsgNo
                text: ""
                Layout.fillWidth: true
                readOnly: tfMsgNo_cb.checked
                function autoFillData() {
                    var curMsgNo = Number(dataExch.memGetData("MsgNo"))
                    tfMsgNo.text = "%1=%2+1".arg(curMsgNo+1).arg(curMsgNo)
                }
            }
            Controls1.CheckBox {
                id: tfMsgNo_cb
                text: "AutoFill"
                checked: true
                onClicked: {
                    tfMsgNo.readOnly = tfMsgNo_cb.checked
                    tfMsgNo.autoFillData()
                }
            }
            Label {
                text: "SeqNo"
            }
            TextField {
                id: tfSeqNo
                text: "0"
                Layout.fillWidth: true
                readOnly: tfSeqNo_cb.checked
            }
            Controls1.CheckBox {
                id: tfSeqNo_cb
                text: "AutoFill"
                checked: true
                onClicked: {
                    tfSeqNo.readOnly = tfSeqNo_cb.checked
                    tfSeqNo.text = "%1".arg(0)
                }
            }
            Label {
                text: "BatchNo"
            }
            TextField {
                id:tfBatchNo
                text: "0"
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "RefNum"
            }
            TextField {
                id: tfRefNum
                text: "0"
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "RefText"
            }
            TextField {
                id: tfRefText
                text: ""
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "SenderID"
            }
            TextField {
                id: tfSenderID
                text: ""
                Layout.fillWidth: true
                readOnly: tfSenderID_cb.checked
                function autoFillData() {
                    tfSenderID.text = dataExch.memGetInfo("myself", ["UserID"])
                }
            }
            Controls1.CheckBox {
                id: tfSenderID_cb
                text: "AutoFill"
                checked: true
                onClicked: {
                    tfSenderID.readOnly = tfSenderID_cb.checked
                    tfSenderID.autoFillData()
                }
            }
            Label {
                text: "RecverID"
            }
            Controls1.ComboBox {
                id: cbRecverID
                property string text2: cbRecverID.editable ? cbRecverID.editText : cbRecverID.currentText
                Layout.fillWidth: true
                editable: false
                model: []
                function autoFillData() {
                    var data = JSON.parse(dataExch.memGetData("PathwayInfo"))["Info"]
                    cbRecverID.model = (typeof(data) === "undefined") ? [] : Object.keys(data)
                }
            }
            Controls1.CheckBox {
                id: cbRecverID_cb
                text: "editable"
                checked: false
                onClicked: {
                    cbRecverID.editable = cbRecverID_cb.checked
                    cbRecverID.autoFillData()
                }
            }
            Label {
                text: "ToRoot"
            }
            Controls1.CheckBox {
                id: cbToRoot
                text: "(叶子节点=>子节点=>根节点)"
                Layout.fillWidth: true
                checked: true
                onClicked: cbToRoot.checked = cbToRoot_cb.checked || cbToRoot.checked
            }
            Controls1.CheckBox {
                id: cbToRoot_cb
                text: "AutoFill"
                checked: true
                onClicked: if(cbToRoot_cb.checked){ cbToRoot.checked = true }
            }
            Label {
                text: "IsLog"
            }
            Controls1.CheckBox {
                id: cbIsLog
                text: "IsLog"
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "IsSafe"
            }
            Controls1.CheckBox {
                id: cbIsSafe
                text: "IsSafe"
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "IsPush"
            }
            Controls1.CheckBox {
                id: cbIsPush
                text: "IsPush"
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "UpCache"
            }
            Controls1.CheckBox {
                id: cbUpCache
                text: "UpCache"
                Layout.fillWidth: true
                Layout.columnSpan: 2
            }
            Label {
                text: "ReqTime"
            }
            TextField {
                id: tfReqTime
                property string text2: tfReqTime.readOnly ? calcText(false) : tfReqTime.text
                text: calcText(true)
                Layout.fillWidth: true
                readOnly: tfReqTime_cb.checked
                function calcText(isFormatter) {
                    return isFormatter ? "yyyy-MM-dd hh:mm:ss" : (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss")
                }
            }
            Controls1.CheckBox {
                id: tfReqTime_cb
                text: "AutoFill"
                checked: true
                onClicked: {
                    tfReqTime.readOnly = tfReqTime_cb.checked
                    tfReqTime.text = tfReqTime.calcText(tfReqTime.readOnly)
                }
            }
            Label {
                text: "ReqType"
            }
            Controls1.ComboBox {
                id: cbReqType
                property string text2: cbReqType.editable ? cbReqType.editText : cbReqType.currentText
                Layout.fillWidth: true
                editable: false
                model: []
                function autoFillData() {
                    cbReqType.model = dataExch.getTxMsgTypeNameList()
                }
            }
            Controls1.CheckBox {
                id: cbReqType_cb
                text: "editable"
                onClicked: cbReqType.editable = cbReqType_cb.checked
            }
            TextArea {
                id: idTextArea
                Layout.fillWidth: true
                Layout.fillHeight: true
                Layout.columnSpan: 3
                background: Rectangle {
                    border.width: 2
                    border.color: "blue"
                }
            }
            RowLayout {
                Layout.fillWidth: true
                Layout.columnSpan: 3
                Button {
                    text: qsTr("填充JSON")
                    Layout.fillWidth: true
                    onClicked: idTextArea.text = dataExch.jsonExample(cbReqType.text2)
                }
                Button {
                    text: qsTr("发送")
                    Layout.fillWidth: true
                    onClicked: {
                        var paramList =  selectCommonReq()
                        var message = dataExch.testSendC1Req(paramList)
                        ToolTip.show("send: "+message, 5000)
                    }
                }
            }
            Button {
                Layout.fillWidth: true
                Layout.columnSpan: 3
            }
        }
        function selectCommonReq() {
            var paramList = []
            paramList.push('UserID'); paramList.push(tfUserID.text);
            paramList.push('MsgNo'); paramList.push(tfMsgNo.text.split("=")[0]);
            paramList.push('SeqNo'); paramList.push(tfSeqNo.text);
            paramList.push('BatchNo'); paramList.push(tfBatchNo.text);
            paramList.push('RefNum'); paramList.push(tfRefNum.text);
            paramList.push('RefText'); paramList.push(tfRefText.text);
            paramList.push('SenderID'); paramList.push(tfSenderID.text);
            paramList.push('RecverID'); paramList.push(cbRecverID.text2);
            paramList.push('ToRoot'); paramList.push(cbToRoot.checked?1:0);
            paramList.push('IsLog'); paramList.push(cbIsLog.checked?1:0);
            paramList.push('IsPush'); paramList.push(cbIsPush.checked?1:0);
            paramList.push('UpCache'); paramList.push(cbUpCache.checked?1:0);
            paramList.push('ReqTime'); paramList.push(tfReqTime.text2);
            paramList.push('ReqType'); paramList.push(cbReqType.text2);
            paramList.push('ReqData'); paramList.push(idTextArea.text);
            return paramList
        }
    }
}
