package capi

import (
	"log"
	"os"
	"reflect"
	"runtime"
	"unsafe"
)

// #cgo pkg-config: cefingo
// #include "cefingo.h"
import "C"

type Settings struct {
	LogSeverity              CLogSeverityT
	NoSandbox                int
	MultiThreadedMessageLoop int
	RemoteDebuggingPort      int
}

// type CLangSizeT C.size_t
type CEventHandleT C.cef_event_handle_t
type CCursorHandleT C.cef_cursor_handle_t

func init() {
	// Check cef library version
	cefVersionMajor := C.cef_version_info(0)
	cefVersionMinor := C.cef_version_info(1)
	cefVersionPatch := C.cef_version_info(2)
	cefCommitNumber := C.cef_version_info(3)
	chromeVersionMajor := C.cef_version_info(4)
	chromeVersionMinor := C.cef_version_info(5)
	chromeVersionBuild := C.cef_version_info(6)
	chromeVersionPatch := C.cef_version_info(7)
	if cefVersionMajor != C.CEF_VERSION_MAJOR || chromeVersionMajor != C.CHROME_VERSION_MAJOR {
		Logger = log.New(os.Stdout, "init", log.LstdFlags)
		Logf("build lib: cef_%d.%d.%d.%d (chrome:%d.%d.%d.%d)",
			C.CEF_VERSION_MAJOR, C.CEF_VERSION_MINOR, C.CEF_VERSION_PATCH, C.CEF_COMMIT_NUMBER,
			C.CHROME_VERSION_MAJOR, C.CHROME_VERSION_MINOR, C.CHROME_VERSION_BUILD, C.CHROME_VERSION_PATCH)
		Logf("load  lib: cef_%d.%d.%d.%d (chrome:%d.%d.%d.%d)",
			cefVersionMajor, cefVersionMinor, cefVersionPatch, cefCommitNumber,
			chromeVersionMajor, chromeVersionMinor, chromeVersionBuild, chromeVersionPatch)
		Panicf("L195: Cef Library mismatch!")
	}
	C.cefingo_init()

}

func DumpInfo() {
	var i int
	const maxUint = ^uint(0)
	const minUint = 0
	const maxInt = int(maxUint >> 1)
	const minInt = -maxInt - 1
	Logf("Size of var (reflect.TypeOf.Size): %d\n", reflect.TypeOf(i).Size())
	Logf("Size of var (unsafe.Sizeof): %d\n", unsafe.Sizeof(i))
	Logf("maxUint: %d\n", maxUint)
	Logf("maxInt: %d\n", maxInt)
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) UNlock() {}

type unrefedBy int

const (
	byApp unrefedBy = iota + 1 // alloc returned takeOvered
	byApi                      // inParam
	byCef                      // passed
	self
	unrefed
)

type cef_string_t struct {
	noCopy         noCopy
	p_cef_string_t *C.cef_string_t
}

func set_cef_string(cs *C.cef_string_t, s string) {
	c_string := C.CString(s)
	defer C.free(unsafe.Pointer(c_string))

	status := C.cef_string_from_utf8(c_string, C.strlen(c_string), cs)
	if status == 0 {
		Panicf("L346: Error cef_string_from_utf8")
	}
}

func set_cef_string_from_byte_array(cs *C.cef_string_t, b []byte) {

	status := C.cef_string_from_utf8((*C.char)(unsafe.Pointer(&b[0])), (C.size_t)(len(b)), cs)
	if status == 0 {
		Panicf("L354: Error cef_string_from_utf8")
	}
}

func create_cef_string(s string) *cef_string_t {
	cs := (*C.cef_string_t)(C.memset(C.malloc(C.sizeof_cef_string_t), 0, C.sizeof_cef_string_t))
	goCs := &cef_string_t{noCopy{}, cs}
	set_cef_string(cs, s)

	runtime.SetFinalizer(goCs, func(goCs *cef_string_t) {
		clear_cef_string(goCs.p_cef_string_t)
		C.free(unsafe.Pointer(goCs.p_cef_string_t))
	})

	return goCs
}

func create_cef_string_from_byte_array(b []byte) *cef_string_t {
	cs := (*C.cef_string_t)(C.memset(C.malloc(C.sizeof_cef_string_t), 0, C.sizeof_cef_string_t))
	goCs := &cef_string_t{noCopy{}, cs}
	set_cef_string_from_byte_array(cs, b)

	runtime.SetFinalizer(goCs, func(goCs *cef_string_t) {
		clear_cef_string(goCs.p_cef_string_t)
		C.free(unsafe.Pointer(goCs.p_cef_string_t))
	})

	return goCs
}

func string_from_cef_string(s *C.cef_string_t) (str string) {
	if s != nil {
		cs := C.cef_string_utf8_t{}
		C.cef_string_to_utf8(s.str, s.length, &cs)
		str = C.GoString(cs.str)
		C.cef_string_utf8_clear(&cs)
	}
	return str
}

func clear_cef_string(s *C.cef_string_t) {
	C.cef_string_clear(s)
}

func ColorSetARGB(a, r, g, b int) CColorT {
	return CColorT(a<<24 | r<<16 | g<<8 | b)
}

func ColorGetA(c CColorT) uint32 {
	return (uint32(c) >> 24) & 0xff
}

func ColorGetR(c CColorT) uint32 {
	return (uint32(c) >> 16) & 0xff
}

func ColorGetG(c CColorT) uint32 {
	return (uint32(c) >> 8) & 0xff
}

func ColorGetB(c CColorT) uint32 {
	return uint32(c) & 0xff
}

func (rect *CRectT) IsEmpty() bool {
	return rect.width <= 0 || rect.height <= 0
}
