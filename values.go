package cefingo

import (
	"unsafe"
)

// #include "cefingo.h"
import "C"

func ValueCreate() *CValueT {
	return (*CValueT)(C.cef_value_create())
}

func BinaryValueCreate(data []byte) *CBinaryValueT {
	return (*CBinaryValueT)(C.cef_binary_value_create(unsafe.Pointer(&data[0]), C.size_t(len(data))))
}

func DictionaryValueCreate() *CDictionaryValueT {
	return (*CDictionaryValueT)(C.cef_dictionary_value_create())
}

func ListValueCreate() *CListValueT {
	return (*CListValueT)(C.cef_list_value_create())
}
