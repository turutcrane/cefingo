package cefingo

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

func RefCountLogOutput(enable bool) {
	if enable {
		C.REF_COUNT_LOG_OUTPUT = C.TRUE
	} else {
		C.REF_COUNT_LOG_OUTPUT = C.FALSE
	}
}

func cast_to_base_ref_counted_t(ptr interface{}) (refp *C.cef_base_ref_counted_t) {
	var up unsafe.Pointer
	switch p := ptr.(type) {
	case *CAppT:
		up = unsafe.Pointer(p)
	case *CBrowserProcessHandlerT:
		up = unsafe.Pointer(p)
	case *CClientT:
		up = unsafe.Pointer(p)
	case *CLifeSpanHandlerT:
		up = unsafe.Pointer(p)
	case *CRenderProcessHandlerT:
		up = unsafe.Pointer(p)
	case *CV8valueT:
		up = unsafe.Pointer(p)
	case *CV8contextT:
		up = unsafe.Pointer(p)
	case *CV8arrayBufferReleaseCallbackT:
		up = unsafe.Pointer(p)
	case *CV8handlerT:
		up = unsafe.Pointer(p)
	case *CLoadHandlerT:
		up = unsafe.Pointer(p)
	case *CSchemeHandlerFactoryT:
		up = unsafe.Pointer(p)
	case *CResourceHandlerT:
		up = unsafe.Pointer(p)
	case *CRequestT:
		up = unsafe.Pointer(p)
	case *CV8exceptionT:
		up = unsafe.Pointer(p)
	default:
		Panicf("Not Refcounted Object: T: %T V: %v", p, p)
	}
	if up == nil {
		Logf("L21: Null passed!")
	}
	refp = (*C.cef_base_ref_counted_t)(up)
	return refp
}

func BaseAddRef(ptr interface{}) {
	C.cefingo_base_add_ref(cast_to_base_ref_counted_t(ptr))
}

///
// Called to decrement the reference count for the object. If the reference
// count falls to 0 the object should self-delete. Returns true (1) if the
// resulting reference count is 0.
///
func BaseRelease(ptr interface{}) (b bool) {
	status := C.cefingo_base_release(cast_to_base_ref_counted_t(ptr))
	return status == 1
}

func BaseHasOneRef(ptr interface{}) bool {
	status := C.cefingo_base_has_one_ref(cast_to_base_ref_counted_t(ptr))
	return status == 1
}
