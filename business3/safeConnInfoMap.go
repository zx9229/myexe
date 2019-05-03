package main

import (
	"encoding/json"
	"errors"
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

func (thls *safeConnInfoMap) sendDataToSon(data ProtoMessage) {
	var byteSlice []byte
	thls.Lock()
	for _, val := range thls.M {
		if len(val.Pathway) == 1 {
			if byteSlice == nil {
				byteSlice = msg2package(data)
			}
			val.conn.Send(byteSlice)
		}
	}
	thls.Unlock()
}

func (thls *safeConnInfoMap) sendDataToUser(data ProtoMessage, userID string) (err error) {
	var isExist bool
	var cInfoEx *connInfoEx

	thls.Lock()
	cInfoEx, isExist = thls.M[userID]
	thls.Unlock()

	if !isExist {
		return errors.New("user if offline")
	}
	return cInfoEx.conn.Send(msg2package(data))
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeConnInfoMap) MarshalJSON() (byteSlice []byte, err error) {
	thls.Lock()
	byteSlice, err = json.Marshal(thls.M)
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) tmpF1() (data map[string]*txdata.ConnectReq) {
	data = make(map[string]*txdata.ConnectReq)
	thls.Lock()
	for k, v := range thls.M {
		data[k] = &txdata.ConnectReq{InfoReq: &v.Info, Pathway: v.Pathway}
	}
	thls.Unlock()
	return
}
