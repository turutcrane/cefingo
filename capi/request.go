package capi

// #include "cefingo.h"
import "C"

func (r *CRequestT) GetUrl() string {

	ufs := C.cefingo_request_get_url(r.p_request)
	s := string_from_cef_string((*C.cef_string_t)(ufs))
	C.cef_string_userfree_free(ufs)
	return s
}
