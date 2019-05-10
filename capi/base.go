package capi

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

func c_calloc(n C.size_t, s C.size_t, msg string, v ...interface{}) (p unsafe.Pointer) {
	p = C.calloc(n, s)
	if ref_count_log.trace {
		ref_count_log.traceSet[p] = true
		traceuf(1, p, msg, v...)
	}
	if p == nil {
		Panicf("Can not Allocate"+msg, v...)
	}
	return p
}

// func cast_to_base_ref_counted_t(ptr interface{}) (refp *C.cef_base_ref_counted_t) {
// 	var up unsafe.Pointer
// 	switch p := ptr.(type) {
// 	case *CAppT:
// 		up = unsafe.Pointer(p)
// 	case *CBinaryValueT:
// 		up = unsafe.Pointer(p)
// 	case *CBrowserT:
// 		up = unsafe.Pointer(p)
// 	case *CBrowserHostT:
// 		up = unsafe.Pointer(p)
// 	case *CBrowserProcessHandlerT:
// 		up = unsafe.Pointer(p)
// 	case *CClientT:
// 		up = unsafe.Pointer(p)
// 	case *CDictionaryValueT:
// 		up = unsafe.Pointer(p)
// 	case *CFrameT:
// 		up = unsafe.Pointer(p)
// 	case *CLifeSpanHandlerT:
// 		up = unsafe.Pointer(p)
// 	case *CListValueT:
// 		up = unsafe.Pointer(p)
// 	case *CLoadHandlerT:
// 		up = unsafe.Pointer(p)
// 	case *CProcessMessageT:
// 		up = unsafe.Pointer(p)
// 	case *CSchemeHandlerFactoryT:
// 		up = unsafe.Pointer(p)
// 	case *CRequestT:
// 		up = unsafe.Pointer(p)
// 	case *CRenderProcessHandlerT:
// 		up = unsafe.Pointer(p)
// 	case *CResourceHandlerT:
// 		up = unsafe.Pointer(p)
// 	case *CRunFileDialogCallbackT:
// 		up = unsafe.Pointer(p)
// 	case *CValueT:
// 		up = unsafe.Pointer(p)
// 	case *CV8valueT:
// 		up = unsafe.Pointer(p)
// 	case *CV8contextT:
// 		up = unsafe.Pointer(p)
// 	case *CV8arrayBufferReleaseCallbackT:
// 		up = unsafe.Pointer(p)
// 	case *CV8handlerT:
// 		up = unsafe.Pointer(p)
// 	case *CV8exceptionT:
// 		up = unsafe.Pointer(p)
// 	default:
// 		Panicf("Not Refcounted Object: T: %T V: %v", p, p)
// 	}
// 	if up == nil {
// 		Logf("L21: Null passed!")
// 	}
// 	refp = (*C.cef_base_ref_counted_t)(up)
// 	return refp
// }

type refCounted interface {
	cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t
}

func BaseAddRef(rc refCounted) {
	C.cefingo_base_add_ref(rc.cast_to_p_base_ref_counted_t())
}

///
// Called to decrement the reference count for the object. If the reference
// count falls to 0 the object should self-delete. Returns true (1) if the
// resulting reference count is 0.
///
func BaseRelease(rc refCounted) (b bool) {
	status := C.cefingo_base_release(rc.cast_to_p_base_ref_counted_t())
	return status == 1
}

func BaseHasOneRef(rc refCounted) bool {
	status := C.cefingo_base_has_one_ref(rc.cast_to_p_base_ref_counted_t())
	return status == 1
}

func BaseHasAtLeastOneRef(rc refCounted) bool {
	status := C.cefingo_base_has_at_least_one_ref(rc.cast_to_p_base_ref_counted_t())
	return status == 1
}

type Deassocer interface {
	Deassoc()
}

var deassocers = map[unsafe.Pointer][]Deassocer{}

func registerDeassocer(cstruct unsafe.Pointer, d Deassocer) {
	entry, ok := deassocers[cstruct]
	if ok {
		deassocers[cstruct] = append(entry, d)
	} else {
		entry = []Deassocer{d}
		deassocers[cstruct] = entry
	}
}

type DeassocFunc func()

func (f DeassocFunc) Deassoc() {
	f()
}

//export cefingo_base_deassoc
func cefingo_base_deassoc(cstruct unsafe.Pointer) {
	e, ok := deassocers[cstruct]
	if ok {
		Tracef(cstruct, "L160:")
		for _, d := range e {
			d.Deassoc()
		}
		delete(deassocers, cstruct)
	}
}
