```
curl -d "{ \"UID\":\"a1\" , \"Cmd\":\"help\" , \"Timeout\":5 }"    http://127.0.0.1:10083/executeCommand
curl -d "{ \"Cache\":true , \"Timeout\":5 , \"Topic\":\"test\" , \"Data\":\"testData\" }"    http://127.0.0.1:10065/reportData
```
我想写一个程序，这个程序只有一个exe，再无依赖，所以使用sqlite，启用MySQL等。  
这个程序不会丢失数据，所以需要缓存和REQ+RSP模式。  
不管什么东西，只要扔到这程序里面就可以，所以它是一个归集程序。就好像日志装箱工具，例如windows中的Nxlog，linux中的Rsyslog等。  
我可以控制它。所以它好像一个监控程序。  
需要添加"发送邮件"功能。  
要么对接微信(等)，要么一劳永逸的写一个APP。  
把report和push合并，增加一个选项（尽力到达和发送即可）。  
