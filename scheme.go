package cefingo

// #include "cefingo.h"
import "C"

type SchemeHandlerFactory interface {
	///
	// Return a new resource handler instance to handle the request or an NULL
	// reference to allow default handling of the request. |browser| and |frame|
	// will be the browser window and frame respectively that originated the
	// request or NULL if the request did not originate from a browser window (for
	// example, if the request came from cef_urlrequest_t). The |request| object
	// passed to this function will not contain cookie data.
	///
	Create(
		self *CSchemeHandlerFactoryT,
		browser *CBrowserT,
		frame *CFrameT,
		scheme_name string,
		request *CRequestT,
	) *CResourceHandlerT
}

///
// Register a scheme handler factory with the global request context. An NULL
// |domain_name| value for a standard scheme will cause the factory to match all
// domain names. The |domain_name| value will be ignored for non-standard
// schemes. If |scheme_name| is a built-in scheme and no handler is returned by
// |factory| then the built-in scheme handler factory will be called. If
// |scheme_name| is a custom scheme then you must also implement the
// cef_app_t::on_register_custom_schemes() function in all processes. This
// function may be called multiple times to change or remove the factory that
// matches the specified |scheme_name| and optional |domain_name|. Returns false
// (0) if an error occurs. This function may be called on any thread in the
// browser process. Using this function is equivalent to calling cef_request_tCo
// ntext::cef_request_context_get_global_context()->register_scheme_handler_fact
// ory().
///
func RegisterSchemeHandlerFactory(
	scheme_name string,
	domain_name string,
	factory *CSchemeHandlerFactoryT,
) int {
	s := create_cef_string(scheme_name)
	defer clear_cef_string(s)

	var d *C.cef_string_t
	if len(domain_name) > 0 {
		d := create_cef_string(domain_name)
		defer clear_cef_string(d)
	}

	return (int)(C.cef_register_scheme_handler_factory(
		s, d,
		(*C.cef_scheme_handler_factory_t)(factory)))
}

var scheme_handler_factory_method = map[*CSchemeHandlerFactoryT]SchemeHandlerFactory{}

func AllocCSchemeHandlerFactoryT(f SchemeHandlerFactory) (factory *CSchemeHandlerFactoryT) {
	p := C.calloc(1, C.sizeof_cefingo_scheme_handler_factory_wrapper_t)

	fp := (*C.cefingo_scheme_handler_factory_wrapper_t)(p)
	C.cefingo_construct_scheme_handler_factory(fp)

	cFactory := (*CSchemeHandlerFactoryT)(p)
	BaseAddRef(cFactory)
	scheme_handler_factory_method[cFactory] = f
	return cFactory
}

//export cefingo_scheme_handler_factory_create
func cefingo_scheme_handler_factory_create(
	self *CSchemeHandlerFactoryT,
	browser *CBrowserT,
	frame *CFrameT,
	scheme_name *CStringT,
	request *CRequestT,
) *CResourceHandlerT {
	f := scheme_handler_factory_method[self]
	if f == nil {
		Logf("L70: No Scheme Factory ")
	}
	s := string_from_cef_string((*C.cef_string_t)(scheme_name))
	return f.Create(self, browser, frame, s, request)
}

func (self *CSchemeRegistrarT) AddCustomScheme(
	scheme_name string,
	is_standard bool,
	is_local bool,
	is_display_isolated bool,
	is_secure bool,
	is_cors_enabled bool,
	is_csp_bypassing bool,
) int {
	s := create_cef_string(scheme_name)
	defer clear_cef_string(s)
	var standard C.int
	if is_standard {
		standard = 1
	}
	var local C.int
	if is_local {
		local = 1
	}
	var display_isolated C.int
	if is_display_isolated {
		display_isolated = 1
	}
	var secure C.int
	if is_secure {
		secure = 1
	}
	var cors_enabled C.int
	if is_cors_enabled {
		cors_enabled = 1
	}
	var csp_bypassing C.int
	if is_csp_bypassing {
		csp_bypassing = 1
	}
	return (int)(C.cefingo_scheme_registrar_add_custom_scheme(
		(*C.cef_scheme_registrar_t)(self),
		s,
		standard, local, display_isolated, secure, cors_enabled, csp_bypassing))

}
