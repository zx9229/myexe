```
curl -d '{ "Timeout":999 , "Recver":"n3" , "Data":"[echo]" , "Mode":1, "IsLog":true }'    http://127.0.0.1:40078/echo
```
在日志中，我们应当按照Key进行检索。日志中的Key一般如下所示：
```
msgData=Key:<UserID:"n4" MsgNo:8042223290400000001 >
msgData=Key:<UserID:"n4" MsgNo:8042223290400000001 SeqNo:1 >
msgData=Key:<UserID:"n4" MsgNo:8042223290400000001 SeqNo:2 >
```
所以，我们应当先找到对应的Key，然后然后将所有的日志文件汇总成一个日志文件：  
`grep -r 'Key:<UserID:"n4" MsgNo:8042223290400000001' .  > ../t1.log`  
然后对这个日志文件，按照时间戳排序：  
`sort -t" " -k2 t1.log > t2.log`  
然后就可以慢慢的分析日志了。  
`grep -r 'Key:<UserID:"n4" MsgNo:8042223290400000001' . | sort -t" " -k2`
