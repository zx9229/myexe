import QtQuick 2.0
import QtQuick.Controls 1.4 as Controls1
import QtQuick.Layouts 1.3

Item {
    ColumnLayout {
        anchors.fill: parent
        GridLayout {
            id: gridLayout
            columns: 2
            rows: 9
            Controls1.Button {
                Layout.fillWidth: true
                text: qsTr("载入")
                onClicked: {
                    cbTtsSubject.checked = (Number(dataExch.dbLoadValue(cbTtsSubject.dbKey))!==0)
                    cbTtsContent.checked = (Number(dataExch.dbLoadValue(cbTtsContent.dbKey))!==0)
                }
            }
            Controls1.Button {
                Layout.fillWidth: true
                text: qsTr("保存")
                onClicked: {
                    dataExch.dbSaveValue(cbTtsSubject.dbKey, cbTtsSubject.checked ? 1 : 0)
                    dataExch.dbSaveValue(cbTtsContent.dbKey, cbTtsContent.checked ? 1 : 0)
                }
            }
            Controls1.CheckBox {
                property string dbKey: "TtsSubject"
                id: cbTtsSubject
                Layout.fillWidth: true
                text: qsTr("TTS播报主题")
            }
            Controls1.CheckBox {
                property string dbKey: "TtsContent"
                id: cbTtsContent
                Layout.fillWidth: true
                text: qsTr("TTS播报内容")
            }
        }
        Rectangle {
            Layout.fillHeight: true
        }
    }
}
