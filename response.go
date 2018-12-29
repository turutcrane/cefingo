package cefingo

// #include "cefingo.h"
import "C"

///
// Set the response error code. This can be used by custom scheme handlers to
// return errors during initial request processing.
///
func (self *CResponseT) SetError(
	errorCode CErrorcodeT,
) {
	C.cefingo_response_set_error((*C.cef_response_t)(self), C.cef_errorcode_t(errorCode))
}

///
// Set the response status code.
///
func (self *CResponseT) SetStatus(
	status int,
) {
	C.cefingo_response_set_status((*C.cef_response_t)(self), C.int(status))
}

///
// Set the response status text.
///
func (self *CResponseT) SetStatusText(
	statusText string,
) {
	s := create_cef_string(statusText)
	defer clear_cef_string(s)
	C.cefingo_response_set_status_text((*C.cef_response_t)(self), s)
}

///
// Set the response mime type.
///
func (self *CResponseT) SetMimeType(mimeType string) {
	s := create_cef_string(mimeType)
	defer clear_cef_string(s)
	C.cefingo_response_set_mime_type((*C.cef_response_t)(self), s)
}
