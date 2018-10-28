package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

//-json        设置json文件的内容到sqlite数据库.
//-sqlite      指定要使用的sqlite数据库.
//-log_dir     启用(golang/glog)
//-logtostderr 启用(golang/glog)
func main() {
	if true {
		flag.Set("stderrthreshold", "WARNING")
		flag.Set("logtostderr", "true")
	}
	var (
		argHelp   bool
		argJSON   string
		argSqlite string
	)
	//我的参数和(glog)的参数混合了,所以用[M]标识这个参数是我的参数,不是(glog)的参数.
	flag.BoolVar(&argHelp, "help", false, "[M] display this help and exit.")
	flag.StringVar(&argJSON, "json", "", "[M] set the contents of the json file to the sqlite database.")
	flag.StringVar(&argSqlite, "sqlite", "", "[M] specify the sqlite database to use.")
	flag.Parse()
	if argHelp {
		flag.Usage()
		fmt.Println()
		fmt.Println("NOTE: please use (-logtostderr) to view details.")
		fmt.Println()
		fmt.Println(exampleConfigData("AGENT"))
		fmt.Println()
		fmt.Println(exampleConfigData("SERVER"))
		return
	}
	defer glog.Flush()
	for range "1" {
		cfgA, cfgS, err := handleArgs(argJSON, argSqlite)
		if err != nil {
			glog.Infoln(err)
			break
		}
		if (cfgA == nil) && (cfgS == nil) {
			glog.Infoln("save json to the database successfully.")
			break
		}
		if cfgA != nil {
			runAgent(cfgA)
		} else if cfgS != nil {
			runServer(cfgS)
		} else {
			panic("logical_error")
		}
	}
	glog.Infoln("the program is about to quit ...")
	log.Println("the program is about to quit ...")
}

//handleArgs 处理程序的入参.
//(err != nil)表示:函数报错,程序应当直接退出.
//(cfgA == nil && cfgS == nil)表示:将文件argJSON的内容写入sqlite中,此时工作已完成,程序应当直接退出.
//(cfgA != nil)表示:执行cfgA的逻辑.
//(cfgS != nil)表示:执行cfgS的逻辑.
func handleArgs(argJSON, argSqlite string) (cfgA *configAgent, cfgS *configServer, err error) {
	for range "1" {
		var content string
		if len(argJSON) != 0 {
			var byteSlice []byte
			if byteSlice, err = ioutil.ReadFile(argJSON); err != nil {
				break
			}
			content = string(byteSlice)
			if cfgA, cfgS, err = parseContent(content); err != nil {
				break
			}
			if (cfgA == nil && cfgS == nil) || (cfgA != nil && cfgS != nil) {
				panic("logical_error")
			}
			var dataSourceName string
			if cfgA != nil {
				dataSourceName = cfgA.DataSourceName
			}
			if cfgS != nil {
				dataSourceName = cfgS.DataSourceName
			}
			err = saveConfigToDb(dataSourceName, content)
			cfgA = nil
			cfgS = nil
		} else if len(argSqlite) != 0 {
			if content, err = loadConfigFromDb(argSqlite); err != nil {
				break
			}
			if cfgA, cfgS, err = parseContent(content); err != nil {
				break
			}
			if (cfgA == nil && cfgS == nil) || (cfgA != nil && cfgS != nil) {
				panic("logical_error")
			}
			if cfgA != nil && cfgA.DataSourceName != argSqlite {
				err = fmt.Errorf("not equal, argSqlite=%v, cfgA.DataSourceName=%v", argSqlite, cfgA.DataSourceName)
			}
			if cfgS != nil && cfgS.DataSourceName != argSqlite {
				err = fmt.Errorf("not equal, argSqlite=%v, cfgS.DataSourceName=%v", argSqlite, cfgS.DataSourceName)
			}
		} else {
			err = fmt.Errorf("is empty, argJSON=%v, argSqlite=%v", argJSON, argSqlite)
		}
	}
	return
}
