package capi

// #include "cefingo.h"
import "C"
import "errors"

// func StringListAlloc() (sl CStringListT) {
// 	sl = CStringListT(C.cef_string_list_alloc())
// 	return sl
// }

// func StringListSize(list CStringListT) int {
// 	return int(C.cef_string_list_size(C.cef_string_list_t(list)))
// }

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

// ///
// // Append a new value at the end of the string list.
// ///
// func StringListAppend(list CStringListT, value string) {
// 	v := create_cef_string(value)
// 	defer clear_cef_string(v)

// 	C.cef_string_list_append(C.cef_string_list_t(list), (*C.cef_string_t)(v))

// }

// ///
// // Clear the string list.
// ///
// func StringListClear(list CStringListT) {
// 	C.cef_string_list_clear(C.cef_string_list_t(list))
// }

// ///
// // Free the string list.
// ///
// func StringListFree(list CStringListT) {
// 	C.cef_string_list_free(C.cef_string_list_t(list))
// }

// ///
// // Creates a copy of an existing string list.
// ///
// func StringListCopy(list CStringListT) CStringListT {
// 	return CStringListT(C.cef_string_list_copy(C.cef_string_list_t(list)))
// }

///
// Return the value assigned to the specified key.
///
func StringMapFind(smap CStringMapT, key string) (value string, err error) {
	go_key := create_cef_string(key)
	defer clear_cef_string(go_key)

	cstr := C.cef_string_t{}

	status := C.cef_string_map_find(
		C.cef_string_map_t(smap),
		(*C.cef_string_t)(go_key),
		&cstr,
	)
	if status == 1 {
		value = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Not find value specified by key in string map")
	}

	return value, err
}

///
// Return the key at the specified zero-based string map index.
///
func StringMapKey(
	cmap CStringMapT,
	index int,
) (key string, err error) {

	cstr := C.cef_string_t{}

	status := C.cef_string_map_key((C.cef_string_map_t)(cmap), (C.size_t)(index), &cstr)
	if status == 1 {
		key = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Not find key by index in string map")
	}

	return key, err
}

///
// Return the value at the specified zero-based string map index.
///
func StringMapValue(
	cmap CStringMapT,
	index int,
) (value string, err error) {
	cstr := C.cef_string_t{}

	status := C.cef_string_map_value((C.cef_string_map_t)(cmap), (C.size_t)(index), &cstr)

	if status == 1 {
		value = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Not find value by index in string map")
	}
	return value, err
}

///
// Return the value_index-th value with the specified key.
///
func StringMultimapEnumerate(
	cmap CStringMultimapT,
	key string,
	value_index int,
) (value string, err error) {
	c_key := create_cef_string(key)
	defer clear_cef_string(c_key)

	cstr := C.cef_string_t{}

	status := C.cef_string_multimap_enumerate((C.cef_string_multimap_t)(cmap), (*C.cef_string_t)(c_key), (C.size_t)(value_index), &cstr)
	if status == 1 {
		value = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Not find value by index in string_multi_map")
	}
	return value, err
}

///
// Return the key at the specified zero-based string multimap index.
///
func StringMultimapKey(
	cmap CStringMultimapT,
	index int,
) (key string, err error) {
	cstr := C.cef_string_t{}

	status := C.cef_string_multimap_key((C.cef_string_multimap_t)(cmap), (C.size_t)(index), &cstr)
	if status == 1 {
		key = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Not find key by index in string multi map")
	}
	return key, err
}

///
// Return the value at the specified zero-based string multimap index.
///
func StringMultimapValue(
	cmap CStringMultimapT,
	index int,
) (value string, err error) {
	cstr := C.cef_string_t{}

	status := C.cef_string_multimap_value((C.cef_string_multimap_t)(cmap), (C.size_t)(index), &cstr)
	if status == 1 {
		value = string_from_cef_string(&cstr)
	} else {
		err = errors.New("Not find key by index in string multi map")
	}
	return value, err
}
