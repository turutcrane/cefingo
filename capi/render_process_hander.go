package capi

import (
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

///
// Class used to implement render process callbacks. The methods of this class
// will be called on the render process main thread (TID_RENDERER) unless
// otherwise indicated.
///

///
// Called after the render process main thread has been created. |extra_info|
// is a read-only value originating from
// cef_browser_process_handler_t::on_render_process_thread_created(). Do not
// keep a reference to |extra_info| outside of this function.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L67-L75
///
type OnRenderThreadCreatedHandler interface {
	OnRenderThreadCreated(
		self *CRenderProcessHandlerT,
		extre_info *CListValueT)
}

///
// Called after WebKit has been initialized.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L77-L81
///
type OnWebKitInitializedHandler interface {
	OnWebKitInitialized(self *CRenderProcessHandlerT)
}

///
// Called after a browser has been created. When browsing cross-origin a new
// browser will be created before the old browser with the same identifier is
// destroyed.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L83-L90
///
type OnBrowserCreatedHandler interface {
	OnBrowserCreated(self *CRenderProcessHandlerT, browser *CBrowserT)
}

///
// Called before a browser is destroyed.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L92-L97
///
type OnBrowserDestroyedHandler interface {
	OnBrowserDestroyed(self *CRenderProcessHandlerT, browser *CBrowserT)
}

///
// Return the handler for browser load status events.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L99-L103
///
// GetLoadHandler(self *CRenderProcessHandlerT) *CLoadHandlerT

///
// Called immediately after the V8 context for a frame has been created. To
// retrieve the JavaScript 'window' object use the
// cef_v8context_t::get_global() function. V8 handles can only be accessed
// from the thread on which they are created. A task runner for posting tasks
// on the associated thread can be retrieved via the
// cef_v8context_t::get_task_runner() function.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L105-L117
///
type OnContextCreatedHandler interface {
	OnContextCreated(self *CRenderProcessHandlerT,
		brower *CBrowserT,
		frame *CFrameT,
		context *CV8contextT)
}

///
// Called immediately before the V8 context for a frame is released. No
// references to the context should be kept after this function is called.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L119-L127
///
type OnContextReleasedHandler interface {
	OnContextReleased(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		context *CV8contextT)
}

///
// Called for global uncaught exceptions in a frame. Execution of this
// callback is disabled by default. To enable set
// CefSettings.uncaught_exception_stack_size > 0.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L129-L140
///
type OnUncaughtExceptionHandler interface {
	OnUncaughtException(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		context *CV8contextT,
		exception *CV8exceptionT,
		stackTrace *CV8stackTraceT,
	)
}

///
// Called when a new node in the the browser gets focus. The |node| value may
// be NULL if no specific node has gained focus. The node object passed to
// this function represents a snapshot of the DOM at the time this function is
// executed. DOM objects are only valid for the scope of this function. Do not
// keep references to or attempt to access any DOM objects outside the scope
// of this function.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L142-L154
///
type OnFocusedNodeChangedHandler interface {
	OnFocusedNodeChanged(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		node *CDomnodeT,
	)
}

///
// Called when a new message is received from a different process. Return true
// (1) if the message was handled or false (0) otherwise. Do not keep a
// reference to or attempt to access the message outside of this callback.
// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L156-L165
///
type OnProcessMessageReceivedHandler interface {
	OnProcessMessageReceived(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		source_process CProcessIdT,
		message *CProcessMessageT) bool
}

var on_render_thread_created_handler = map[*C.cef_render_process_handler_t]OnRenderThreadCreatedHandler{}
var on_web_kit_initialized_handler = map[*C.cef_render_process_handler_t]OnWebKitInitializedHandler{}
var on_browser_created_handler = map[*C.cef_render_process_handler_t]OnBrowserCreatedHandler{}
var on_browser_destroyed_handler = map[*C.cef_render_process_handler_t]OnBrowserDestroyedHandler{}
var on_context_created_handler = map[*C.cef_render_process_handler_t]OnContextCreatedHandler{}
var on_context_released_handler = map[*C.cef_render_process_handler_t]OnContextReleasedHandler{}
var on_uncaught_exception_handler = map[*C.cef_render_process_handler_t]OnUncaughtExceptionHandler{}
var on_focused_node_changed_handler = map[*C.cef_render_process_handler_t]OnFocusedNodeChangedHandler{}
var on_process_message_received_handler = map[*C.cef_render_process_handler_t]OnProcessMessageReceivedHandler{}

var rphLoadHandlers = map[*C.cef_render_process_handler_t]*CLoadHandlerT{}

func newCRenderProcessHandlerT(cef *C.cef_render_process_handler_t) *CRenderProcessHandlerT {
	Tracef(unsafe.Pointer(cef), "L122:")
	BaseAddRef(cef)
	handler := CRenderProcessHandlerT{cef}
	runtime.SetFinalizer(&handler, func(h *CRenderProcessHandlerT) {
		Tracef(unsafe.Pointer(h.p_render_process_handler), "L126:")
		BaseRelease(h.p_render_process_handler)
	})
	return &handler
}

func AllocCRenderProcessHandlerT() *CRenderProcessHandlerT {
	p := c_calloc(1, C.sizeof_cefingo_render_process_handler_wrapper_t, "L133:")
	C.cefingo_construct_render_process_handler((*C.cefingo_render_process_handler_wrapper_t)(p))

	return newCRenderProcessHandlerT((*C.cef_render_process_handler_t)(p))
}

func (rph *CRenderProcessHandlerT) Bind(handler interface{}) *CRenderProcessHandlerT {
	crph := rph.p_render_process_handler

	if h, ok := handler.(OnRenderThreadCreatedHandler); ok {
		on_render_thread_created_handler[crph] = h
	}
	if h, ok := handler.(OnWebKitInitializedHandler); ok {
		on_web_kit_initialized_handler[crph] = h
	}
	if h, ok := handler.(OnBrowserCreatedHandler); ok {
		on_browser_created_handler[crph] = h
	}
	if h, ok := handler.(OnBrowserDestroyedHandler); ok {
		on_browser_destroyed_handler[crph] = h
	}
	if h, ok := handler.(OnContextCreatedHandler); ok {
		on_context_created_handler[crph] = h
	}
	if h, ok := handler.(OnContextReleasedHandler); ok {
		on_context_released_handler[crph] = h
	}
	if h, ok := handler.(OnUncaughtExceptionHandler); ok {
		on_uncaught_exception_handler[crph] = h
	}
	if h, ok := handler.(OnFocusedNodeChangedHandler); ok {
		on_focused_node_changed_handler[crph] = h
	}
	if h, ok := handler.(OnProcessMessageReceivedHandler); ok {
		on_process_message_received_handler[crph] = h
	}

	registerDeassocer(unsafe.Pointer(crph), DeassocFunc(func() {
		Tracef(unsafe.Pointer(crph), "L201:")
		delete(on_render_thread_created_handler, crph)
		delete(on_web_kit_initialized_handler, crph)
		delete(on_browser_created_handler, crph)
		delete(on_browser_destroyed_handler, crph)
		delete(on_context_created_handler, crph)
		delete(on_context_released_handler, crph)
		delete(on_uncaught_exception_handler, crph)
		delete(on_focused_node_changed_handler, crph)
		delete(on_process_message_received_handler, crph)
	}))

	return rph
}

func (h *C.cef_render_process_handler_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(h))
}

