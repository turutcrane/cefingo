package capi

// #include "cefingo.h"
import "C"

func (self *CResponseT) DumpHeaders() {
	m := C.cef_string_multimap_alloc()
	defer C.cef_string_multimap_free(m)

	C.cefingo_response_get_header_map(self.p_response, m)
	size := C.cef_string_multimap_size(m)
	Logf("T49: size:%d", size)
	for i := C.size_t(0); i < size; i++ {
		k := C.cef_string_t{}
		C.cef_string_multimap_key(m, i, &k)
		key := string_from_cef_string(&k)
		Logf("T53: %d key:%s", i, key)
	}

}
