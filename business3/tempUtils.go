package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zx9229/myexe/txdata"
)

const (
	//EMPTYSTR 空字符串
	EMPTYSTR = ""
	//FORONCE 仅执行一次for循环.
	FORONCE = "1"
)

func toConfigNode(filename string) (cfg *configNode) {
	var err error
	for range FORONCE {
		var byteSlice []byte
		if byteSlice, err = ioutil.ReadFile(filename); err != nil {
			break
		}
		cfg = new(configNode)
		err = json.Unmarshal(byteSlice, cfg)
		if err != nil {
			break
		}
	}
	assert4true(err == nil)
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