//export cefingo_render_process_handler_on_render_thread_created
func cefingo_render_process_handler_on_render_thread_created(
	self *C.cef_render_process_handler_t,
	extra_info *C.cef_list_value_t) {
	Logf("L122: self: %p", self)

	f := on_render_thread_created_handler[self]
	if f != nil {
		f.OnRenderThreadCreated(
			newCRenderProcessHandlerT(self),
			newCListValueT(extra_info),
		)
		Logf("L168: %b", BaseHasOneRef(self))
	} else {
		Logf("L209: Noo!")
	}

}

//export cefingo_render_process_handler_on_web_kit_initialized
func cefingo_render_process_handler_on_web_kit_initialized(self *C.cef_render_process_handler_t) {
	Logf("L134: self: %p", self)

	f := on_web_kit_initialized_handler[self]
	if f != nil {
		f.OnWebKitInitialized(newCRenderProcessHandlerT(self))
	} else {
		Logf("L219: Noo!")
	}
}

//export cefingo_render_process_handler_on_browser_created
func cefingo_render_process_handler_on_browser_created(self *C.cef_render_process_handler_t, browser *C.cef_browser_t) {
	Logf("L147: self: %p", self)

	f := on_browser_created_handler[self]
	if f != nil {
		f.OnBrowserCreated(newCRenderProcessHandlerT(self), newCBrowserT(browser))
	} else {
		Logf("L251: Noo!")
	}
}

