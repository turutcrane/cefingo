package capi

import (
	"sync"
	"unsafe"
)

// #include "cefingo.h"
import "C"

///
// Begin processing the request. To handle the request return true (1) and
// call cef_callback_t::cont() once the response header information is
// available (cef_callback_t::cont() can also be called from inside this
// function if header information is available immediately). To cancel the
// request return false (0).
///
type ProcessRequestHandler interface {
	ProcessRequest(
		self *CResourceHandlerT,
		request *CRequestT,
		callback *CCallbackT,
	) bool
}

///
// Retrieve response header information. If the response length is not known
// set |response_length| to -1 and read_response() will be called until it
// returns false (0). If the response length is known set |response_length| to
// a positive value and read_response() will be called until it returns false
// (0) or the specified number of bytes have been read. Use the |response|
// object to set the mime type, http status code and other optional header
// values. To redirect the request to a new URL set |redirectUrl| to the new
// URL. If an error occured while setting up the request you can call
// set_error() on |response| to indicate the error condition.
///
type GetResponseHeadersHandler interface {
	GetResponseHeaders(
		self *CResourceHandlerT,
		response *CResponseT,
		response_length *int64,
		redirectUrl *string)
}

///
// Read response data. If data is available immediately copy up to
// |bytes_to_read| bytes into |data_out|, set |bytes_read| to the number of
// bytes copied, and return true (1). To read the data at a later time set
// |bytes_read| to 0, return true (1) and call cef_callback_t::cont() when the
// data is available. To indicate response completion return false (0).
///
type ReadResponseHandler interface {
	ReadResponse(
		self *CResourceHandlerT,
		data_out []byte,
		bytes_to_read int,
		bytes_read *int,
		callback *CCallbackT,
	) bool
}

///
// Return true (1) if the specified cookie can be sent with the request or
// false (0) otherwise. If false (0) is returned for any cookie then no
// cookies will be sent with the request.
///
type CanGetCookieHandler interface {
	CanGetCookie(
		self *CResourceHandlerT,
		cookie *CCookieT,
	) bool
}

///
// Return true (1) if the specified cookie returned with the response can be
// set or false (0) otherwise.
///
type CanSetCookieHandler interface {
	CanSetCookie(
		self *CResourceHandlerT,
		cookie *CCookieT,
	) bool
}

///
// Request processing has been canceled.
///
type CancelHandler interface {
	Cancel(self *CResourceHandlerT)
}

var resourceHandlers = struct {
	m                            sync.Mutex
	process_request_handler      map[*C.cef_resource_handler_t]ProcessRequestHandler
	get_response_headers_handler map[*C.cef_resource_handler_t]GetResponseHeadersHandler
	read_response_handler        map[*C.cef_resource_handler_t]ReadResponseHandler
	can_get_cookie_handler       map[*C.cef_resource_handler_t]CanGetCookieHandler
	can_set_cookie_handler       map[*C.cef_resource_handler_t]CanSetCookieHandler
	cancel_handler               map[*C.cef_resource_handler_t]CancelHandler
}{
	sync.Mutex{},
	map[*C.cef_resource_handler_t]ProcessRequestHandler{},
	map[*C.cef_resource_handler_t]GetResponseHeadersHandler{},
	map[*C.cef_resource_handler_t]ReadResponseHandler{},
	map[*C.cef_resource_handler_t]CanGetCookieHandler{},
	map[*C.cef_resource_handler_t]CanSetCookieHandler{},
	map[*C.cef_resource_handler_t]CancelHandler{},
}

func AllocCResourceHanderT() *CResourceHandlerT {
	up := c_calloc(1, C.sizeof_cefingo_resource_handler_wrapper_t, "T102:")
	cefp := C.cefingo_construct_resource_handler((*C.cefingo_resource_handler_wrapper_t)(up))

	registerDeassocer(up, DeassocFunc(func() {
		Tracef(up, "T108:")
		resourceHandlers.m.Lock()
		defer resourceHandlers.m.Unlock()

		delete(resourceHandlers.process_request_handler, cefp)
		delete(resourceHandlers.get_response_headers_handler, cefp)
		delete(resourceHandlers.read_response_handler, cefp)
		delete(resourceHandlers.can_get_cookie_handler, cefp)
		delete(resourceHandlers.can_set_cookie_handler, cefp)
		delete(resourceHandlers.cancel_handler, cefp)
	}))

	return newCResourceHandlerT(cefp)
}

func (rh *CResourceHandlerT) Bind(handler interface{}) *CResourceHandlerT {
	cefp := rh.p_resource_handler
	resourceHandlers.m.Lock()
	defer resourceHandlers.m.Unlock()

	if h, ok := handler.(ProcessRequestHandler); ok {
		resourceHandlers.process_request_handler[cefp] = h
	}

	if h, ok := handler.(GetResponseHeadersHandler); ok {
		resourceHandlers.get_response_headers_handler[cefp] = h
	}

	if h, ok := handler.(ReadResponseHandler); ok {
		resourceHandlers.read_response_handler[cefp] = h
	}

	if h, ok := handler.(CanGetCookieHandler); ok {
		resourceHandlers.can_get_cookie_handler[cefp] = h
	}

	if h, ok := handler.(CanSetCookieHandler); ok {
		resourceHandlers.can_set_cookie_handler[cefp] = h
	}

	if h, ok := handler.(CancelHandler); ok {
		resourceHandlers.cancel_handler[cefp] = h
	}

	if accessor, ok := handler.(CResourceHandlerTAccessor); ok {
		accessor.SetCResourceHandlerT(rh)
		Logf("T180:")
	}

	return rh
}

