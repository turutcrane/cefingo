package cefingo

import (
	"runtime"
	"unsafe"

	"github.com/pkg/errors"
)

// #include "cefingo.h"
// #include "cefingo_v8.h"
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
		exception *string) bool
}

var v8handlerMap = map[*CV8handlerT]V8handler{}

// AllocCV8arrayBufferReleaseCallbackT allocates CV8arrayBufferReleaseCallbackT and construct it
func AllocCV8arrayBufferReleaseCallbackT() (cv8ab_release_callback *CV8arrayBufferReleaseCallbackT) {
	p := C.calloc(1, C.sizeof_cefingo_v8array_buffer_release_callback_wrapper_t)
	Logf("L22: p: %v", p)

	C.cefingo_construct_v8array_buffer_release_callback((*C.cefingo_v8array_buffer_release_callback_wrapper_t)(p))

	cv8ab_release_callback = (*CV8arrayBufferReleaseCallbackT)(p)
	BaseAddRef(cv8ab_release_callback)

	return cv8ab_release_callback
}

//export cefingo_v8array_buffer_release_callback_release_buffer
func cefingo_v8array_buffer_release_callback_release_buffer(self *C.cef_v8array_buffer_release_callback_t, buffer unsafe.Pointer) {
	Logf("L46: p:%v", buffer)
	// C.free(buffer)
}

func (self *CV8contextT) GetGlobal() *CV8valueT {
	g := (*CV8valueT)(C.cefingo_v8context_get_global((*C.cef_v8context_t)(self)))
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

func (self *CV8valueT) SetValueBykey(key string, value *CV8valueT) bool {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	BaseAddRef(value)
	status := C.cefingo_v8context_set_value_bykey((*C.cef_v8value_t)(self),
		key_string, (*C.cef_v8value_t)(value), C.V8_PROPERTY_ATTRIBUTE_NONE)
	return status == 1
}

func (self *CV8valueT) HasValueBykey(key string) bool {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	status := C.cefingo_v8value_has_value_bykey((*C.cef_v8value_t)(self), key_string)
	return (status == 1)
}

func (self *CV8valueT) GetValueBykey(key string) (value *CV8valueT) {
	key_string := create_cef_string(key)
	defer clear_cef_string(key_string)

	value = (*CV8valueT)(C.cefingo_v8value_get_value_bykey((*C.cef_v8value_t)(self), key_string))
	BaseAddRef(value)
	return value
}

func (self *CV8valueT) IsValid() bool {
	status := C.cefingo_v8value_is_valid((*C.cef_v8value_t)(self))
	return status == 1
}

func (self *CV8valueT) IsFunction() bool {
	status := C.cefingo_v8value_is_function((*C.cef_v8value_t)(self))
	return status == 1
}

func (self *CV8valueT) IsString() bool {
	status := C.cefingo_v8value_is_string((*C.cef_v8value_t)(self))
	return status == 1
}

func (self *CV8valueT) IsObject() bool {
	status := C.cefingo_v8value_is_object((*C.cef_v8value_t)(self))
	return status == 1
}

func (self *CV8valueT) GetString() string {
	usfs := C.cefingo_v8value_get_string((*C.cef_v8value_t)(self))
	s := string_from_cef_string((*C.cef_string_t)(usfs))
	C.cef_string_userfree_free(usfs)
	return s
}

func (self *CV8valueT) GetFunctionName() string {
	usfs := C.cefingo_v8value_get_function_name((*C.cef_v8value_t)(self))
	s := string_from_cef_string((*C.cef_string_t)(usfs))
	C.cef_string_userfree_free(usfs)
	return s
}

func (self *CV8valueT) ExecuteFunction(
	this *CV8valueT,
	argumentsCount int,
	arguments []*CV8valueT,
) (v *CV8valueT, err error) {

	if !self.IsFunction() {
		cause := errors.Errorf("Object is Not Function")
		return nil, cause
	}
	ca := C.calloc((C.size_t)(argumentsCount), (C.size_t)(unsafe.Sizeof(this)))
	slice := (*[1 << 30]*CV8valueT)(ca)[:argumentsCount:argumentsCount]

	BaseAddRef(this)
	for i, v := range arguments {
		BaseAddRef(v)
		slice[i] = v
	}
	v = (*CV8valueT)(C.cefingo_v8value_execute_function(
		(*C.cef_v8value_t)(self),
		(*C.cef_v8value_t)(this),
		(C.size_t)(argumentsCount),
		(**C.cef_v8value_t)(ca)))
	if v == nil {
		name := self.GetFunctionName()
		err = errors.Errorf("%s returns NULL", name)
	}
	return v, err
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
	C.cefingo_construct_v8handler(hp)

	v8handler = (*CV8handlerT)(p)
	BaseAddRef(v8handler)
	v8handlerMap[v8handler] = handler

	return v8handler
}

// Excute Handler
//export cefingo_v8handler_execute
func cefingo_v8handler_execute(self *CV8handlerT,
	name *C.cef_string_t,
	object *CV8valueT,
	argumentsCount C.size_t,
	arguments **CV8valueT,
	retval **CV8valueT,
	exception *C.cef_string_t,
) (ret C.int) {
	goname := string_from_cef_string(name)
	handler := v8handlerMap[self]

	if handler == nil {
		Logf("L121: No V8 Execute Handler")
		ret = 0
	} else {
		var slice []*CV8valueT
		if arguments != nil {
			slice = (*[1 << 30]*CV8valueT)(unsafe.Pointer(arguments))[:argumentsCount:argumentsCount]
		}
		runtime.LockOSThread()
		var exc string
		if handler.Execute(self, goname, object, (int)(argumentsCount), slice, retval, &exc) {
			// Is required release of member of arguments ?
			set_cef_string(exception, exc)
			ret = 1
		} else {
			ret = 0
		}
		runtime.UnlockOSThread()
	}
	return ret
}

func V8contextInContext() bool {
	// Returns true (1) if V8 is currently inside a context.
	inContext := C.cef_v8context_in_context()
	Logf("L150: %d", inContext)
	return (inContext == 1)
}

func V8contextGetEnterdContext() (context *CV8contextT) {
	c := C.cef_v8context_get_entered_context()
	context = (*CV8contextT)(c)
	BaseAddRef(context)
	return context
}

func (self *CV8contextT) Enter() bool {
	runtime.LockOSThread()
	c := C.cefingo_v8context_enter((*C.cef_v8context_t)(self))
	return (c == 1)
}

func (self *CV8contextT) Exit() bool {
	c := C.cefingo_v8context_exit((*C.cef_v8context_t)(self))
	runtime.UnlockOSThread()
	return (c == 1)
}

func (self *CV8contextT) IsSame(that *CV8contextT) bool {
	BaseAddRef(that)
	s := C.cefingo_v8context_is_same(
		(*C.cef_v8context_t)(self),
		(*C.cef_v8context_t)(that),
	)
	return s == 1
}

func (self *CV8contextT) EvalString(code string, retval **CV8valueT, e **CV8exceptionT) bool {
	s := create_cef_string(code)
	defer clear_cef_string(s)
	var r *C.cef_v8value_t
	var exc *C.cef_v8exception_t
	status := C.cefingo_v8context_eval(
		(*C.cef_v8context_t)(self), s, nil, 0,
		&r, &exc)
	*retval = (*CV8valueT)(r)
	return status == 1
}
