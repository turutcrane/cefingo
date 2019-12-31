// Code generated by "gen_cef_types.go" DO NOT EDIT.

package capi

// "runtime"
// "unsafe"

// #include "cefingo.h"
import "C"

// Go type for cef_v8array_buffer_release_callback_t
// type CV8arrayBufferReleaseCallbackT struct {
// 	p_v8array_buffer_release_callback *C.cef_v8array_buffer_release_callback_t
// }

// type RefToCV8arrayBufferReleaseCallbackT struct {
// 	p_v8array_buffer_release_callback *CV8arrayBufferReleaseCallbackT
// }

// type CV8arrayBufferReleaseCallbackTAccessor interface {
// 	GetCV8arrayBufferReleaseCallbackT() *CV8arrayBufferReleaseCallbackT
// 	SetCV8arrayBufferReleaseCallbackT(*CV8arrayBufferReleaseCallbackT)
// }

// func (r RefToCV8arrayBufferReleaseCallbackT) GetCV8arrayBufferReleaseCallbackT() *CV8arrayBufferReleaseCallbackT {
// 	return r.p_v8array_buffer_release_callback
// }

// func (r *RefToCV8arrayBufferReleaseCallbackT) SetCV8arrayBufferReleaseCallbackT(p *CV8arrayBufferReleaseCallbackT) {
// 	r.p_v8array_buffer_release_callback = p
// }

// // Go type CV8arrayBufferReleaseCallbackT wraps cef type *C.cef_v8array_buffer_release_callback_t
// func newCV8arrayBufferReleaseCallbackT(p *C.cef_v8array_buffer_release_callback_t) *CV8arrayBufferReleaseCallbackT {
// 	Tracef(unsafe.Pointer(p), "T778:")
// 	BaseAddRef(p)
// 	go_v8array_buffer_release_callback := CV8arrayBufferReleaseCallbackT{p}
// 	runtime.SetFinalizer(&go_v8array_buffer_release_callback, func(g *CV8arrayBufferReleaseCallbackT) {
// 		Tracef(unsafe.Pointer(g.p_v8array_buffer_release_callback), "T782:")
// 		BaseRelease(g.p_v8array_buffer_release_callback)
// 	})
// 	return &go_v8array_buffer_release_callback
// }

// // *C.cef_v8array_buffer_release_callback_t has refCounted interface
// func (p *C.cef_v8array_buffer_release_callback_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
// 	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
// }

// Go type for cef_v8context_t
// type CV8contextT struct {
// 	p_v8context *C.cef_v8context_t
// }

// type RefToCV8contextT struct {
// 	p_v8context *CV8contextT
// }

// type CV8contextTAccessor interface {
// 	GetCV8contextT() *CV8contextT
// 	SetCV8contextT(*CV8contextT)
// }

// func (r RefToCV8contextT) GetCV8contextT() *CV8contextT {
// 	return r.p_v8context
// }

// func (r *RefToCV8contextT) SetCV8contextT(p *CV8contextT) {
// 	r.p_v8context = p
// }

// // Go type CV8contextT wraps cef type *C.cef_v8context_t
// func newCV8contextT(p *C.cef_v8context_t) *CV8contextT {
// 	Tracef(unsafe.Pointer(p), "T817:")
// 	BaseAddRef(p)
// 	go_v8context := CV8contextT{p}
// 	runtime.SetFinalizer(&go_v8context, func(g *CV8contextT) {
// 		Tracef(unsafe.Pointer(g.p_v8context), "T821:")
// 		BaseRelease(g.p_v8context)
// 	})
// 	return &go_v8context
// }

// // *C.cef_v8context_t has refCounted interface
// func (p *C.cef_v8context_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
// 	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
// }

// Go type for cef_v8exception_t
// type CV8exceptionT struct {
// 	p_v8exception *C.cef_v8exception_t
// }

// type RefToCV8exceptionT struct {
// 	p_v8exception *CV8exceptionT
// }

// type CV8exceptionTAccessor interface {
// 	GetCV8exceptionT() *CV8exceptionT
// 	SetCV8exceptionT(*CV8exceptionT)
// }

// func (r RefToCV8exceptionT) GetCV8exceptionT() *CV8exceptionT {
// 	return r.p_v8exception
// }

// func (r *RefToCV8exceptionT) SetCV8exceptionT(p *CV8exceptionT) {
// 	r.p_v8exception = p
// }

