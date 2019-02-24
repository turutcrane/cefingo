package cefingo

import (
	"runtime"
	"unsafe"
)

// #include "cefingo.h"
import "C"

func newCBrowserT(cef *C.cef_browser_t) *CBrowserT {
	Logf("L42: %p", cef)
	BaseAddRef(cef)
	browser := CBrowserT{cef}
	runtime.SetFinalizer(&browser, func(b *CBrowserT) {
		if ref_count_log.output {
			Logf("L47: %p", b.p_browser)
		}
		BaseRelease(b.p_browser)
	})
	return &browser
}

func (bh *CBrowserHostT) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(bh))
}

func (b *C.cef_browser_t) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(b))
}

func (self *CBrowserT) GetHost() (h *CBrowserHostT) {
	h = (*CBrowserHostT)(C.cefingo_browser_get_host(self.p_browser))
	BaseAddRef(h)
	return h
}

///
// Returns the focused frame for the browser window.
///
func (b *CBrowserT) GetFocusedFrame() (f *CFrameT) {
	f = (*CFrameT)(C.cefingo_browser_get_focused_frame(b.p_browser))
	BaseAddRef(f)
	return f
}

///
// Send a message to the specified |target_process|. Returns true (1) if the
// message was sent successfully.
///
func (self *CBrowserT) SendProcessMessage(
	target_process CProcessIdT,
	message *CProcessMessageT) bool {

	status := C.cefingo_browser_send_process_message(
		self.p_browser,
		C.cef_process_id_t(target_process),
		message.p_process_message,
	)
	return status == 1
}

///
// Called asynchronously after the file dialog is dismissed.
// |selected_accept_filter| is the 0-based index of the value selected from
// the accept filters array passed to cef_browser_host_t::RunFileDialog.
// |file_paths| will be a single value or a list of values depending on the
// dialog mode. If the selection was cancelled |file_paths| will be NULL.
///
type RunFileDialogCallback interface {
	OnFileDialogDismissed(
		self *CRunFileDialogCallbackT,
		selected_accept_filter int,
		file_paths CStringListT,
	)
}

var run_file_dialog_callback = map[*CRunFileDialogCallbackT]RunFileDialogCallback{}

func AllocRunFileDialogCallbackT(callback RunFileDialogCallback) (crun_file_dialog_callback *CRunFileDialogCallbackT) {
	p := C.calloc(1, C.sizeof_cefingo_run_file_dialog_callback_wrapper_t)
	Logf("L43: p: %v", p)

	C.cefingo_construct_run_file_dialog_callback((*C.cefingo_run_file_dialog_callback_wrapper_t)(p))

	crun_file_dialog_callback = (*CRunFileDialogCallbackT)(p)
	BaseAddRef(crun_file_dialog_callback)
	run_file_dialog_callback[crun_file_dialog_callback] = callback

	registerDeassocer(p, DeassocFunc(func() {
		Logf("L56: Deassoc of *CRunFileDialogCallbackT %p", p)
		delete(run_file_dialog_callback, crun_file_dialog_callback)
	}))

	return crun_file_dialog_callback
}

func (c *CRunFileDialogCallbackT) cast_to_p_base_ref_counted_t() *C.cef_base_ref_counted_t {
	return (*C.cef_base_ref_counted_t)(unsafe.Pointer(c))
}

//export cefingo_run_file_dialog_callback_on_file_dialog_dismissed
func cefingo_run_file_dialog_callback_on_file_dialog_dismissed(
	self *CRunFileDialogCallbackT,
	selected_accept_filter C.int,
	file_paths CStringListT) {

	c := run_file_dialog_callback[self]
	if c == nil {
		Panicf("L62: on_file_dialog_dismissed: Noo!")
	}
	c.OnFileDialogDismissed(self, int(selected_accept_filter), file_paths)
}

///
// Call to run a file chooser dialog. Only a single file chooser dialog may be
// pending at any given time. |mode| represents the type of dialog to display.
// |title| to the title to be used for the dialog and may be NULL to show the
// default title ("Open" or "Save" depending on the mode). |default_file_path|
// is the path with optional directory and/or file name component that will be
// initially selected in the dialog. |accept_filters| are used to restrict the
// selectable file types and may any combination of (a) valid lower-cased MIME
// types (e.g. "text/*" or "image/*"), (b) individual file extensions (e.g.
// ".txt" or ".png"), or (c) combined description and file extension delimited
// using "|" and ";" (e.g. "Image Types|.png;.gif;.jpg").
// |selected_accept_filter| is the 0-based index of the filter that will be
// selected by default. |callback| will be executed after the dialog is
// dismissed or immediately if another dialog is already pending. The dialog
// will be initiated asynchronously on the UI thread.
///
func (h *CBrowserHostT) RunFileDialog(
	mode CFileDialogModeT,
	title string,
	default_file_path string,
	accept_filters CStringListT,
	selected_accept_filter int,
	callback *CRunFileDialogCallbackT,
) {
	t := create_cef_string(title)
	defer clear_cef_string(t)

	dfp := create_cef_string(default_file_path)
	defer clear_cef_string(dfp)

	BaseAddRef(callback)
	C.cefingo_browser_host_run_file_dialog(
		(*C.cef_browser_host_t)(h),
		C.cef_file_dialog_mode_t(mode),
		t, dfp,
		C.cef_string_list_t(accept_filters), C.int(selected_accept_filter),
		(*C.cef_run_file_dialog_callback_t)(callback),
	)
}
