package capi

import (
	"sync"
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

var loadHandlers = struct {
	m                               sync.Mutex
	on_loading_state_change_handler map[*C.cef_load_handler_t]OnLoadingStateChangeHandler
	on_load_start_handler           map[*C.cef_load_handler_t]OnLoadStartHandler
	on_load_end_handler             map[*C.cef_load_handler_t]OnLoadEndHandler
	on_load_error_handler           map[*C.cef_load_handler_t]OnLoadErrorHandler
}{
	sync.Mutex{},
	map[*C.cef_load_handler_t]OnLoadingStateChangeHandler{},
	map[*C.cef_load_handler_t]OnLoadStartHandler{},
	map[*C.cef_load_handler_t]OnLoadEndHandler{},
	map[*C.cef_load_handler_t]OnLoadErrorHandler{},
}

// AllocCLoadHandlerT allocates CLoadHandlerT and construct it
func AllocCLoadHandlerT() *CLoadHandlerT {
	up := c_calloc(1, C.sizeof_cefingo_load_handler_wrapper_t, "T106:")
	cefp := C.cefingo_construct_load_handler((*C.cefingo_load_handler_wrapper_t)(up))

	registerDeassocer(up, DeassocFunc(func() {
		Tracef(up, "T108:")
		loadHandlers.m.Lock()
		defer loadHandlers.m.Unlock()

		delete(loadHandlers.on_loading_state_change_handler, cefp)
		delete(loadHandlers.on_load_start_handler, cefp)
		delete(loadHandlers.on_load_end_handler, cefp)
		delete(loadHandlers.on_load_error_handler, cefp)
	}))

	return newCLoadHandlerT(cefp)
}

func (lh *CLoadHandlerT) Bind(handler interface{}) *CLoadHandlerT {
	cefp := lh.p_load_handler
	loadHandlers.m.Lock()
	defer loadHandlers.m.Unlock()

	if h, ok := handler.(OnLoadingStateChangeHandler); ok {
		loadHandlers.on_loading_state_change_handler[cefp] = h
	}
	if h, ok := handler.(OnLoadStartHandler); ok {
		loadHandlers.on_load_start_handler[cefp] = h
	}
	if h, ok := handler.(OnLoadEndHandler); ok {
		loadHandlers.on_load_end_handler[cefp] = h
	}
	if h, ok := handler.(OnLoadErrorHandler); ok {
		loadHandlers.on_load_error_handler[cefp] = h
	}

	if accessor, ok := handler.(CLoadHandlerTAccessor); ok {
		accessor.SetCLoadHandlerT(lh)
		Logf("T109:")
	}

	return lh
}

//export cefingo_load_handler_on_loading_state_change
func cefingo_load_handler_on_loading_state_change(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	isLoading C.int,
	canGoBack C.int,
	canGoForward C.int,
) {
	loadHandlers.m.Lock()
	h := loadHandlers.on_loading_state_change_handler[self]
	loadHandlers.m.Unlock()

	if h != nil {
		handler := newCLoadHandlerT(self)
		b := newCBrowserT(browser)
		h.OnLoadingStateChange(handler, b, (int)(isLoading), (int)(canGoBack), (int)(canGoForward))
	} else {
		Logf("T139: on_loading_state_change: Noo!")
	}
}

//export cefingo_load_handler_on_load_start
func cefingo_load_handler_on_load_start(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	transitionType CTransitionTypeT,
) {
	loadHandlers.m.Lock()
	h := loadHandlers.on_load_start_handler[self]
	loadHandlers.m.Unlock()

	if h != nil {
		h.OnLoadStart(newCLoadHandlerT(self),
			newCBrowserT(browser), newCFrameT(frame), transitionType)
	} else {
		Logf("T159: on_load_start: Noo!")
	}
}

//export cefingo_load_handler_on_load_end
func cefingo_load_handler_on_load_end(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	httpStatusCode C.int,
) {
	loadHandlers.m.Lock()
	h := loadHandlers.on_load_end_handler[self]
	loadHandlers.m.Unlock()

	if h != nil {
		h.OnLoadEnd(newCLoadHandlerT(self),
			newCBrowserT(browser),
			newCFrameT(frame),
			(int)(httpStatusCode))
	} else {
		Logf("T177: on_load_end: Noo!")
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
	loadHandlers.m.Lock()
	h := loadHandlers.on_load_error_handler[self]
	loadHandlers.m.Unlock()

	if h != nil {
		t := string_from_cef_string(errorText)
		u := string_from_cef_string(failedUrl)
		h.OnLoadError(newCLoadHandlerT(self),
			newCBrowserT(browser),
			newCFrameT(frame),
			errorCode, t, u)
	} else {
		Logf("T192: on_load_error: Noo!")
	}
}
