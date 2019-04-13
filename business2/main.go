package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

func main() {
	glog.Infoln(os.Args)
	cfgNode := toConfigNode(os.Args[1])
	globalNode := newBusinessNode(cfgNode)
	cs := wsnet.NewWsCliSrv()
	cs.CbConnected = globalNode.onConnected
	cs.CbDisconnected = globalNode.onDisconnected
	cs.CbReceive = globalNode.onMessage
	cs.Init(cfgNode.ClientURL, cfgNode.ServerURL)

	nodeCache := func(w http.ResponseWriter, r *http.Request) {
		jsonContent := calcNodeCache(globalNode)
		fmt.Fprintf(w, jsonContent)
	}
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/cache", nodeCache)
	cs.Run()
}

func calcNodeCache(node *businessNode) (jsonContent string) {
	tmpObj := new(struct {
		OwnInfo    *txdata.ConnectionInfo
		ParentInfo *safeFatherData
		RootOnline bool
		CacheUser  *safeConnInfoMap
		CacheSock  *safeWsSocketMap
	})
	tmpObj.OwnInfo = &node.ownInfo
	tmpObj.ParentInfo = &node.parentInfo
	tmpObj.RootOnline = node.rootOnline
	tmpObj.CacheUser = node.cacheUser
	tmpObj.CacheSock = node.cacheSock

	tmpObj.ParentInfo.Lock()
	defer tmpObj.ParentInfo.Unlock()
	tmpObj.CacheUser.Lock()
	defer tmpObj.CacheUser.Unlock()
	tmpObj.CacheSock.Lock()
	defer tmpObj.CacheSock.Unlock()

	if byteSlice, err := json.Marshal(tmpObj); err != nil {
		glog.Fatalln(err, tmpObj)
	} else {
		jsonContent = string(byteSlice)
	}

	return
}
