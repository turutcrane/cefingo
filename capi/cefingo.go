package capi

import (
	"log"
	"os"
	"unsafe"

	"github.com/alexbrainman/gowingui/winapi"
)

//	#include "cefingo.h"
import "C"

// type Cint C.int
// type CSizeT C.size_t

// type CStringT C.cef_string_t
type CTimeT C.cef_time_t

type CErrorcodeT C.cef_errorcode_t
type CLogSeverityT C.cef_log_severity_t
type CStringListT C.cef_string_list_t
type CTransitionTypeT C.cef_transition_type_t
type CValueTypeT C.cef_value_type_t
type CSchemeOptionsT C.cef_scheme_options_t

const (
	ErrNone            CErrorcodeT = C.ERR_NONE
	ErrFailed          CErrorcodeT = C.ERR_FAILED
	ErrAborted         CErrorcodeT = C.ERR_ABORTED
	ErrInvalidArgument CErrorcodeT = C.ERR_INVALID_ARGUMENT
	ErrInvalidHandle   CErrorcodeT = C.ERR_INVALID_HANDLE
	ErrFileNotFound    CErrorcodeT = C.ERR_FILE_NOT_FOUND
	ErrTimedOut        CErrorcodeT = C.ERR_TIMED_OUT
	ErrFileTooBig      CErrorcodeT = C.ERR_FILE_TOO_BIG
	ErrUnexpected      CErrorcodeT = C.ERR_UNEXPECTED
	ErrAccessDenied    CErrorcodeT = C.ERR_ACCESS_DENIED
	ErrNotImplemented  CErrorcodeT = C.ERR_NOT_IMPLEMENTED

// ERR_CONNECTION_CLOSED = -100,
// ERR_CONNECTION_RESET = -101,
// ERR_CONNECTION_REFUSED = -102,
// ERR_CONNECTION_ABORTED = -103,
// ERR_CONNECTION_FAILED = -104,
// ERR_NAME_NOT_RESOLVED = -105,
// ERR_INTERNET_DISCONNECTED = -106,
// ERR_SSL_PROTOCOL_ERROR = -107,
// ERR_ADDRESS_INVALID = -108,
// ERR_ADDRESS_UNREACHABLE = -109,
// ERR_SSL_CLIENT_AUTH_CERT_NEEDED = -110,
// ERR_TUNNEL_CONNECTION_FAILED = -111,
// ERR_NO_SSL_VERSIONS_ENABLED = -112,
// ERR_SSL_VERSION_OR_CIPHER_MISMATCH = -113,
// ERR_SSL_RENEGOTIATION_REQUESTED = -114,
// ERR_CERT_COMMON_NAME_INVALID = -200,
// ERR_CERT_BEGIN = ERR_CERT_COMMON_NAME_INVALID,
// ERR_CERT_DATE_INVALID = -201,
// ERR_CERT_AUTHORITY_INVALID = -202,
// ERR_CERT_CONTAINS_ERRORS = -203,
// ERR_CERT_NO_REVOCATION_MECHANISM = -204,
// ERR_CERT_UNABLE_TO_CHECK_REVOCATION = -205,
// ERR_CERT_REVOKED = -206,
// ERR_CERT_INVALID = -207,
// ERR_CERT_WEAK_SIGNATURE_ALGORITHM = -208,
// // -209 is available: was ERR_CERT_NOT_IN_DNS.
// ERR_CERT_NON_UNIQUE_NAME = -210,
// ERR_CERT_WEAK_KEY = -211,
// ERR_CERT_NAME_CONSTRAINT_VIOLATION = -212,
// ERR_CERT_VALIDITY_TOO_LONG = -213,
// ERR_CERT_END = ERR_CERT_VALIDITY_TOO_LONG,
// ERR_INVALID_URL = -300,
// ERR_DISALLOWED_URL_SCHEME = -301,
// ERR_UNKNOWN_URL_SCHEME = -302,
// ERR_TOO_MANY_REDIRECTS = -310,
// ERR_UNSAFE_REDIRECT = -311,
// ERR_UNSAFE_PORT = -312,
// ERR_INVALID_RESPONSE = -320,
// ERR_INVALID_CHUNKED_ENCODING = -321,
// ERR_METHOD_NOT_SUPPORTED = -322,
// ERR_UNEXPECTED_PROXY_AUTH = -323,
// ERR_EMPTY_RESPONSE = -324,
// ERR_RESPONSE_HEADERS_TOO_BIG = -325,
// ERR_CACHE_MISS = -400,
// ERR_INSECURE_RESPONSE = -501,
)

