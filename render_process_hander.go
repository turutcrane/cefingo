package cefingo

import (
	"log"
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

type RenderProcessHandler interface {
	///
	// Called after the render process main thread has been created. |extra_info|
	// is a read-only value originating from
	// cef_browser_process_handler_t::on_render_process_thread_created(). Do not
	// keep a reference to |extra_info| outside of this function.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L67-L75
	///
	OnRenderThreadCreated(
		self *CRenderProcessHandlerT,
		extre_info *CListValueT)

	///
	// Called after WebKit has been initialized.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L77-L81
	///
	OnWebKitInitialized(self *CRenderProcessHandlerT)

	///
	// Called after a browser has been created. When browsing cross-origin a new
	// browser will be created before the old browser with the same identifier is
	// destroyed.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L83-L90
	///
	OnBrowserCreated(self *CRenderProcessHandlerT, browser *CBrowserT)

	///
	// Called before a browser is destroyed.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L92-L97
	///
	OnBrowserDestroyed(self *CRenderProcessHandlerT, browser *CBrowserT)

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
	OnContextCreated(self *CRenderProcessHandlerT,
		brower *CBrowserT,
		frame *CFrameT,
		context *CV8contextT)

	///
	// Called immediately before the V8 context for a frame is released. No
	// references to the context should be kept after this function is called.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L119-L127
	///
	OnContextReleased(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		context *CV8contextT)
	///
	// Called for global uncaught exceptions in a frame. Execution of this
	// callback is disabled by default. To enable set
	// CefSettings.uncaught_exception_stack_size > 0.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L129-L140
	///
	OnUncaughtException(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		context *CV8contextT,
		exception *CV8exceptionT,
		stackTrace *CV8stackTraceT,
	)

	///
	// Called when a new node in the the browser gets focus. The |node| value may
	// be NULL if no specific node has gained focus. The node object passed to
	// this function represents a snapshot of the DOM at the time this function is
	// executed. DOM objects are only valid for the scope of this function. Do not
	// keep references to or attempt to access any DOM objects outside the scope
	// of this function.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L142-L154
	///
	OnFocusedNodeChanged(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		node *CDomnodeT,
	)

	///
	// Called when a new message is received from a different process. Return true
	// (1) if the message was handled or false (0) otherwise. Do not keep a
	// reference to or attempt to access the message outside of this callback.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_render_process_handler_capi.h#L156-L165
	///
	OnProcessMessageReceived(
		self *CRenderProcessHandlerT,
		browser *CBrowserT,
		source_process CProcessIdT,
		message *CProcessMessageT) bool
}

var renderProcessHandlers = map[*C.cef_render_process_handler_t]RenderProcessHandler{}
var rphLoadHandlers = map[*C.cef_render_process_handler_t]*CLoadHandlerT{}

func newCRenderProcessHandlerT(cef *C.cef_render_process_handler_t) *CRenderProcessHandlerT {
	Logf("L30: %p", cef)
	BaseAddRef(cef)
	handler := CRenderProcessHandlerT{cef}
	runtime.SetFinalizer(&handler, func(h *CRenderProcessHandlerT) {
		if ref_count_log.output {
			Logf("L35: %p", h.p_render_process_handler)
		}
		BaseRelease(h.p_render_process_handler)
	})
	return &handler
}

func AllocCRenderProcessHandlerT(handler RenderProcessHandler) (cHandler *CRenderProcessHandlerT) {
	p := C.calloc(1, C.sizeof_cefingo_render_process_handler_wrapper_t)
	Logf("L120: p: %v", p)
	C.cefingo_construct_render_process_handler((*C.cefingo_render_process_handler_wrapper_t)(p))

	rph := newCRenderProcessHandlerT((*C.cef_render_process_handler_t)(p))
	crph := rph.p_render_process_handler
	renderProcessHandlers[crph] = handler
	registerDeassocer(unsafe.Pointer(crph), DeassocFunc(func() {
		delete(renderProcessHandlers, crph)
	}))

	return rph
}

func (h *C.cef_render_process_handler_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(h))
}

