package cef

import (
	"log"
	"unsafe"

	"github.com/turutcrane/cefingo/capi"
	"github.com/turutcrane/win32api"
)

func CMainArgsTSetInstance(mainArgs *capi.CMainArgsT) {
	h, err := win32api.GetModuleHandle(nil)
	if err != nil {
		log.Panicf("T92: %v", err)
	}
	mainArgs.SetInstance(capi.WinHinstance(unsafe.Pointer(h)))
}
