import QtQuick 2.4

LoginForm {
    id: loginForm

    signal sigSignIn() //声明一个自定义的信号, https://www.jianshu.com/p/442f461ee62b

    Timer {
        interval: 500
        running: true
        repeat: true
        onTriggered: {
            var dts = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss.zzz ddd")
            loginForm.labelDT.text = dts
        }
    }

    Connections {
        //因为【void sigStatusError(const QString& errMessage, int errType);】所以可以如下书写:
        target: dataExch
        onSigStatusError: {
            var localDT = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss")
            loginForm.textAreaMessage.text = localDT + '\n' + errMessage
        }
    }

    Connections {
        //因为【void sigReady();】所以可以如下书写:
        target: dataExch
        onSigReady: {
            var localDT = (new Date()).toLocaleString(Qt.locale(), "yyyy-MM-dd hh:mm:ss")
            loginForm.textAreaMessage.text = localDT + '\n' + "SUCCESS"
            sigSignIn()
        }
    }
}
