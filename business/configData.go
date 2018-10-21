package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

/*
大写字母A的10进制是65，所以Agent 端暂定端口号10065
{
    "ConfType": "AGENT",
    "UniqueID": "a1",
    "BelongID": "s1",
    "ServerURL":   { "Scheme": "ws", "Host": "localhost:10065", "Path": "/websocket" },
    "ClientURL": [ { "Scheme": "ws", "Host": "localhost:10083", "Path": "/websocket" } ],
    "DataSourceName": "database_a1.db",
    "LocationName": "Asia/Shanghai"
}
大写字母S的10进制是83，所以Server端暂定端口号10083
{
    "ConfType": "SERVER",
    "UniqueID": "a1",
    "ServerURL": { "Scheme": "ws", "Host": "localhost:10083", "Path": "/websocket" },
    "DataSourceName": "database_s1.db",
    "LocationName": "Asia/Shanghai"
}
*/
type configBase struct {
	ConfType string //可选值为(AGENT/SERVER)
}

type configAgent struct {
	configBase
	UniqueID       string
	BelongID       string
	ServerURL      url.URL
	ClientURL      []url.URL
	DataSourceName string //数据源的名字.
	LocationName   string //数据源的时区的名字.
}

//isValid 仅做基本检查(类似URL是否合法这种详细的检查,不做).
func (thls *configAgent) isValid() error {
	PREFIX := "configAgent: "
	var err error
	for range "1" {
		if len(thls.UniqueID) == 0 {
			err = fmt.Errorf(PREFIX + "UniqueID is empty")
			break
		}
		if len(thls.BelongID) == 0 {
			err = fmt.Errorf(PREFIX + "BelongID is empty")
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

type configServer struct {
	configBase
	UniqueID       string
	ServerURL      url.URL
	ClientURL      []url.URL
	DataSourceName string //数据源的名字.
	LocationName   string //数据源的时区的名字.
}

//isValid 仅做基本检查(类似URL是否合法这种详细的检查,不做).
func (thls *configServer) isValid() error {
	PREFIX := "configServer: "
	var err error
	for range "1" {
		if len(thls.UniqueID) == 0 {
			err = fmt.Errorf(PREFIX + "UniqueID is empty")
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
func parseContent(content string) (cfgA *configAgent, cfgS *configServer, err error) {
	byteSlice := []byte(content)
	var cfgB configBase
	if err = json.Unmarshal(byteSlice, &cfgB); err == nil {
		switch cfgB.ConfType {
		case "AGENT":
			cfgA = new(configAgent)
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
		var cfgA configAgent
		cfgA.ConfType = "AGENT"
		cfgA.UniqueID = "a1"
		cfgA.BelongID = "s1"
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
		cfgS.UniqueID = "a1"
		cfgS.ServerURL = url.URL{Scheme: "ws", Host: "localhost:10083", Path: "/websocket"}
		cfgS.ClientURL = []url.URL{}
		cfgS.DataSourceName = "database_s1.db"
		cfgS.LocationName = "Asia/Shanghai"
		byteSlice, _ := json.Marshal(cfgS)
		return string(byteSlice)
	}
	switch confType {
	case "AGENT":
		return exampleA()
	case "SERVER":
		return exampleS()
	default:
		return ""
	}
}
