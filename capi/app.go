package capi

import (
	"sync"
	"unsafe"
)

// #include "cefingo.h"
import "C"

// Client is Go interface of C.cef_app_t
type OnBeforeCommandLineProcessingHandler interface {
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
}

type OnRegisterCustomSchemesHandler interface {
	///
	// Provides an opportunity to register custom schemes. Do not keep a reference
	// to the |registrar| object. This function is called on the main thread for
	// each process and the registered schemes should be the same across all
	// processes.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#82-L90
	///
	OnRegisterCustomSchemes(self *CAppT, registrar *CSchemeRegistrarT)
}

var appHandlers = struct {
	m                                         sync.Mutex
	on_register_custom_schemes_handler        map[*C.cef_app_t]OnRegisterCustomSchemesHandler
	on_before_command_line_processing_handler map[*C.cef_app_t]OnBeforeCommandLineProcessingHandler
	browser_process_handler                   map[*C.cef_app_t]*CBrowserProcessHandlerT
	render_process_handler                    map[*C.cef_app_t]*CRenderProcessHandlerT
}{
	sync.Mutex{},
	map[*C.cef_app_t]OnRegisterCustomSchemesHandler{},
	map[*C.cef_app_t]OnBeforeCommandLineProcessingHandler{},
	map[*C.cef_app_t]*CBrowserProcessHandlerT{},
	map[*C.cef_app_t]*CRenderProcessHandlerT{},
}

// AllocCAppT allocates CAppT and construct it
func AllocCAppT() *CAppT {
	up := c_calloc(1, C.sizeof_cefingo_app_wrapper_t, "T58:")
	cefp := C.cefingo_construct_app((*C.cefingo_app_wrapper_t)(up))

	registerDeassocer(up, DeassocFunc(func() {
		// Do not have reference to capp itself in DeassocFunc,
		// or app is never GCed.
		Tracef(up, "T67:")

		appHandlers.m.Lock()
		defer appHandlers.m.Unlock()

		delete(appHandlers.on_before_command_line_processing_handler, cefp)
		delete(appHandlers.on_register_custom_schemes_handler, cefp)
		delete(appHandlers.browser_process_handler, cefp)
		delete(appHandlers.render_process_handler, cefp)
	}))

	return newCAppT(cefp)
}

func (capp *CAppT) Bind(a interface{}) *CAppT {
	cp := capp.p_app

	appHandlers.m.Lock()
	defer appHandlers.m.Unlock()
	if h, ok := a.(OnBeforeCommandLineProcessingHandler); ok {
		appHandlers.on_before_command_line_processing_handler[cp] = h
	}

	if h, ok := a.(OnRegisterCustomSchemesHandler); ok {
		appHandlers.on_register_custom_schemes_handler[cp] = h
	}

	if accessor, ok := a.(CAppTAccessor); ok {
		accessor.SetCAppT(capp)
		Logf("T109:")
	}

	return capp
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
	appHandlers.m.Lock()
	defer appHandlers.m.Unlock()

	appHandlers.browser_process_handler[ap] = handler
}

///
// Return the handler for functionality specific to the browser process. This
// function is called on multiple threads in the browser process.
///
//export cefing_app_get_browser_process_handler
func cefing_app_get_browser_process_handler(self *C.cef_app_t) (ch *C.cef_browser_process_handler_t) {

	appHandlers.m.Lock()
	handler := appHandlers.browser_process_handler[self]
	appHandlers.m.Unlock()

	if handler == nil {
		Logf("T77: No Browser Process Handler")
	} else {
		BaseAddRef(handler.p_browser_process_handler)
		ch = handler.p_browser_process_handler
	}
	return ch
}

// AssocRenderProcessHandler associate a hander to app
func (app *CAppT) AssocRenderProcessHandler(handler *CRenderProcessHandlerT) {
	ap := app.p_app
	appHandlers.m.Lock()
	defer appHandlers.m.Unlock()

	appHandlers.render_process_handler[ap] = handler
}

///
// Return the handler for functionality specific to the render process. This
// function is called on the render process main thread.
///
//export cefing_app_get_render_process_handler
func cefing_app_get_render_process_handler(self *C.cef_app_t) (h *C.cef_render_process_handler_t) {
	appHandlers.m.Lock()
	handler := appHandlers.render_process_handler[self]
	appHandlers.m.Unlock()

	if handler == nil {
		Logf("T77: No Render Process Handler")
	} else {
		h = handler.p_render_process_handler
		BaseAddRef(h) // ??
	}
	return h
}

//on_process_mesage_received call OnProcessMessageRecived method
//export cefing_app_on_before_command_line_processing
func cefing_app_on_before_command_line_processing(self *C.cef_app_t, process_type *C.cef_string_t, command_line *CCommandLineT) {
	Tracef(unsafe.Pointer(self), "T36:")
	appHandlers.m.Lock()
	f := appHandlers.on_before_command_line_processing_handler[self]
	appHandlers.m.Unlock()

	if f != nil {
		pt := string_from_cef_string(process_type)
		app := newCAppT(self)
		f.OnBeforeCommandLineProcessing(app, pt, command_line)
	} else {
		Logf("T48: on_before_command_line_processing: Noo!")
	}
}

//on_pregiser_custom_schemes call OnRegisterCustomSchemes method
//export cefing_app_on_register_custom_schemes
func cefing_app_on_register_custom_schemes(self *C.cef_app_t, registrar *CSchemeRegistrarT) {
	Tracef(unsafe.Pointer(self), "T36:")
	appHandlers.m.Lock()
	f := appHandlers.on_register_custom_schemes_handler[self]
	appHandlers.m.Unlock()

	if f != nil {
		app := newCAppT(self)
		f.OnRegisterCustomSchemes(app, registrar)
	} else {
		Logf("T48: on_before_command_line_processing: Noo!")
	}
}
