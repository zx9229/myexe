import json
import urllib
import threading
import http.server
import sys
import getopt
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
        try:
            data = self.rfile.read(int(self.headers['content-length']))
            data = MyResquestHandler.bytes2str(data)
            data = json.loads(data)
            splitResult = urllib.parse.urlsplit(self.path)
            #
            if splitResult.path != "/send":
                raise Exception("unknown url")
            #
            if data.get("To") == "user":
                send_to_user(data["Msg"], data.get("Name"),
                             data.get("NickName"), data.get("RemarkName"))
            elif data.get("To") == "group":
                send_to_group(data["Msg"], data.get("Name"))
            else:
                raise Exception("unknown value of field To")
            #
            result = {'ErrNo': 0, 'ErrMsg': ''}
        except Exception as ex:
            result = {'ErrNo': -1, 'ErrMsg': ex.__str__()}
        self.wfile.write(json.dumps(result).encode())


# #############################################################################


def myLoginCallback():
    print(datetime.datetime.now(), "login")


def myExitCallback():
    print(datetime.datetime.now(), "exit")


@itchat.msg_register(itchat.content.TEXT, isFriendChat=True, isGroupChat=True)
def text_reply(msg):
    # print(msg)
    d = {
        "MsgId": msg.msgId,
        "IsAt": msg.get("IsAt"),
        "Text": msg.text,
        "Content": msg.content,
        "CreateTime": msg.createTime,
        "NickName": msg.user.nickName,
        "RemarkName": msg.user.remarkName,
        "PYQuanPin": msg.user.pYQuanPin,
        "UserName": msg.user.UserName
    }
    return str(d)


def send_to_user(msg, name=None, nName=None, rName=None):
    '''
    https://itchat.readthedocs.io/zh/latest/intro/contact/
    '''
    users = itchat.search_friends(name=name, nickName=nName, remarkName=rName)
    if not users:
        raise Exception("user not found")
    user = users[0]
    user.send(msg)


def send_to_group(msg, name=None):
    '''
    itchat 转发指定的微信群、用户的发言到指定的群
    https://blog.csdn.net/zhizunyu2009/article/details/79126375
    该程序的主要功能是监控撤回消息，并且如果有消息撤回就会撤回的消息发送给你，以后再也不用担心看不到好友的撤回的消息了
    https://www.cnblogs.com/ouyangping/p/8453920.html
    '''
    groups = itchat.search_chatrooms(name=name)
    if not groups:
        raise Exception("group not found")
    group = groups[0]
    group.send(msg)


# #############################################################################


def calc_args():
    host = 'localhost'
    port = 65535
    opts, args = getopt.getopt(sys.argv[1:], "h:p:")
    for op, value in opts:
        if op == "-h":
            host = value
        if op == "-p":
            port = int(value)
    return host, port


# curl -d '{"To":"group","Msg":"测试消息","Name":"群名"}'  http://localhost:65535/send
# curl -d '{"To":"user" ,"Msg":"测试消息","Name":"模糊搜索的名字","NickName":"昵称","RemarkName":"备注姓名"}' http://localhost:65535/send
if __name__ == "__main__":
    host, port = calc_args()
    myServer = MyHttpServer(host, port)
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
