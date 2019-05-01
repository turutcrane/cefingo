package capi

// #include "cefingo.h"
import "C"

func ProcessMessageCreate(name string) *CProcessMessageT {
	cef_name := create_cef_string(name)
	defer clear_cef_string(cef_name)

	return newCProcessMessageT(C.cef_process_message_create(cef_name))
}

func (self *CProcessMessageT) IsValid() bool {
	status := C.cefingo_process_message_is_valid(self.p_process_message)
	return status == 1
}

func (self *CProcessMessageT) IsReadOnly() bool {
	status := C.cefingo_process_message_is_read_only(self.p_process_message)
	return status == 1
}

func (self *CProcessMessageT) Copy() *CProcessMessageT {
	copy := C.cefingo_process_message_copy(self.p_process_message)
	return newCProcessMessageT(copy)
}

func (self *CProcessMessageT) GetName() string {
	usfs := C.cefingo_process_message_get_name(self.p_process_message)
	name := string_from_cef_string((*C.cef_string_t)(usfs))
	C.cef_string_userfree_free(usfs)
	return name
}

func (self *CProcessMessageT) GetArgumentList() *CListValueT {
	l := C.cefingo_process_message_get_argument_list(self.p_process_message)
	return newCListValueT(l)
}
