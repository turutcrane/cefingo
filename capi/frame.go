package capi

// #include "cefingo.h"
import "C"

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
