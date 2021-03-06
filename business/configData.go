package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zx9229/myexe/txdata"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

/*
大写字母N的10进制是78，所以Node 端暂定端口号10078
{
    "ConfType": "NODE",
    "UserKey":   { "ZoneName": "", "NodeName": "n1", "ExecType": "NODE",   "ExecName": "" },
    "BelongKey": { "ZoneName": "", "NodeName": "s1", "ExecType": "SERVER", "ExecName": "" },
    "ServerURL":   { "Scheme": "ws", "Host": "localhost:10078", "Path": "/websocket" },
    "ClientURL": [ { "Scheme": "ws", "Host": "localhost:10083", "Path": "/websocket" } ],
    "DataSourceName": "database_n1.db",
    "LocationName": "Asia/Shanghai"
}
大写字母S的10进制是83，所以Server端暂定端口号10083
{
    "ConfType": "SERVER",
    "UserKey": { "ZoneName": "", "NodeName": "s1", "ExecType": "SERVER", "ExecName": "" },
    "ServerURL": { "Scheme": "ws", "Host": "localhost:10083", "Path": "/websocket" },
    "DataSourceName": "database_s1.db",
	"LocationName": "Asia/Shanghai",
	"MailCfg": {
        "Username": "用户名@126.com",
        "Password": "用户名的密码",
        "SMTPAddr": "smtp.126.com:25"
    }
}
*/
type configAtomic struct {
	ZoneName string
	NodeName string
	ExecType string
	ExecName string
}

func (thls *configAtomic) toTxAtomicKey() *txdata.AtomicKey {
	return &txdata.AtomicKey{ZoneName: thls.ZoneName, NodeName: thls.NodeName, ExecType: str2ProgramType(thls.ExecType), ExecName: thls.ExecName}
}

type configBase struct {
	ConfType string //可选值为(NODE/SERVER)
}

type configNode struct {
	configBase
	UserKey        configAtomic
	BelongKey      configAtomic
	ServerURL      url.URL
	ClientURL      []url.URL
	DataSourceName string //数据源的名字.
	LocationName   string //数据源的时区的名字.
}

//isValid 仅做基本检查(类似URL是否合法这种详细的检查,不做).
func (thls *configNode) isValid() error {
	PREFIX := "configNode: "
	var err error
	for range "1" {
		if !atomicKeyIsValid(thls.UserKey.toTxAtomicKey()) {
			err = fmt.Errorf(PREFIX + "UserKey is invalid")
			break
		}
		if len(thls.UserKey.NodeName) == 0 {
			err = fmt.Errorf(PREFIX + "UserKey.NodeName must not empty")
			break
		}
		if str2ProgramType(thls.UserKey.ExecType) != txdata.ProgramType_NODE {
			err = fmt.Errorf(PREFIX + "UserKey.ExecType must be NODE")
			break
		}
		if len(thls.UserKey.ExecName) != 0 {
			err = fmt.Errorf(PREFIX + "UserKey.ExecName must be empty")
			break
		}
		if !atomicKeyIsValid(thls.BelongKey.toTxAtomicKey()) {
			err = fmt.Errorf(PREFIX + "BelongKey is invalid")
			break
		}
		if len(thls.BelongKey.NodeName) == 0 {
			err = fmt.Errorf(PREFIX + "BelongKey.NodeName must not empty")
			break
		}
		if str2ProgramType(thls.BelongKey.ExecType) != txdata.ProgramType_NODE &&
			str2ProgramType(thls.BelongKey.ExecType) != txdata.ProgramType_SERVER {
			err = fmt.Errorf(PREFIX + "UserKey.ExecType must be NODE or SERVER")
			break
		}
		if len(thls.BelongKey.ExecName) != 0 {
			err = fmt.Errorf(PREFIX + "BelongKey.ExecName must be empty")
			break
		}
		if len(thls.ServerURL.String()) == 0 {
			err = fmt.Errorf(PREFIX + "ServerURL is empty")
			break
		}
		if thls.ClientURL != nil {
			for _, u := range thls.ClientURL {
				if len(u.String()) == 0 {
					err = fmt.Errorf(PREFIX + "ClientURL contains empty elements")
					break
				}
			}
		}
		if len(thls.DataSourceName) == 0 {
			err = fmt.Errorf(PREFIX + "DataSourceName is empty")
			break
		}
	}
	return err
}

type config4mail struct {
	Username string //邮箱的用户名
	Password string //邮箱的密码
	SmtpAddr string //邮箱的SMTP地址
}

type configServer struct {
	configBase
	UserKey        configAtomic
	ServerURL      url.URL
	ClientURL      []url.URL
	DataSourceName string //数据源的名字.
	LocationName   string //数据源的时区的名字.
	MailCfg        config4mail
}

