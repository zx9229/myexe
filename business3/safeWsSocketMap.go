package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zx9229/myexe/wsnet"
)

type safeWsSocketMap struct {
	sync.Mutex
	M map[*wsnet.WsSocket]bool
}

func newSafeWsSocketMap() *safeWsSocketMap {
	return &safeWsSocketMap{M: make(map[*wsnet.WsSocket]bool)}
}

func (thls *safeWsSocketMap) insertData(k *wsnet.WsSocket, v bool) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[k]; !isSuccess {
		thls.M[k] = v
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeWsSocketMap) deleteData(k *wsnet.WsSocket) (v bool, isSuccess bool) {
	thls.Lock()
	if v, isSuccess = thls.M[k]; isSuccess {
		delete(thls.M, k)
	}
	thls.Unlock()
	return
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeWsSocketMap) MarshalJSON() (byteSlice []byte, err error) {
	thls.Lock()
	tmpMap := make(map[string]bool)
	for k, v := range thls.M {
		tmpMap[fmt.Sprintf("%v", k)] = v
	}
	byteSlice, err = json.Marshal(tmpMap)
	thls.Unlock()
	return
}
