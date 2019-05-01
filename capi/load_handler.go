package capi

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

///
// Called when the loading state has changed. This callback will be executed
// twice -- once when loading is initiated either programmatically or by user
// action, and once when loading is terminated due to completion, cancellation
// of failure. It will be called before any calls to OnLoadStart and after all
// calls to OnLoadError and/or OnLoadEnd.
///
type OnLoadingStateChangeHandler interface {
	OnLoadingStateChange(
		self *CLoadHandlerT,
		browser *CBrowserT,
		isLoading int,
		canGoBack int,
		canGoForward int,
	)
}

///
// Called after a navigation has been committed and before the browser begins
// loading contents in the frame. The |frame| value will never be NULL -- call
// the is_main() function to check if this frame is the main frame.
// |transition_type| provides information about the source of the navigation
// and an accurate value is only available in the browser process. Multiple
// frames may be loading at the same time. Sub-frames may start or continue
// loading after the main frame load has ended. This function will not be
// called for same page navigations (fragments, history state, etc.) or for
// navigations that fail or are canceled before commit. For notification of
// overall browser load status use OnLoadingStateChange instead.
///
type OnLoadStartHandler interface {
	OnLoadStart(
		self *CLoadHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		transitionType CTransitionTypeT,
	)
}

///
// Called when the browser is done loading a frame. The |frame| value will
// never be NULL -- call the is_main() function to check if this frame is the
// main frame. Multiple frames may be loading at the same time. Sub-frames may
// start or continue loading after the main frame load has ended. This
// function will not be called for same page navigations (fragments, history
// state, etc.) or for navigations that fail or are canceled before commit.
// For notification of overall browser load status use OnLoadingStateChange
// instead.
///
type OnLoadEndHandler interface {
	OnLoadEnd(
		self *CLoadHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		httpStatusCode int,
	)
}

///
// Called when a navigation fails or is canceled. This function may be called
// by itself if before commit or in combination with OnLoadStart/OnLoadEnd if
// after commit. |errorCode| is the error code number, |errorText| is the
// error text and |failedUrl| is the URL that failed to load. See
// net\base\net_error_list.h for complete descriptions of the error codes.
///
type OnLoadErrorHandler interface {
	OnLoadError(
		self *CLoadHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		errorCode CErrorcodeT,
		errorText string,
		failedUrl string,
	)
}

var on_loading_state_change_handler = map[*C.cef_load_handler_t]OnLoadingStateChangeHandler{}
var on_load_start_handler = map[*C.cef_load_handler_t]OnLoadStartHandler{}
var on_load_end_handler = map[*C.cef_load_handler_t]OnLoadEndHandler{}
var on_load_error_handler = map[*C.cef_load_handler_t]OnLoadErrorHandler{}

// AllocCLoadHandlerT allocates CLoadHandlerT and construct it
func AllocCLoadHandlerT() *CLoadHandlerT {
	p := (*C.cefingo_load_handler_wrapper_t)(
		c_calloc(1, C.sizeof_cefingo_load_handler_wrapper_t, "L106:"))

	C.cefingo_construct_load_handler(p)

	return newCLoadHandlerT((*C.cef_load_handler_t)(unsafe.Pointer(p)))
}

func (loadHandler *CLoadHandlerT) Bind(handler interface{}) *CLoadHandlerT {
	cefp := loadHandler.p_load_handler

	if h, ok := handler.(OnLoadingStateChangeHandler); ok {
		on_loading_state_change_handler[cefp] = h
	}
	if h, ok := handler.(OnLoadStartHandler); ok {
		on_load_start_handler[cefp] = h
	}
	if h, ok := handler.(OnLoadEndHandler); ok {
		on_load_end_handler[cefp] = h
	}
	if h, ok := handler.(OnLoadErrorHandler); ok {
		on_load_error_handler[cefp] = h
	}

	registerDeassocer(unsafe.Pointer(cefp), DeassocFunc(func() {
		Tracef(unsafe.Pointer(cefp), "L108:")
		delete(on_loading_state_change_handler, cefp)
		delete(on_load_start_handler, cefp)
		delete(on_load_end_handler, cefp)
		delete(on_load_error_handler, cefp)
	}))

	if accessor, ok := handler.(CLoadHandlerTAccessor); ok {
		accessor.SetCLoadHandlerT(loadHandler)
		Logf("L109:")
	}

	return loadHandler
}

//export cefingo_load_handler_on_loading_state_change
func cefingo_load_handler_on_loading_state_change(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	isLoading C.int,
	canGoBack C.int,
	canGoForward C.int,
) {
	h := on_loading_state_change_handler[self]
	if h != nil {
		handler := newCLoadHandlerT(self)
		b := newCBrowserT(browser)
		h.OnLoadingStateChange(handler, b, (int)(isLoading), (int)(canGoBack), (int)(canGoForward))
	} else {
		Logf("L139: on_loading_state_change: Noo!")
	}
}

//export cefingo_load_handler_on_load_start
func cefingo_load_handler_on_load_start(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	transitionType CTransitionTypeT,
) {
	h := on_load_start_handler[self]
	if h != nil {
		h.OnLoadStart(newCLoadHandlerT(self),
			newCBrowserT(browser), newCFrameT(frame), transitionType)
	} else {
		Logf("L159: on_load_start: Noo!")
	}
}

//export cefingo_load_handler_on_load_end
func cefingo_load_handler_on_load_end(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	httpStatusCode C.int,
) {
	h := on_load_end_handler[self]
	if h != nil {
		h.OnLoadEnd(newCLoadHandlerT(self),
			newCBrowserT(browser),
			newCFrameT(frame),
			(int)(httpStatusCode))
	} else {
		Logf("L177: on_load_end: Noo!")
	}
}

//export cefingo_load_handler_on_load_error
func cefingo_load_handler_on_load_error(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	errorCode CErrorcodeT,
	errorText *C.cef_string_t,
	failedUrl *C.cef_string_t,
) {
	h := on_load_error_handler[self]
	if h != nil {
		t := string_from_cef_string(errorText)
		u := string_from_cef_string(failedUrl)
		h.OnLoadError(newCLoadHandlerT(self),
			newCBrowserT(browser),
			newCFrameT(frame),
			errorCode, t, u)
	} else {
		Logf("L192: on_load_error: Noo!")
	}
}
