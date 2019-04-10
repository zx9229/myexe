```
curl -d '{ "UID":"/n1/3/" , "Cmd":"help" , "Timeout":5 }'    http://127.0.0.1:10083/executeCommand
curl -d '{ "Cache":true , "Timeout":9 , "Topic":"test" , "Data":"testData" }'    http://127.0.0.1:10083/reportData
curl -d '{ "Cache":true , "Timeout":9 , "To":"略" , "Subject":"标题", "Content":"内容" }'    http://127.0.0.1:10083/sendMail
curl -d '{ "Cache":false , "Timeout":9 , "Recver":"//n1/3/" , "Data":"test_echo" }'    http://127.0.0.1:30078/echo
```
我想写一个程序，这个程序只有一个exe，再无依赖，所以使用sqlite，启用MySQL等。  
这个程序不会丢失数据，所以需要缓存和REQ+RSP模式。  
不管什么东西，只要扔到这程序里面就可以，所以它是一个归集程序。就好像日志装箱工具，例如windows中的Nxlog，linux中的Rsyslog等。  
我可以控制它。所以它好像一个监控程序。  
发送请求之后，可能响应要1分钟之后才能处理完，然后返回来，所以需要1请求n响应（收到请求了，结果1，结果2，响应结束）。  
被控者（执行命令者）执行期间如果和上手断线，就缓存到内存中，等待与上手重连后再发给上手，期间如果重启了就算了。

为什么要用websocket：想找一个没有粘包的tcp协议。
为什么要用protobuf：想找一个跨语言的通信协议。
API：websocket通信。
仅内存做了缓存。
可以执行命令。
