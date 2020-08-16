package capi

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

type v8arrayBufferRelaseCallback func(self *CV8arrayBufferReleaseCallbackT, buffer unsafe.Pointer)

func (f v8arrayBufferRelaseCallback) ReleaseBuffer(
	self *CV8arrayBufferReleaseCallbackT,
	buffer unsafe.Pointer,
) {
	f(self, buffer)
}

func CreateArrayBuffer(
	buffer []byte,
) *CV8valueT {
	releaseCallback := AllocCV8arrayBufferReleaseCallbackT().Bind(
		v8arrayBufferRelaseCallback(func(self *CV8arrayBufferReleaseCallbackT, buffer unsafe.Pointer) {
			Logf("T34: %v\n", buffer)
			C.free(buffer)
		}),
	)

	length := len(buffer)
	tmpbuffer := C.CBytes(buffer) // TODO: should be freed in CV8arrayBufferReleaseCallbackT
	Logf("T40: %v\n", tmpbuffer)

	ab := V8valueCreateArrayBuffer(tmpbuffer, int64(length), releaseCallback)

	return ab
}

func V8valueCreateStringFromByteArray(b []byte) (val *CV8valueT) {
	cef_string := create_cef_string_from_byte_array(b)

	v := C.cef_v8value_create_string(cef_string.p_cef_string_t)
	return newCV8valueT(v)
}

// func (self *CV8valueT) GetFunctionName() (s string) {
// 	if self.IsFunction() {
// 		usfs := C.cefingo_v8value_get_function_name(self.p_v8value)
// 		if usfs != nil {
// 			s = string_from_cef_string((*C.cef_string_t)(usfs))
// 			C.cef_string_userfree_free(usfs)
// 		}
// 	}
// 	return s
// }

// func (self *CV8valueT) ExecuteFunction(
// 	this *CV8valueT,
// 	argumentsCount int,
// 	arguments []*CV8valueT,
// ) (val *CV8valueT, err error) {

// 	if !self.IsFunction() {
// 		cause := errors.Errorf("Object is Not Function")
// 		return nil, cause
// 	}
// 	cargs := C.calloc((C.size_t)(argumentsCount), (C.size_t)(unsafe.Sizeof(this)))
// 	slice := (*[1 << 30]*C.cef_v8value_t)(cargs)[:argumentsCount:argumentsCount]

// 	var pThis *C.cef_v8value_t

// 	if this != nil {
// 		BaseAddRef(this.p_v8value)
// 		pThis = this.p_v8value
// 	}
// 	for i, v := range arguments {
// 		BaseAddRef(v.p_v8value)
// 		slice[i] = v.p_v8value
// 	}
// 	v := C.cefingo_v8value_execute_function(
// 		self.p_v8value,
// 		pThis,
// 		(C.size_t)(argumentsCount),
// 		(**C.cef_v8value_t)(cargs))
// 	if v == nil {
// 		name := self.GetFunctionName()
// 		if self.HasException() {
// 			e := self.GetException()
// 			m := e.GetMessage()
// 			err = errors.Errorf("E363: %s returns NULL and %s has Exception: %s", name, name, m)
// 		} else if this != nil && this.HasException() {
// 			e := this.GetException()
// 			m := e.GetMessage()
// 			err = errors.Errorf("E367: %s returns NULL and (this) has Exception: %s", name, m)
// 		} else {
// 			err = errors.Errorf("E369: %s returns NULL", name)
// 		}
// 	} else if status := C.cefingo_v8value_is_valid(v); status == 1 {
// 		val = newCV8valueT(v)
// 	}
// 	return val, err
// }

// // Excute Handler
// //export cefingo_v8handler_execute
// func cefingo_v8handler_execute(self *C.cef_v8handler_t,
// 	name *C.cef_string_t,
// 	object *C.cef_v8value_t,
// 	argumentsCount C.size_t,
// 	arguments **C.cef_v8value_t,
// 	retval **C.cef_v8value_t,
// 	exception *C.cef_string_t,
// ) (ret C.int) {
// 	goname := string_from_cef_string(name)
// 	v8handlers.m.Lock()
// 	handler := v8handlers.v8handler[self]
// 	v8handlers.m.Unlock()
// 	if handler == nil {
// 		Logf("L121: No V8 Execute Handler")
// 		ret = 0
// 	} else {
// 		var slice []*CV8valueT
// 		if arguments != nil {
// 			s := (*[1 << 30]*C.cef_v8value_t)(unsafe.Pointer(arguments))[:argumentsCount:argumentsCount]
// 			for _, v := range s {
// 				slice = append(slice, newCV8valueT(v))
// 			}
// 		}

// 		runtime.LockOSThread()
// 		defer runtime.UnlockOSThread()

// 		var exc string
// 		var v8v *CV8valueT
// 		if handler.Execute(newCV8handlerT(self), goname,
// 			newCV8valueT(object),
// 			(int)(argumentsCount), slice, &v8v, &exc) {
// 			// Is required release of member of arguments ?
// 			if v8v != nil {
// 				*retval = v8v.p_v8value
// 				BaseAddRef(*retval)
// 			}
// 			set_cef_string(exception, exc)
// 			ret = 1
// 		} else {
// 			ret = 0
// 		}
// 	}
// 	return ret
// }

// func (self *CV8contextT) Eval(code string, retval **CV8valueT, e **CV8exceptionT) (ret bool) {
// 	s := create_cef_string(code)
// 	defer clear_cef_string(s)
// 	var r *C.cef_v8value_t
// 	var exc *C.cef_v8exception_t
// 	status := C.cefingo_v8context_eval(
// 		self.p_v8context, s, nil, 0,
// 		&r, &exc)
// 	ret = (status == 1)
// 	if ret {
// 		*retval = newCV8valueT(r)
// 	} else {
// 		// if exc != nil {
// 		*e = newCV8exceptionT(exc)
// 		// BaseAddRef(*e)
// 		// }
// 	}

// 	return ret
// }
