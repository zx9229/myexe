package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/zxgo/zxxorm"
)

type safeSynchCacheRoot struct {
	sync.Mutex
	dbChan  chan *DbOp
	xEngine *xorm.Engine
	M       map[UniSym]*node4sync
}

func newSafeSynchCacheRoot(engine *xorm.Engine, dbc chan *DbOp) *safeSynchCacheRoot {
	return &safeSynchCacheRoot{dbChan: dbc, xEngine: engine, M: make(map[UniSym]*node4sync)}
}

//insertData 入参uniKey是入参pm的一个数据成员.
func (thls *safeSynchCacheRoot) insertData(uniKey *txdata.UniKey, toR bool, rID string, pm ProtoMessage, tmpTag int32) (isExist, isInsert bool) {
	//isExist, isInsert
	//true   , false    (已经存在了,肯定插不进去)
	//false  , true
	//false  , false (异常情况)
	if thls.dbChan == nil {
		node := new(node4sync)
		node.Key.fromUniKey(uniKey)
		node.ToRoot = toR
		node.RecverID = rID
		node.data = pm
		node.TmpFlag = tmpTag
		thls.Lock()
		if _, isExist = thls.M[node.Key]; !isExist {
			thls.M[node.Key] = node
			isInsert = true
		}
		thls.Unlock()
	} else {
		action := &DbOp{}
		action.pm = pm
		action.handler = func(session *xorm.Session, dbop *DbOp) {
			temp := &DbCommonReqRspRoot{}
			if c2req, isOk := dbop.pm.(*txdata.Common2Req); isOk {
				temp.FromC2Req(c2req, tmpTag)
				dbop.affected, dbop.err = session.Insert(temp)
			} else if c2rsp, isOk := dbop.pm.(*txdata.Common2Rsp); isOk {
				temp.FromC2Rsp(c2rsp, tmpTag)
				dbop.affected, dbop.err = session.Insert(temp)
			} else {
				dbop.affected = 0
				dbop.err = errors.New("xxxx")
			}
		}
		action.wg.Add(1)
		thls.dbChan <- action
		action.wg.Wait()
		isExist, isInsert = (*DbCommonReqRspRoot)(nil).insertOneResult(action.affected, action.err)
	}
	return
}

func (thls *safeSynchCacheRoot) deleteData(uniKey *txdata.UniKey) {
	if thls.dbChan == nil {
		var sym UniSym
		sym.fromUniKey(uniKey)
		thls.Lock()
		delete(thls.M, sym)
		thls.Unlock()
	} else {
		action := &DbOp{}
		action.pm = uniKey
		action.handler = func(session *xorm.Session, dbop *DbOp) {
			key := dbop.pm.(*txdata.UniKey)
			//Id(interface{})传入(一个)主键字段的值，作为查询条件:engine.Id(1).Get(&user)
			//如果是复合主键，则可以:engine.Id(core.PK{1, "name"}).Get(&user)
			//传入的两个参数按照struct中pk标记字段出现的顺序赋值.
			dbop.affected, dbop.err = session.Id(core.PK{key.UserID, key.MsgNo, key.SeqNo}).Delete(&DbCommonReqRspRoot{})
		}
		action.wg.Add(1)
		thls.dbChan <- action
		action.wg.Wait()
	}
	return
}

func (thls *safeSynchCacheRoot) queryData(toR bool, rID string) (results []ProtoMessage) {
	if thls.dbChan == nil {
		thls.Lock()
		for _, node := range thls.M {
			if node.ToRoot == toR && node.RecverID == rID {
				if results == nil {
					results = make([]ProtoMessage, 0)
				}
				results = append(results, node.data)
			}
		}
		thls.Unlock()
	} else {
		var err error
		tmpResults := make([]*DbCommonReqRspRoot, 0)
		data4qry := &DbCommonReqRspRoot{}
		fnToRoot := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.ToRoot), true)
		fnRecverID := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.RecverID), true)
		if err = thls.xEngine.Where(builder.Eq{fnToRoot: toR}.And(builder.Eq{fnRecverID: rID})).Find(&tmpResults); err == nil {
			for _, tmp := range tmpResults {
				if pmTMP := tmp.To(); pmTMP != nil {
					if results == nil {
						results = make([]ProtoMessage, 0)
					}
					results = append(results, pmTMP)
				}
			}
		}
	}
	return
}

func (thls *safeSynchCacheRoot) queryCount(uID string, msgNo int64) (cnt int) {
	if thls.xEngine == nil {
		thls.Lock()
		for _, node := range thls.M {
			if node.Key.UserID == uID && node.Key.MsgNo == msgNo {
				cnt++
			}
		}
		thls.Unlock()
	} else {
		var num int64
		var err error
		data4qry := &DbCommonReqRspRoot{}
		fnUserID := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.UserID), true)
		fnMsgNo := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.MsgNo), true)
		if num, err = thls.xEngine.Where(builder.Eq{fnUserID: uID}.And(builder.Eq{fnMsgNo: msgNo})).Count(data4qry); err == nil {
			cnt = int(num)
		}
	}
	return
}

func (thls *safeSynchCacheRoot) queryDataByToRoot(toR bool) (results []ProtoMessage) {
	if thls.xEngine == nil {
		thls.Lock()
		for _, node := range thls.M {
			if node.ToRoot == toR {
				if results == nil {
					results = make([]ProtoMessage, 0)
				}
				results = append(results, node.data)
			}
		}
		thls.Unlock()
	} else {
		var err error
		tmpResults := make([]*DbCommonReqRspRoot, 0)
		data4qry := &DbCommonReqRspRoot{}
		fnToRoot := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.ToRoot), true)
		if err = thls.xEngine.Where(builder.Eq{fnToRoot: toR}).Find(&tmpResults); err == nil {
			for _, tmp := range tmpResults {
				if pmTMP := tmp.To(); pmTMP != nil {
					if results == nil {
						results = make([]ProtoMessage, 0)
					}
					results = append(results, pmTMP)
				}
			}
		}
	}
	return
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeSynchCacheRoot) MarshalJSON() ([]byte, error) {
	tmpMap := make(map[string]string)
	thls.Lock()
	for k, v := range thls.M {
		tmpK := fmt.Sprintf("(%v|%v|%v)", k.UserID, k.MsgNo, k.SeqNo)
		tmpV := fmt.Sprintf("(%v|%v|%v)", v.RecverID, v.ToRoot, CalcMessageType(v.data))
		tmpMap[tmpK] = tmpV
	}
	thls.Unlock()
	return json.Marshal(tmpMap)
}
