# myexe
我的程序。  

# 下载代码
```
go get -u -v github.com/zx9229/myexe/business
```

# 编译
```
cd ./myexe/business
go build
```
你有可能看到下面的错误提示：
```
\> go build
# github.com/mattn/go-sqlite3
cc1.exe: sorry, unimplemented: 64-bit mode not compiled in
```
它大概表示Golang是64位的，MinGW是32位的，无法编译。此时你需要替换编译器为MinGW-W64。  
