package cefingo

import (
	"unsafe"

	// #include "cefingo.h"
	"C"
)

type CBaseRefCountedT C.cef_base_ref_counted_t

func BaseAddRef(p unsafe.Pointer) {
	C.cefingo_add_ref((*C.cef_base_ref_counted_t)(p))
}

func BaseRelease(p unsafe.Pointer) {
	C.cefingo_release((*C.cef_base_ref_counted_t)(p))
}

func BaseHasOneRef(p unsafe.Pointer) {
	C.cefingo_has_one_ref((*C.cef_base_ref_counted_t)(p))
}
