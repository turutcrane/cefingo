package cefingo

import (
	"log"
	"os"
	"unsafe"

	"github.com/alexbrainman/gowingui/winapi"

	/*
		#include "cefingo.h"
	*/
	"C"
)

type Cint C.int
type LogSeverityT C.cef_log_severity_t

const (
	LogSeverityDefault = C.LOGSEVERITY_DEFAULT
	LogSeverityVerbose = C.LOGSEVERITY_VERBOSE
	LogSeverityDebug   = C.LOGSEVERITY_DEBUG
	LogSeverityInfo    = C.LOGSEVERITY_INFO
	LogSeverityWarning = C.LOGSEVERITY_WARNING
	LogSeverityError   = C.LOGSEVERITY_ERROR
	LogSeverityDisable = C.LOGSEVERITY_DISABLE
)

type Settings struct {
	LogSeverity              LogSeverityT
	NoSandbox                Cint
	MultiThreadedMessageLoop Cint
}

// type CefSettingsT struct {
// }

// Go Equivalent Type of C.cef_xxx
type CAppT C.cef_app_t
type CBrowserT C.cef_browser_t
type CClientT C.cef_client_t
type CFrameT C.cef_frame_t
type CDomnodeT C.cef_domnode_t
type CListValueT C.cef_list_value_t
type CProcessIdT C.cef_process_id_t
type CProcessMessageT C.cef_process_message_t
type CStringT C.cef_string_t
type CCommandLineT C.cef_command_line_t
type CSchemeRegistrarT C.cef_scheme_registrar_t
type CV8accessorT C.cef_v8accessor_t
type CV8contextT C.cef_v8context_t
type CV8exceptionT C.cef_v8exception_t
type CV8interceptorT C.cef_v8interceptor_t
type CV8stackTraceT C.cef_v8stack_trace_t
type CV8valueT C.cef_v8value_t

type CBrowserProcessHandlerT C.cef_browser_process_handler_t
type CContextMenuHandlerT C.cef_context_menu_handler_t
type CDialogHandlerT C.cef_dialog_handler_t
type CDisplayHandlerT C.cef_display_handler_t
type CDownloaddHanderT C.cef_download_handler_t
type CDragHandlerT C.cef_drag_handler_t
type CFindHandlerT C.cef_find_handler_t
type CFocusHanderT C.cef_focus_handler_t
type CJsdialogHandlerT C.cef_jsdialog_handler_t
type CKeyboardHandlerT C.cef_keyboard_handler_t
type CLifeSpanHandlerT C.cef_life_span_handler_t
type CLoadHandlerT C.cef_load_handler_t
type CRenderHandlerT C.cef_render_handler_t
type CRequestHandlerT C.cef_request_handler_t
type CResourceBundleHanderT C.cef_resource_bundle_handler_t
type CRenderProcessHandlerT C.cef_render_process_handler_t

func init() {
	// Check cef library version
	cefVersionMajor := C.cef_version_info(0)
	cefCommitNumber := C.cef_version_info(1)
	chromeVersionMajor := C.cef_version_info(2)
	// chromeVersionMinor := C.cef_version_info(3)
	chromeVersionBuild := C.cef_version_info(4)
	// chromeVersionPatch := C.cef_version_info(5)
	if cefVersionMajor != C.CEF_VERSION_MAJOR || chromeVersionMajor != C.CHROME_VERSION_MAJOR {
		Logf("build lib: cef_binary_%d.%d.%d (chrome:%d)", C.CEF_VERSION_MAJOR, C.CHROME_VERSION_BUILD, C.CEF_COMMIT_NUMBER, C.CHROME_VERSION_MAJOR)
		Logf("load  lib: cef_binary_%d.%d.%d (chrome:%d)", cefVersionMajor, chromeVersionBuild, cefCommitNumber, chromeVersionMajor)
		log.Panicln("Cef Library mismatch!")
	}
}

