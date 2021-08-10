package capi

// #include "cefingo.h"
import "C"

///
// Retrieve the post data elements.
///
// func (self *CPostDataT) GetElements() []*CPostDataElementT {

// 	n := self.GetElementCount()
// 	tmpelements := c_calloc(C.size_t(n), C.sizeof_ptr_cef_post_data_element_t, "T10.1:cef_post_data_t::get_elements::elements")
// 	defer C.free(tmpelements)
// 	slice := (*[1 << 30]*C.cef_post_data_element_t)(tmpelements)[:n:n]

// 	np := unsafe.Pointer(&n)
// 	C.cefingo_post_data_get_elements((*C.cef_post_data_t)(self.pc_post_data), (*C.size_t)(np), (**C.cef_post_data_element_t)(tmpelements))

// 	ret := []*CPostDataElementT{}
// 	for _, e := range slice {
// 		ret = append(ret, newCPostDataElementT(e, false))
// 	}
// 	return ret
// }
