package capi

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

func c_calloc(n C.size_t, s C.size_t, msg string, v ...interface{}) (p unsafe.Pointer) {
	p = C.malloc(n * s) // never returns nil
	C.memset(p, 0, n * s)

	if ref_count_log.trace {
		ref_count_log.traceSet[p] = true
		traceuf(1, p, msg, v...)
	}
	return p
}

type refCounted interface {
	cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t
}

func BaseAddRef(rc refCounted) {
	if rc != nil {
		C.cefingo_base_add_ref(rc.cast_to_p_base_ref_counted_t())
	}
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
		Tracef(cstruct, "T135:")
		for _, d := range e {
			d.Deassoc()
		}
		delete(deassocers, cstruct)
	}
}
