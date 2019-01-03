package cefingo

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

type ResourceHandler interface {
	///
	// Begin processing the request. To handle the request return true (1) and
	// call cef_callback_t::cont() once the response header information is
	// available (cef_callback_t::cont() can also be called from inside this
	// function if header information is available immediately). To cancel the
	// request return false (0).
	///
	ProcessRequest(
		self *CResourceHandlerT,
		request *CRequestT,
		callback *CCallbackT,
	) bool

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
	GetResponseHeaders(
		self *CResourceHandlerT,
		response *CResponseT,
		response_length *int64,
		redirectUrl *string)

	///
	// Read response data. If data is available immediately copy up to
	// |bytes_to_read| bytes into |data_out|, set |bytes_read| to the number of
	// bytes copied, and return true (1). To read the data at a later time set
	// |bytes_read| to 0, return true (1) and call cef_callback_t::cont() when the
	// data is available. To indicate response completion return false (0).
	///
	ReadResponse(
		self *CResourceHandlerT,
		data_out []byte,
		bytes_to_read int,
		bytes_read *int,
		callback *CCallbackT,
	) bool

	///
	// Return true (1) if the specified cookie can be sent with the request or
	// false (0) otherwise. If false (0) is returned for any cookie then no
	// cookies will be sent with the request.
	///
	CanGetCookie(
		self *CResourceHandlerT,
		cookie *CCookieT,
	) bool

	///
	// Return true (1) if the specified cookie returned with the response can be
	// set or false (0) otherwise.
	///
	CanSetCookie(
		self *CResourceHandlerT,
		cookie *CCookieT,
	) bool

	///
	// Request processing has been canceled.
	///
	Cancel(self *CResourceHandlerT)
}

var resource_handler_method = map[*CResourceHandlerT]ResourceHandler{}

func AllocCResourceHanderT(h ResourceHandler) (cHandler *CResourceHandlerT) {
	p := C.calloc(1, C.sizeof_cefingo_resource_handler_wrapper_t)

	hp := (*C.cefingo_resource_handler_wrapper_t)(p)
	C.cefingo_construct_resource_handler(hp)

	cHandler = (*CResourceHandlerT)(p)
	BaseAddRef(cHandler)
	resource_handler_method[cHandler] = h

	return cHandler
}

//export cefingo_resource_handler_process_request
func cefingo_resource_handler_process_request(
	self *CResourceHandlerT,
	request *CRequestT,
	callback *CCallbackT,
) (ret C.int) {
	h := resource_handler_method[self]
	if h == nil {
		Logger.Panicln("L99: No Handler")
	}
	if h.ProcessRequest(self, request, callback) {
		ret = 1
	}
	return ret
}

//export cefingo_resource_handler_get_response_headers
func cefingo_resource_handler_get_response_headers(
	self *CResourceHandlerT,
	response *CResponseT,
	response_length *int64,
	redirectUrl *C.cef_string_t,
) {
	h := resource_handler_method[self]
	if h == nil {
		Logger.Panicln("L113: No Handler")
	}
	var r string
	h.GetResponseHeaders(self, response, response_length, &r)
	set_cef_string(redirectUrl, r)
}

//export cefingo_resource_handler_read_response
func cefingo_resource_handler_read_response(
	self *CResourceHandlerT,
	data_out unsafe.Pointer,
	bytes_to_read C.int,
	bytes_read *C.int,
	callback *CCallbackT,
) (ret C.int) {
	h := resource_handler_method[self]
	if h == nil {
		Logger.Panicln("L128: No Handler")
	}
	buff := (*[1 << 30]byte)(data_out)[:bytes_to_read:bytes_to_read]
	var i int
	if h.ReadResponse(self, buff, (int)(bytes_to_read), &i, callback) {
		ret = 1
	}
	*bytes_read = C.int(i)
	return ret
}

//export cefingo_resource_handler_can_get_cookie
func cefingo_resource_handler_can_get_cookie(
	self *CResourceHandlerT,
	cookie *CCookieT,
) (ret C.int) {
	h := resource_handler_method[self]
	if h == nil {
		Logger.Panicln("L142: No Handler")
	}
	if h.CanGetCookie(self, cookie) {
		ret = 1
	}
	return ret
}

//export cefingo_resource_handler_can_set_cookie
func cefingo_resource_handler_can_set_cookie(
	self *CResourceHandlerT,
	cookie *CCookieT,
) (ret C.int) {
	h := resource_handler_method[self]
	if h == nil {
		Logger.Panicln("L154: No Handler")
	}
	if h.CanSetCookie(self, cookie) {
		ret = 1
	}
	return ret
}

//export cefingo_resource_handler_cancel
func cefingo_resource_handler_cancel(
	self *CResourceHandlerT,
) {
	h := resource_handler_method[self]
	if h == nil {
		Logger.Panicln("L154: No Handler")
	}
	h.Cancel(self)
}

type DefaultResourceHandler struct {
}

func (rh *DefaultResourceHandler) ProcessRequest(
	self *CResourceHandlerT,
	request *CRequestT,
	callback *CCallbackT,
) bool {

	callback.Cont()
	return true
}

const sample_text = `
<html>
<head>
</head>
<body>
Sample Text: Hello Cefingo!!"
<br>
<a href="http://cefingo.internal/aaa">Go</a>
</body>
</html>
`

func (rh *DefaultResourceHandler) GetResponseHeaders(
	self *CResourceHandlerT,
	response *CResponseT,
	response_length *int64,
	redirectUrl *string,
) {
	response.SetMimeType("text/html")
	response.SetStatus(200)
	response.SetStatusText("OK")
	*response_length = int64(len(sample_text))
}

func (fh *DefaultResourceHandler) ReadResponse(
	self *CResourceHandlerT,
	data_out []byte,
	bytes_to_read int,
	bytes_read *int,
	callback *CCallbackT,
) bool {
	l := len(sample_text)
	buf := []byte(sample_text)
	l = min(l, len(buf))
	for i, b := range buf[:l] {
		data_out[i] = b
	}
	*bytes_read = l

	return true
}

func (*DefaultResourceHandler) CanGetCookie(
	self *CResourceHandlerT,
	cookie *CCookieT,
) bool {
	return true
}

func (*DefaultResourceHandler) CanSetCookie(
	self *CResourceHandlerT,
	cookie *CCookieT,
) bool {
	return true
}

func (*DefaultResourceHandler) Cancel(self *CResourceHandlerT) {
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
