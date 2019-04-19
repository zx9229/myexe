package main

import (
	"sync"

	"github.com/golang/glog"
	"github.com/zx9229/myexe/txdata"
)

//PshKey omit
type PshKey struct {
	UserID string
	SeqNo  int64
}

//PshVal omit
type PshVal struct {
	sync.Mutex
	insIdx int64 //插入缓存时的序号(ins:insert,从1开始递增,被赋值之后,不允许再变动)
	slcIdx int   //处于切片中的序号(slc:slice,随着slice的变化而变化,负数(-1)表示不在slice中)
	Psh    *txdata.DataPsh
	Ack    *txdata.DataAck
	Resp   map[PshKey]*txdata.DataPsh //如果Psh是请求的话那么这个缓存就是请求的响应.
}

func (thls *PshVal) relateAck(dataA *txdata.DataAck) (isSuccess bool) {
	thls.Lock()
	if thls.Psh.SendNo == dataA.SendNo && thls.Psh.SendUID == dataA.SendUID {
		if thls.Ack == nil || thls.Ack.ErrNo != 0 {
			thls.Ack = dataA
			isSuccess = true
		} else { //thls.Ack.ErrNo为0
			if dataA.ErrNo == 0 {
				thls.Ack = dataA
				isSuccess = true
			} else {
				//已经收到了一个正确的回复,然后又收到了一个失败的回复.
				glog.Errorf("relateAck, prevA=%v, dataA=%v", thls.Ack, dataA)
				isSuccess = false
			}
		}
	} else {
		isSuccess = false
	}
	thls.Unlock()
	return
}

func (thls *PshVal) relatePsh4Rsp(dataP *txdata.DataPsh, replaceIt bool) (isSuccess bool) {
	thls.Lock()
	if thls.Psh.SendNo == dataP.RecvNo && thls.Psh.SendUID == dataP.RecvUID {
		key := PshKey{UserID: dataP.SendUID, SeqNo: dataP.SendNo}
		if thls.Resp == nil {
			thls.Resp = make(map[PshKey]*txdata.DataPsh)
		}
		if replaceIt {
			thls.Resp[key] = dataP
			isSuccess = true
		} else {
			if _, isSuccess = thls.Resp[key]; !isSuccess {
				thls.Resp[key] = dataP
			}
			isSuccess = !isSuccess
		}
	} else {
		isSuccess = false
	}
	thls.Unlock()
	return
}

func (thls *PshVal) deletePsh4Rsp(key *PshKey) (isSuccess bool) {
	thls.Lock()
	if thls.Resp != nil {
		if _, isSuccess = thls.Resp[*key]; isSuccess {
			delete(thls.Resp, *key)
		}
	}
	thls.Unlock()
	return
}

type safeDataPshCache struct {
	sync.Mutex
	isRoot bool  //现在是根节点(ROOT)在使用缓存,(一经赋值,禁止修改)
	maxIdx int64 //因为是int64所以假定它无法溢出.
	M      map[PshKey]*PshVal
	Slc    []*PshVal //尚未同步成功的数据存放在这里.
}

func newSafeDataPshCache(isR bool) *safeDataPshCache {
	curData := new(safeDataPshCache)
	curData.isRoot = isR
	curData.maxIdx = 0
	curData.M = make(map[PshKey]*PshVal)
	curData.Slc = make([]*PshVal, 0)
	return curData
}

