package capi

// #include "cefingo.h"
import "C"

func (self *CCallbackT) Cont() {
	C.cefingo_callback_cont((*C.cef_callback_t)(self))
}

func (self *CCallbackT) Cancel() {
	C.cefingo_callback_cancel((*C.cef_callback_t)(self))
}