const (
	LogSeverityDefault CLogSeverityT = C.LOGSEVERITY_DEFAULT
	LogSeverityVerbose CLogSeverityT = C.LOGSEVERITY_VERBOSE
	LogSeverityDebug   CLogSeverityT = C.LOGSEVERITY_DEBUG
	LogSeverityInfo    CLogSeverityT = C.LOGSEVERITY_INFO
	LogSeverityWarning CLogSeverityT = C.LOGSEVERITY_WARNING
	LogSeverityError   CLogSeverityT = C.LOGSEVERITY_ERROR
	LogSeverityDisable CLogSeverityT = C.LOGSEVERITY_DISABLE
)

const (
	TtLink               CTransitionTypeT = C.TT_LINK
	TtExplicit           CTransitionTypeT = C.TT_EXPLICIT
	TtAutoSubframe       CTransitionTypeT = C.TT_AUTO_SUBFRAME
	TtManualSubframe     CTransitionTypeT = C.TT_MANUAL_SUBFRAME
	TtFormSubmit         CTransitionTypeT = C.TT_FORM_SUBMIT
	TtReload             CTransitionTypeT = C.TT_RELOAD
	TtSourceMask         CTransitionTypeT = C.TT_SOURCE_MASK
	TtBlockedFlag        CTransitionTypeT = C.TT_BLOCKED_FLAG
	TtForwardBackFlag    CTransitionTypeT = C.TT_FORWARD_BACK_FLAG
	TtChainStartFlag     CTransitionTypeT = C.TT_CHAIN_START_FLAG
	TtChainEndFlag       CTransitionTypeT = C.TT_CHAIN_END_FLAG
	TtClientRedirectFlag CTransitionTypeT = C.TT_CLIENT_REDIRECT_FLAG
	TtServerRedirectFlag CTransitionTypeT = C.TT_SERVER_REDIRECT_FLAG
	TtIsRedirectMask     CTransitionTypeT = C.TT_IS_REDIRECT_MASK
	TtQualifierMask      CTransitionTypeT = C.TT_QUALIFIER_MASK
)

const (
	PidBrowser  CProcessIdT = C.PID_BROWSER
	PidRenderer CProcessIdT = C.PID_RENDERER
)

const (
	FileDialogOpen                CFileDialogModeT = C.FILE_DIALOG_OPEN
	FileDialogOpenMultiple        CFileDialogModeT = C.FILE_DIALOG_OPEN_MULTIPLE
	FileDialogOpenFolder          CFileDialogModeT = C.FILE_DIALOG_OPEN_FOLDER
	FileDialogSave                CFileDialogModeT = C.FILE_DIALOG_SAVE
	FileDialogTypeMask            CFileDialogModeT = C.FILE_DIALOG_TYPE_MASK
	FileDialogOverwritepromptFlag CFileDialogModeT = C.FILE_DIALOG_OVERWRITEPROMPT_FLAG
	FileDialogHidereadonlyFlag    CFileDialogModeT = C.FILE_DIALOG_HIDEREADONLY_FLAG
)

const (
	VtypeInvalid    CValueTypeT = C.VTYPE_INVALID
	VtypeNull       CValueTypeT = C.VTYPE_NULL
	VtypeBool       CValueTypeT = C.VTYPE_BOOL
	VtypeInt        CValueTypeT = C.VTYPE_INT
	VtypeDouble     CValueTypeT = C.VTYPE_DOUBLE
	VtypeString     CValueTypeT = C.VTYPE_STRING
	VtypeBinary     CValueTypeT = C.VTYPE_BINARY
	VtypeDictionary CValueTypeT = C.VTYPE_DICTIONARY
	VtypeList       CValueTypeT = C.VTYPE_LIST
)