//①没有这个消息,缓存成功.
//②存在这个消息,返回成功的警告.
//③没有这个消息,但是硬盘满了,sqlite插入数据库失败,(缓存失败).
func (thls *safeDataPshCache) feedDataPsh(dataP *txdata.DataPsh, ackOut *txdata.DataAck) {
	assert4true(dataP.SendNo != 0) //必须SendNo非0.

	var isExist bool
	key := PshKey{UserID: dataP.SendUID, SeqNo: dataP.SendNo}
	thls.Lock()
	if _, isExist = thls.M[key]; !isExist {
		val := &PshVal{insIdx: thls.maxIdx + 1, slcIdx: len(thls.Slc), Psh: dataP}
		thls.M[key] = val
		thls.Slc = append(thls.Slc, val)
		thls.maxIdx = val.insIdx
		//////////////////////////////////////////////////////////////////////////
		if thls.isRoot {
			if dataP.RecvNo != 0 { //(rsp的Psh) //因为SendNo非0所以RecvNo为0是找不到数据的.
				key4req := PshKey{UserID: dataP.RecvUID, SeqNo: dataP.RecvNo}
				if val4req, isOk := thls.M[key4req]; isOk {
					isOk = val4req.relatePsh4Rsp(dataP, true)
					assert4true(isOk)
				}
			} else { //(req的Psh)
				//TODO:可能rsp的Psh先进来,然后req的Psh后进来,此时需要,查出将所有rsp的Psh,然后关联到req的Psh,这里先不实现它.
				//这里先不实现它,因为,对于这个类来说,逻辑上是可能出现的,但是对于整个程序来说,是不应当出现的.
			}
		}
		//////////////////////////////////////////////////////////////////////////
	}
	thls.Unlock()
	if isExist {
		ackOut.ErrNo = 0
		ackOut.ErrMsg = "already exist"
	}
	return
}

func (thls *safeDataPshCache) feedDataAck(dataA *txdata.DataAck, delByAck bool) (isSuccess bool) {
	assert4true(dataA.SendNo != 0) //必须SendNo非0.

	key := PshKey{UserID: dataA.SendUID, SeqNo: dataA.SendNo}
	thls.Lock()
	if val, isOk := thls.M[key]; isOk {
		isOk = val.relateAck(dataA)
		assert4true(isOk)
		isSuccess = true
		if dataA.ErrNo == 0 && delByAck { //从Ack知道了这个Psh投递成功,其使命已经完成,现在又要执行del操作.
			delete(thls.M, key)
			if 0 <= val.slcIdx { //它现在处于切片中.
				assert4true(thls.Slc[val.slcIdx] == val)
				thls.Slc = append(thls.Slc[:val.slcIdx], thls.Slc[val.slcIdx+1:]...)
				for i := val.slcIdx; i < len(thls.Slc); i++ {
					thls.Slc[i].slcIdx--
				}
			}
			//////////////////////////////////////////////////////////////////////////
			if thls.isRoot {
				if val.Psh.RecvNo != 0 { //这个消息传输了Rsp消息,所以要找到Req消息,然后继续清理.
					key4req := PshKey{UserID: val.Psh.RecvUID, SeqNo: val.Psh.RecvNo}
					if val4req, isOk := thls.M[key4req]; isOk {
						isOk = val4req.deletePsh4Rsp(&key)
						assert4true(isOk) //逻辑上一定能执行成功.
					} else {
						assert4true(isOk) //逻辑上一定能找到对应的请求包,所以这里让它崩溃.
					}
				} else {
					//已经在上面执行了清理操作.
					//这个消息传输了Req消息,要清理的话,直接清理掉它自己就行了,无需清理掉关联的Rsp消息,因为Rsp消息可能未传输成功.
				}
			}
			//////////////////////////////////////////////////////////////////////////
		}
	}
	thls.Unlock()
	return
}

func (thls *safeDataPshCache) calcDataNeedSync() (dst map[string][]*txdata.DataPsh) {
	//TODO:发送某消息后,若一段时间都没有回应的话,需要间隔一段时间(比如30秒)再发送.
	thls.Lock()
	for idx, val := range thls.Slc {
		assert4true(idx == val.slcIdx)
		if val.Ack == nil || val.Ack.ErrNo != 0 { //尚未同步成功.
			if dst == nil {
				dst = make(map[string][]*txdata.DataPsh)
			}
			if dataSlice, isExist := dst[val.Psh.RecverID]; isExist {
				dst[val.Psh.RecverID] = append(dataSlice, val.Psh)
			} else {
				dataSlice = make([]*txdata.DataPsh, 0)
				dst[val.Psh.RecverID] = append(dataSlice, val.Psh)
			}
		} else {
			//TODO:移除出slice.
			//TODO:移除时,脱离slice的那些数据,其slcIdx要设置为负数(比如-1)
		}
	}
	thls.Unlock()
	return
}

//DataPsh2DataAck omit
func DataPsh2DataAck(src *txdata.DataPsh) (dst *txdata.DataAck) {
	return &txdata.DataAck{SenderID: src.RecverID, RecverID: src.SenderID, SendUID: src.SendUID, SendNo: src.SendNo}
}
