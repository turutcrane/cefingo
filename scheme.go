package cefingo

// #include "cefingo.h"
import "C"

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

// func RegisterSchemeHandlerFactory(schemeName string, domain_name string)
// 	return C.cef_register_scheme_handler_factory(
// 	const cef_string_t* scheme_name,
// 	const cef_string_t* domain_name,
// 	cef_scheme_handler_factory_t* factory);