// // Go type CV8exceptionT wraps cef type *C.cef_v8exception_t
// func newCV8exceptionT(p *C.cef_v8exception_t) *CV8exceptionT {
// 	Tracef(unsafe.Pointer(p), "T856:")
// 	BaseAddRef(p)
// 	go_v8exception := CV8exceptionT{p}
// 	runtime.SetFinalizer(&go_v8exception, func(g *CV8exceptionT) {
// 		Tracef(unsafe.Pointer(g.p_v8exception), "T860:")
// 		BaseRelease(g.p_v8exception)
// 	})
// 	return &go_v8exception
// }

// // *C.cef_v8exception_t has refCounted interface
// func (p *C.cef_v8exception_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
// 	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
// }

// Go type for cef_v8handler_t
// type CV8handlerT struct {
// 	p_v8handler *C.cef_v8handler_t
// }

// type RefToCV8handlerT struct {
// 	p_v8handler *CV8handlerT
// }

// type CV8handlerTAccessor interface {
// 	GetCV8handlerT() *CV8handlerT
// 	SetCV8handlerT(*CV8handlerT)
// }

// func (r RefToCV8handlerT) GetCV8handlerT() *CV8handlerT {
// 	return r.p_v8handler
// }

// func (r *RefToCV8handlerT) SetCV8handlerT(p *CV8handlerT) {
// 	r.p_v8handler = p
// }

// // Go type CV8handlerT wraps cef type *C.cef_v8handler_t
// func newCV8handlerT(p *C.cef_v8handler_t) *CV8handlerT {
// 	Tracef(unsafe.Pointer(p), "T895:")
// 	BaseAddRef(p)
// 	go_v8handler := CV8handlerT{p}
// 	runtime.SetFinalizer(&go_v8handler, func(g *CV8handlerT) {
// 		Tracef(unsafe.Pointer(g.p_v8handler), "T899:")
// 		BaseRelease(g.p_v8handler)
// 	})
// 	return &go_v8handler
// }

// // *C.cef_v8handler_t has refCounted interface
// func (p *C.cef_v8handler_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
// 	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
// }

// // Go type for cef_v8value_t
// type CV8valueT struct {
// 	p_v8value *C.cef_v8value_t
// }

// type RefToCV8valueT struct {
// 	p_v8value *CV8valueT
// }

// type CV8valueTAccessor interface {
// 	GetCV8valueT() *CV8valueT
// 	SetCV8valueT(*CV8valueT)
// }

// func (r RefToCV8valueT) GetCV8valueT() *CV8valueT {
// 	return r.p_v8value
// }

// func (r *RefToCV8valueT) SetCV8valueT(p *CV8valueT) {
// 	r.p_v8value = p
// }

// // Go type CV8valueT wraps cef type *C.cef_v8value_t
// func newCV8valueT(p *C.cef_v8value_t) *CV8valueT {
// 	Tracef(unsafe.Pointer(p), "T934:")
// 	BaseAddRef(p)
// 	go_v8value := CV8valueT{p}
// 	runtime.SetFinalizer(&go_v8value, func(g *CV8valueT) {
// 		Tracef(unsafe.Pointer(g.p_v8value), "T938:")
// 		BaseRelease(g.p_v8value)
// 	})
// 	return &go_v8value
// }

// // *C.cef_v8value_t has refCounted interface
// func (p *C.cef_v8value_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
// 	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
// }

// Go type for cef_v8stack_trace_t
// type CV8stackTraceT struct {
// 	p_v8stack_trace *C.cef_v8stack_trace_t
// }

// type RefToCV8stackTraceT struct {
// 	p_v8stack_trace *CV8stackTraceT
// }

// type CV8stackTraceTAccessor interface {
// 	GetCV8stackTraceT() *CV8stackTraceT
// 	SetCV8stackTraceT(*CV8stackTraceT)
// }

// func (r RefToCV8stackTraceT) GetCV8stackTraceT() *CV8stackTraceT {
// 	return r.p_v8stack_trace
// }

// func (r *RefToCV8stackTraceT) SetCV8stackTraceT(p *CV8stackTraceT) {
// 	r.p_v8stack_trace = p
// }

// // Go type CV8stackTraceT wraps cef type *C.cef_v8stack_trace_t
// func newCV8stackTraceT(p *C.cef_v8stack_trace_t) *CV8stackTraceT {
// 	Tracef(unsafe.Pointer(p), "T973:")
// 	BaseAddRef(p)
// 	go_v8stack_trace := CV8stackTraceT{p}
// 	runtime.SetFinalizer(&go_v8stack_trace, func(g *CV8stackTraceT) {
// 		Tracef(unsafe.Pointer(g.p_v8stack_trace), "T977:")
// 		BaseRelease(g.p_v8stack_trace)
// 	})
// 	return &go_v8stack_trace
// }

// // *C.cef_v8stack_trace_t has refCounted interface
// func (p *C.cef_v8stack_trace_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
// 	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(p))
// }
