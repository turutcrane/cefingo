package capi

import (
	"fmt"
	"syscall"
	"unsafe"
)

//	#include "cefingo.h"
import "C"

type CWindowHandleT C.cef_window_handle_t

const (
	WinWsOverlapped         = C.WS_OVERLAPPED
	WinWsPopup              = C.WS_POPUP
	WinWsChild              = C.WS_CHILD
	WinWsMinimize           = C.WS_MINIMIZE
	WinWsVisible            = C.WS_VISIBLE
	WinWsDisabled           = C.WS_DISABLED
	WinWsClipsiblings       = C.WS_CLIPSIBLINGS
	WinWsClipchildren       = C.WS_CLIPCHILDREN
	WinWsMaximize           = C.WS_MAXIMIZE
	WinWsCaption            = C.WS_CAPTION
	WinWsBorder             = C.WS_BORDER
	WinWsDlgframe           = C.WS_DLGFRAME
	WinWsVscroll            = C.WS_VSCROLL
	WinWsHscroll            = C.WS_HSCROLL
	WinWsSysmenu            = C.WS_SYSMENU
	WinWsThickframe         = C.WS_THICKFRAME
	WinWsGroup              = C.WS_GROUP
	WinWsTabstop            = C.WS_TABSTOP
	WinWsMinimizebox        = C.WS_MINIMIZEBOX
	WinWsMaximizebox        = C.WS_MAXIMIZEBOX
	WinWsTiled              = C.WS_TILED
	WinWsIconic             = C.WS_ICONIC
	WinWsSizebox            = C.WS_SIZEBOX
	WinWsTiledwindow        = C.WS_TILEDWINDOW
	WinWsOverlappedwindow   = C.WS_OVERLAPPEDWINDOW
	WinWsPopupwindow        = C.WS_POPUPWINDOW
	WinWsChildwindow        = C.WS_CHILDWINDOW
	WinWsExDlgmodalframe    = C.WS_EX_DLGMODALFRAME
	WinWsExNoparentnotify   = C.WS_EX_NOPARENTNOTIFY
	WinWsExTopmost          = C.WS_EX_TOPMOST
	WinWsExAcceptfiles      = C.WS_EX_ACCEPTFILES
	WinWsExTransparent      = C.WS_EX_TRANSPARENT
	WinWsExMdichild         = C.WS_EX_MDICHILD
	WinWsExToolwindow       = C.WS_EX_TOOLWINDOW
	WinWsExWindowedge       = C.WS_EX_WINDOWEDGE
	WinWsExClientedge       = C.WS_EX_CLIENTEDGE
	WinWsExContexthelp      = C.WS_EX_CONTEXTHELP
	WinWsExRight            = C.WS_EX_RIGHT
	WinWsExLeft             = C.WS_EX_LEFT
	WinWsExRtlreading       = C.WS_EX_RTLREADING
	WinWsExLtrreading       = C.WS_EX_LTRREADING
	WinWsExLeftscrollbar    = C.WS_EX_LEFTSCROLLBAR
	WinWsExRightscrollbar   = C.WS_EX_RIGHTSCROLLBAR
	WinWsExControlparent    = C.WS_EX_CONTROLPARENT
	WinWsExStaticedge       = C.WS_EX_STATICEDGE
	WinWsExAppwindow        = C.WS_EX_APPWINDOW
	WinWsExOverlappedwindow = C.WS_EX_OVERLAPPEDWINDOW
	WinWsExPalettewindow    = C.WS_EX_PALETTEWINDOW
	WinWsExLayered          = C.WS_EX_LAYERED
	WinWsExNoinheritlayout  = C.WS_EX_NOINHERITLAYOUT
	// WinWsExNoredirectionbitmap = C.WS_EX_NOREDIRECTIONBITMAP #if WINVER >= 0x0602
	WinWsExLayoutrtl   = C.WS_EX_LAYOUTRTL
	WinWsExComposited  = C.WS_EX_COMPOSITED
	WinWsExNoactivate  = C.WS_EX_NOACTIVATE
	WinWsActivecaption = C.WS_ACTIVECAPTION
)

type WinHinstance C.HINSTANCE
type WinDword C.DWORD
type WinHmenu C.HMENU

func GetWinHandle() WinHinstance {
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

	return WinHinstance(unsafe.Pointer(moduleHandle))
}

func (mainArg *CMainArgsT) SetWinHandle() {
	mainArg.instance = C.HINSTANCE(GetWinHandle())
}

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
