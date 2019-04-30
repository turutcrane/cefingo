package capi

import (
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

type CResourceHandlerT struct {
	p_resource_handler *C.cef_resource_handler_t
}

type RefToCResourceHandlerT struct {
	rh *CResourceHandlerT
}

type CResourceHandlerTAccessor interface {
	GetCResourceHandlerT() *CResourceHandlerT
	SetCResourceHandlerT(*CResourceHandlerT)
}

func (r RefToCResourceHandlerT) GetCResourceHandlerT() *CResourceHandlerT {
	return r.rh
}

func (r *RefToCResourceHandlerT) SetCResourceHandlerT(c *CResourceHandlerT) {
	r.rh = c
}

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

var process_request_handler = map[*C.cef_resource_handler_t]ProcessRequestHandler{}
var get_response_headers_handler = map[*C.cef_resource_handler_t]GetResponseHeadersHandler{}
var read_response_handler = map[*C.cef_resource_handler_t]ReadResponseHandler{}
var can_get_cookie_handler = map[*C.cef_resource_handler_t]CanGetCookieHandler{}
var can_set_cookie_handler = map[*C.cef_resource_handler_t]CanSetCookieHandler{}
var cancel_handler = map[*C.cef_resource_handler_t]CancelHandler{}

func newCResourceHanderT(cef *C.cef_resource_handler_t) *CResourceHandlerT {
	Tracef(unsafe.Pointer(cef), "L88:")
	BaseAddRef(cef)
	handler := CResourceHandlerT{cef}

	runtime.SetFinalizer(&handler, func(h *CResourceHandlerT) {
		Tracef(unsafe.Pointer(h.p_resource_handler), "L92:")
		BaseRelease(h.p_resource_handler)
	})
	return &handler
}

func AllocCResourceHanderT() *CResourceHandlerT {
	p := (*C.cefingo_resource_handler_wrapper_t)(
		c_calloc(1, C.sizeof_cefingo_resource_handler_wrapper_t, "L102:"))
	C.cefingo_construct_resource_handler(p)

	return newCResourceHanderT(
		(*C.cef_resource_handler_t)(unsafe.Pointer(p)))
}

func (rh *CResourceHandlerT) Bind(handler interface{}) *CResourceHandlerT {
	cefp := rh.p_resource_handler

	if h, ok := handler.(ProcessRequestHandler); ok {
		process_request_handler[cefp] = h
	}

	if h, ok := handler.(GetResponseHeadersHandler); ok {
		get_response_headers_handler[cefp] = h
	}

	if h, ok := handler.(ReadResponseHandler); ok {
		read_response_handler[cefp] = h
	}

	if h, ok := handler.(CanGetCookieHandler); ok {
		can_get_cookie_handler[cefp] = h
	}

	if h, ok := handler.(CanSetCookieHandler); ok {
		can_set_cookie_handler[cefp] = h
	}

	if h, ok := handler.(CancelHandler); ok {
		cancel_handler[cefp] = h
	}

	registerDeassocer(unsafe.Pointer(cefp), DeassocFunc(func() {
		Tracef(unsafe.Pointer(cefp), "L108:")
		delete(process_request_handler, cefp)
		delete(get_response_headers_handler, cefp)
		delete(read_response_handler, cefp)
		delete(can_get_cookie_handler, cefp)
		delete(can_set_cookie_handler, cefp)
		delete(cancel_handler, cefp)
	}))

	if accessor, ok := handler.(CResourceHandlerTAccessor); ok {
		accessor.SetCResourceHandlerT(rh)
		Logf("L180:")
	}

	return rh
}

func (h *C.cef_resource_handler_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(h))
}

//export cefingo_resource_handler_process_request
func cefingo_resource_handler_process_request(
	self *C.cef_resource_handler_t,
	request *CRequestT,
	callback *CCallbackT,
) (ret C.int) {
	h := process_request_handler[self]
	if h != nil {
		if h.ProcessRequest(
			newCResourceHanderT(self),
			request,
			callback,
		) {
			ret = 1
		}
	} else {
		Logf("L139: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_get_response_headers
func cefingo_resource_handler_get_response_headers(
	self *C.cef_resource_handler_t,
	response *CResponseT,
	response_length *int64,
	redirectUrl *C.cef_string_t,
) {
	h := get_response_headers_handler[self]
	if h != nil {
		var r string
		h.GetResponseHeaders(
			newCResourceHanderT(self), response, response_length, &r)
		set_cef_string(redirectUrl, r)
	} else {
		Logf("L156: No Handler")
	}
}

//export cefingo_resource_handler_read_response
func cefingo_resource_handler_read_response(
	self *C.cef_resource_handler_t,
	data_out unsafe.Pointer,
	bytes_to_read C.int,
	bytes_read *C.int,
	callback *CCallbackT,
) (ret C.int) {
	h := read_response_handler[self]
	if h != nil {
		buff := (*[1 << 30]byte)(data_out)[:bytes_to_read:bytes_to_read]
		var i int
		if h.ReadResponse(
			newCResourceHanderT(self), buff, (int)(bytes_to_read), &i, callback) {
			ret = 1
		}
		*bytes_read = C.int(i)
	} else {
		Logf("L128: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_can_get_cookie
func cefingo_resource_handler_can_get_cookie(
	self *C.cef_resource_handler_t,
	cookie *CCookieT,
) (ret C.int) {
	h := can_get_cookie_handler[self]
	if h != nil {
		if h.CanGetCookie(newCResourceHanderT(self), cookie) {
			ret = 1
		}
	} else {
		Logf("L142: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_can_set_cookie
func cefingo_resource_handler_can_set_cookie(
	self *C.cef_resource_handler_t,
	cookie *CCookieT,
) (ret C.int) {
	h := can_set_cookie_handler[self]
	if h != nil {
		if h.CanSetCookie(newCResourceHanderT(self), cookie) {
			ret = 1
		}
	} else {
		Logf("L154: No Handler")
	}
	return ret
}

//export cefingo_resource_handler_cancel
func cefingo_resource_handler_cancel(
	self *C.cef_resource_handler_t,
) {
	h := cancel_handler[self]
	if h != nil {
		h.Cancel(newCResourceHanderT(self))
	} else {
		Logger.Panicln("L154: No Handler")
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
// 	callback *CCallbackT,
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