//export cefingo_render_process_handler_on_render_thread_created
func cefingo_render_process_handler_on_render_thread_created(
	self *C.cef_render_process_handler_t,
	extra_info *CListValueT) {
	Logf("L122: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}

	f.OnRenderThreadCreated(
		newCRenderProcessHandlerT(self),
		extra_info,
	)
	Logf("L168: %b", BaseHasOneRef(self))
}

//export cefingo_render_process_handler_on_web_kit_initialized
func cefingo_render_process_handler_on_web_kit_initialized(self *C.cef_render_process_handler_t) {
	Logf("L134: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}

	f.OnWebKitInitialized(newCRenderProcessHandlerT(self))

}

//export cefingo_render_process_handler_on_browser_created
func cefingo_render_process_handler_on_browser_created(self *C.cef_render_process_handler_t, browser *C.cef_browser_t) {
	Logf("L147: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}
	f.OnBrowserCreated(newCRenderProcessHandlerT(self), newCBrowserT(browser))
}

//export cefingo_render_process_handler_on_browser_destroyed
func cefingo_render_process_handler_on_browser_destroyed(self *C.cef_render_process_handler_t, browser *C.cef_browser_t) {
	Logf("L160: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}

	f.OnBrowserDestroyed(newCRenderProcessHandlerT(self), newCBrowserT(browser))
}

//export cefingo_render_process_handler_get_load_handler
func cefingo_render_process_handler_get_load_handler(
	self *C.cef_render_process_handler_t,
) (ch *C.cef_load_handler_t) {
	h := rphLoadHandlers[self]
	if h == nil {
		Logf("L188: No Handler %v", self)
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
	frame *CFrameT,
	context *CV8contextT,
) {
	Logf("L191: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}
	f.OnContextCreated(newCRenderProcessHandlerT(self), newCBrowserT(browser),
		frame, context)
}

//export cefingo_render_process_handler_on_context_released
func cefingo_render_process_handler_on_context_released(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *CFrameT,
	context *CV8contextT) {
	Logf("L207: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}

	f.OnContextReleased(newCRenderProcessHandlerT(self), newCBrowserT(browser),
		frame, context)

}

//export cefingo_render_process_handler_on_uncaught_exception
func cefingo_render_process_handler_on_uncaught_exception(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *CFrameT,
	context *CV8contextT,
	exception *CV8exceptionT,
	stackTrace *CV8stackTraceT,
) {
	Logf("L227: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}

	f.OnUncaughtException(newCRenderProcessHandlerT(self), newCBrowserT(browser),
		frame, context, exception, stackTrace)

}

//export cefingo_render_process_handler_on_focused_node_changed
func cefingo_render_process_handler_on_focused_node_changed(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	frame *CFrameT,
	node *CDomnodeT,
) {
	Logf("L245: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("36: Noo!")
	}

	f.OnFocusedNodeChanged(newCRenderProcessHandlerT(self), newCBrowserT(browser),
		frame, node)
}

//export cefingo_render_process_handler_on_process_message_received
func cefingo_render_process_handler_on_process_message_received(
	self *C.cef_render_process_handler_t,
	browser *C.cef_browser_t,
	source_process CProcessIdT,
	message *C.cef_process_message_t,
) (ret C.int) {
	Logf("L261: self: %p", self)

	f := renderProcessHandlers[self]
	if f == nil {
		log.Panicln("285: Noo!")
	}

	if f.OnProcessMessageReceived(newCRenderProcessHandlerT(self),
		newCBrowserT(browser), source_process, newCProcessMessageT(message)) {
		ret = 1
	} else {
		ret = 0
	}
	return ret
}

type DefaultRenderProcessHander struct {
}

func (*DefaultRenderProcessHander) OnRenderThreadCreated(
	self *CRenderProcessHandlerT,
	extre_info *CListValueT,
) {
	Logf("L278:")
}

func (*DefaultRenderProcessHander) OnWebKitInitialized(self *CRenderProcessHandlerT) {
	Logf("L282:")
}

func (*DefaultRenderProcessHander) OnBrowserCreated(self *CRenderProcessHandlerT, browser *CBrowserT) {
	Logf("L286:")
}

func (*DefaultRenderProcessHander) OnBrowserDestroyed(self *CRenderProcessHandlerT, browser *CBrowserT) {
	Logf("L290:")
}

func (*DefaultRenderProcessHander) OnContextCreated(self *CRenderProcessHandlerT,
	brower *CBrowserT,
	frame *CFrameT,
	context *CV8contextT,
) {
	Logf("L303:")
}

func (*DefaultRenderProcessHander) OnContextReleased(
	self *CRenderProcessHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	context *CV8contextT,
) {
	Logf("L312:")
}

func (*DefaultRenderProcessHander) OnUncaughtException(
	self *CRenderProcessHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	context *CV8contextT,
	exception *CV8exceptionT,
	stackTrace *CV8stackTraceT,
) {
	Logf("L323:")
}

func (*DefaultRenderProcessHander) OnFocusedNodeChanged(
	self *CRenderProcessHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	node *CDomnodeT,
) {
	Logf("L332:")
}

func (*DefaultRenderProcessHander) OnProcessMessageReceived(
	self *CRenderProcessHandlerT,
	browser *CBrowserT,
	source_process CProcessIdT,
	message *CProcessMessageT,
) bool {
	Logf("L341:")

	return true
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
