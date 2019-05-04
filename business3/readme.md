1. 绝对安全模式: 同步写数据库,写成功了,才进行下一步,随时可以杀进程.
2. 安全退出模式: 东西都在内存中,退出时,需要让程序自己退出,不可以杀进程,不可以宕机.
插入数据库可能的情况：
1. 存在该主键，插入失败。
2. 没有该主键，插入成功。
3. 没有该主键，插入失败（硬盘满等原因）。
此时，硬盘满的Node不回应任何消息(返回ACK等操作)。那么，续传对话就没有完成，其他Node还在缓存着数据。
然后异常的Node向ROOT报警，人工介入处理，处理完成后，ROOT广播一个“某Node连接成功/某Node在线了/等”的消息，
然后其他Node就认为这个Node又在线了，然后接着续传，然后一切就正常了。
```
curl -d '{ "Timeout":999 , "Recver":"n3" , "Data":"[echo]" , "Mode":1, "IsLog":true, "IsC2":true  }'    http://127.0.0.1:40078/echo
curl -d '{ "Timeout":999 , "Recver":"n1" , "Data":"[echo]" , "Mode":0, "IsLog":true, "IsC2":false }'    http://127.0.0.1:40078/echo
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