//isValid 仅做基本检查(类似URL是否合法这种详细的检查,不做).
func (thls *configServer) isValid() error {
	PREFIX := "configServer: "
	var err error
	for range "1" {
		if !atomicKeyIsValid(thls.UserKey.toTxAtomicKey()) {
			err = fmt.Errorf(PREFIX + "UserKey is invalid")
			break
		}
		if len(thls.ServerURL.String()) == 0 {
			err = fmt.Errorf(PREFIX + "ServerURL is empty")
			break
		}
		if thls.ClientURL != nil {
			for _, u := range thls.ClientURL {
				if len(u.String()) == 0 {
					err = fmt.Errorf(PREFIX + "ClientURL contains empty elements")
					break
				}
			}
		}
		if len(thls.DataSourceName) == 0 {
			err = fmt.Errorf(PREFIX + "DataSourceName is empty")
			break
		}
	}
	return err
}

//parseContent 解析内容到配置结构体
func parseContent(content string) (cfgA *configNode, cfgS *configServer, err error) {
	byteSlice := []byte(content)
	var cfgB configBase
	if err = json.Unmarshal(byteSlice, &cfgB); err == nil {
		switch cfgB.ConfType {
		case "NODE":
			cfgA = new(configNode)
			err = json.Unmarshal(byteSlice, cfgA)
		case "SERVER":
			cfgS = new(configServer)
			err = json.Unmarshal(byteSlice, cfgS)
		default:
			err = fmt.Errorf("unknown ConfType=%v", cfgB.ConfType)
		}
	}
	if err == nil && cfgA != nil {
		err = cfgA.isValid()
	}
	if err == nil && cfgS != nil {
		err = cfgS.isValid()
	}
	if err != nil {
		cfgA = nil
		cfgS = nil
	}
	return
}

func saveConfigToDb(dataSourceName, content string) error {
	var err error
	var tmpEngine *xorm.Engine
	for range "1" {
		if tmpEngine, err = xorm.NewEngine("sqlite3", dataSourceName); err != nil {
			break
		}
		//支持struct为驼峰式命名,表结构为下划线命名之间的转换,同时对于特定词支持更好.
		tmpEngine.SetMapper(core.GonicMapper{})
		//应该是:只要存在这个tablename,就跳过它.
		if err = tmpEngine.CreateTables(&KeyValue{}); err != nil {
			break
		}
		var affected int64
		if affected, err = tmpEngine.ID("json").Update(&KeyValue{Value: content}); err != nil {
			break
		}
		if affected != 0 && affected != 1 { //此判断可有可无.
			panic(affected)
		}
		if affected == 1 { //找到了key,更新数据成功.
			break
		}
		if affected == 0 {
			if affected, err = tmpEngine.InsertOne(&KeyValue{Key: "json", Value: content}); err != nil {
				break
			}
			if affected != 1 {
				err = fmt.Errorf("InsertOne with err=%v and affected=%v", err, affected)
			}
		}
	}
	if tmpEngine != nil {
		tmpEngine.Close()
	}
	return err
}

func loadConfigFromDb(dataSourceName string) (content string, err error) {
	var tmpEngine *xorm.Engine
	for range "1" {
		if tmpEngine, err = xorm.NewEngine("sqlite3", dataSourceName); err != nil {
			break
		}
		//支持struct为驼峰式命名,表结构为下划线命名之间的转换,同时对于特定词支持更好.
		tmpEngine.SetMapper(core.GonicMapper{})
		//应该是:只要存在这个tablename,就跳过它.
		if err = tmpEngine.CreateTables(&KeyValue{}); err != nil {
			break
		}
		kv := KeyValue{Key: "json"}
		var isOk bool
		if isOk, err = tmpEngine.Get(&kv); err != nil {
			break
		}
		if !isOk {
			err = fmt.Errorf("can not find Key=%v in dataSourceName", kv.Key)
			break
		}
		content = kv.Value
	}
	if tmpEngine != nil {
		tmpEngine.Close()
	}
	return
}

func exampleConfigData(confType string) string {
	exampleA := func() string {
		var cfgA configNode
		cfgA.ConfType = "NODE"
		cfgA.UserKey = configAtomic{NodeName: "n1", ExecType: "NODE"}
		cfgA.BelongKey = configAtomic{NodeName: "s1", ExecType: "SERVER"}
		cfgA.ServerURL = url.URL{Scheme: "ws", Host: "localhost:10065", Path: "/websocket"}
		cfgA.ClientURL = []url.URL{url.URL{Scheme: "ws", Host: "localhost:10083", Path: "/websocket"}}
		cfgA.DataSourceName = "database_a1.db"
		cfgA.LocationName = "Asia/Shanghai"
		byteSlice, _ := json.Marshal(cfgA)
		return string(byteSlice)
	}
	exampleS := func() string {
		var cfgS configServer
		cfgS.ConfType = "SERVER"
		cfgS.UserKey = configAtomic{NodeName: "s1", ExecType: "SERVER"}
		cfgS.ServerURL = url.URL{Scheme: "ws", Host: "localhost:10083", Path: "/websocket"}
		cfgS.ClientURL = []url.URL{}
		cfgS.DataSourceName = "database_s1.db"
		cfgS.LocationName = "Asia/Shanghai"
		byteSlice, _ := json.Marshal(cfgS)
		return string(byteSlice)
	}
	switch confType {
	case "NODE":
		return exampleA()
	case "SERVER":
		return exampleS()
	default:
		return ""
	}
}
