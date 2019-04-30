package capi

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

///
// Called on the browser process UI thread immediately after the CEF context
// has been initialized.
///
type OnContextInitializedHandler interface {
	OnContextInitialized(self *CBrowserProcessHandlerT)
}

///
// Called before a child process is launched. Will be called on the browser
// process UI thread when launching a render process and on the browser
// process IO thread when launching a GPU or plugin process. Provides an
// opportunity to modify the child process command line. Do not keep a
// reference to |command_line| outside of this function.
///
type OnBeforeChildProcessLaunch interface {
	OnBeforeChildProcessLaunch(
		self *CBrowserProcessHandlerT,
		command_line *CCommandLineT,
	)
}

///
// Called on the browser process IO thread after the main thread has been
// created for a new render process. Provides an opportunity to specify extra
// information that will be passed to
// cef_render_process_handler_t::on_render_thread_created() in the render
// process. Do not keep a reference to |extra_info| outside of this function.
///
type OnRenderProcessThreadCreatedHandler interface {
	OnRenderProcessThreadCreated(
		self *CBrowserProcessHandlerT,
		extra_info *CListValueT)
}

///
// Return the handler for printing on Linux. If a print handler is not
// provided then printing will not be supported on the Linux platform.
///
//struct _cef_print_handler_t*(CEF_CALLBACK* get_print_handler)(
//	struct _cef_browser_process_handler_t* self);

///
// Called from any thread when work has been scheduled for the browser process
// main (UI) thread. This callback is used in combination with CefSettings.
// external_message_pump and cef_do_message_loop_work() in cases where the CEF
// message loop must be integrated into an existing application message loop
// (see additional comments and warnings on CefDoMessageLoopWork). This
// callback should schedule a cef_do_message_loop_work() call to happen on the
// main (UI) thread. |delay_ms| is the requested delay in milliseconds. If
// |delay_ms| is <= 0 then the call should happen reasonably soon. If
// |delay_ms| is > 0 then the call should be scheduled to happen after the
// specified delay and any currently pending scheduled call should be
// cancelled.
///
//void(CEF_CALLBACK* on_schedule_message_pump_work)(
//	struct _cef_browser_process_handler_t* self,
//	int64 delay_ms);

var on_context_initialized_handler = map[*C.cef_browser_process_handler_t]OnContextInitializedHandler{}
var on_render_process_thread_created_handler = map[*C.cef_browser_process_handler_t]OnRenderProcessThreadCreatedHandler{}

//export cefingo_browser_process_handler_on_context_initialized
func cefingo_browser_process_handler_on_context_initialized(self *C.cef_browser_process_handler_t) {

	f := on_context_initialized_handler[self]
	if f != nil {
		h := newCBrowserProcessHandlerT(self)
		f.OnContextInitialized(h)
	} else {
		Logf("75: Noo!")
	}
}

//export cefingo_browser_process_handler_on_render_process_thread_created
func cefingo_browser_process_handler_on_render_process_thread_created(
	self *C.cef_browser_process_handler_t,
	extra_info *C.cef_list_value_t,
) {
	f := on_render_process_thread_created_handler[self]
	if f != nil {
		f.OnRenderProcessThreadCreated(newCBrowserProcessHandlerT(self),
			newCListValueT(extra_info))
	} else {
		Logf("L109: Noo!")
	}
}

func AllocCBrowserProcessHandlerT() *CBrowserProcessHandlerT {
	p := (*C.cefingo_browser_process_handler_wrapper_t)(
		c_calloc(1, C.sizeof_cefingo_browser_process_handler_wrapper_t, "L112:"))
	C.cefingo_construct_browser_process_handler(p)

	return newCBrowserProcessHandlerT(
		(*C.cef_browser_process_handler_t)(unsafe.Pointer(p)))
}

func (bph *CBrowserProcessHandlerT) Bind(handler interface{}) *CBrowserProcessHandlerT {

	cefp := bph.p_browser_process_handler

	if h, ok := handler.(OnContextInitializedHandler); ok {
		on_context_initialized_handler[cefp] = h
	}

	if h, ok := handler.(OnRenderProcessThreadCreatedHandler); ok {
		on_render_process_thread_created_handler[cefp] = h
	}

	registerDeassocer(unsafe.Pointer(cefp), DeassocFunc(func() {
		Tracef(unsafe.Pointer(cefp), "L126:")
		delete(on_context_initialized_handler, cefp)
		delete(on_render_process_thread_created_handler, cefp)
	}))

	if accessor, ok := handler.(CBrowserProcessHandlerTAccessor); ok {
		accessor.SetCBrowserProcessHandlerT(bph)
		Logf("L161:")
	}

	return bph
}
