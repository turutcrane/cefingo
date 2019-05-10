package capi

// #include "cefingo.h"
import "C"
import "sync"

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
		d = create_cef_string(domain_name)
		defer clear_cef_string(d)
	}

	BaseAddRef(factory.p_scheme_handler_factory)
	return (int)(C.cef_register_scheme_handler_factory(
		s, d,
		factory.p_scheme_handler_factory))
}

var scheme_handler_factory_methods = struct {
	m                      sync.Mutex
	scheme_handler_factory map[*C.cef_scheme_handler_factory_t]SchemeHandlerFactory
}{
	sync.Mutex{},
	map[*C.cef_scheme_handler_factory_t]SchemeHandlerFactory{},
}

func AllocCSchemeHandlerFactoryT() *CSchemeHandlerFactoryT {
	up := c_calloc(1, C.sizeof_cefingo_scheme_handler_factory_wrapper_t, "L81:")
	cefp := C.cefingo_construct_scheme_handler_factory(
		(*C.cefingo_scheme_handler_factory_wrapper_t)(up))

	registerDeassocer(up, DeassocFunc(func() {
		scheme_handler_factory_methods.m.Lock()
		defer scheme_handler_factory_methods.m.Unlock()

		delete(scheme_handler_factory_methods.scheme_handler_factory, cefp)
	}))

	return newCSchemeHandlerFactoryT(cefp)
}

func (factory *CSchemeHandlerFactoryT) Bind(f SchemeHandlerFactory) *CSchemeHandlerFactoryT {
	cefp := factory.p_scheme_handler_factory
	scheme_handler_factory_methods.m.Lock()
	scheme_handler_factory_methods.scheme_handler_factory[cefp] = f
	scheme_handler_factory_methods.m.Unlock()

	if accessor, ok := f.(CSchemeHandlerFactoryTAccessor); ok {
		accessor.SetCSchemeHandlerFactoryT(factory)
		Logf("L109:")
	}

	return factory
}

//export cefingo_scheme_handler_factory_create
func cefingo_scheme_handler_factory_create(
	self *C.cef_scheme_handler_factory_t,
	browser *C.cef_browser_t,
	frame *C.cef_frame_t,
	scheme_name *C.cef_string_t,
	request *C.cef_request_t,
) *C.cef_resource_handler_t {
	scheme_handler_factory_methods.m.Lock()
	f := scheme_handler_factory_methods.scheme_handler_factory[self]
	scheme_handler_factory_methods.m.Unlock()

	if f == nil {
		Logf("L70: No Scheme Factory ")
	}
	s := string_from_cef_string(scheme_name)
	h := f.Create(
		newCSchemeHandlerFactoryT(self),
		newCBrowserT(browser),
		newCFrameT(frame),
		s,
		newCRequestT(request))
	BaseAddRef(h.p_resource_handler)
	return h.p_resource_handler
}

func (self *CSchemeRegistrarT) AddCustomScheme(
	scheme_name string,
	options CSchemeOptionsT,
) int {
	s := create_cef_string(scheme_name)
	defer clear_cef_string(s)
	return (int)(C.cefingo_scheme_registrar_add_custom_scheme(
		(*C.cef_scheme_registrar_t)(self),
		s,
		C.cef_scheme_options_t(options),
	))
}
