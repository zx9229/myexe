package main

import (
	"sync"
	"unsafe"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/zxgo/zxxorm"
)

type safePushCache struct {
	sync.Mutex
	dbChan chan *DbOp
	engine *xorm.Engine
	Slc    []*txdata.PushWrap
	idx    int64
}

func newSafePushCache(dbc chan *DbOp, eng *xorm.Engine) *safePushCache {
	return &safePushCache{dbChan: dbc, engine: eng, Slc: make([]*txdata.PushWrap, 0)}
}

func (thls *safePushCache) Insert(data *txdata.PushWrap) (isInsert bool) {
	data.MsgNo = 0 //重新分配MsgNo
	if thls.dbChan == nil {
		thls.Lock()
		data.MsgNo = thls.idx + 1
		temp := DeepCopy(data).(*txdata.PushWrap)
		thls.Slc = append(thls.Slc, temp)
		assert4true(data.MsgNo == temp.MsgNo)
		thls.idx = data.MsgNo
		isInsert = true
		thls.Unlock()
	} else {
		action := &DbOp{}
		action.pm = data
		action.handler = func(session *xorm.Session, dbop *DbOp) {
			temp := &DbPushWrap{}
			temp.From(dbop.pm.(*txdata.PushWrap))
			dbop.affected, dbop.err = session.Insert(temp)
		}
		action.wg.Add(1)
		thls.dbChan <- action
		action.wg.Wait()
		var nilobj *DbPushWrap
		_, isInsert = nilobj.insertOneResult(action.affected, action.err)
	}
	return
}

//Select 筛选条件(msgNoBeg<MsgNo AND MsgNo<msgNoEnd)
func (thls *safePushCache) Select(msgNoBeg, msgNoEnd int64) (results []*txdata.PushWrap) {
	results = make([]*txdata.PushWrap, 0)
	if thls.engine == nil {
		thls.Lock()
		for _, node := range thls.Slc {
			if msgNoBeg < node.MsgNo {
				if 0 < msgNoEnd {
					if node.MsgNo < msgNoEnd {
						results = append(results, node)
					}
				} else {
					results = append(results, node)
				}
			}
		}
		thls.Unlock()
	} else {
		//在Xorm中，构建相对复杂的查询条件
		//https://www.golangtc.com/t/57b5e0e9b09ecc163500000e
		var err error
		tmpResults := make([]*DbPushWrap, 0)
		data4qry := &DbPushWrap{}
		fnMsgNo := zxxorm.GuessColName(thls.engine, data4qry, unsafe.Offsetof(data4qry.MsgNo), true)
		if 0 < msgNoEnd {
			err = thls.engine.Where(builder.Gt{fnMsgNo: msgNoBeg}.And(builder.Lt{fnMsgNo: msgNoEnd})).Find(&tmpResults)
		} else {
			err = thls.engine.Where(builder.Gt{fnMsgNo: msgNoBeg}).Find(&tmpResults)
		}
		if err == nil {
			for _, tmp := range tmpResults {
				results = append(results, tmp.To())
			}
		}
	}
	if results == nil || len(results) == 0 {
		results = nil
	}
	return
}
