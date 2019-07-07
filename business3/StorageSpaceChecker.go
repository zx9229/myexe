package main

import (
	"time"

	"github.com/zx9229/myexe/business3/ostemp"
)

//StorageSpaceChecker omit
type StorageSpaceChecker struct {
	path      string
	freeBytes int64
	threshold int64
}

func newStorageSpaceChecker(path string, thresholdBytes int64) *StorageSpaceChecker {
	curData := new(StorageSpaceChecker)
	curData.path = path
	curData.threshold = thresholdBytes
	curData.refreshData()
	return curData
}

func (thls *StorageSpaceChecker) refreshData() {
	go func() {
		for {
			thls.freeBytes, _, _ = ostemp.GetStorageSpace(thls.path)
			time.Sleep(time.Second * 60)
		}
	}()
}

func (thls *StorageSpaceChecker) getFreeBytes() int64 {
	return thls.freeBytes
}

func (thls *StorageSpaceChecker) getFreeKB() int64 {
	return thls.freeBytes / 1024
}

func (thls *StorageSpaceChecker) available() bool {
	return thls.threshold < thls.freeBytes
}