//export cefingo_render_process_handler_on_browser_destroyed
func cefingo_render_process_handler_on_browser_destroyed(self *C.cef_render_process_handler_t, browser *C.cef_browser_t) {
	Logf("L160: self: %p", self)

	f := on_browser_destroyed_handler[self]
	if f != nil {
		f.OnBrowserDestroyed(newCRenderProcessHandlerT(self), newCBrowserT(browser))
	} else {
		Logf("L263: Noo!")
	}

}

//export cefingo_render_process_handler_get_load_handler
func cefingo_render_process_handler_get_load_handler(
	self *C.cef_render_process_handler_t,
) (ch *C.cef_load_handler_t) {
	h := rphLoadHandlers[self]
	if h == nil {
		Logf("L274: No Handler %v", self)
	} else {
		ch = h.p_load_handler
		BaseAddRef(ch)
	}
	return ch
}

//export cefingo_render_process_handler_on_context_created
func cefingo_render_process_handler_on_context_created(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	context *C.cef_v8context_t,
) {
	Logf("L191: self: %p", self)

	f := on_context_created_handler[self]
	if f != nil {
		f.OnContextCreated(newCRenderProcessHandlerT(self), newCBrowserT(browser),
			newCFrameT(frame), newCV8contextT(context))
	} else {
		Logf("L296: Noo!")
	}
}

//export cefingo_render_process_handler_on_context_released
func cefingo_render_process_handler_on_context_released(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	context *C.cef_v8context_t) {
	Logf("L207: self: %p", self)

	f := on_context_released_handler[self]
	if f != nil {
		f.OnContextReleased(newCRenderProcessHandlerT(self), newCBrowserT(browser),
			newCFrameT(frame), newCV8contextT(context))
	} else {
		Logf("L313: Noo!")
	}
}

//export cefingo_render_process_handler_on_uncaught_exception
func cefingo_render_process_handler_on_uncaught_exception(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	context *C.cef_v8context_t,
	exception *CV8exceptionT,
	stackTrace *CV8stackTraceT,
) {
	Logf("L227: self: %p", self)

	f := on_uncaught_exception_handler[self]
	if f != nil {
		f.OnUncaughtException(newCRenderProcessHandlerT(self), newCBrowserT(browser),
			newCFrameT(frame), newCV8contextT(context), exception, stackTrace)
	} else {
		Logf("L333: Noo!")
	}
}

//export cefingo_render_process_handler_on_focused_node_changed
func cefingo_render_process_handler_on_focused_node_changed(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	node *CDomnodeT,
) {
	Logf("L245: self: %p", self)

	f := on_focused_node_changed_handler[self]
	if f != nil {
		f.OnFocusedNodeChanged(newCRenderProcessHandlerT(self), newCBrowserT(browser),
			newCFrameT(frame), node)
	} else {
		Logf("L358: Noo!")
	}
}

//export cefingo_render_process_handler_on_process_message_received
func cefingo_render_process_handler_on_process_message_received(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	source_process CProcessIdT,
	message *C.cef_process_message_t,
) (ret C.int) {
	Logf("L261: self: %p", self)

	f := on_process_message_received_handler[self]
	if f != nil {
		if f.OnProcessMessageReceived(newCRenderProcessHandlerT(self),
			newCBrowserT(browser), source_process, newCProcessMessageT(message)) {
			ret = 1
		}
	} else {
		Logf("285: Noo!")
	}

	return ret
}

// AssocLoadHandler associate hander to CRenderProcessHandlerT
func (rph *CRenderProcessHandlerT) AssocLoadHandler(h *CLoadHandlerT) {

	crph := rph.p_render_process_handler
	rphLoadHandlers[crph] = h
	registerDeassocer(unsafe.Pointer(crph), DeassocFunc(func() {
		// Do not have reference to rph itself in DeassocFunc,
		// or rph is never GCed.
		Logf("L397: %p", crph)
		delete(rphLoadHandlers, crph)
	}))
}
