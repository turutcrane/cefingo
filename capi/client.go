package capi

import (
	"log"
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

// Client is Go interface of C.cef_client_t
type Client interface {
	///
	// Called when a new message is received from a different process. Return true
	// (1) if the message was handled or false (0) otherwise. Do not keep a
	// reference to or attempt to access the message outside of this callback.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_client_capi.h#L154-L164
	///
	OnProcessMessageRecived(self *CClientT,
		browser *CBrowserT,
		source_process CProcessIdT,
		message *CProcessMessageT,
	) bool
}

var client_method = map[*C.cef_client_t]Client{}
var life_span_handler = map[*C.cef_client_t]*CLifeSpanHandlerT{}

func newCClientT(cef *C.cef_client_t) *CClientT {
	Tracef(unsafe.Pointer(cef), "L31:")
	BaseAddRef(cef)
	client := CClientT{cef}
	runtime.SetFinalizer(&client, func(c *CClientT) {
		Tracef(unsafe.Pointer(c.p_client), "L35:")
		BaseRelease(c.p_client)
	})
	return &client
}

// AllocCClient allocates CClientT and construct it
func AllocCClient(c Client) (cClient *CClientT) {
	p := c_calloc(1, C.sizeof_cefingo_client_wrapper_t, "L43:")
	C.cefingo_construct_client((*C.cefingo_client_wrapper_t)(p))

	cClient = newCClientT((*C.cef_client_t)(p))
	cp := cClient.p_client
	client_method[cp] = c
	registerDeassocer(unsafe.Pointer(cp), DeassocFunc(func() {
		Tracef(unsafe.Pointer(cp), "L50:")
		delete(client_method, cp)
	}))

	return cClient
}

func (c *C.cef_client_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(c))
}

///
// Return the handler for context menus. If no handler is
// provided the default implementation will be used.
///
//export cefingo_client_get_context_menu_handler
func cefingo_client_get_context_menu_handler(self *C.cef_client_t) *CContextMenuHandlerT {
	return nil
}

///
// Return the handler for dialogs. If no handler is provided the default
// implementation will be used.
///
//export cefingo_client_get_dialog_handler
func cefingo_client_get_dialog_handler(self *C.cef_client_t) *CDialogHandlerT {
	return nil
}

///
// Return the handler for browser display state events.
///
//export cefingo_client_get_display_handler
func cefingo_client_get_display_handler(self *C.cef_client_t) *CDisplayHandlerT {
	return nil
}

///
// Return the handler for download events. If no handler is returned downloads
// will not be allowed.
///
//export cefingo_client_get_download_handler
func cefingo_client_get_download_handler(self *C.cef_client_t) *CDownloaddHanderT {
	return nil
}

///
// Return the handler for drag events.
///
//export cefingo_client_get_drag_handler
func cefingo_client_get_drag_handler(self *C.cef_client_t) *CDragHandlerT {
	return nil
}

///
// Return the handler for find result events.
///
//export cefingo_client_get_find_handler
func cefingo_client_get_find_handler(self *C.cef_client_t) *CFindHandlerT {
	return nil
}

///
// Return the handler for focus events.
///
//export cefingo_client_get_focus_handler
func cefingo_client_get_focus_handler(self *C.cef_client_t) *CFocusHanderT {
	return nil
}

///
// Return the handler for JavaScript dialogs. If no handler is provided the
// default implementation will be used.
///
//export cefingo_client_get_jsdialog_handler
func cefingo_client_get_jsdialog_handler(self *C.cef_client_t) *CJsdialogHandlerT {
	return nil
}

///
// Return the handler for keyboard events.
///
//export cefingo_client_get_keyboard_handler
func cefingo_client_get_keyboard_handler(self *C.cef_client_t) *CKeyboardHandlerT {
	return nil
}

// AssocLifeSpanHandler associate hander to client
func (client *CClientT) AssocLifeSpanHandler(handler *CLifeSpanHandlerT) {

	cp := client.p_client
	life_span_handler[cp] = handler
	registerDeassocer(unsafe.Pointer(cp), DeassocFunc(func() {
		// Do not have reference to client itself in DeassocFunc,
		// or client is never GCed.
		Tracef(unsafe.Pointer(cp), "L145:")
		delete(life_span_handler, cp)
	}))
}

///
// Return the handler for browser life span events.
///
//export cefingo_client_get_life_span_handler
func cefingo_client_get_life_span_handler(self *C.cef_client_t) *C.cef_life_span_handler_t {
	Logf("L70:")

	handler := life_span_handler[self]
	if handler == nil {
		Logf("L77: No Life Span Handler")
	} else {
		BaseAddRef(handler.p_life_span_handler)
	}
	return handler.p_life_span_handler
}

///
// Return the handler for browser load status events.
///
//export cefingo_client_client_get_load_handler
func cefingo_client_client_get_load_handler(self *C.cef_client_t) *C.cef_load_handler_t {
	return nil
}

///
// Return the handler for off-screen rendering events.
///
//export cefingo_client_get_render_handler
func cefingo_client_get_render_handler(self *C.cef_client_t) *CRenderHandlerT {
	return nil
}

///
// Return the handler for browser request events.
///
//export cefingo_client_get_request_handler
func cefingo_client_get_request_handler(self *C.cef_client_t) *CRequestHandlerT {
	return nil
}

//on_process_mesage_received call OnProcessMessageRecived method
//export cefingo_client_on_process_message_received
func cefingo_client_on_process_message_received(
	self *C.cef_client_t,
	browser *C.cef_browser_t,
	source_process CProcessIdT,
	message *C.cef_process_message_t,
) (ret C.int) {
	Tracef(unsafe.Pointer(self), "L46:")
	f := client_method[self]
	if f == nil {
		log.Panicln("L48: on_process_message_received: Noo!")
	}

	client := newCClientT(self)
	b := newCBrowserT(browser)
	m := newCProcessMessageT(message)
	if f.OnProcessMessageRecived(client, b, source_process, m) {
		ret = 1
	} else {
		ret = 0
	}
	return ret
}

// DefaultClient is dummy implementation of CClientT
type DefaultClient struct {
}

func (*DefaultClient) OnProcessMessageRecived(self *CClientT,
	browser *CBrowserT,
	source_process CProcessIdT,
	message *CProcessMessageT,
) bool {
	return false
}
