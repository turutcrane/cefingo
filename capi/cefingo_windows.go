package capi

import (
	"syscall"
	"unsafe"
)

//	#include "cefingo.h"
import "C"

type CWindowHandleT C.cef_window_handle_t

type WinHinstance C.HINSTANCE
type WinDword C.DWORD
type WinHmenu C.HMENU

func ToCWindowHandleT(hwnd syscall.Handle) CWindowHandleT {
	return CWindowHandleT(unsafe.Pointer(uintptr(hwnd)))
}
