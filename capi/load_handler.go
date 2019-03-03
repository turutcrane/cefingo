package capi

import (
	"log"
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

type LoadHandler interface {
	///
	// Called when the loading state has changed. This callback will be executed
	// twice -- once when loading is initiated either programmatically or by user
	// action, and once when loading is terminated due to completion, cancellation
	// of failure. It will be called before any calls to OnLoadStart and after all
	// calls to OnLoadError and/or OnLoadEnd.
	///
	OnLoadingStateChange(
		self *CLoadHandlerT,
		browser *CBrowserT,
		isLoading int,
		canGoBack int,
		canGoForward int,
	)

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
	OnLoadStart(
		self *CLoadHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		transitionType CTransitionTypeT,
	)

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
	OnLoadEnd(
		self *CLoadHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		httpStatusCode int,
	)

	///
	// Called when a navigation fails or is canceled. This function may be called
	// by itself if before commit or in combination with OnLoadStart/OnLoadEnd if
	// after commit. |errorCode| is the error code number, |errorText| is the
	// error text and |failedUrl| is the URL that failed to load. See
	// net\base\net_error_list.h for complete descriptions of the error codes.
	///
	OnLoadError(
		self *CLoadHandlerT,
		browser *CBrowserT,
		frame *CFrameT,
		errorCode CErrorcodeT,
		errorText string,
		failedUrl string,
	)
}

var loadHandlerMap = map[*C.cef_load_handler_t]LoadHandler{}

func newCLoadHandlerT(cef *C.cef_load_handler_t) *CLoadHandlerT {
	Tracef(unsafe.Pointer(cef), "L84:")
	BaseAddRef(cef)
	handler := CLoadHandlerT{cef}

	runtime.SetFinalizer(&handler, func(h *CLoadHandlerT) {
		if ref_count_log.output {
			Tracef(unsafe.Pointer(h.p_load_handler), "L133:")
		}
		BaseRelease(h.p_load_handler)
	})
	return &handler
}

// AllocCLoadHandlerT allocates CLoadHandlerT and construct it
func AllocCLoadHandlerT(h LoadHandler) (loadHandler *CLoadHandlerT) {
	p := c_calloc(1, C.sizeof_cefingo_load_handler_wrapper_t, "L99:")

	C.cefingo_construct_load_handler((*C.cefingo_load_handler_wrapper_t)(p))

	loadHandler = newCLoadHandlerT((*C.cef_load_handler_t)(p))

	cefp := loadHandler.p_load_handler
	loadHandlerMap[cefp] = h
	registerDeassocer(unsafe.Pointer(cefp), DeassocFunc(func() {
		Tracef(unsafe.Pointer(cefp), "L108:")
		delete(loadHandlerMap, cefp)
	}))

	return loadHandler
}

func (h *C.cef_load_handler_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(h))
}

//export cefingo_load_handler_on_loading_state_change
func cefingo_load_handler_on_loading_state_change(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	isLoading C.int,
	canGoBack C.int,
	canGoForward C.int,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_loading_state_change: Noo!")
	}
	handler := newCLoadHandlerT(self)
	b := newCBrowserT(browser)
	h.OnLoadingStateChange(handler, b, (int)(isLoading), (int)(canGoBack), (int)(canGoForward))
}

//export cefingo_load_handler_on_load_start
func cefingo_load_handler_on_load_start(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	transitionType CTransitionTypeT,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_load_start: Noo!")
	}
	h.OnLoadStart(newCLoadHandlerT(self),
		newCBrowserT(browser), newCFrameT(frame), transitionType)
}

//export cefingo_load_handler_on_load_end
func cefingo_load_handler_on_load_end(
	self *C.cef_load_handler_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	httpStatusCode C.int,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_load_end: Noo!")
	}
	h.OnLoadEnd(newCLoadHandlerT(self),
		newCBrowserT(browser),
		newCFrameT(frame),
		(int)(httpStatusCode))
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
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_load_error: Noo!")
	}
	t := string_from_cef_string(errorText)
	u := string_from_cef_string(failedUrl)
	h.OnLoadError(newCLoadHandlerT(self),
		newCBrowserT(browser),
		newCFrameT(frame),
		errorCode, t, u)
}

// Default LoadHander is dummy implementaion of CLoadHanderT
type DefaultLoadHandler struct {
}

func (*DefaultLoadHandler) OnLoadingStateChange(
	self *CLoadHandlerT,
	browser *CBrowserT,
	isLoading int,
	canGoBack int,
	canGoForward int,

) {
	// No Operation
}

func (*DefaultLoadHandler) OnLoadStart(
	self *CLoadHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	transitionType CTransitionTypeT,
) {
	// No Operation
}

func (*DefaultLoadHandler) OnLoadEnd(
	self *CLoadHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	httpStatusCode int,
) {
	// No Operation
}

func (*DefaultLoadHandler) OnLoadError(
	self *CLoadHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	errorCode CErrorcodeT,
	errorText string,
	failedUrl string,
) {
	// No Operaion
}
