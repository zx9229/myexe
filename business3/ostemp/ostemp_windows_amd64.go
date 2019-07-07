package ostemp

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

//GetStorageSpace omit
func GetStorageSpace(path string) (freeBytes, totalBytes int64, err error) {
	//在Golang中获取系统的磁盘空间内存占用
	//http://wendal.net/2012/1224.html
	for range "1" {
		if path, err = filepath.Abs(path); err != nil {
			break
		}
		path = path[:strings.Index(path, string(os.PathSeparator))]
		var kernel32_dll *syscall.DLL
		if kernel32_dll, err = syscall.LoadDLL("kernel32.dll"); err != nil {
			break
		}
		defer kernel32_dll.Release()
		var GetDiskFreeSpaceExW *syscall.Proc
		if GetDiskFreeSpaceExW, err = kernel32_dll.FindProc("GetDiskFreeSpaceExW"); err != nil {
			break
		}
		var r1, r2 uintptr
		var lpFreeBytesAvailableToCaller, lpTotalNumberOfBytes, lpTotalNumberOfFreeBytes int64
		r1, r2, err = GetDiskFreeSpaceExW.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
			uintptr(unsafe.Pointer(&lpFreeBytesAvailableToCaller)),
			uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
			uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
		if r1 == r2 { //为了规避(declared and not used)而这样写.
		}
		freeBytes = lpFreeBytesAvailableToCaller
		totalBytes = lpTotalNumberOfBytes
	}
	return
}