//export cefingo_resource_handler_process_request
func cefingo_resource_handler_process_request(
	self *C.cef_resource_handler_t,
	request *C.cef_request_t,
	callback *C.cef_callback_t,
) (ret C.int) {
	resourceHandlers.m.Lock()
	h := resourceHandlers.process_request_handler[self]
	resourceHandlers.m.Unlock()

	if h != nil {
		if h.ProcessRequest(
			newCResourceHandlerT(self),
			newCRequestT(request),
			newCCallbackT(callback),
		) {
			ret = 1
		}
	} else {
		Logf("T139: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_get_response_headers
func cefingo_resource_handler_get_response_headers(
	self *C.cef_resource_handler_t,
	response *C.cef_response_t,
	response_length *int64,
	redirectUrl *C.cef_string_t,
) {
	resourceHandlers.m.Lock()
	h := resourceHandlers.get_response_headers_handler[self]
	resourceHandlers.m.Unlock()

	if h != nil {
		var r string
		h.GetResponseHeaders(
			newCResourceHandlerT(self), newCResponseT(response),
			response_length, &r,
		)
		set_cef_string(redirectUrl, r)
	} else {
		Logf("T156: No Handler")
	}
}

//export cefingo_resource_handler_read_response
func cefingo_resource_handler_read_response(
	self *C.cef_resource_handler_t,
	data_out unsafe.Pointer,
	bytes_to_read C.int,
	bytes_read *C.int,
	callback *C.cef_callback_t,
) (ret C.int) {
	resourceHandlers.m.Lock()
	h := resourceHandlers.read_response_handler[self]
	resourceHandlers.m.Unlock()

	if h != nil {
		buff := (*[1 << 30]byte)(data_out)[:bytes_to_read:bytes_to_read]
		var i int
		if h.ReadResponse(
			newCResourceHandlerT(self),
			buff,
			int(bytes_to_read),
			&i,
			newCCallbackT(callback),
		) {
			ret = 1
		}
		*bytes_read = C.int(i)
	} else {
		Logf("T128: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_can_get_cookie
func cefingo_resource_handler_can_get_cookie(
	self *C.cef_resource_handler_t,
	cookie *CCookieT,
) (ret C.int) {
	resourceHandlers.m.Lock()
	h := resourceHandlers.can_get_cookie_handler[self]
	resourceHandlers.m.Unlock()
	if h != nil {
		if h.CanGetCookie(newCResourceHandlerT(self), cookie) {
			ret = 1
		}
	} else {
		Logf("T142: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_can_set_cookie
func cefingo_resource_handler_can_set_cookie(
	self *C.cef_resource_handler_t,
	cookie *CCookieT,
) (ret C.int) {
	resourceHandlers.m.Lock()
	h := resourceHandlers.can_set_cookie_handler[self]
	resourceHandlers.m.Unlock()

	if h != nil {
		if h.CanSetCookie(newCResourceHandlerT(self), cookie) {
			ret = 1
		}
	} else {
		Logf("T154: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_cancel
func cefingo_resource_handler_cancel(
	self *C.cef_resource_handler_t,
) {
	resourceHandlers.m.Lock()
	h := resourceHandlers.cancel_handler[self]
	resourceHandlers.m.Unlock()

	if h != nil {
		h.Cancel(newCResourceHandlerT(self))
	} else {
		Logf("T154: No Handler")
	}
}

// type DefaultResourceHandler struct {
// }

// func (rh *DefaultResourceHandler) ProcessRequest(
// 	self *CResourceHandlerT,
// 	request *CRequestT,
// 	callback *CCallbackT,
// ) bool {
// 	callback.Cont()
// 	return true
// }

// const sample_text = `
// <html>
// <head>
// </head>
// <body>
// Sample Text: Hello Cefingo!!"
// <br>
// <a href="http://cefingo.internal/aaa">Go</a>
// </body>
// </html>
// `

// func (rh *DefaultResourceHandler) GetResponseHeaders(
// 	self *CResourceHandlerT,
// 	response *CResponseT,
// 	response_length *int64,
// 	redirectUrl *string,
// ) {
// 	response.SetMimeType("text/html")
// 	response.SetStatus(200)
// 	response.SetStatusText("OK")
// 	*response_length = int64(len(sample_text))
// }

// func (fh *DefaultResourceHandler) ReadResponse(
// 	self *CResourceHandlerT,
// 	data_out []byte,
// 	bytes_to_read int,
// 	bytes_read *int,
// 	callback *C.cef_callback_t,
// ) bool {
// 	l := len(sample_text)
// 	buf := []byte(sample_text)
// 	l = min(l, len(buf))
// 	for i, b := range buf[:l] {
// 		data_out[i] = b
// 	}
// 	*bytes_read = l

// 	return true
// }

// func (*DefaultResourceHandler) CanGetCookie(
// 	self *CResourceHandlerT,
// 	cookie *CCookieT,
// ) bool {
// 	return true
// }

// func (*DefaultResourceHandler) CanSetCookie(
// 	self *CResourceHandlerT,
// 	cookie *CCookieT,
// ) bool {
// 	return true
// }

// func (*DefaultResourceHandler) Cancel(self *CResourceHandlerT) {
// }

// func min(x, y int) int {
// 	if x < y {
// 		return x
// 	}
// 	return y
// }
