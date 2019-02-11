package cefingo

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

var app_method = map[*CAppT]App{}
var browser_process_handler = map[*CAppT]*CBrowserProcessHandlerT{}
var render_process_handler = map[*CAppT]*CRenderProcessHandlerT{}

// Client is Go interface of C.cef_app_t
type App interface {
	///
	// Provides an opportunity to view and/or modify command-line arguments before
	// processing by CEF and Chromium. The |process_type| value will be NULL for
	// the browser process. Do not keep a reference to the cef_command_line_t
	// object passed to this function. The CefSettings.command_line_args_disabled
	// value can be used to start with an NULL command-line object. Any values
	// specified in CefSettings that equate to command-line arguments will be set
	// before this function is called. Be cautious when using this function to
	// modify command-line arguments for non-browser processes as this may result
	// in undefined behavior including crashes.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#L66-L80
	///
	OnBeforeCommandLineProcessing(self *CAppT, process_type string, command_line *CCommandLineT)

	///
	// Provides an opportunity to register custom schemes. Do not keep a reference
	// to the |registrar| object. This function is called on the main thread for
	// each process and the registered schemes should be the same across all
	// processes.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#82-L90
	///
	OnRegisterCustomSchemes(self *CAppT, registrar *CSchemeRegistrarT)
}

// AllocCClient allocates CAppT and construct it
func AllocCAppT(a App) (cApp *CAppT) {
	p := C.calloc(1, C.sizeof_cefingo_app_wrapper_t)
	Logf("L22: p: %v", p)

	ap := (*C.cefingo_app_wrapper_t)(p)
	C.cefingo_construct_app(ap)

	cApp = (*CAppT)(p)
	BaseAddRef(cApp)
	app_method[cApp] = a

	return cApp
}

func (a *CAppT) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(a))
}

///
// Return the handler for resource bundle events. If
// CefSettings.pack_loading_disabled is true (1) a handler must be returned.
// If no handler is returned resources will be loaded from pack files. This
// function is called by the browser and render processes on multiple threads.
///
//export cefing_app_get_resource_bundle_handler
func cefing_app_get_resource_bundle_handler(self *CAppT) *CResourceBundleHanderT {
	return nil
}

// AssocBrowserProcessHandler associate a hander to app
func AssocBrowserProcessHandler(app *CAppT, handler *CBrowserProcessHandlerT) {
	BaseAddRef(handler)
	browser_process_handler[app] = handler
}

///
// Return the handler for functionality specific to the browser process. This
// function is called on multiple threads in the browser process.
///
//export cefing_app_get_browser_process_handler
func cefing_app_get_browser_process_handler(self *CAppT) *CBrowserProcessHandlerT {
	Logf("L48:")

	handler := browser_process_handler[self]
	if handler == nil {
		Logf("L77: No Browser Process Handler")
	} else {
		BaseAddRef(handler)
	}
	return handler

}

// AssocRenderProcessHandler associate a hander to app
func (app *CAppT) AssocRenderProcessHandler(handler *CRenderProcessHandlerT) {
	BaseAddRef(handler)
	render_process_handler[app] = handler
}

///
// Return the handler for functionality specific to the render process. This
// function is called on the render process main thread.
///
//export cefing_app_get_render_process_handler
func cefing_app_get_render_process_handler(self *CAppT) *CRenderProcessHandlerT {
	Logf("L97:")

	handler := render_process_handler[self]
	if handler == nil {
		Logf("L77: No Render Process Handler")
	} else {
		BaseAddRef(handler)
	}
	return handler
}

//on_process_mesage_received call OnProcessMessageRecived method
//export cefing_app_on_before_command_line_processing
func cefing_app_on_before_command_line_processing(self *CAppT, process_type *C.cef_string_t, command_line *CCommandLineT) {
	Logf("L36: app: %p", self)

	f := app_method[self]
	if f == nil {
		Logger.Panicln("L48: on_before_command_line_processing: Noo!")
	}
	pt := string_from_cef_string(process_type)
	f.OnBeforeCommandLineProcessing(self, pt, command_line)
}

//on_pregiser_custom_schemes call OnRegisterCustomSchemes method
//export cefing_app_on_register_custom_schemes
func cefing_app_on_register_custom_schemes(self *CAppT, registrar *CSchemeRegistrarT) {
	Logf("L36: app: %p", self)

	f := app_method[self]
	if f == nil {
		Logger.Panicln("L48: on_before_command_line_processing: Noo!")
	}
	f.OnRegisterCustomSchemes(self, registrar)
}

// DefaultApp is dummy implementation of CClientT
type DefaultApp struct {
}

func (*DefaultApp) OnBeforeCommandLineProcessing(self *CAppT, process_type string, command_line *CCommandLineT) {

}

func (*DefaultApp) OnRegisterCustomSchemes(self *CAppT, registrar *CSchemeRegistrarT) {

}
