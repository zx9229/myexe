```
curl -d '{ "Timeout":999 , "Recver":"n3" , "Data":"[echo]" , "Mode":1, "IsLog":true }'    http://127.0.0.1:40078/echo
```
在日志中，我们应当按照Key进行检索。日志中的Key一般如下所示：
```
msgData=Key:<UserID:"n4" MsgNo:8042223290400000001 >
msgData=Key:<UserID:"n4" MsgNo:8042223290400000001 SeqNo:1 >
msgData=Key:<UserID:"n4" MsgNo:8042223290400000001 SeqNo:2 >
```
分析日志的命令的例子和命令如下：  
ff  'Key:<UserID:"n4" MsgNo:8042223290400000001'  
ff(){ grep -r --include=*.log.INFO.* "$*" . | sort -t" " -k2 ; }