const (
	CSchemeOptionNone            CSchemeOptionsT = C.CEF_SCHEME_OPTION_NONE
	CSchemeOptionStandard        CSchemeOptionsT = C.CEF_SCHEME_OPTION_STANDARD
	CSchemeOptionLocal           CSchemeOptionsT = C.CEF_SCHEME_OPTION_LOCAL
	CSchemeOptionDisplayIsolated CSchemeOptionsT = C.CEF_SCHEME_OPTION_DISPLAY_ISOLATED
	CSchemeOptionSecure          CSchemeOptionsT = C.CEF_SCHEME_OPTION_SECURE
	CSchemeOptionCorsEnabled     CSchemeOptionsT = C.CEF_SCHEME_OPTION_CORS_ENABLED
	CSchemeOptionCspBypassing    CSchemeOptionsT = C.CEF_SCHEME_OPTION_CSP_BYPASSING
	CSchemeOptionFetchEnabled    CSchemeOptionsT = C.CEF_SCHEME_OPTION_FETCH_ENABLED
)

type Settings struct {
	LogSeverity              CLogSeverityT
	NoSandbox                int
	MultiThreadedMessageLoop int
	RemoteDebuggingPort      int
}

type CAppT struct {
	p_app *C.cef_app_t
}

type CBrowserT struct {
	p_browser *C.cef_browser_t
}
type CBrowserHostT C.cef_browser_host_t
type CBinaryValueT C.cef_binary_value_t
type CDictionaryValueT C.cef_dictionary_value_t
type CCallbackT C.cef_callback_t
type CClientT struct {
	p_client *C.cef_client_t
}
type CFrameT struct {
	p_frame *C.cef_frame_t
}
type CCookieT C.cef_cookie_t
type CCommandLineT C.cef_command_line_t
type CDomnodeT C.cef_domnode_t
type CFileDialogModeT C.cef_file_dialog_mode_t
type CListValueT struct {
	p_list_value *C.cef_list_value_t
}
type CProcessIdT C.cef_process_id_t
type CProcessMessageT struct {
	p_process_message *C.cef_process_message_t
}
type CRequestT C.cef_request_t
type CResourceHandlerT C.cef_resource_handler_t
type CResponseT C.cef_response_t
type CSchemeHandlerFactoryT struct {
	p_scheme_handler_factory *C.cef_scheme_handler_factory_t
}
type CSchemeRegistrarT C.cef_scheme_registrar_t
type CV8accessorT C.cef_v8accessor_t
type CV8arrayBufferReleaseCallbackT struct {
	p_v8array_buffer_release_callback *C.cef_v8array_buffer_release_callback_t
}
type CV8contextT struct {
	p_v8context *C.cef_v8context_t
}
type CV8exceptionT C.cef_v8exception_t
type CV8handlerT struct {
	p_v8handler *C.cef_v8handler_t
}
type CV8interceptorT C.cef_v8interceptor_t
type CV8stackTraceT C.cef_v8stack_trace_t
type CV8valueT struct {
	p_v8value *C.cef_v8value_t
}
type CValueT C.cef_value_t

type CBrowserProcessHandlerT struct {
	p_browser_process_handler *C.cef_browser_process_handler_t
}
type CContextMenuHandlerT C.cef_context_menu_handler_t
type CDialogHandlerT C.cef_dialog_handler_t
type CDisplayHandlerT C.cef_display_handler_t
type CDownloaddHanderT C.cef_download_handler_t
type CDragHandlerT C.cef_drag_handler_t
type CFindHandlerT C.cef_find_handler_t
type CFocusHanderT C.cef_focus_handler_t
type CJsdialogHandlerT C.cef_jsdialog_handler_t
type CKeyboardHandlerT C.cef_keyboard_handler_t
type CLifeSpanHandlerT struct {
	p_life_span_handler *C.cef_life_span_handler_t
}
type CLoadHandlerT struct {
	p_load_handler *C.cef_load_handler_t
}
type CRenderHandlerT C.cef_render_handler_t
type CRequestHandlerT C.cef_request_handler_t
type CResourceBundleHanderT C.cef_resource_bundle_handler_t
type CRenderProcessHandlerT struct {
	p_render_process_handler *C.cef_render_process_handler_t
}

type CRunFileDialogCallbackT struct {
	p_run_file_dialog_callback *C.cef_run_file_dialog_callback_t
}

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
}

var main_args = C.cef_main_args_t{}

func ExecuteProcess(app *CAppT) {

	instance, err := winapi.GetModuleHandle(nil)
	if err != nil {
		Panicf("L205: %v", err)
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
	// BaseAddRef(client.p_client) ?? Need not? clinet.p_client is CefClientCToCpp::Wrap-ed
	C.cef_browser_host_create_browser(
		&window_info,
		client.p_client,
		cef_url,
		&browser_settings, nil,
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
