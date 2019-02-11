package cefingo

// #include "cefingo.h"
import "C"
import "unsafe"

func (m *CProcessMessageT) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(m))
}

func ProcessMessageCreate(name string) *CProcessMessageT {
	cef_name := create_cef_string(name)
	defer clear_cef_string(cef_name)

	msg := (*CProcessMessageT)(C.cef_process_message_create(cef_name))
	BaseAddRef(msg)
	return msg
}

func (self *CProcessMessageT) IsValid() bool {
	status := C.cefingo_process_message_is_valid((*C.cef_process_message_t)(self))
	return status == 1
}

func (self *CProcessMessageT) IsReadOnly() bool {
	status := C.cefingo_process_message_is_read_only((*C.cef_process_message_t)(self))
	return status == 1
}

func (self *CProcessMessageT) Copy() *CProcessMessageT {
	copy := (*CProcessMessageT)(C.cefingo_process_message_copy((*C.cef_process_message_t)(self)))
	BaseAddRef(copy)
	return copy
}

func (self *CProcessMessageT) GetName() string {
	usfs := C.cefingo_process_message_get_name((*C.cef_process_message_t)(self))
	name := string_from_cef_string((*C.cef_string_t)(usfs))
	C.cef_string_userfree_free(usfs)
	return name
}

func (self *CProcessMessageT) GetArgumentList() *CListValueT {
	l := (*CListValueT)(C.cefingo_process_message_get_argument_list((*C.cef_process_message_t)(self)))
	BaseAddRef(l)
	return l
}
