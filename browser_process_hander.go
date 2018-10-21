package cefingo

import (
	"log"

	// #include "cefingo.h"
	"C"
)

// BrowserProcessHandler is Go interface of C.cef_browser_process_handler_t
type BrowserProcessHandler interface {
	OnContextInitialized(self *CBrowserProcessHandlerT)
}

var browserProcessHandlers = map[*CBrowserProcessHandlerT]BrowserProcessHandler{}

//export browser_process_on_context_initialized
func browser_process_on_context_initialized(self *CBrowserProcessHandlerT) {
	// Logf("L19: self: %p", self)

	f := browserProcessHandlers[self]
	if f == nil {
		log.Panicln("37: Noo!")
	}

	f.OnContextInitialized(self)
}

type DefaultBrowserProcessHandler struct {
}

func (*DefaultBrowserProcessHandler) OnBeforeClose(self *CBrowserProcessHandlerT) {
	Logf("L79:")
}

func AllocCBrowserProcessHandler(handler BrowserProcessHandler) (cHandler *CBrowserProcessHandlerT) {
	p := C.calloc(1, C.sizeof_cefingo_browser_process_handler_wrapper_t)
	Logf("L39: p: %v", p)
	C.construct_cefingo_browser_process_handler((*C.cefingo_browser_process_handler_wrapper_t)(p))

	cHandler = (*CBrowserProcessHandlerT)(p)
	BaseAddRef(cHandler)
	browserProcessHandlers[cHandler] = handler

	return cHandler
}
