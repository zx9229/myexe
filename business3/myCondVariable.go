package main

import (
	"sync"
	"time"
)

type myCondVariable struct {
	sync.Mutex
	M map[*struct {
		sync.Mutex             //Lock两次产生阻塞的效果
		timer      *time.Timer //超时之用
		isTimeout  bool        //是否超时
	}]bool
}

func newMyCondVariable() *myCondVariable {
	return &myCondVariable{M: make(map[*struct {
		sync.Mutex
		timer     *time.Timer
		isTimeout bool
	}]bool)}
}

func (thls *myCondVariable) wait() {
	elem := new(struct {
		sync.Mutex
		timer     *time.Timer
		isTimeout bool
	})

	thls.Lock()
	thls.M[elem] = true
	elem.Lock()
	thls.Unlock()

	elem.Lock()
	elem.Unlock()
}

func (thls *myCondVariable) waitFor(d time.Duration) bool {
	elem := new(struct {
		sync.Mutex
		timer     *time.Timer
		isTimeout bool
	})

	timeoutFun := func() {
		THLS := thls
		THLS.Lock()
		if _, isOk := THLS.M[elem]; isOk {
			delete(THLS.M, elem)
			elem.isTimeout = true
			elem.Unlock()
		}
		thls.Unlock()
	}
	thls.Lock()
	thls.M[elem] = true
	elem.timer = time.AfterFunc(d, timeoutFun)
	elem.Lock()
	thls.Unlock()
	elem.Lock()
	elem.Unlock()
	return elem.isTimeout
}

func (thls *myCondVariable) notifyOne() {
	thls.Lock()
	for elem := range thls.M {
		delete(thls.M, elem)
		elem.Unlock()
		break
	}
	thls.Unlock()
}

func (thls *myCondVariable) notifyAll() {
	thls.Lock()
	for elem := range thls.M {
		delete(thls.M, elem)
		elem.Unlock()
	}
	thls.Unlock()
}
