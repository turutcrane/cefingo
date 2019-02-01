package cefingo

import (
	"log"

	// #include "cefingo.h"
	"C"
)

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

var browserProcessHandlers = map[*CBrowserProcessHandlerT]BrowserProcessHandler{}

//export cefingo_browser_process_handler_on_context_initialized
func cefingo_browser_process_handler_on_context_initialized(self *CBrowserProcessHandlerT) {
	// Logf("L19: self: %p", self)

	f := browserProcessHandlers[self]
	if f == nil {
		log.Panicln("37: Noo!")
	}

	f.OnContextInitialized(self)
}

type DefaultBrowserProcessHandler struct {
}

func AllocCBrowserProcessHandlerT(handler BrowserProcessHandler) (cHandler *CBrowserProcessHandlerT) {
	p := C.calloc(1, C.sizeof_cefingo_browser_process_handler_wrapper_t)
	Logf("L39: p: %v", p)
	C.cefingo_construct_browser_process_handler((*C.cefingo_browser_process_handler_wrapper_t)(p))

	cHandler = (*CBrowserProcessHandlerT)(p)
	BaseAddRef(cHandler)
	browserProcessHandlers[cHandler] = handler

	return cHandler
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
