package cefingo

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

func (f *CFrameT) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(f))
}

func (self *CFrameT) GetV8context() (context *CV8contextT) {
	c := C.cefingo_frame_get_v8context((*C.cef_frame_t)(self))
	context = (*CV8contextT)(c)
	BaseAddRef(context)
	return context
}

///
// Returns the URL currently loaded in this frame.
///
// The resulting string must be freed by calling cef_string_userfree_free().
func (self *CFrameT) GetUrl() (s string) {
	usfs := C.cefingo_frame_get_url((*C.cef_frame_t)(self))
	s = string_from_cef_string((*C.cef_string_t)(usfs))
	C.cef_string_userfree_free(usfs)
	return s
}
