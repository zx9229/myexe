import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4

ColumnLayout {
    signal sigClicked()
    property string txtText
    property string btnText
    TextEdit {
        readOnly: true
        Layout.fillWidth: true
        text: txtText
    }
    Button {
        onClicked: sigClicked()
        Layout.fillWidth: true
        text: btnText
    }
}
