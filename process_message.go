package cefingo

// #include "cefingo.h"
import "C"
import (
	"runtime"
	"unsafe"
)

func (m *C.cef_process_message_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(m))
}

func newCProcessMessageT(cef *C.cef_process_message_t) *CProcessMessageT {
	Logf("L42: %p", cef)
	BaseAddRef(cef)
	message := CProcessMessageT{cef}
	runtime.SetFinalizer(&message, func(m *CProcessMessageT) {
		if ref_count_log.output {
			Logf("L47: %p", m.p_process_message)
		}
		BaseRelease(m.p_process_message)
	})
	return &message
}

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
	l := (*CListValueT)(C.cefingo_process_message_get_argument_list(self.p_process_message))
	BaseAddRef(l)
	return l
}
