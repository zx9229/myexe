import json
import urllib
import threading
import http.server
import datetime
import itchat


# #############################################################################
class MyHttpServer(object):
    def __init__(self, host, port):
        self._server = http.server.HTTPServer((host, port), MyResquestHandler)

    def serverBackground(self):
        print("listen {}".format(self._server.server_address))
        self._serverThread = threading.Thread(
            target=self._server.serve_forever)
        self._serverThread.setDaemon(True)
        self._serverThread.start()


class MyResquestHandler(http.server.BaseHTTPRequestHandler):
    @staticmethod
    def bytes2str(src):
        dst = None
        for encoding in ('utf_8', 'gbk'):
            try:
                dst = src.decode(encoding='gbk', errors='strict')
                break
            except Exception as ex:
                continue
        return str(src) if dst is None else dst

    def do_GET(self):
        return

    def do_POST(self):
        data = self.rfile.read(int(self.headers['content-length']))
        data = MyResquestHandler.bytes2str(data)
        data = json.loads(data)
        splitResult = urllib.parse.urlsplit(self.path)
        if splitResult.path == "/sendToGroup":
            result = self.sendToGroup(data)
        elif splitResult.path == "/sendToUser":
            result = self.sendToUser(data)
        else:
            result = {'ErrNo': -1, 'ErrMsg': 'unknown url'}
        self.wfile.write(json.dumps(result).encode())

    def sendToGroup(self, data):
        try:
            send_to_group(data["Msg"], data["Name"])
            result = {'ErrNo': 0, 'ErrMsg': ''}
        except Exception as ex:
            result = {'ErrNo': -1, 'ErrMsg': ex.__str__()}
        return result

    def sendToUser(self, data):
        try:
            send_to_user(data["Msg"], data["Name"])
            result = {'ErrNo': 0, 'ErrMsg': ''}
        except Exception as ex:
            result = {'ErrNo': -1, 'ErrMsg': ex.__str__()}
        return result


# #############################################################################


def myLoginCallback():
    print(datetime.datetime.now(), "login")


def myExitCallback():
    print(datetime.datetime.now(), "exit")


@itchat.msg_register(itchat.content.TEXT, isFriendChat=True, isGroupChat=True)
def text_reply(msg):
    print(msg)
    d = {
        "Text": msg.text,
        "CreateTime": msg.createTime,
        "IsAt": msg.isAt,
        "MsgId": msg.msgId,
        "NickName": msg.user.nickName,
        "PYQuanPin": msg.user.pYQuanPin
    }
    return str(d)


def send_to_user(message, name):
    raise NotImplementedError()


def send_to_group(message, name):
    '''
    itchat 转发指定的微信群、用户的发言到指定的群
    https://blog.csdn.net/zhizunyu2009/article/details/79126375
    该程序的主要功能是监控撤回消息，并且如果有消息撤回就会撤回的消息发送给你，以后再也不用担心看不到好友的撤回的消息了
    https://www.cnblogs.com/ouyangping/p/8453920.html
    '''
    toGroup = None
    itchat.search_chatrooms()
    groups = itchat.get_chatrooms(update=True)
    for group in groups:
        if group.nickName == name:
            toGroup = group.userName
    if toGroup is None:
        print("is_none")
    else:
        itchat.send(msg=message, toUserName=toGroup)


if __name__ == "__main__":
    # curl -d '{"Msg":"测试消息","Name":"ZX自用"}' http://127.0.0.1:8080/sendToGroup
    # curl -d '{"Msg":"测试消息","Name":"ZX自用"}' http://127.0.0.1:8080/sendToUser
    myServer = MyHttpServer('localhost', 8080)
    myServer.serverBackground()
    #
    enableCmdQR = -1
    enableCmdQR = 2
    enableCmdQR = True
    itchat.auto_login(
        enableCmdQR=enableCmdQR,
        hotReload=True,
        loginCallback=myLoginCallback,
        exitCallback=myExitCallback)
    itchat.run()
