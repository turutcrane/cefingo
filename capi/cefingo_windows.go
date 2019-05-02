package capi

import (
	"fmt"
	"syscall"
	"unsafe"
)

//	#include "cefingo.h"
import "C"

func setup_main_args(ma *C.cef_main_args_t) {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	// psapi := syscall.MustLoadDLL("psapi.dll")
	getModuleHandle := kernel32.MustFindProc("GetModuleHandleW")
	// getCurrentProcess := kernel32.MustFindProc("GetCurrentProcess")
	// getModuleInformation := psapi.MustFindProc("GetModuleInformation")

	// procHandle, _, _ := getCurrentProcess.Call()
	moduleHandle, _, err := getModuleHandle.Call(0)
	if moduleHandle == 0 {
		panic(fmt.Sprintf("GetModuleHandle() failed: %d", err))
	}

	ma.instance = C.HINSTANCE(unsafe.Pointer(moduleHandle))
}
