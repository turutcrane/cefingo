package cefingo

import (
	"log"
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
		errorText *CStringT,
		failedUrl *CStringT,
	)
}

var loadHandlerMap = map[*CLoadHandlerT]LoadHandler{}

// AllocCLoadHandlerT allocates CLoadHandlerT and construct it
func AllocCLoadHandlerT(h LoadHandler) (loadHandler *CLoadHandlerT) {
	p := C.calloc(1, C.sizeof_cefingo_load_handler_wrapper_t)
	Logf("L84: p: %v", p)

	C.construct_cefingo_load_handler((*C.cefingo_load_handler_wrapper_t)(p))

	loadHandler = (*CLoadHandlerT)(p)
	BaseAddRef(loadHandler)

	loadHandlerMap[loadHandler] = h

	return loadHandler
}

//export on_loading_state_change
func on_loading_state_change(
	self *CLoadHandlerT,
	browser *CBrowserT,
	isLoading C.int,
	canGoBack C.int,
	canGoForward C.int,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_loading_state_change: Noo!")
	}
	h.OnLoadingStateChange(self, browser, (int)(isLoading), (int)(canGoBack), (int)(canGoForward))
}

//export on_load_start
func on_load_start(
	self *CLoadHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	transitionType CTransitionTypeT,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_load_start: Noo!")
	}
	h.OnLoadStart(self, browser, frame, transitionType)
}

//export on_load_end
func on_load_end(
	self *CLoadHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	httpStatusCode C.int,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_load_end: Noo!")
	}
	h.OnLoadEnd(self, browser, frame, (int)(httpStatusCode))
}

//export on_load_error
func on_load_error(
	self *CLoadHandlerT,
	browser *CBrowserT,
	frame *CFrameT,
	errorCode CErrorcodeT,
	errorText *CStringT,
	failedUrl *CStringT,
) {
	h := loadHandlerMap[self]
	if h == nil {
		log.Panicln("L100: on_load_error: Noo!")
	}
	h.OnLoadError(self, browser, frame, errorCode, errorText, failedUrl)
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
	errorText *CStringT,
	failedUrl *CStringT,
) {
	// No Operaion
}
