package cefingo

// import (
// 	"log"
// )

// #include "cefingo.h"
import "C"

func (self *CFrameT) GetV8context() (context *CV8contextT) {
	c := C.cefingo_frame_get_v8context((*C.cef_frame_t)(self))
	context = (*CV8contextT)(c)
	BaseAddRef(context)
	return context
}
