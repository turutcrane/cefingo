package cefingo

import (
	"log"
	"unsafe"
)

// #include "cefingo.h"
import "C"

// AllocCV8arrayBufferReleaseCallbackT allocates CV8arrayBufferReleaseCallbackT and construct it
func AllocCV8arrayBufferReleaseCallbackT() (cv8ab_release_callback *CV8arrayBufferReleaseCallbackT) {
	p := C.calloc(1, C.sizeof_cefingo_v8array_buffer_release_callback_wrapper_t)
	Logf("L22: p: %v", p)

	C.construct_cefingo_v8array_buffer_release_callback((*C.cefingo_v8array_buffer_release_callback_wrapper_t)(p))

	cv8ab_release_callback = (*CV8arrayBufferReleaseCallbackT)(p)
	BaseAddRef(cv8ab_release_callback)

	return cv8ab_release_callback
}

//export v8array_buffer_release_buffer
func v8array_buffer_release_buffer(self *C.cef_v8array_buffer_release_callback_t, buffer unsafe.Pointer) {
	Logf("L26: p:%", buffer)
	// C.free(buffer)
}

func GetGlobal(self *CV8contextT) *CV8valueT {
	return (*CV8valueT)(C.v8context_get_global((*C.cef_v8context_t)(self)))
}

func V8valueCreateString(s string) *CV8valueT {
	cef_string := create_cef_string(s)
	defer clear_cef_string(cef_string)

	return (*CV8valueT)(C.cef_v8value_create_string(cef_string))
}

func V8valueCreateObject(accessor *CV8accessorT, interceptor *CV8interceptorT) *CV8valueT {
	return (*CV8valueT)(C.cef_v8value_create_object(
		(*C.cef_v8accessor_t)(accessor), (*C.cef_v8interceptor_t)(interceptor)))
}

func V8valueCreateArrayBuffer(buffer []byte) *CV8valueT {
	release_callback := AllocCV8arrayBufferReleaseCallbackT()

	// buf := [100]byte{}
	cbytes := C.CBytes(buffer[:])
	buffer_len := (C.size_t)(len(buffer[:]))
	Logf("L31: %v, %v, %v", cbytes, len(buffer[:]), buffer_len)
	v8array_buffer := C.cef_v8value_create_array_buffer(
		cbytes,
		buffer_len,
		(*C.cef_v8array_buffer_release_callback_t)(release_callback),
	)
	return (*CV8valueT)(v8array_buffer)
}

func SetValueBykey(self *CV8valueT, key string, value *CV8valueT) {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	BaseAddRef(value)
	status := C.v8context_set_value_bykey((*C.cef_v8value_t)(self),
		key_string, (*C.cef_v8value_t)(value), C.V8_PROPERTY_ATTRIBUTE_NONE)
	if status == 0 {
		log.Panicln("can not set_value_bykey")
	}
}
