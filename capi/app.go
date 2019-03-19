package capi

import (
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

var app_method = map[*C.cef_app_t]IApp{}
var browser_process_handler = map[*C.cef_app_t]*CBrowserProcessHandlerT{}
var render_process_handler = map[*C.cef_app_t]*CRenderProcessHandlerT{}

// Client is Go interface of C.cef_app_t
type IApp interface {
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

func newCAppT(cef *C.cef_app_t) *CAppT {
	Tracef(unsafe.Pointer(cef), "L42:")
	BaseAddRef(cef)
	app := CAppT{cef}
	runtime.SetFinalizer(&app, func(a *CAppT) {
		Tracef(unsafe.Pointer(a.p_app), "L47:")
		BaseRelease(a.p_app)
	})
	return &app
}

// AllocCAppT allocates CAppT and construct it
func AllocCAppT(a IApp) (cApp *CAppT) {
	p := c_calloc(1, C.sizeof_cefingo_app_wrapper_t, "L56:")

	C.cefingo_construct_app((*C.cefingo_app_wrapper_t)(p))

	cApp = newCAppT((*C.cef_app_t)(p))
	cp := cApp.p_app
	app_method[cp] = a
	registerDeassocer(unsafe.Pointer(cp), DeassocFunc(func() {
		// Do not have reference to cApp itself in DeassocFunc,
		// or app is never GCed.
		Tracef(unsafe.Pointer(cp), "L67:")
		delete(app_method, cp)
	}))

	return cApp
}

func (a *C.cef_app_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(a))
}

///
// Return the handler for resource bundle events. If
// CefSettings.pack_loading_disabled is true (1) a handler must be returned.
// If no handler is returned resources will be loaded from pack files. This
// function is called by the browser and render processes on multiple threads.
///
//export cefing_app_get_resource_bundle_handler
func cefing_app_get_resource_bundle_handler(self *C.cef_app_t) *CResourceBundleHanderT {
	return nil
}

// AssocBrowserProcessHandler associate a hander to app
func (app *CAppT) AssocBrowserProcessHandler(handler *CBrowserProcessHandlerT) {
	ap := app.p_app
	browser_process_handler[ap] = handler
	registerDeassocer(unsafe.Pointer(ap), DeassocFunc(func() {
		// Do not have reference to app itself in DeassocFunc,
		// or app is never GCed.
		Tracef(unsafe.Pointer(ap), "L95:")
		delete(browser_process_handler, ap)
	}))
}

///
// Return the handler for functionality specific to the browser process. This
// function is called on multiple threads in the browser process.
///
//export cefing_app_get_browser_process_handler
func cefing_app_get_browser_process_handler(self *C.cef_app_t) (ch *C.cef_browser_process_handler_t) {
	handler := browser_process_handler[self]
	if handler == nil {
		Logf("L77: No Browser Process Handler")
	} else {
		BaseAddRef(handler.p_browser_process_handler) // ??
		ch = handler.p_browser_process_handler
	}
	return ch
}

// AssocRenderProcessHandler associate a hander to app
func (app *CAppT) AssocRenderProcessHandler(handler *CRenderProcessHandlerT) {
	ap := app.p_app
	render_process_handler[ap] = handler
	registerDeassocer(unsafe.Pointer(ap), DeassocFunc(func() {
		// Do not have reference to app itself in DeassocFunc,
		// or app is never GCed.
		Tracef(unsafe.Pointer(ap), "L125:")
		delete(browser_process_handler, ap)
	}))
}

///
// Return the handler for functionality specific to the render process. This
// function is called on the render process main thread.
///
//export cefing_app_get_render_process_handler
func cefing_app_get_render_process_handler(self *C.cef_app_t) (h *C.cef_render_process_handler_t) {
	handler := render_process_handler[self]
	if handler == nil {
		Logf("L77: No Render Process Handler")
	} else {
		h = handler.p_render_process_handler
		BaseAddRef(h) // ??
	}
	return h
}

//on_process_mesage_received call OnProcessMessageRecived method
//export cefing_app_on_before_command_line_processing
func cefing_app_on_before_command_line_processing(self *C.cef_app_t, process_type *C.cef_string_t, command_line *CCommandLineT) {
	Tracef(unsafe.Pointer(self), "L36:")

	f := app_method[self]
	if f == nil {
		Logger.Panicln("L48: on_before_command_line_processing: Noo!")
	}
	pt := string_from_cef_string(process_type)
	app := newCAppT(self)
	f.OnBeforeCommandLineProcessing(app, pt, command_line)
}

//on_pregiser_custom_schemes call OnRegisterCustomSchemes method
//export cefing_app_on_register_custom_schemes
func cefing_app_on_register_custom_schemes(self *C.cef_app_t, registrar *CSchemeRegistrarT) {
	Tracef(unsafe.Pointer(self), "L36:")

	f := app_method[self]
	if f == nil {
		Logger.Panicln("L48: on_before_command_line_processing: Noo!")
	}
	app := newCAppT(self)
	f.OnRegisterCustomSchemes(app, registrar)
}

// DefaultApp is dummy implementation of CAppT
type DefaultApp struct {
}

func (*DefaultApp) OnBeforeCommandLineProcessing(self *CAppT, process_type string, command_line *CCommandLineT) {

}

func (*DefaultApp) OnRegisterCustomSchemes(self *CAppT, registrar *CSchemeRegistrarT) {

}
