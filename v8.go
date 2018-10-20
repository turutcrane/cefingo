package cefingo

// #include "cefingo.h"
import "C"
import "log"

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

func SetValueBykey(self *CV8valueT, key string, value *CV8valueT) {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)
	status := C.v8context_set_value_bykey((*C.cef_v8value_t)(self),
		key_string, (*C.cef_v8value_t)(value), C.V8_PROPERTY_ATTRIBUTE_NONE)
	if status == 0 {
		log.Panicln("can not set_value_bykey")
	}
}
