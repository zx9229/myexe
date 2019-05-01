import QtQuick 2.4

LoginForm {
    id: loginForm
    Timer {
        interval: 500
        running: true
        repeat: true
        onTriggered: {
            var dts = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss.zzz ddd")
            loginForm.labelDT.text = dts
        }
    }
}
