package cefingo

// #include "cefingo.h"
import "C"
import "errors"

func StringListAlloc() (sl CStringListT) {
	sl = CStringListT(C.cef_string_list_alloc())
	return sl
}

func StringListSize(list CStringListT) int {
	return int(C.cef_string_list_size(C.cef_string_list_t(list)))
}

///
// Retrieve the value at the specified zero-based string list index. Returns
// true (1) if the value was successfully retrieved.
///
func StringListValue(list CStringListT, index int) (s string, err error) {
	cstr := C.cef_string_t{}

	status := C.cef_string_list_value(C.cef_string_list_t(list), C.size_t(index), &cstr)
	if status == 1 {
		s = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Can not retieve string from string list")
	}
	return s, err
}

///
// Append a new value at the end of the string list.
///
func StringListAppend(list CStringListT, value string) {
	v := create_cef_string(value)
	defer clear_cef_string(v)

	C.cef_string_list_append(C.cef_string_list_t(list), (*C.cef_string_t)(v))

}

///
// Clear the string list.
///
func StringListClear(list CStringListT) {
	C.cef_string_list_clear(C.cef_string_list_t(list))
}

///
// Free the string list.
///
func StringListFree(list CStringListT) {
	C.cef_string_list_free(C.cef_string_list_t(list))
}

///
// Creates a copy of an existing string list.
///
func StringListCopy(list CStringListT) CStringListT {
	return CStringListT(C.cef_string_list_copy(C.cef_string_list_t(list)))
}
