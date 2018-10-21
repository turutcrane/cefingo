package cefingo

import (
	"log"
	"unsafe"

	// #include "cefingo.h"
	"C"
)

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
	) Cint
}

var client_method = map[*CClientT]Client{}
var life_span_handler = map[*CClientT]*CLifeSpanHandlerT{}

// AllocCClient allocates CClientT and construct it
func AllocCClient(c Client) (cClient *CClientT) {
	p := C.calloc(1, C.sizeof_cefingo_client_wrapper_t)
	Logf("L26: p: %v", p)
	C.construct_cefingo_client((*C.cefingo_client_wrapper_t)(p))

	cClient = (*CClientT)(p)
	BaseAddRef(cClient)
	client_method[cClient] = c

	return cClient
}

///
// Return the handler for context menus. If no handler is
// provided the default implementation will be used.
///
//export get_context_menu_handler
func get_context_menu_handler(self *CClientT) *CContextMenuHandlerT {
	return nil
}

///
// Return the handler for dialogs. If no handler is provided the default
// implementation will be used.
///
//export get_dialog_handler
func get_dialog_handler(self *CClientT) *CDialogHandlerT {
	return nil
}

///
// Return the handler for browser display state events.
///
//export get_display_handler
func get_display_handler(self *CClientT) *CDisplayHandlerT {
	return nil
}

///
// Return the handler for download events. If no handler is returned downloads
// will not be allowed.
///
//export get_download_handler
func get_download_handler(self *CClientT) *CDownloaddHanderT {
	return nil
}

///
// Return the handler for drag events.
///
//export get_drag_handler
func get_drag_handler(self *CClientT) *CDragHandlerT {
	return nil
}

///
// Return the handler for find result events.
///
//export get_find_handler
func get_find_handler(self *CClientT) *CFindHandlerT {
	return nil
}

///
// Return the handler for focus events.
///
//export get_focus_handler
func get_focus_handler(self *CClientT) *CFocusHanderT {
	return nil
}

///
// Return the handler for JavaScript dialogs. If no handler is provided the
// default implementation will be used.
///
//export get_jsdialog_handler
func get_jsdialog_handler(self *CClientT) *CJsdialogHandlerT {
	return nil
}

///
// Return the handler for keyboard events.
///
//export get_keyboard_handler
func get_keyboard_handler(self *CClientT) *CKeyboardHandlerT {
	return nil
}

// AssocLifeSpanHandler associate hander to client
func AssocLifeSpanHandler(client *CClientT, handler *CLifeSpanHandlerT) {
	p := (unsafe.Pointer)(handler)
	C.cefingo_add_ref((*C.cef_base_ref_counted_t)(p))
	life_span_handler[client] = handler
}

///
// Return the handler for browser life span events.
///
//export get_life_span_handler
func get_life_span_handler(self *CClientT) *CLifeSpanHandlerT {
	Logf("L70:")

	handler := life_span_handler[self]
	if handler == nil {
		Logf("L77: No Life Span Handler")
	} else {
		p := (unsafe.Pointer)(handler)
		C.cefingo_add_ref((*C.cef_base_ref_counted_t)(p))
	}
	return handler
}

///
// Return the handler for browser load status events.
///
//export client_get_load_handler
func client_get_load_handler(self *CClientT) *CLoadHandlerT {
	return nil
}

///
// Return the handler for off-screen rendering events.
///
//export get_render_handler
func get_render_handler(self *CClientT) *CRenderHandlerT {
	return nil
}

///
// Return the handler for browser request events.
///
//export get_request_handler
func get_request_handler(self *CClientT) *CRequestHandlerT {
	return nil
}

//on_process_mesage_received call OnProcessMessageRecived method
//export client_on_process_message_received
func client_on_process_message_received(
	self *CClientT, browser *CBrowserT, source_process CProcessIdT, message *CProcessMessageT) Cint {

	Logf("L46: client: %p", self)
	f := client_method[self]
	if f == nil {
		log.Panicln("L48: on_process_message_received: Noo!")
	}

	return f.OnProcessMessageRecived(self, browser, source_process, message)
}

// DefaultClient is dummy implementation of CClientT
type DefaultClient struct {
}

func (*DefaultClient) OnProcessMessageRecived(self *CClientT,
	browser *CBrowserT,
	source_process CProcessIdT,
	message *CProcessMessageT) Cint {

	return 0
}
