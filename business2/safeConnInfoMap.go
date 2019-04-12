package main

import (
	"encoding/json"
	"sync"

	"github.com/golang/glog"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

type connInfoEx struct {
	conn    *wsnet.WsSocket
	Info    txdata.ConnectionInfo
	Pathway []string
}

type safeConnInfoMap struct {
	sync.Mutex
	M map[string]*connInfoEx
}

func newSafeConnInfoMap() *safeConnInfoMap {
	return &safeConnInfoMap{M: make(map[string]*connInfoEx)}
}

func (thls *safeConnInfoMap) humanReadable() (jsonContent string) {
	thls.Lock()
	if byteSlice, err := json.Marshal(thls.M); err != nil {
		glog.Fatalln(err, thls.M)
	} else {
		jsonContent = string(byteSlice)
	}
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) queryData(key string) (connEx *connInfoEx, isExist bool) {
	thls.Lock()
	connEx, isExist = thls.M[key]
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) isValidData(data *connInfoEx) bool {
	var isOk bool
	for range FORONCE {
		if data.conn == nil {
			break
		}
		if data.Info.UserID == EMPTYSTR {
			break
		}
		if data.Pathway == nil {
			break
		}
		isOk = true
	}
	return isOk
}

func (thls *safeConnInfoMap) insertData(data *connInfoEx) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[data.Info.UserID]; !isSuccess {
		thls.M[data.Info.UserID] = data
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeConnInfoMap) deleteData(key string) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[key]; isSuccess {
		delete(thls.M, key)
	}
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) deleteDataByConn(conn *wsnet.WsSocket) []*connInfoEx {
	var dataSlice []*connInfoEx
	thls.Lock()
	for key, val := range thls.M {
		if val.conn == conn {
			if dataSlice == nil {
				dataSlice = make([]*connInfoEx, 0)
			}
			dataSlice = append(dataSlice, val)
			delete(thls.M, key)
		}
	}
	thls.Unlock()
	return dataSlice
}
