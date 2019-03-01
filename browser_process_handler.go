package cefingo

import (
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

// BrowserProcessHandler is Go interface of C.cef_browser_process_handler_t
type BrowserProcessHandler interface {
	///
	// Called on the browser process UI thread immediately after the CEF context
	// has been initialized.
	///
	OnContextInitialized(self *CBrowserProcessHandlerT)

	///
	// Called before a child process is launched. Will be called on the browser
	// process UI thread when launching a render process and on the browser
	// process IO thread when launching a GPU or plugin process. Provides an
	// opportunity to modify the child process command line. Do not keep a
	// reference to |command_line| outside of this function.
	///
	OnBeforeChildProcessLaunch(
		self *CBrowserProcessHandlerT,
		command_line *CCommandLineT,
	)

	///
	// Called on the browser process IO thread after the main thread has been
	// created for a new render process. Provides an opportunity to specify extra
	// information that will be passed to
	// cef_render_process_handler_t::on_render_thread_created() in the render
	// process. Do not keep a reference to |extra_info| outside of this function.
	///
	OnRenderProcessThreadCreated(
		self *CBrowserProcessHandlerT,
		extra_info *CListValueT)

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
}

var browserProcessHandlers = map[*C.cef_browser_process_handler_t]BrowserProcessHandler{}

func newCBrowserProcessHandlerT(cef *C.cef_browser_process_handler_t) *CBrowserProcessHandlerT {
	Tracef(unsafe.Pointer(cef), "L127:")
	BaseAddRef(cef)
	handler := CBrowserProcessHandlerT{cef}

	runtime.SetFinalizer(&handler, func(h *CBrowserProcessHandlerT) {
		Tracef(unsafe.Pointer(h.p_browser_process_handler), "L133:")
		BaseRelease(h.p_browser_process_handler)
	})
	return &handler
}

//export cefingo_browser_process_handler_on_context_initialized
func cefingo_browser_process_handler_on_context_initialized(self *C.cef_browser_process_handler_t) {

	f := browserProcessHandlers[self]
	if f == nil {
		Panicf("75: Noo!")
	}
	h := newCBrowserProcessHandlerT(self)
	f.OnContextInitialized(h)
}

//export cefingo_browser_process_handler_on_render_process_thread_created
func cefingo_browser_process_handler_on_render_process_thread_created(
	self *C.cef_browser_process_handler_t,
	extra_info *C.cef_list_value_t,
) {
	f := browserProcessHandlers[self]
	if f == nil {
		Panicf("88: Noo!")
	}
	f.OnRenderProcessThreadCreated(newCBrowserProcessHandlerT(self),
		newCListValueT(extra_info))
}

type DefaultBrowserProcessHandler struct {
}

func AllocCBrowserProcessHandlerT(handler BrowserProcessHandler) (cHandler *CBrowserProcessHandlerT) {
	p := c_calloc(1, C.sizeof_cefingo_browser_process_handler_wrapper_t, "L112:")
	C.cefingo_construct_browser_process_handler((*C.cefingo_browser_process_handler_wrapper_t)(p))

	cHandler = newCBrowserProcessHandlerT((*C.cef_browser_process_handler_t)(p))

	cefp := cHandler.p_browser_process_handler
	browserProcessHandlers[cefp] = handler
	registerDeassocer(unsafe.Pointer(cefp), DeassocFunc(func() {
		delete(browserProcessHandlers, cefp)
	}))

	return cHandler
}

func (h *C.cef_browser_process_handler_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(h))
}

func (*DefaultBrowserProcessHandler) OnBeforeChildProcessLaunch(
	self *CBrowserProcessHandlerT,
	command_line *CCommandLineT,
) {
	return
}

func (*DefaultBrowserProcessHandler) OnRenderProcessThreadCreated(
	self *CBrowserProcessHandlerT,
	extra_info *CListValueT) {
	return
}
