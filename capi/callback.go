package capi

// #include "cefingo.h"
import "C"

func (self *CCallbackT) Cont() {
	C.cefingo_callback_cont(self.p_callback)
}

func (self *CCallbackT) Cancel() {
	C.cefingo_callback_cancel(self.p_callback)
}
