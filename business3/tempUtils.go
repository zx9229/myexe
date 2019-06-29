package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zx9229/myexe/txdata"
)

const (
	//EMPTYSTR 空字符串
	EMPTYSTR = ""
	//FORONCE 仅执行一次for循环.
	FORONCE = "1"
	//ROOTNODE 根节点.
	ROOTNODE = "ROOT"
)

func toConfigNode(filename string) (cfg *configNode, err error) {
	for range FORONCE {
		var byteSlice []byte
		if byteSlice, err = ioutil.ReadFile(filename); err != nil {
			break
		}
		cfg = new(configNode)
		if err = json.Unmarshal(byteSlice, cfg); err != nil {
			break
		}
	}
	return
}

func cloneUniKey(src *txdata.UniKey) *txdata.UniKey {
	if src == nil {
		return nil
	}
	return &txdata.UniKey{UserID: src.UserID, MsgNo: src.MsgNo, SeqNo: src.SeqNo}
}

func toUniSym(src *txdata.UniKey) *UniSym {
	return &UniSym{UserID: src.UserID, MsgNo: src.MsgNo, SeqNo: src.SeqNo}
}

func int2mode(src int) (isPush bool, isSafe bool) {
	//0 不推送,不安全
	//1 不推送,要安全
	//2 要推送,不安全
	//3 要推送,要安全
	isSafe = ((src & 1) == 1)
	isPush = (((src >> 1) & 1) == 1)
	return
}

func filename_line_funcname() (filename string, line int, funcname string) {
	filename, line, funcname = "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(2)
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	if ok {
		filename = filepath.Base(filename)           // /full/path/basename.go => basename.go
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
	}
	return
}

func funName() (funcname string) {
	if pc, _, _, ok := runtime.Caller(2); ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
	} else {
		funcname = "???"
	}
	return
}
