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

//UniSym omit
type UniSym struct {
	UserID string
	MsgNo  int64
	SeqNo  int32
}

func (thls *UniSym) fromUniKey(src *txdata.UniKey) {
	thls.UserID = src.UserID
	thls.MsgNo = src.MsgNo
	thls.SeqNo = src.SeqNo
	assert4true(0 <= thls.MsgNo)
	assert4true(0 <= thls.SeqNo)
}

type node4sync struct {
	Key      UniSym       //(一经设置,不再修改)
	ToRoot   bool         //(一经设置,不再修改)
	RecverID string       //(一经设置,不再修改)
	data     ProtoMessage //(一经设置,不再修改)
}

type safeSynchCache struct {
	sync.Mutex
	dbChan  chan *DbOp
	xEngine *xorm.Engine
	M       map[UniSym]*node4sync
}

func newSafeSynchCache(dbc chan *DbOp, engine *xorm.Engine) *safeSynchCache {
	return &safeSynchCache{dbChan: dbc, xEngine: engine, M: make(map[UniSym]*node4sync)}
}

//insertData 入参uniKey是入参pm的一个数据成员.
func (thls *safeSynchCache) insertData(uniKey *txdata.UniKey, toR bool, rID string, pm ProtoMessage) (isExist, isInsert bool) {
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
			temp := &DbCommonReqRsp{}
			if c2req, isOk := dbop.pm.(*txdata.Common2Req); isOk {
				temp.FromC2Req(c2req)
				dbop.affected, dbop.err = session.Insert(temp)
			} else if c2rsp, isOk := dbop.pm.(*txdata.Common2Rsp); isOk {
				temp.FromC2Rsp(c2rsp)
				dbop.affected, dbop.err = session.Insert(temp)
			} else {
				dbop.affected = 0
				dbop.err = errors.New("xxxx")
			}
		}
		action.wg.Add(1)
		thls.dbChan <- action
		action.wg.Wait()
		isExist, isInsert = (*DbCommonReqRsp)(nil).insertOneResult(action.affected, action.err)
	}
	return
}

func (thls *safeSynchCache) deleteData(uniKey *txdata.UniKey) {
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
			dbop.affected, dbop.err = session.Id(core.PK{key.UserID, key.MsgNo, key.SeqNo}).Delete(&DbCommonReqRsp{})
		}
		action.wg.Add(1)
		thls.dbChan <- action
		action.wg.Wait()
	}
	return
}

func (thls *safeSynchCache) queryData(toR bool, rID string) (results []ProtoMessage) {
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
		tmpResults := make([]*DbCommonReqRsp, 0)
		data4qry := &DbCommonReqRsp{}
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

func (thls *safeSynchCache) queryCount(uID string, msgNo int64) (cnt int) {
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
		data4qry := &DbCommonReqRsp{}
		fnUserID := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.UserID), true)
		fnMsgNo := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.MsgNo), true)
		if num, err = thls.xEngine.Where(builder.Eq{fnUserID: uID}.And(builder.Eq{fnMsgNo: msgNo})).Count(data4qry); err == nil {
			cnt = int(num)
		}
	}
	return
}

func (thls *safeSynchCache) queryDataByToRoot(toR bool) (results []ProtoMessage) {
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
		tmpResults := make([]*DbCommonReqRsp, 0)
		data4qry := &DbCommonReqRsp{}
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
func (thls *safeSynchCache) MarshalJSON() ([]byte, error) {
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
