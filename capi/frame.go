package capi

import (
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

func (f *C.cef_frame_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(f))
}

func newCFrameT(cef *C.cef_frame_t) *CFrameT {
	Tracef(unsafe.Pointer(cef), "L15:")
	BaseAddRef(cef)
	frame := CFrameT{cef}
	runtime.SetFinalizer(&frame, func(f *CFrameT) {
		Tracef(unsafe.Pointer(f.p_frame), "L19:")
		BaseRelease(f.p_frame)
	})
	return &frame
}

func (self *CFrameT) GetV8context() (context *CV8contextT) {
	c := C.cefingo_frame_get_v8context(self.p_frame)
	context = newCV8contextT(c)
	return context
}

///
// Returns the URL currently loaded in this frame.
///
// The resulting string must be freed by calling cef_string_userfree_free().
func (self *CFrameT) GetUrl() (s string) {
	usfs := C.cefingo_frame_get_url(self.p_frame)
	s = string_from_cef_string((*C.cef_string_t)(usfs))
	C.cef_string_userfree_free(usfs)
	return s
}
