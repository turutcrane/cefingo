package capi

//go:generate go run ../tools/gen_cef_types.go

import (
	"log"
	"os"
	"reflect"
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

type CLangSizeT C.size_t
type CWindowHandleT C.cef_window_handle_t
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

	var i int
	const maxUint = ^uint(0)
	const minUint = 0
	const maxInt = int(maxUint >> 1)
	const minInt = -maxInt - 1
	Logger = log.New(os.Stdout, "init", log.LstdFlags)
	Logf("Size of var (reflect.TypeOf.Size): %d\n", reflect.TypeOf(i).Size())
	Logf("Size of var (unsafe.Sizeof): %d\n", unsafe.Sizeof(i))
	Logf("maxUint: %d\n", maxUint)
	Logf("maxInt: %d\n", maxInt)
}

var main_args = C.cef_main_args_t{}

func ExecuteProcess(app *CAppT) {

	setup_main_args(&main_args)

	///
	// This function should be called from the application entry point function to
	// execute a secondary process. It can be used to run secondary processes from
	// the browser client executable (default behavior) or from a separate
	// executable specified by the CefSettings.browser_subprocess_path value. If
	// called for the browser process (identified by no "type" command-line value)
	// it will return immediately with a value of -1. If called for a recognized
	// secondary process it will block until the process should exit and then return
	// the process exit code. The |application| parameter may be NULL. The
	// |windows_sandbox_info| parameter is only used on Windows and may be NULL (see
	// cef_sandbox_win.h for details).
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#L116-L130
	///
	BaseAddRef(app.p_app) // ??
	Logf("L243: %t", BaseHasOneRef(app.p_app))
	code := C.cef_execute_process(&main_args, app.p_app, nil)
	Logf("L245: code: %d: %t", code, BaseHasOneRef(app.p_app))
	if code >= 0 {
		os.Exit(int(code))
	}
}

func construct_settings(s Settings) C.cef_settings_t {
	settings := C.cef_settings_t{}
	// Application settings. It is mandatory to set the
	// "size" member.
	settings.size = C.sizeof_cef_settings_t

	settings.log_severity = (C.cef_log_severity_t)(s.LogSeverity) // C.LOGSEVERITY_WARNING // Show only warnings/errors
	settings.no_sandbox = (C.int)(s.NoSandbox)
	settings.multi_threaded_message_loop = (C.int)(s.MultiThreadedMessageLoop)
	settings.remote_debugging_port = (C.int)(s.RemoteDebuggingPort) //8088

	return settings
}

func Initialize(s Settings, app *CAppT) {
	settings := construct_settings(s)

	var c_status C.int
	// resource_path := "C:\\DiskC\\dev\\cef2go\\cef_binary_3.3282.1742.g96f907e_windows64\\Resources"
	// c_resource_path := C.CString(resource_path)
	// defer C.free(c_resource_paht)
	// status := C.cef_string_from_utf8(c_resource_path, C.strlen(c_resource_path), &C.resource_path)
	// log.Println("cef_string_from_utf8", status)
	// settings.resources_dir_path = C.resource_path

	///
	// This function should be called on the main application thread to initialize
	// the CEF browser process. The |application| parameter may be NULL. A return
	// value of true (1) indicates that it succeeded and false (0) indicates that it
	// failed. The |windows_sandbox_info| parameter is only used on Windows and may
	// be NULL (see cef_sandbox_win.h for details).
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#L131-L142
	///
	BaseAddRef(app.p_app) // ??
	c_status = C.cef_initialize(&main_args, &settings, app.p_app, nil)
	Logf("L51: cef_initialize: %d", c_status)

}

// func RunMessageLoop() {
// 	///
// 	// Run the CEF message loop. Use this function instead of an application-
// 	// provided message loop to get the best balance between performance and CPU
// 	// usage. This function should only be called on the main application thread and
// 	// only if cef_initialize() is called with a
// 	// CefSettings.multi_threaded_message_loop value of false (0). This function
// 	// will block until a quit message is received by the system.
// 	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#L167-L175
// 	///
// 	C.cef_run_message_loop()
// }

// // QuitMessageLoop
// func QuitMessageLoop() {
// 	Logf("L166:")
// 	C.cef_quit_message_loop()
// }

// // Shutdown CEF
// func Shutdown() {
// 	Logf("L118:")
// 	C.cef_shutdown()
// }

func BrowserHostCreateBrowser(window_name, url_string string, client *CClientT) {
	Logf("L330:")
	// Window info
	window_info := C.cef_window_info_t{}
	window_info.style = C.WS_OVERLAPPEDWINDOW | C.WS_CLIPCHILDREN |
		C.WS_CLIPSIBLINGS | C.WS_VISIBLE
	window_info.parent_window = nil
	window_info.x = C.CW_USEDEFAULT
	window_info.y = C.CW_USEDEFAULT
	window_info.width = C.CW_USEDEFAULT
	window_info.height = C.CW_USEDEFAULT

	cef_window_name := create_cef_string(window_name)
	defer clear_cef_string(cef_window_name)

	window_info.window_name = *cef_window_name // Do not clear window_info.window_name

	// Initial url
	cef_url := create_cef_string(url_string)
	defer clear_cef_string(cef_url)

	// Browser settings. It is mandatory to set the
	// "size" member.
	browser_settings := C.cef_browser_settings_t{}
	browser_settings.size = C.sizeof_cef_browser_settings_t

	///
	// Create a new browser window using the window parameters specified by
	// |windowInfo|. All values will be copied internally and the actual window will
	// be created on the UI thread. If |request_context| is NULL the global request
	// context will be used. This function can be called on any browser process
	// thread and will not block.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_browser_capi.h#L842-L854
	///
	BaseAddRef(client.p_client) // ?? Need not? clinet.p_client is CefClientCToCpp::Wrap-ed
	C.cef_browser_host_create_browser(
		&window_info,
		client.p_client,
		cef_url,
		&browser_settings, nil, nil,
	)
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

func create_cef_string(s string) *C.cef_string_t {
	cs := C.cef_string_t{}
	set_cef_string(&cs, s)
	return &cs
}

func create_cef_string_from_byte_array(b []byte) *C.cef_string_t {
	cs := C.cef_string_t{}
	set_cef_string_from_byte_array(&cs, b)
	return &cs
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
