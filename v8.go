package cefingo

import (
	"log"
	"unsafe"
)

// #include "cefingo.h"
import "C"

type V8handler interface {
	///
	// Handle execution of the function identified by |name|. |object| is the
	// receiver ('this' object) of the function. |arguments| is the list of
	// arguments passed to the function. If execution succeeds set |retval| to the
	// function return value. If execution fails set |exception| to the exception
	// that will be thrown. Return true (1) if execution was handled.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_v8_capi.h#L158-L183
	///
	Execute(self *CV8handlerT,
		name string,
		object *CV8valueT,
		argumentsCount int,
		arguments []*CV8valueT,
		retval **CV8valueT,
		exception *CStringT) bool
}

var v8handlerMap = map[*CV8handlerT]V8handler{}

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
	Logf("L46: p:%v", buffer)
	// C.free(buffer)
}

func (self *CV8contextT) GetGlobal() *CV8valueT {
	g := (*CV8valueT)(C.v8context_get_global((*C.cef_v8context_t)(self)))
	BaseAddRef(g)
	return g
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

func (self *CV8valueT) SetValueBykey(key string, value *CV8valueT) {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	BaseAddRef(value)
	status := C.v8context_set_value_bykey((*C.cef_v8value_t)(self),
		key_string, (*C.cef_v8value_t)(value), C.V8_PROPERTY_ATTRIBUTE_NONE)
	if status == 0 {
		log.Panicln("can not set_value_bykey")
	}
}

func (self *CV8valueT) HasValueBykey(key string) bool {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	status := C.v8context_has_value_bykey((*C.cef_v8value_t)(self), key_string)
	return (status == 1)
}

func (self *CV8valueT) GetValueBykey(key string) (value *CV8valueT) {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	value = (*CV8valueT)(C.v8context_get_value_bykey((*C.cef_v8value_t)(self), key_string))
	BaseAddRef(value)
	return value
}

func (self *CV8valueT) IsFunction() bool {
	status := C.cefingo_v8value_is_function((*C.cef_v8value_t)(self))
	return status == 1
}

// V8valueCreateFunction create V8 function
///
// Create a new cef_v8value_t object of type function. This function should only
// be called from within the scope of a cef_render_process_handler_t,
// cef_v8handler_t or cef_v8accessor_t callback, or in combination with calling
// enter() and exit() on a stored cef_v8context_t reference.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_v8_capi.h#L812-L819
///
func V8valueCreateFunction(name string, handler *CV8handlerT) (function *CV8valueT) {
	cef_name := create_cef_string(name)
	defer clear_cef_string(cef_name)

	BaseAddRef(handler)
	return (*CV8valueT)(C.cef_v8value_create_function(cef_name, (*C.cef_v8handler_t)(handler)))
}

// AllocCV8handlerT allocates CV8handlerT and construct it
func AllocCV8handlerT(handler V8handler) (v8handler *CV8handlerT) {
	p := C.calloc(1, C.sizeof_cefingo_v8handler_wrapper_t)
	Logf("L22: p: %v", p)

	hp := (*C.cefingo_v8handler_wrapper_t)(p)
	C.construct_cefingo_v8handler(hp)

	v8handler = (*CV8handlerT)(p)
	BaseAddRef(v8handler)
	v8handlerMap[v8handler] = handler

	return v8handler
}

// Excute Handler
//export execute
func execute(self *CV8handlerT,
	name *CStringT,
	object *CV8valueT,
	argumentsCount C.size_t,
	arguments **CV8valueT,
	retval **CV8valueT,
	exception *CStringT,
) (ret C.int) {

	goname := C.GoString((*C.char)(unsafe.Pointer(name.str)))
	handler := v8handlerMap[self]
	if handler == nil {
		Logf("L121: No V8 Execute Handler")
		ret = 0
	} else {
		var slice []*CV8valueT
		if arguments != nil {
			slice = (*[1<<30]*CV8valueT)(unsafe.Pointer(arguments))[:argumentsCount:argumentsCount]
		}
		if handler.Execute(self, goname, object, (int)(argumentsCount), slice, retval, exception) {
			// Is required release of member of arguments ?
			ret = 1
		} else {
			ret = 0
		}
	}
	return ret
}

func V8contextInContext() bool {
	// Returns true (1) if V8 is currently inside a context.
	inContext := C.cef_v8context_in_context()
	Logf("L150: %d", inContext)
	return (inContext == 1)
}

func (self *CV8contextT) Enter() bool {
	c := C.cefingo_v8context_enter((*C.cef_v8context_t)(self))
	return (c == 1)
}

func (self *CV8contextT) Exit() bool {
	c := C.cefingo_v8context_exit((*C.cef_v8context_t)(self))
	return (c == 1)
}