var main_args = C.cef_main_args_t{}

func ExecuteProcess(app *CAppT) {

	instance, err := winapi.GetModuleHandle(nil)
	if err != nil {
		log.Panicln(err)
	}

	main_args.instance = C.HINSTANCE(unsafe.Pointer(instance))
	Logf("L33: %T: %#v :: %T: %#v", instance, instance, main_args.instance, main_args.instance)

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
	BaseAddRef(app)
	code := C.cef_execute_process(&main_args, (*C.cef_app_t)(unsafe.Pointer(app)), nil)
	Logf("L37: code: %d", code)
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
	BaseAddRef(app)
	c_status = C.cef_initialize(&main_args, &settings, (*C.cef_app_t)(unsafe.Pointer(app)), nil)
	Logf("L51: cef_initialize: %d", c_status)

}

func RunMessageLoop() {
	///
	// Run the CEF message loop. Use this function instead of an application-
	// provided message loop to get the best balance between performance and CPU
	// usage. This function should only be called on the main application thread and
	// only if cef_initialize() is called with a
	// CefSettings.multi_threaded_message_loop value of false (0). This function
	// will block until a quit message is received by the system.
	// https://github.com/chromiumembedded/cef/blob/3497/include/capi/cef_app_capi.h#L167-L175
	///
	C.cef_run_message_loop()
}

// QuitMessageLoop
func QuitMessageLoop() {
	Logf("L166:")
	C.cef_quit_message_loop()
}

// Shutdown CEF
func Shutdown() {
	Logf("L118:")
	C.cef_shutdown()
}

func BrowserHostCreateBrowser(window_name, url_string string, client *CClientT) {
	Logf("L97:")
	// Window info
	window_info := C.cef_window_info_t{}
	window_info.style = C.WS_OVERLAPPEDWINDOW | C.WS_CLIPCHILDREN |
		C.WS_CLIPSIBLINGS | C.WS_VISIBLE
	window_info.parent_window = nil
	window_info.x = C.CW_USEDEFAULT
	window_info.y = C.CW_USEDEFAULT
	window_info.width = C.CW_USEDEFAULT
	window_info.height = C.CW_USEDEFAULT

	cef_window_name := C.cef_string_t{}

	c_window_name := C.CString(window_name)
	defer C.free(unsafe.Pointer(c_window_name))
	_ = C.cef_string_utf8_to_utf16(c_window_name, C.strlen(c_window_name),
		&cef_window_name)
	window_info.window_name = cef_window_name

	// Initial url
	url := C.CString(url_string)
	defer C.free(unsafe.Pointer(url))

	cef_url := C.cef_string_t{} // (*C.cef_string_t)(calloc(1, C.sizeof_cef_string_t))
	// defer func() {
	// 	C.cef_string_clear(cef_url)
	// 	C.free(unsafe.Pointer(cef_url))
	// }()

	_ = C.cef_string_utf8_to_utf16(url, C.strlen(url), &cef_url)
	// Logf("L88: cef_string_utf8_to_utf16: %d", c_status)

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
	C.cef_browser_host_create_browser(
		&window_info,
		(*C.cef_client_t)(client),
		&cef_url,
		&browser_settings, nil,
	)

}

func calloc(num C.size_t, size C.size_t) unsafe.Pointer {
	p := C.calloc(num, size)
	if p == nil {
		log.Panicln("L58: Cannot Allocated.")
	}
	return p
}

func create_cef_string(s string) *C.cef_string_t {
	c_string := C.CString(s)
	defer C.free(unsafe.Pointer(c_string))
	cs := C.cef_string_t{}

	status := C.cef_string_from_utf8(c_string, C.strlen(c_string), &cs)
	if status == 0 {
		log.Panicln("Error cef_string_from_utf8")
	}
	return &cs
}

func clear_cef_string(s *C.cef_string_t) {
	C.cef_string_clear(s)
}
