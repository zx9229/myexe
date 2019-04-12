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

func (thls *safeWsSocketMap) MarshalJSON() ([]byte, error) {
	//thls.Lock()
	//defer thls.Unlock()
	tmpObj := new(struct {
		M map[string]bool
	})
	tmpObj.M = make(map[string]bool)
	for k, v := range thls.M {
		tmpObj.M[fmt.Sprintf("%v", k)] = v
	}
	return json.Marshal(tmpObj)
}
