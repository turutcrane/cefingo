package cefingo

import (
	"log"
	"unsafe"

	// #include "cefingo.h"
	"C"
)

func RefCountLogOutput(enable bool) {
	if enable {
		C.REF_COUNT_LOG_OUTPUT = C.TRUE
	} else {
		C.REF_COUNT_LOG_OUTPUT = C.FALSE
	}
};

func cast_to_base_ref_counted_t(ptr interface{}) (cp *C.cef_base_ref_counted_t) {
	switch p := ptr.(type) {
	case *CAppT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	case *CBrowserProcessHandlerT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	case *CClientT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	case *CLifeSpanHandlerT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	case *CRenderProcessHandlerT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	case *CV8valueT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	case *CV8contextT:
		cp = (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
	default:
		log.Panicf("Not Refcounted Object: T: %t V: %v", p, p)
	}
	return cp
}

func BaseAddRef(ptr interface{}) {
	C.cefingo_base_add_ref(cast_to_base_ref_counted_t(ptr))
}

func BaseRelease(ptr interface{}) Cint {
	status := C.cefingo_base_release(cast_to_base_ref_counted_t(ptr))
	return Cint(status)
}

func BaseHasOneRef(ptr interface{}) Cint {
	status := C.cefingo_base_has_one_ref(cast_to_base_ref_counted_t(ptr))
	return Cint(status)
}
