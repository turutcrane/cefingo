package cefingo

//
import (
	"log"
)

// #include "cefingo.h"
import "C"

// LifeSpanHandler is Go interface of C.cef_life_span_handler_t
type LifeSpanHandler interface {
	///
	// Called after a new browser is created. This callback will be the first
	// notification that references |browser|.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_life_span_handler_capi.h#L99-#L102
	///
	OnAfterCreated(self *CLifeSpanHandlerT, brwoser *CBrowserT)

	///
	// Called just before a browser is destroyed. Release all references to the
	// browser object and do not attempt to execute any functions on the browser
	// object after this callback returns. This callback will be the last
	// notification that references |browser|. See do_close() documentation for
	// additional usage information.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_life_span_handler_capi.h#L198-#L206
	///
	OnBeforeClose(self *CLifeSpanHandlerT, brwoser *CBrowserT)

	///
	// Called when a browser has recieved a request to close. This may result
	// directly from a call to cef_browser_host_t::*close_browser() or indirectly
	// if the browser is parented to a top-level window created by CEF and the
	// user attempts to close that window (by clicking the 'X', for example). The
	// do_close() function will be called after the JavaScript 'onunload' event
	// has been fired.
	//
	// An application should handle top-level owner window close notifications by
	// calling cef_browser_host_t::try_close_browser() or
	// cef_browser_host_t::CloseBrowser(false (0)) instead of allowing the window
	// to close immediately (see the examples below). This gives CEF an
	// opportunity to process the 'onbeforeunload' event and optionally cancel the
	// close before do_close() is called.
	//
	// When windowed rendering is enabled CEF will internally create a window or
	// view to host the browser. In that case returning false (0) from do_close()
	// will send the standard close notification to the browser's top-level owner
	// window (e.g. WM_CLOSE on Windows, performClose: on OS X, "delete_event" on
	// Linux or cef_window_delegate_t::can_close() callback from Views). If the
	// browser's host window/view has already been destroyed (via view hierarchy
	// tear-down, for example) then do_close() will not be called for that browser
	// since is no longer possible to cancel the close.
	//
	// When windowed rendering is disabled returning false (0) from do_close()
	// will cause the browser object to be destroyed immediately.
	//
	// If the browser's top-level owner window requires a non-standard close
	// notification then send that notification from do_close() and return true
	// (1).
	//
	// The cef_life_span_handler_t::on_before_close() function will be called
	// after do_close() (if do_close() is called) and immediately before the
	// browser object is destroyed. The application should only exit after
	// on_before_close() has been called for all existing browsers.
	//
	// The below examples describe what should happen during window close when the
	// browser is parented to an application-provided top-level window.
	//
	// Example 1: Using cef_browser_host_t::try_close_browser(). This is
	// recommended for clients using standard close handling and windows created
	// on the browser process UI thread. 1.  User clicks the window close button
	// which sends a close notification to
	//     the application's top-level window.
	// 2.  Application's top-level window receives the close notification and
	//     calls TryCloseBrowser() (which internally calls CloseBrowser(false)).
	//     TryCloseBrowser() returns false so the client cancels the window close.
	// 3.  JavaScript 'onbeforeunload' handler executes and shows the close
	//     confirmation dialog (which can be overridden via
	//     CefJSDialogHandler::OnBeforeUnloadDialog()).
	// 4.  User approves the close. 5.  JavaScript 'onunload' handler executes. 6.
	// CEF sends a close notification to the application's top-level window
	//     (because DoClose() returned false by default).
	// 7.  Application's top-level window receives the close notification and
	//     calls TryCloseBrowser(). TryCloseBrowser() returns true so the client
	//     allows the window close.
	// 8.  Application's top-level window is destroyed. 9.  Application's
	// on_before_close() handler is called and the browser object
	//     is destroyed.
	// 10. Application exits by calling cef_quit_message_loop() if no other
	// browsers
	//     exist.
	//
	// Example 2: Using cef_browser_host_t::CloseBrowser(false (0)) and
	// implementing the do_close() callback. This is recommended for clients using
	// non-standard close handling or windows that were not created on the browser
	// process UI thread. 1.  User clicks the window close button which sends a
	// close notification to
	//     the application's top-level window.
	// 2.  Application's top-level window receives the close notification and:
	//     A. Calls CefBrowserHost::CloseBrowser(false).
	//     B. Cancels the window close.
	// 3.  JavaScript 'onbeforeunload' handler executes and shows the close
	//     confirmation dialog (which can be overridden via
	//     CefJSDialogHandler::OnBeforeUnloadDialog()).
	// 4.  User approves the close. 5.  JavaScript 'onunload' handler executes. 6.
	// Application's do_close() handler is called. Application will:
	//     A. Set a flag to indicate that the next close attempt will be allowed.
	//     B. Return false.
	// 7.  CEF sends an close notification to the application's top-level window.
	// 8.  Application's top-level window receives the close notification and
	//     allows the window to close based on the flag from #6B.
	// 9.  Application's top-level window is destroyed. 10. Application's
	// on_before_close() handler is called and the browser object
	//     is destroyed.
	// 11. Application exits by calling cef_quit_message_loop() if no other
	// browsers
	//     exist.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_life_span_handler_capi.h#L106-#L194
	///
	DoClose(self *CLifeSpanHandlerT, brwoser *CBrowserT) bool
}

var lifeSpanHandlers = map[*CLifeSpanHandlerT]LifeSpanHandler{}

func AllocCLifeSpanHandlerT(handler LifeSpanHandler) (cHandler *CLifeSpanHandlerT) {
	p := C.calloc(1, C.sizeof_cefingo_life_span_handler_wrapper_t)
	Logf("L23: p: %v", p)
	C.cefingo_construct_life_span_handler((*C.cefingo_life_span_handler_wrapper_t)(p))

	ch := (*CLifeSpanHandlerT)(p)
	BaseAddRef(ch)

	lifeSpanHandlers[ch] = handler

	return ch
}

//export life_span_on_before_close
func life_span_on_before_close(self *CLifeSpanHandlerT, browser *CBrowserT) {
	Logf("L39:")

	f := lifeSpanHandlers[self]
	if f == nil {
		log.Panicln("L44: life_span_on_before_close: Noo!")
	}

	f.OnBeforeClose(self, browser)
}

//export life_span_do_close
func life_span_do_close(self *CLifeSpanHandlerT, brwoser *CBrowserT) (ret C.int) {
	Logf("L50:")
	f := lifeSpanHandlers[self]
	if f == nil {
		log.Panicln("L58: life_span_do_close: Noo!")
	}

	if f.DoClose(self, brwoser) {
		ret = 1
	} else {
		ret = 0
	}
	return ret
}

// void on_after_created(struct _cef_life_span_handler_t* self,
//                               struct _cef_browser_t* browser) {
// }
//export life_span_on_after_created
func life_span_on_after_created(self *CLifeSpanHandlerT, brwoser *CBrowserT) {
	Logf("L60:")
	f := lifeSpanHandlers[self]
	if f == nil {
		log.Panicln("L70: life_span_on_after_created Noo!")
	}

	f.OnAfterCreated(self, brwoser)
}

type DefaultLifeSpanHandler struct {
}

func (*DefaultLifeSpanHandler) OnBeforeClose(self *CLifeSpanHandlerT, brwoser *CBrowserT) {
	Logf("L79:")
}

func (*DefaultLifeSpanHandler) DoClose(self *CLifeSpanHandlerT, brwoser *CBrowserT) bool {
	Logf("L83:")
	return false
}

func (*DefaultLifeSpanHandler) OnAfterCreated(self *CLifeSpanHandlerT, brwoser *CBrowserT) {
	Logf("L88:")
}
