package parser

//go:generate go run golang.org/x/tools/cmd/stringer -output parse_string.go -type IdentKind,Ty,StructType,DefKind,TypeQualifier parse.go

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"modernc.org/cc"
)

type void struct{}

var setElement void

var targetFileList = map[string]void{
	"cef_string_list.h":     setElement,
	"cef_string_map.h":      setElement,
	"cef_string_multimap.h": setElement,
	"cef_time.h":            setElement,
	"cef_types.h":           setElement,
	"cef_types_win.h":       setElement,
}

var handlerClasses = map[string]void{
	"cef_accessibility_handler_t":           setElement,
	"cef_app_t":                             setElement,
	"cef_audio_handler_t":                   setElement,
	"cef_browser_process_handler_t":         setElement,
	"cef_client_t":                          setElement,
	"cef_cookie_visitor_t":                  setElement,
	"cef_cookie_access_filter_t":            setElement,
	"cef_dev_tools_message_observer_t":      setElement,
	"cef_delete_cookies_callback_t":         setElement,
	"cef_context_menu_handler_t":            setElement,
	"cef_dialog_handler_t":                  setElement,
	"cef_display_handler_t":                 setElement,
	"cef_domvisitor_t":                      setElement,
	"cef_download_handler_t":                setElement,
	"cef_drag_handler_t":                    setElement,
	"cef_end_tracing_callback_t":            setElement,
	"cef_extension_handler_t":               setElement,
	"cef_find_handler_t":                    setElement,
	"cef_focus_handler_t":                   setElement,
	"cef_jsdialog_handler_t":                setElement,
	"cef_keyboard_handler_t":                setElement,
	"cef_life_span_handler_t":               setElement,
	"cef_load_handler_t":                    setElement,
	"cef_menu_model_delegate_t":             setElement,
	"cef_media_observer_t":                  setElement,
	"cef_print_handler_t":                   setElement,
	"cef_read_handler_t":                    setElement,
	"cef_register_cdm_callback_t":           setElement,
	"cef_render_handler_t":                  setElement,
	"cef_render_process_handler_t":          setElement,
	"cef_request_context_handler_t":         setElement,
	"cef_request_handler_t":                 setElement,
	"cef_resource_bundle_handler_t":         setElement,
	"cef_resource_handler_t":                setElement,
	"cef_resource_request_handler_t":        setElement,
	"cef_response_filter_t":                 setElement,
	"cef_run_file_dialog_callback_t":        setElement,
	"cef_scheme_handler_factory_t":          setElement,
	"cef_server_handler_t":                  setElement,
	"cef_set_cookie_callback_t":             setElement,
	"cef_string_visitor_t":                  setElement,
	"cef_v8accessor_t":                      setElement,
	"cef_v8handler_t":                       setElement,
	"cef_v8array_buffer_release_callback_t": setElement,
	"cef_task_t":                            setElement,
	"cef_urlrequest_client_t":               setElement,
	"cef_web_plugin_info_visitor_t":         setElement,
	"cef_web_plugin_unstable_callback_t":    setElement,
	"cef_write_handler_t":                   setElement,

	"cef_browser_view_delegate_t": setElement,
	"cef_button_delegate_t":       setElement,
	"cef_menu_button_delegate_t":  setElement,
	"cef_panel_delegate_t":        setElement,
	"cef_textfield_delegate_t":    setElement,
	"cef_view_delegate_t":         setElement,
	"cef_window_delegate_t":       setElement,
}

var unGenerateMethod = map[string]void{
	"cef_command_line_t::init_from_argv":          setElement, // use char *
	"cef_browser_t::get_frame_identifiers":        setElement, // parameter identifiers should be **int64 instead of *int64
	"cef_post_data_t::get_elements":               setElement, // pre allocate slice with modifiable count currently can not support
	"cef_print_settings_t::get_page_ranges":       setElement, // parameter ranges should be cef_range_t** instead of cef_range_t* (PageRangeList& ranges)
	"cef_audio_handler_t::on_audio_stream_packet": setElement, // data parameter type **float shoud be *float (float[frame * length of channel_layout])
	"cef_v8value_t::get_user_data":                setElement, // use struct _cef_base_ref_counted_t*
	"cef_v8value_t::set_user_data":                setElement, // use struct _cef_base_ref_counted_t*
}

var notBoolValueMethod = map[string]void{
	"::cef_execute_process":         setElement,
	"::cef_string_list_size":        setElement,
	"cef_list_value_t::get_int":     setElement,
	"cef_browser_t::get_identifier": setElement,
}

var duplicatedHandler = map[string]void{
	"can_set_cookie":               setElement,
	"execute":                      setElement,
	"get_auth_credentials":         setElement,
	"get_load_handler":             setElement,
	"get_resource_request_handler": setElement,
	"may_block":                    setElement,
	"on_browser_created":           setElement,
	"on_browser_destroyed":         setElement,
	"on_complete":                  setElement,
	"on_key_event":                 setElement,
	"on_process_message_received":  setElement,
	"read":                         setElement,
	"seek":                         setElement,
	"tell":                         setElement,
	"visit":                        setElement,
}

var outParameter = map[string]void{
	"cef_cookie_visitor_t::visit::deleteCookie":                                             setElement,
	"cef_image_t::get_representation_info::actual_scale_factor":                             setElement,
	"cef_image_t::get_representation_info::pixel_height":                                    setElement,
	"cef_image_t::get_representation_info::pixel_width":                                     setElement,
	"cef_image_t::get_as_bitmap::pixel_height":                                              setElement,
	"cef_image_t::get_as_bitmap::pixel_width":                                               setElement,
	"cef_image_t::get_as_png::pixel_height":                                                 setElement,
	"cef_image_t::get_as_png::pixel_width":                                                  setElement,
	"cef_image_t::get_as_jpeg::pixel_height":                                                setElement,
	"cef_image_t::get_as_jpeg::pixel_width":                                                 setElement,
	"cef_jsdialog_handler_t::on_jsdialog::suppress_message":                                 setElement,
	"cef_keyboard_handler_t::on_pre_key_event::is_keyboard_shortcut":                        setElement,
	"cef_life_span_handler_t::on_before_popup::extra_info":                                  setElement,
	"cef_menu_model_t::get_accelerator::key_code":                                           setElement,
	"cef_menu_model_t::get_accelerator::shift_pressed":                                      setElement,
	"cef_menu_model_t::get_accelerator::ctrl_pressed":                                       setElement,
	"cef_menu_model_t::get_accelerator::alt_pressed":                                        setElement,
	"cef_menu_model_t::get_accelerator_at::key_code":                                        setElement,
	"cef_menu_model_t::get_accelerator_at::shift_pressed":                                   setElement,
	"cef_menu_model_t::get_accelerator_at::ctrl_pressed":                                    setElement,
	"cef_menu_model_t::get_accelerator_at::alt_pressed":                                     setElement,
	"cef_render_handler_t::get_view_rect::rect":                                             setElement,
	"cef_render_handler_t::get_root_screen_rect::rect":                                      setElement,
	"cef_render_handler_t::get_screen_point::screenX":                                       setElement,
	"cef_render_handler_t::get_screen_point::screenY":                                       setElement,
	"cef_request_context_t::set_preference::error":                                          setElement,
	"cef_request_context_handler_t::get_resource_request_handler::disable_default_handling": setElement,
	"cef_request_context_handler_t::on_before_plugin_load::plugin_policy":                   setElement,
	"cef_request_handler_t::get_resource_request_handler::disable_default_handling":         setElement,
	"cef_request_handler_t::on_protocol_execution::allow_os_execution":                      setElement,
	"cef_request_handler_t::on_resource_redirect::new_url":                                  setElement,
	"cef_response_filter_t::filter::data_in_read":                                           setElement,
	"cef_response_filter_t::filter::data_out_written":                                       setElement,
	"cef_resource_bundle_handler_t::get_data_resource::data":                                setElement,
	"cef_resource_bundle_handler_t::get_data_resource_for_scale::data":                      setElement,
	"cef_resource_handler_t::get_response_headers::response_length":                         setElement,
	"cef_resource_handler_t::get_response_headers::redirectUrl":                             setElement,
	"cef_resource_handler_t::open::handle_request":                                          setElement,
	"cef_resource_handler_t::read::bytes_read":                                              setElement,
	"cef_resource_handler_t::read_response::bytes_read":                                     setElement,
	"cef_resource_handler_t::skip::bytes_skipped":                                           setElement,
	"cef_resource_bundle_handler_t::get_localized_string::string":                           setElement,
	"cef_resource_request_handler_t::on_resource_redirect::new_url":                         setElement,
	"cef_resource_request_handler_t::on_protocol_execution::allow_os_execution":             setElement,
	"cef_v8accessor_t::get::exception":                                                      setElement,
	"cef_v8accessor_t::get::retval":                                                         setElement,
	"cef_v8accessor_t::set::exception":                                                      setElement,
	"cef_v8context_t::eval::exception":                                                      setElement,
	"cef_v8context_t::eval::retval":                                                         setElement,
	"cef_v8handler_t::execute::exception":                                                   setElement,
	"cef_v8handler_t::execute::retval":                                                      setElement,
	"cef_v8interceptor_t::get_byname::exception":                                            setElement,
	"cef_v8interceptor_t::get_byname::retval":                                               setElement,
	"cef_v8interceptor_t::set_byname::exception":                                            setElement,
	"cef_v8interceptor_t::get_byindex::exception":                                           setElement,
	"cef_v8interceptor_t::get_byindex::retval":                                              setElement,
	"cef_v8interceptor_t::set_byindex::exception":                                           setElement,
	"cef_x509certificate_t::get_derencoded_issuer_chain::chain":                             setElement,
	"cef_x509certificate_t::get_pemencoded_issuer_chain::chain":                             setElement,

	"::cef_time_to_timet::time": setElement,

	"cef_window_delegate_t::get_parent_window::is_menu":           setElement,
	"cef_window_delegate_t::get_parent_window::can_activate_menu": setElement,
	"::cef_display_get_alls::displays":                            setElement,
	"::cef_string_list_value::value":                              setElement,
	"::cef_string_map_find::value":                                setElement,
	"::cef_string_map_key::key":                                   setElement,
	"::cef_string_map_value::value":                               setElement,
	"::cef_string_multimap_enumerate::value":                      setElement,
	"::cef_string_multimap_key::key":                              setElement,
	"::cef_string_multimap_value::value":                          setElement,
}

var inOutParameter = map[string]void{
	"cef_display_handler_t::on_tooltip::text":                         setElement,
	"cef_extension_handler_t::on_before_background_browser::client":   setElement,
	"cef_extension_handler_t::on_before_background_browser::settings": setElement,
	"cef_extension_handler_t::on_before_browser::windowInfo":          setElement,
	"cef_extension_handler_t::on_before_browser::client":              setElement,
	"cef_extension_handler_t::on_before_browser::settings":            setElement,
	"cef_life_span_handler_t::on_before_popup::windowInfo":            setElement,
	"cef_life_span_handler_t::on_before_popup::client":                setElement,
	"cef_life_span_handler_t::on_before_popup::settings":              setElement,
	"cef_life_span_handler_t::on_before_popup::no_javascript_access":  setElement,
	"cef_menu_model_delegate_t::format_label::label":                  setElement,
}

var byteSliceParameter = map[string]string{
	"cef_browser_host_t::send_dev_tools_message::message":                  "message_size",
	"cef_image_t::add_bitmap::pixel_data":                                  "pixel_data_size",
	"cef_image_t::add_png::png_data":                                       "png_data_size",
	"cef_image_t::add_jpeg::jpeg_data":                                     "jpeg_data_size",
	"cef_dev_tools_message_observer_t::on_dev_tools_message::message":      "message_size",
	"cef_dev_tools_message_observer_t::on_dev_tools_method_result::result": "result_size",
	"cef_dev_tools_message_observer_t::on_dev_tools_event::params":         "params_size",
	"cef_media_observer_t::on_route_message_received::message":             "message_size",
	"cef_media_route_t::send_route_message::message":                       "message_size",
	"cef_post_data_element_t::set_to_bytes::bytes":                         "size",
	"cef_read_handler_t::read::ptr":                                        "size*n",
	"cef_render_handler_t::on_paint::buffer":                               "width*height*4",
	"cef_response_filter_t::filter::data_in":                               "data_in_size",
	"cef_response_filter_t::filter::data_out":                              "data_out_size",
	"cef_resource_bundle_handler_t::get_data_resource::data":               "data_size",
	"cef_resource_bundle_handler_t::get_data_resource_for_scale::data":     "data_size",
	"cef_resource_handler_t::read::data_out":                               "bytes_to_read",
	"cef_resource_handler_t::read_response::data_out":                      "bytes_to_read",
	"cef_urlrequest_client_t::on_download_data::data":                      "data_length",
	"cef_write_handler_t::write::ptr":                                      "size*n",
	"::cef_binary_value_create::data":                                      "data_size",
	"::cef_stream_reader_create_for_data::data":                            "size",
	// "::cef_v8value_create_array_buffer::buffer":       "length", buffer parameter should be shared with cef_v8array_buffer_release_callback_t
}
var byteSliceLengthParameter = map[string]string{}

var sliceParameter = map[string]string{
	"cef_media_observer_t::on_sinks::sinks":                             "sinksCount",
	"cef_media_observer_t::on_routes::routes":                           "routesCount",
	"cef_render_handler_t::on_paint::dirtyRects":                        "dirtyRectsCount",
	"cef_request_handler_t::on_select_client_certificate::certificates": "certificatesCount",
	"cef_v8handler_t::execute::arguments":                               "argumentsCount",
	"cef_v8value_t::execute_function::arguments":                        "argumentsCount",
	"cef_v8value_t::execute_function_with_context::arguments":           "argumentsCount",
	"cef_x509certificate_t::get_derencoded_issuer_chain::chain":         "chainCount",
	"cef_x509certificate_t::get_pemencoded_issuer_chain::chain":         "chainCount",
	"::cef_display_get_alls::displays":                                  "displaysCount",
}

var sliceLengthParameter = map[string]string{}

func init() {
	for bs, length := range byteSliceParameter {
		c := strings.Split(bs, "::")
		byteSliceLengthParameter[strings.Join([]string{c[0], c[1], length}, "::")] = c[2]
	}

	for s, length := range sliceParameter {
		c := strings.Split(s, "::")
		sliceLengthParameter[strings.Join([]string{c[0], c[1], length}, "::")] = c[2]
	}
}

var boolParameter = map[string]void{
	"::cef_set_osmodal_loop::osModalLoop":                                                   setElement,
	"cef_browser_host_t::close_browser::force_close":                                        setElement,
	"cef_browser_host_t::find::forward":                                                     setElement,
	"cef_browser_host_t::find::matchCase":                                                   setElement,
	"cef_browser_host_t::find::findNext":                                                    setElement,
	"cef_browser_host_t::set_focus::focus":                                                  setElement,
	"cef_browser_host_t::stop_finding::clearSelection":                                      setElement,
	"cef_browser_host_t::download_image::is_favicon":                                        setElement,
	"cef_browser_host_t::download_image::bypass_cache":                                      setElement,
	"cef_browser_host_t::get_navigation_entries::current_only":                              setElement,
	"cef_browser_view_t::set_prefer_accelerators::prefer_accelerators":                      setElement,
	"cef_browser_view_delegate_t::get_delegate_for_popup_browser_view::is_devtools":         setElement,
	"cef_browser_view_delegate_t::on_popup_browser_view_created::is_devtools":               setElement,
	"cef_dictionary_value_t::set_bool::value":                                               setElement,
	"cef_life_span_handler_t::on_before_popup::no_javascript_access":                        setElement,
	"cef_load_handler_t::on_loading_state_change::isLoading":                                setElement,
	"cef_load_handler_t::on_loading_state_change::canGoBack":                                setElement,
	"cef_load_handler_t::on_loading_state_change::canGoForward":                             setElement,
	"cef_list_value_t::set_bool::value":                                                     setElement,
	"cef_navigation_entry_visitor_t::visit::current":                                        setElement,
	"cef_pdf_print_callback_t::on_pdf_print_finished::ok":                                   setElement,
	"cef_request_handler_t::get_resource_request_handler::disable_default_handling":         setElement,
	"cef_request_handler_t::on_open_urlfrom_tab::user_gesture":                              setElement,
	"cef_request_context_handler_t::get_resource_request_handler::disable_default_handling": setElement,
	"cef_resource_handler_t::open::handle_request":                                          setElement,
	"cef_v8value_t::create_bool::value":                                                     setElement,
	"cef_value_t::set_bool::value":                                                          setElement,
	"cef_view_t::set_focusable::focusable":                                                  setElement,
	"cef_window_delegate_t::get_parent_window::is_menu":                                     setElement,
	"cef_window_delegate_t::get_parent_window::can_activate_menu":                           setElement,
}

var boolMember = map[string]void{
	"cef_box_layout_settings_t::horizontal":                     setElement,
	"cef_composition_underline_t::thick":                        setElement,
	"cef_settings_t::no_sandbox":                                setElement,
	"cef_settings_t::multi_threaded_message_loop":               setElement,
	"cef_settings_t::external_message_pump":                     setElement,
	"cef_settings_t::windowless_rendering_enabled":              setElement,
	"cef_settings_t::command_line_args_disabled":                setElement,
	"cef_settings_t::persist_session_cookies":                   setElement,
	"cef_settings_t::persist_user_preferences":                  setElement,
	"cef_settings_t::pack_loading_disabled":                     setElement,
	"cef_settings_t::ignore_certificate_errors":                 setElement,
	"cef_draggable_region_t::draggable":                         setElement,
	"cef_pdf_print_settings_t::header_footer_enabled":           setElement,
	"cef_pdf_print_settings_t::selection_only":                  setElement,
	"cef_pdf_print_settings_t::landscape":                       setElement,
	"cef_pdf_print_settings_t::backgrounds_enabled":             setElement,
	"cef_request_context_settings_t::persist_session_cookies":   setElement,
	"cef_request_context_settings_t::persist_user_preferences":  setElement,
	"cef_request_context_settings_t::ignore_certificate_errors": setElement,
	"cef_popup_features_t::xSet":                                setElement,
	"cef_popup_features_t::ySet":                                setElement,
	"cef_popup_features_t::widthSet":                            setElement,
	"cef_popup_features_t::heightSet":                           setElement,
	"cef_window_info_t::windowless_rendering_enabled":           setElement,
	"cef_window_info_t::shared_texture_enabled":                 setElement,
	"cef_window_info_t::external_begin_frame_enabled":           setElement,
}

var cefdir string

type IdentKind int

const (
	IkNone IdentKind = iota
	IkName
	IkFunc
	IkArray
	IkStructTag
)

type Ty int

const (
	TyUnknown Ty = iota
	TyVoid
	TyChar
	TyUnsigned
	TyInt
	TyLong
	TyLongLong
	TyULong
	TyULongLong
	TyFloat
	TyDouble
	TySizeT
	TyHWND
	TyStructUnhandled
	TyStructSimple
	TyStructRefCounted
	TyStructScoped
	TyStructNotDefined
	TyInt16
	TyInt32
	TyInt64
	TyUint16
	TyUint32
	TyUint64
	TyTimeT
	TyStringT
	TyStringUserfreeT
	TyOther
	TyEnum
	TySimple
	TyMSG
	TyHCURSOR
	TyHINSTANCE
	TyDWORD
	TyHMENU
)

var primitiveTypeDef = map[string]Ty{
	"size_t":    TySizeT,
	"HWND":      TyHWND,
	"int32":     TyInt32,
	"int64":     TyInt64,
	"uint32":    TyUint32,
	"uint64":    TyUint64,
	"time_t":    TyTimeT,
	"int16":     TyInt16,
	"uint16":    TyUint16,
	"char16":    TyUint16,
	"HINSTANCE": TyHINSTANCE,
	"DWORD":     TyDWORD,
	"HMENU":     TyHMENU,
}

type TypeQualifier int

const (
	TqUnknown TypeQualifier = iota
	TqConst
	TqPointer
)

type Type struct {
	Ty
	Pointer int
	Const   bool
	Token   Token
	Typedef bool
	Tq      []TypeQualifier
}

func (t Type) String() string {
	var ptr string
	for i := 0; i < t.Pointer; i++ {
		ptr += "*"
	}
	return ptr + t.Ty.String() + " : " + t.Token.FilePos()
}

type StructType int

const (
	StUnknown StructType = iota
	StRefCounted
	StScoped
	StYetNotDefined
)

var structDefNames = map[string]void{
	"cef_settings_t":                 setElement,
	"cef_main_args_t":                setElement,
	"cef_window_info_t":              setElement,
	"cef_request_context_settings_t": setElement,
	"cef_browser_settings_t":         setElement,
	"cef_urlparts_t":                 setElement,
	"cef_cookie_t":                   setElement,
	"cef_point_t":                    setElement,
	"cef_rect_t":                     setElement,
	"cef_size_t":                     setElement,
	"cef_range_t":                    setElement,
	"cef_insets_t":                   setElement,
	"cef_draggable_region_t":         setElement,
	"cef_screen_info_t":              setElement,
	"cef_mouse_event_t":              setElement,
	"cef_touch_event_t":              setElement,
	"cef_key_event_t":                setElement,
	"cef_popup_features_t":           setElement,
	"cef_cursor_info_t":              setElement,
	"cef_pdf_print_settings_t":       setElement,
	"cef_box_layout_settings_t":      setElement,
	"cef_composition_underline_t":    setElement,
	"cef_audio_parameters_t":         setElement,
	"cef_media_sink_device_info_t":   setElement,
}

var simpleDefNames = map[string]void{
	"cef_color_t":           setElement,
	"cef_time_t":            setElement,
	"cef_string_list_t":     setElement,
	"cef_string_map_t":      setElement,
	"cef_string_multimap_t": setElement,
}

func isSimpleDefName(s string) (b bool) {
	_, b = simpleDefNames[s]
	return b
}

func isStructDefName(s string) (b bool) {
	_, b = structDefNames[s]
	return b
}

func IsHandlerClass(t Token) (c bool) {
	return isHandlerClass(t.Name())
}

type DefKind int

const (
	DkUnknown DefKind = iota
	DkUnhandled
	DkSimple
	DkEnum
	DkCefClass
	DkFunc
	DkStruct
)

type Decl interface {
	SetComment(comments map[int][]string)
	Common() *DeclCommon
}

type DeclCommon struct {
	Dk      DefKind
	d       *cc.Declaration
	Comment []string
}

func (d *DeclCommon) Common() *DeclCommon {
	return d
}

func (decl DeclCommon) Line() int {
	return decl.Token().Line()
}

func (decl DeclCommon) FilePos() string {
	return decl.Token().FilePos()
}

func (s *DeclCommon) SetComment(comments map[int][]string) {
	if c, ok := comments[s.LineOfTypedef()]; ok {
		s.Comment = c
	}
}

func (s *CefClassDecl) SetComment(comments map[int][]string) {
	switch s.DeclCommon.Dk {
	case DkCefClass:
		s.DeclCommon.SetComment(comments)
		for i, _ := range s.Methods {
			m := s.Methods[i]
			if c, ok := comments[(*m).FirstLine()]; ok {
				m.Comment = c
			}
		}
	default:
		log.Panicf("T475: %v\n", s.d)
	}
}

func (s *CefClassDecl) GetBase() (base *CefClassDecl) {
	if c, ok := Defs[s.BaseType]; ok {
		if base, cefClass := c.(*CefClassDecl); cefClass {
			return base
		}
	}
	return nil
}

type UnhandledDecl struct {
	DeclCommon
}

type SimpleDecl struct {
	DeclCommon
}

type Member struct {
	token     Token
	typ       Type
	decl      *cc.StructDeclaration
	container *StructDecl
}

func (m Member) Name() string {
	return m.token.Name()
}

func (m Member) GoName() string {
	return m.token.TitleCase()
}

func (m Member) GoType() string {
	if m.typ.Ty == TyStringT && m.typ.Pointer == 0 {
		return "string"
	}
	if m.IsBoolMember() {
		return "bool"
	}
	return m.typ.GoType()
}

func (m Member) IsBoolMember() bool {
	_, isBool := boolMember[m.container.CefName()+"::"+m.Name()]
	return isBool
}

func (m Member) Type() Type {
	return m.typ
}

type StructDecl struct {
	DeclCommon
	Members []Member
}

type CefClassDecl struct {
	DeclCommon
	St       StructType
	BaseType string
	Methods  []*MethodDecl
}

type EnumDecl struct {
	DeclCommon
	Enums []Token
}

type Callee interface {
	CalleeName() string
}

type FuncDecl struct {
	DeclCommon
	Funcname Token
	params   []Param
}

type MethodDecl struct {
	Funcname Token
	params   []Param
	sd       cc.StructDeclaration // Struct Member Declaration
	Comment  []string
	sdecl    *CefClassDecl // Struct Declaration
}

type Callable interface {
	Params() []Param
	HasReturnValue() bool
	HasOutParam() bool
	IsBoolValueMethod() bool
	ReturnGoType() string
	ReturnType() Type
}

func (m MethodDecl) Params() []Param {
	return m.params
}

func (m MethodDecl) FirstLine() (line int) {
	ts := getTypeSpecifier(m.sd)
	switch ts.Case {
	case 0, 3, 5, 6, 13: // void, int, float, double, TYPEDEFNAME
		line = ts.Token.Position().Line
	case 11: // StructOrUnionSpecifier
		sous := ts.StructOrUnionSpecifier
		if sous.Case == 1 {
			t := sous.StructOrUnion.Token
			line = t.Position().Line
		} else {
			log.Panicf("T150: %v\n", m.sd)
		}
	default:
		log.Panicf("T153: %v\n", m.sd)
	}
	return line
}

func (m MethodDecl) IfName() (ifname string) {
	if _, dupn := duplicatedHandler[m.Funcname.Name()]; dupn {
		return m.sdecl.GoName() + m.Funcname.TitleCase() + "Handler"
	}
	return m.Funcname.TitleCase() + "Handler"
}

func (m MethodDecl) HasConstParams() (has bool) {
	for _, p := range m.Params() {
		if p.Type().Const {
			return true
		}
		for _, tq := range p.Type().Tq {
			if tq == TqConst {
				return true
			}
		}
	}
	return false
}

func (m MethodDecl) HasReturnValue() (has bool) {
	return m.ReturnGoType() != ""
}

func (m MethodDecl) HasOutParam() (has bool) {
	for _, p := range m.Params() {
		if p.IsOutParam() || p.IsInOutParam() {
			has = true
		}
	}
	return has
}

func InTargetFile(t Token) (f bool, fname string) {
	fn := filepath.Base(t.Filename())
	_, f = targetFileList[fn]
	if !f && strings.HasPrefix(fn, "cef_") && strings.HasSuffix(fn, "_capi.h") {
		f = true
	}
	return f, fn
}

func InCefHeader(t Token) bool {
	base := filepath.Base(t.Filename())
	if strings.HasPrefix(base, "cef_") && strings.HasSuffix(base, ".h") &&
		base != "cef_string_types.h" && base != "cef_string.h" &&
		base != "cef_base_capi.h" {
		return true
	}
	return false
}

var Defs = map[string]Decl{}
var FileDefs = map[string][]Decl{}

var hasHandlerClass = map[string]bool{}

func HasHandlerClass(fname string) bool {
	return hasHandlerClass[fname]
}

func addDeclToFile(fname string, decl Decl) {
	if a, ok := FileDefs[fname]; ok {
		FileDefs[fname] = append(a, decl)
	} else {
		FileDefs[fname] = []Decl{decl}
	}
}

// __DIR__
func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func getPredefinedMacros() string {
	nullPath := "/dev/null"
	if runtime.GOOS == "windows" {
		nullPath = "nul:"
	}
	pre, err := exec.Command("cpp", "-dM", nullPath).CombinedOutput()
	if err != nil {
		log.Panicf("E572: %v, %s", string(pre), err)
	}
	return string(pre)
}

func Parse() []*cc.TranslationUnit {
	stripNl := strings.NewReplacer("\n", "", "\r", "")

	out, err := exec.Command("pkg-config", "cefingo", "--variable=includedir").Output()
	if err != nil {
		log.Panicf("T268: %v", err)
	}
	cefdir = filepath.ToSlash(filepath.Clean(stripNl.Replace(string(out))))

	out, err = exec.Command("cygpath", "-w", os.Getenv("MINGW_PREFIX")).Output()
	if err != nil {
		log.Panicf("T684: %v", err)
	}
	mingwPrefix := filepath.ToSlash(filepath.Clean(stripNl.Replace(string(out))))

	// this list build from following command output and cef include path root
	// $ echo | cpp -v
	includePaths := []string{
		filepath.Clean(mingwPrefix + "/lib/gcc/x86_64-w64-mingw32/9.1.0/include"),
		filepath.Clean(mingwPrefix + "/include"),
		filepath.Clean(mingwPrefix + "/lib/gcc/x86_64-w64-mingw32/9.1.0/include-fixed"),
		filepath.Clean(mingwPrefix + "/x86_64-w64-mingw32/include"),
		filepath.Clean(cefdir),
	}
	// log.Printf("T259: CurrentDir: %s\n", GetCurrentDir())
	// log.Printf("T260: Cef dir: %v\n", includePaths)

	sourcePaths := []string{GetCurrentDir() + "/../parser_c/parser_root.c"}
	predefined := strings.Join([]string{
		builtinBase,
		getPredefinedMacros(),
		"#define CEFINGO_BUILD 1",
	}, "\n")
	u0, err := cc.Parse(predefined, sourcePaths, model64,
		cc.SysIncludePaths(includePaths),
		cc.EnableAnonymousStructFields(),
		cc.EnableAsm(),
		cc.EnableAlternateKeywords(),
		cc.EnableIncludeNext(),
		cc.EnableNoreturn(),
		cc.EnableEmptyDeclarations(),
		cc.EnableWideEnumValues(),
		cc.EnableWideBitFieldTypes(),
		cc.KeepComments(),
		cc.AllowCompatibleTypedefRedefinitions(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	tus := []*cc.TranslationUnit{}
	for next := u0.TranslationUnit; next != nil; {
		tu := next
		next = tu.TranslationUnit
		tu.TranslationUnit = nil
		tus = append(tus, tu)
	}
	return tus
}

//	ExternalDeclaration:
//	        FunctionDefinition
//	|       Declaration                  // Case 1
//	|       BasicAssemblerStatement ';'  // Case 2
//	|       ';'                          // Case 3
func ExternalDeclaration(i int, ed *cc.ExternalDeclaration) {
	switch ed.Case {
	case 0: // FunctionDefinition
		funcname := getFuncname(ed)
		if in, _ := InTargetFile(funcname); in {
			log.Printf("T73:Func: i:%d, Case:%d, %s, %s", i, ed.Case, funcname.Name(), funcname.Filename())
		}
	case 1: // Declaration
		d := ed.Declaration
		// log.Printf("T76:Dcl: i:%d, Case:%d", i, ed.Case)
		processDeclaration(d)
	case 2: // BasicAssemblerStatement
		log.Printf("T78:BAS: i:%d, Case:%d, %v", i, ed.Case, ed.BasicAssemblerStatement)
	default:
		log.Printf("T80:OTH: i:%d, Case:%d, %v", i, ed.Case, ed)
	}
}

//	Declaration:
//	        DeclarationSpecifiers InitDeclaratorListOpt ';'
//	|       StaticAssertDeclaration                          // Case 1
func processDeclaration(d *cc.Declaration) {
	switch d.Case {
	case 0:
		if d.InitDeclaratorListOpt == nil {
			tag := getTag(d.DeclarationSpecifiers)
			typedefName := tagToTypdefName(tag.Name())
			if _, ok := Defs[typedefName]; !ok {
				Defs[typedefName] = &CefClassDecl{
					DeclCommon{DkUnknown, d, nil},
					StYetNotDefined,
					"",
					nil,
				}
			}
			log.Printf("T337: No InitDecl %s, %s;\n", tag.Name(), tag.FilePos())
		} else {
			name := getDeclaratorIdent(getFirstDeclarator(d))
			var decl Decl
			if InCefHeader(name) {
				decl = analyzeDecl(d)
				if in, fname := InTargetFile(name); in && decl != nil {
					Defs[name.Name()] = decl
					addDeclToFile(fname, decl)
				}
			}
			if isHandlerClass(name.Name()) {
				hasHandlerClass[filepath.Base(name.Filename())] = true
			}
		}
	default:
		log.Panicf("T94: Declaraion Case:%d\n", d.Case)
	}
}

func analyzeDecl(d *cc.Declaration) (decl Decl) {
	base := DeclCommon{}
	base.d = d
	name := getDeclaratorIdent(getFirstDeclarator(d))
	kind := getKind(d)
	ds := base.d.DeclarationSpecifiers
	switch ds.Case {
	case 0:
		switch ds.StorageClassSpecifier.Case {
		case 0: // typedef
			log.Printf("T280: %s: %s, %s\n", kind.String(), name.Name(), base.FilePos())
			switch kind {
			case IkName:
				decl = handleTypedef(base)
			default:
				log.Panicf("T166: Not Name: %v\n", base.d)
			}
		default:
			log.Panicf("T169: ds.Case:%d\n", ds.Case)
		}
	case 1:
		switch kind {
		case IkFunc:
			decl = handleFunc(base)
		default:
			log.Panicf("T176: Not Func: ds.Case:%d, %s\n", ds.Case, name.FilePos())
		}
	default:
		log.Panicf("T179: ds.Case:%d, %s\n", ds.Case, name.FilePos())
	}
	return decl
}

func (decl DeclCommon) LineOfTypedef() (lineno int) {
	ds := decl.d.DeclarationSpecifiers
	switch ds.Case {
	case 0: // StorageClassSpecifier DeclarationSpecifiersOpt
		if ds.StorageClassSpecifier.Case == 0 { // typedef
			lineno = Token(ds.StorageClassSpecifier.Token).Line()
		} else {
			log.Panicf("T307: Somthing Wrong: %v\n", decl.d)
		}
	case 1: // TypeSpecifier DeclarationSpecifiersOpt
		lineno = Token(ds.TypeSpecifier.Token).Line()
	default:
		log.Panicf("T307: Somthing Wrong: %v\n", decl.d)
	}
	return lineno
}

func getFirstDeclarator(d *cc.Declaration) *cc.Declarator {
	listOpt := d.InitDeclaratorListOpt
	if listOpt == nil {
		log.Panicf("T161: %v\n", d)
	}
	declList := listOpt.InitDeclaratorList
	switch declList.Case {
	case 0:
		return declList.InitDeclarator.Declarator
	default:
		log.Panicf("T104: InitDeclaratorList: mutliple")
	}

	return nil
}

func (decl DeclCommon) Token() Token {
	return getDeclaratorIdent(getFirstDeclarator(decl.d))
}

func getKind(d *cc.Declaration) (ik IdentKind) {
	dd := getFirstDeclarator(d).DirectDeclarator
	switch dd.Case {
	case 0: // IDENTIFIR
		ik = IkName
	case 1: //
		log.Panicf("T185: %v", dd)
	case 2: // Array
		ik = IkArray
	case 6, 7: // Function
		ik = IkFunc
	default:
		log.Panicf("T115: %v", dd)
	}
	return ik
}

//	StructOrUnionSpecifier:
//	        StructOrUnion IdentifierOpt '{' StructDeclarationList '}'
//	|       StructOrUnion IDENTIFIER                                   // Case 1
//	|       StructOrUnion IdentifierOpt '{' '}'                        // Case 2
func getStructTag(su *cc.StructOrUnionSpecifier) Token {
	switch su.Case {
	case 1:
		return Token(su.Token)
	default:
		return Token(su.IdentifierOpt.Token)
		// log.Panicf("T133: Not IDENTIFIER %v", ds)
	}
}

func getTag(ds *cc.DeclarationSpecifiers) Token {
	switch ds.Case {
	case 1:
		ts := ds.TypeSpecifier
		switch ts.Case {
		case 11: // StructOrUnionSpecifier
			return getStructTag(ts.StructOrUnionSpecifier)
		}
	default:
		log.Panicf("T127: Not TypeSpecifier %v\n", ds)
	}
	return noToken
}

func getFuncname(f *cc.ExternalDeclaration) Token {
	if f.Case != 0 {
		log.Panicf("T151: Not Function: %v\n", f)
	}
	d := f.FunctionDefinition.Declarator
	return getDeclaratorIdent(d)

}

// Declarator represents data reduced by production:
//
//	Declarator:
//	        PointerOpt DirectDeclarator
//	DirectDeclarator:
//	        IDENTIFIER
//	|       '(' Declarator ')'                                                 // Case 1
//	|       DirectDeclarator '[' TypeQualifierListOpt ExpressionOpt ']'        // Case 2
//	|       DirectDeclarator '[' "static" TypeQualifierListOpt Expression ']'  // Case 3
//	|       DirectDeclarator '[' TypeQualifierList "static" Expression ']'     // Case 4
//	|       DirectDeclarator '[' TypeQualifierListOpt '*' ']'                  // Case 5
//	|       DirectDeclarator '(' ParameterTypeList ')'                         // Case 6
//	|       DirectDeclarator '(' IdentifierListOpt ')'                         // Case 7
func getDeclaratorIdent(d *cc.Declarator) (identToken Token) {
	return getDirectDeclaratorToken(d.DirectDeclarator)
}

func getDirectDeclaratorToken(dd *cc.DirectDeclarator) (identToken Token) {
	switch dd.Case {
	case 0: // IDENTIFIR
		identToken = Token(dd.Token)
	case 1: // '(' Declarator ')'
		identToken = getDirectDeclaratorToken(dd.Declarator.DirectDeclarator)
	case 2: // Array
		identToken = getDirectDeclaratorToken(dd.DirectDeclarator)
	case 6, 7: // Function
		identToken = getDirectDeclaratorToken(dd.DirectDeclarator)
	default:
		log.Panicf("T115: %v", dd)
	}
	return Token(identToken)
}

func handleFunc(base DeclCommon) (decl Decl) {
	base.Dk = DkFunc
	fd := getFirstDeclarator(base.d)
	dd := fd.DirectDeclarator
	fname := getDirectDeclaratorToken(dd)
	// log.Printf("T600: %v\n", base.d)
	log.Printf("T327: func %v\n", fname.Name())
	log.Printf("T235: Ret: %s\n", fd.Type.Result())

	f := &FuncDecl{base, fname, nil}
	switch dd.Case {
	case 6: // DirectDeclarator '(' ParameterTypeList ')'
		for p := dd.ParameterTypeList.ParameterList; p != nil; p = p.ParameterList {
			f.params = append(f.params, getParam(p.ParameterDeclaration, f))
		}
		if dd.ParameterTypeList.Case == 1 { //ParameterList ',' "..."  // Case 1
			variadic := Param{nil, true, noToken, nil}
			f.params = append(f.params, variadic)
		}
	case 7: // DirectDeclarator '(' IdentifierListOpt ')' No Arguments
	default:
		log.Panicf("T335: %v\n", decl)
	}
	for i, p := range f.params {
		log.Printf("T342:   p%d, %s", i, p)
	}

	if _, nf := unGenerateMethod[f.CalleeName()]; nf {
		log.Printf("T887: Skip: %s\n", fname.Name())
		return nil
	}
	return f
}

func (f *FuncDecl) CalleeName() string {
	return "::" + f.Funcname.Name()
}

func (f *FuncDecl) ReturnType() (retType Type) {
	if f.DeclCommon.d.Case != 0 {
		log.Panicf("T632: %v\n", f.DeclCommon.d)
	}
	ds := f.DeclCommon.d.DeclarationSpecifiers
	idlo := f.DeclCommon.d.InitDeclaratorListOpt
	if ds.Case != 1 {
		log.Panicf("T636: %v\n", f.DeclCommon.d)
	}
	ts := ds.TypeSpecifier
	retType = getTsType(ts)

	if idlo.InitDeclaratorList != nil && idlo.InitDeclaratorList.Case == 0 {
		pointer := idlo.InitDeclaratorList.InitDeclarator.Declarator.PointerOpt
		if pointer != nil {
			retType.Pointer = 1
			retType.Tq = append(retType.Tq, TqPointer)
			if pointer.Pointer.Case != 0 {
				log.Panicf("T991: %v\n", pointer)
			}
		}
	} else {
		log.Panicf("T652: %v\n", f.DeclCommon.d)
	}
	return retType
}

func (f FuncDecl) Params() []Param {
	return f.params
}

func (f FuncDecl) IsBoolValueMethod() (boolMethod bool) {
	_, notBoolMethod := notBoolValueMethod[f.CalleeName()]
	if f.ReturnType().Ty == TyInt && !notBoolMethod {
		boolMethod = true
	}
	return boolMethod
}

func (f FuncDecl) ReturnGoType() string {
	if f.IsBoolValueMethod() {
		return "bool"
	}
	retType := f.ReturnType()
	if retType.Ty == TyStringUserfreeT {
		return "string"
	}
	return retType.GoType()
}

func (f FuncDecl) HasOutParam() (has bool) {
	for _, p := range f.Params() {
		if p.IsOutParam() || p.IsInOutParam() {
			has = true
		}
	}
	return has
}

func (f FuncDecl) HasReturnValue() (has bool) {
	return f.ReturnGoType() != ""
}

func handleTypedef(base DeclCommon) (decl Decl) {
	dso := base.d.DeclarationSpecifiers.DeclarationSpecifiersOpt
	ds := dso.DeclarationSpecifiers
	if ds.Case != 1 {
		log.Panicf("T349: Unhandled: %v\n", ds)
	}
	name := base.Token().Name()
	switch ds.TypeSpecifier.Case {
	case 11: // StructOrUnionSpecifier
		sdecl := handleStruct(base, ds.TypeSpecifier.StructOrUnionSpecifier)
		if s, ok := sdecl.(*CefClassDecl); ok && s.St == StUnknown {
			log.Printf("OUT-S: type %s C.%s", s.GoName(), name)
		}
		decl = sdecl
	case 12: // EnumSpecifier
		base.Dk = DkEnum
		edecl := EnumDecl{base, nil}
		edecl.Enums = getEnumSpecifier(ds.TypeSpecifier.EnumSpecifier)
		log.Printf("OUT-E: type %s C.%s", edecl.GoName(), name)
		decl = &edecl
	case 13: // TYPEDEFNAME
		typeDefName := Token(ds.TypeSpecifier.Token).Name()
		sdecl := SimpleDecl{base}
		if _, ok := primitiveTypeDef[name]; ok {
			sdecl.Dk = DkSimple
			log.Printf("OUT-Pri: type %s C.%s", sdecl.GoName(), name)
		} else if isSimpleDefName(name) {
			sdecl.Dk = DkSimple
			log.Printf("OUT-S: type %s C.%s", sdecl.GoName(), name)
		} else if isStructDefName(name) {
			sdecl.Dk = DkStruct
			log.Printf("OUT-St: type %s C.%s", sdecl.GoName(), name)
		} else {
			log.Panicf("T595: %s :%s, %s\n", name, typeDefName, base.Token().FilePos())
		}
		decl = &sdecl
	default:
		sdecl := SimpleDecl{base}
		if isSimpleDefName(name) {
			sdecl.Dk = DkSimple
		} else {
			if _, ok := primitiveTypeDef[name]; ok {
				sdecl.Dk = DkSimple
			} else {
				log.Panicf("T353: %d. No Struct: %s, %v\n", ds.TypeSpecifier.Case, sdecl.FilePos(), sdecl.d)
			}
		}
		decl = &sdecl
	}
	return decl
}

func handleStruct(base DeclCommon, st *cc.StructOrUnionSpecifier) (decl Decl) {
	name := base.Token().Name()
	tag := getStructTag(st)
	log.Printf("T364: name:%s struct tag: %s\n", name, tag.Name())
	// log.Printf("T355: Struct Def: %v\n", ds.TypeSpecifier)
	sm := []cc.StructDeclaration{}
	for m := st.StructDeclarationList; m != nil; m = m.StructDeclarationList {
		member := *m.StructDeclaration
		// member.StructDeclarationList = nil
		sm = append(sm, member)
	}
	if isStructDefName(name) {
		stDecl := &StructDecl{base, nil}
		decl = stDecl
		stDecl.Common().Dk = DkStruct
		for _, m := range sm {
			ts := getTypeSpecifier(m)
			ty := getTsType(ts)
			d := getDeclarator(m)
			if d.PointerOpt != nil {
				ty.Pointer++
				if d.PointerOpt.Pointer.Case == 1 {
					log.Panicf("T1103: %s, %v\n", ty.Ty, m)
				}
			}
			dd := d.DirectDeclarator
			stDecl.Members = append(stDecl.Members, Member{Token(dd.Token), ty, &m, stDecl})
		}
	} else {
		var stType StructType
		var stBase string
		base.Dk = DkCefClass
		var sdecl *CefClassDecl
		sdecl = &CefClassDecl{base, StUnknown, "", nil}
	MLOOP:
		for i, m0 := range sm {
			if i == 0 {
				stType, stBase = checkBase(m0)
				switch stType {
				case StRefCounted:
					sdecl.St = StRefCounted
					sdecl.BaseType = stBase
					log.Printf("T737: %s\n", name)
				case StScoped:
					sdecl.St = StScoped
					sdecl.BaseType = stBase
					log.Printf("T771: %s\n", name)
				default:
					log.Printf("T743: %s\n", name)
					break MLOOP
				}
				Defs[name] = sdecl // for following method handling
			} else {
				fp := getFuncPointer(sdecl, m0)
				if fp.Funcname != noToken {
					log.Printf("T647: Fn: %s", fp.Funcname.Name())
					ts := getTypeSpecifier(m0)
					ty := getTsType(ts)
					log.Printf("T372:   Case %d, %s %t", ts.Case, ty.Token.Name(), ty.Typedef)
					for i, p := range fp.params {
						log.Printf("T378:   p%d, %s", i, p)
					}
					sdecl.Methods = append(sdecl.Methods, fp)
				}
			}
		}

		switch sdecl.St {
		case StUnknown:
			if isSimpleDefName(name) {
				decl = &SimpleDecl{*sdecl.Common()}
				decl.Common().Dk = DkSimple
			} else {
				log.Panicf("T1163: Can not Handle, %v", sdecl.d)
			}
		default:
			decl = sdecl
		}
		if decl == nil {
			log.Panicf("T733:\n")
		}
	}
	Defs[name] = decl
	return decl

}

func getTypeSpecifier(sd cc.StructDeclaration) *cc.TypeSpecifier {
	sq := sd.SpecifierQualifierList
	if sq.Case != 0 {
		log.Panicf("T415: SpecifierQualifierList.Case != 0, %v\n", sd)
	}
	return sq.TypeSpecifier
}

func getTsType(ts *cc.TypeSpecifier) (ty Type) {
	ty.Token = noToken // No Implemented
	switch ts.Case {
	case 0, 1, 3, 4, 5, 6, 8: // char, void, int, float, double, unsigned
		tm := map[int]Ty{
			0: TyVoid,
			1: TyChar,
			3: TyInt,
			4: TyLong,
			5: TyFloat,
			6: TyDouble,
			8: TyUnsigned,
		}
		ty.Ty = tm[ts.Case]
		ty.Token = Token(ts.Token)
	case 11: // struct
		if ts.StructOrUnionSpecifier.Case != 1 {
			log.Panicf("T426: Not struct identifier, %v", ts)
		}
		ty.Token = Token(ts.StructOrUnionSpecifier.Token)
		name := tagToTypdefName(ty.Token.Name())
		if decl, ok := Defs[name]; ok {
			if _, ok := decl.(*SimpleDecl); ok {
				ty.Ty = TyStructSimple
			} else if s, ok := decl.(*CefClassDecl); ok {
				switch s.St {
				case StRefCounted:
					ty.Ty = TyStructRefCounted
				case StScoped:
					ty.Ty = TyStructScoped
				case StYetNotDefined:
					ty.Ty = TyStructNotDefined
				default:
					log.Panicf("T1271: %s, %s, %v\n", name, decl.Common().Dk, ts)
				}
			} else if _, ok := decl.(*StructDecl); ok {
				ty.Ty = TyStructSimple
			} else if _, ok := decl.(*UnhandledDecl); ok {
				ty.Ty = TyStructUnhandled
			} else {
				log.Panicf("T775: %s, %v\n", name, decl)
			}
		} else {
			log.Panicf("T719: %s: %v", name, ts)
		}
	case 13: // typedef
		ty.Token = Token(ts.Token)
		name := ty.Token.Name()
		ty.Typedef = true
		if t, ok := primitiveTypeDef[name]; ok {
			ty.Ty = t
		} else {
			if decl, ok := Defs[name]; ok {
				switch decl.Common().Dk {
				case DkCefClass:
					s := decl.(*CefClassDecl)
					switch s.St {
					case StRefCounted:
						ty.Ty = TyStructRefCounted
					case StScoped:
						ty.Ty = TyStructScoped
					default:
						log.Panicf("T739: %s, %d, %v\n", name, s.St, ts)
					}
				case DkEnum:
					ty.Ty = TyEnum
				case DkSimple:
					ty.Ty = TySimple
				case DkUnhandled:
					ty.Ty = TyStructUnhandled
				case DkStruct:
					ty.Ty = TyStructSimple
				default:
					log.Panicf("1256: %s, %s, %v\n", name, decl.Common().Dk, ts)
				}
			} else if name == "cef_string_t" {
				ty.Ty = TyStringT
			} else if name == "cef_string_userfree_t" {
				ty.Ty = TyStringUserfreeT
			} else if name == "MSG" {
				ty.Ty = TyMSG
			} else if name == "HCURSOR" {
				ty.Ty = TyHCURSOR
			} else {
				log.Panicf("T690: %s: %v", name, ts)
			}
		}
	default:
		log.Panicf("T428: %v\n", ts)
	}
	return ty
}

func tagToTypdefName(tag string) (name string) {
	if strings.HasPrefix(tag, "_") {
		name = strings.Replace(tag, "_", "", 1)
	} else {
		name = tag
	}
	return name
}

func checkBase(sd cc.StructDeclaration) (stType StructType, stBase string) {
	tq := getTypeSpecifier(sd)
	t := Token(tq.Token)
	name := t.Name()
	if tq.Case != 13 { // TYPEDEFNAME
		log.Printf("T419:   Not Typedef Name, %v\n", name)
		return stType, stBase
	}
	switch name {
	case "cef_base_ref_counted_t":
		stType = StRefCounted
		stBase = name
	case "cef_base_scoped_t":
		stType = StScoped
		stBase = name
	default:
		b := Defs[name]
		if c, ok := b.(*CefClassDecl); ok {
			if c.St == StRefCounted {
				stType = StRefCounted
				stBase = name
				log.Printf("T1362:   %s is base name\n", name)
			} else {
				log.Panicf("T1364: Unpredicted baseL %s, %v\n", name, sd)
			}
		} else {
			log.Printf("T1367: %s, %v\n", name, t.FilePos())
			return stType, stBase
		}
	}
	return stType, stBase
	// c := v.(lex.Char)
	// r := c.Rune
	// s := yySymName(int(r))
}

func getEnumSpecifier(es *cc.EnumSpecifier) (enums []Token) {
	enums = []Token{}
	for el := es.EnumeratorList; el != nil; el = el.EnumeratorList {
		t := Token(el.Enumerator.EnumerationConstant.Token)
		enums = append(enums, t)
	}
	return enums
}

func (decl DeclCommon) GoName() string {
	return decl.Token().GoName()
}

func (decl DeclCommon) BaseName() string {
	return decl.Token().BaseName()
}

func (decl DeclCommon) CefName() string {
	return decl.Token().Name()
}

func (decl DeclCommon) RetStr() (str string) {
	d := getFirstDeclarator(decl.d)
	ret := d.Type.Result().String()
	if ret == "int" {
		str = "int"
	}
	return str
}

func (decl DeclCommon) Call() (call string) {
	d := getFirstDeclarator(decl.d)
	ret := d.Type.Result().String()
	call = "C." + decl.CefName() + "()"
	if ret == "int" {
		call = "return int(" + call + ")"
	}
	return call
}

type Param struct {
	Callee
	Variadic         bool
	paramNameToken   Token
	paramDeclaration *cc.ParameterDeclaration
}

func (p Param) Name() string {
	return p.paramNameToken.Name()
}

func (p Param) GoTypeIn() (t string) {
	if p.IsInOutParam() {
		if p.Type().Ty == TyStringT && p.Type().Pointer == 1 {
			return p.GoType()
		} else {
			if p.IsBoolParam() {
				return "bool"
			}
			return p.Type().Deref().GoType()
		}
	} else {
		return p.GoType()
	}
}

func (p Param) GoType() (t string) {
	if bs, _ := p.IsByteSliceParam(); bs {
		return "[]byte"
	}
	if s, _ := p.IsSliceParam(); s {
		return "[]" + p.Type().Deref().GoType()
	}
	if p.IsBoolParam() {
		return "bool"
	}

	pType := p.Type()

	return pType.GoType()
}

func (p Param) CType() (t string) {
	pType := p.Type()
	t = pType.CType()

	return t
}

func (p Param) GoCType() (t string) {
	return p.Type().GoCType()
}

func (p Param) IsRefCountedClass() bool {
	t := p.Type()
	return t.IsRefCountedClass()
}

func (p Param) IsScopedClass() bool {
	return p.Type().IsScopedClass()
}

func (p Param) IsHandlerClass() bool {
	return p.Type().IsHandlerClass()
}

func (p Param) IsOutParam() (isOut bool) {
	_, isOut = outParameter[p.CalleeName()+"::"+p.Name()]
	return isOut
}

func (p Param) IsInOutParam() (isOut bool) {
	_, isOut = inOutParameter[p.CalleeName()+"::"+p.Name()]
	return isOut
}

func (p Param) IsByteSliceParam() (isByteSlice bool, maxLengthArg string) {
	maxLengthArg, isByteSlice = byteSliceParameter[p.CalleeName()+"::"+p.Name()]
	return isByteSlice, maxLengthArg
}

func (p Param) IsByteSliceLengthParam() (isByteSliceLength bool) {
	_, isByteSliceLength = byteSliceLengthParameter[p.CalleeName()+"::"+p.Name()]
	return isByteSliceLength
}

func (p Param) IsSliceParam() (isSlice bool, maxLengthArg string) {
	maxLengthArg, isSlice = sliceParameter[p.CalleeName()+"::"+p.Name()]
	return isSlice, maxLengthArg
}

func (p Param) IsSliceLengthParam() (isSliceLength bool) {
	_, isSliceLength = sliceLengthParameter[p.CalleeName()+"::"+p.Name()]
	return isSliceLength
}

func (p Param) IsBoolParam() (isBool bool) {
	_, isBool = boolParameter[p.CalleeName()+"::"+p.Name()]
	return isBool
}

func (p Param) String() string {
	var suffix string
	var prefix string
	var pointer string
	paramName := p.paramNameToken.Name()
	pType := p.Type()
	if pType.Const {
		prefix = "const "
	}
	for i := 0; i < pType.Pointer; i++ {
		pointer += "*"
	}

	if pType.Typedef {
		suffix = " (typedef)"
		if d, ok := Defs[pType.Token.Name()]; ok {
			if d.Common().Dk == DkEnum {
				suffix = " (enum)"
			}
			if d.Common().Dk == DkSimple {
				suffix = " (simple)"
			}
		}
	}

	switch pType.Ty {
	case TyStructUnhandled:
		suffix = " (unhandled)"
		return paramName + ": " + prefix + pointer + "struct " + fmt.Sprintf("%s", pType.Token.Name()) + suffix
	case TyStructSimple:
		suffix = " (simple)"
		return paramName + ": " + prefix + pointer + "struct " + fmt.Sprintf("%s", pType.Token.Name()) + suffix
	case TyStructRefCounted, TyStructScoped:
		return paramName + ": " + prefix + pointer + "struct " + fmt.Sprintf("%s", pType.Token.Name()) + suffix
	case TyUnknown:
		log.Panicf("T852: %s, %v, %v\n", paramName, pType.Ty, pType.Token.FilePos())
	}
	return paramName + ": " + prefix + pointer + pType.Ty.String() + suffix
}

func getDeclarator(sd cc.StructDeclaration) *cc.Declarator {
	return sd.StructDeclaratorList.StructDeclarator.Declarator
}

func getDirectDeclarator(sd cc.StructDeclaration) (dd *cc.DirectDeclarator) {
	d := getDeclarator(sd)
	return d.DirectDeclarator
}

func getFuncPointer(sdecl *CefClassDecl, sd cc.StructDeclaration) (methodP *MethodDecl) {
	dd := getDirectDeclarator(sd)
	m := &MethodDecl{noToken, nil, sd, nil, sdecl}
	switch dd.Case {
	case 6: // DirectDeclarator '(' ParameterTypeList ')'
		f := dd.DirectDeclarator
		if f.Case != 1 {
			log.Panicf("T496: Not Function %v\n", dd)
		}
		if f.Declarator.PointerOpt == nil {
			log.Panicf("T493: not pointer %v\n", dd)
		}
		m.Funcname = Token(f.Declarator.DirectDeclarator.Token)
		if _, um := unGenerateMethod[m.CalleeName()]; um {
			log.Printf("T1385: Skip: %s\n", m.CalleeName())
			return &MethodDecl{noToken, nil, sd, nil, sdecl}
		}
		for p := dd.ParameterTypeList.ParameterList; p != nil; p = p.ParameterList {
			m.params = append(m.params, getParam(p.ParameterDeclaration, m))
		}
		if dd.ParameterTypeList.Case == 1 { //ParameterList ',' "..."  // Case 1
			variadic := Param{nil, true, noToken, nil}
			m.params = append(m.params, variadic)
		}
	default:
		log.Panicf("T525: %v\n", dd)
	}
	return m
}

func getDeclType(ds *cc.DeclarationSpecifiers) (ty Type) {
	for ds != nil {
		switch ds.Case {
		case 1: // TypeSpecifier DeclarationSpecifiersOpt
			ty0 := getTsType(ds.TypeSpecifier)
			if ty.Ty == TyUnknown {
				c := ty.Const
				ty = ty0
				ty.Const = c
			} else if ty.Ty == TyUnsigned && ty0.Ty == TyLong {
				ty.Ty = TyULong
			} else if ty.Ty == TyULong && ty0.Ty == TyLong {
				ty.Ty = TyULongLong
			} else if ty.Ty == TyLong && ty0.Ty == TyLong {
				ty.Ty = TyLongLong
			} else if ty.Ty != TyUnknown {
				log.Panicf("T1161: %v\n", ds)
			}
		case 2: // TypeQualifier DeclarationSpecifiersOpt
			tq := ds.TypeQualifier
			if tq.Case == 0 {
				ty.Const = true
			} else {
				log.Panicf("T1181: %v\n", tq)
			}
		default:
			log.Panicf("T526: %v\n", ds)
		}
		if ds.DeclarationSpecifiersOpt != nil {
			ds = ds.DeclarationSpecifiersOpt.DeclarationSpecifiers
		} else {
			ds = nil
		}
	}
	return ty
}

func (p Param) Type() (ty Type) {
	pDec := p.paramDeclaration
	ty = getDeclType(pDec.DeclarationSpecifiers)

	if pDec.Declarator.PointerOpt != nil {
		ty.Pointer = 1
		ty.Tq = append(ty.Tq, TqPointer)
		p := pDec.Declarator.PointerOpt.Pointer
		if p.Case == 1 {
			if p.TypeQualifierListOpt != nil && p.TypeQualifierListOpt.TypeQualifierList != nil {
				if p.TypeQualifierListOpt.TypeQualifierList.Case == 0 &&
					p.TypeQualifierListOpt.TypeQualifierList.TypeQualifier.Case == 0 {
					ty.Tq = append(ty.Tq, TqConst)
				} else {
					log.Panicf("T1410: Too many TypeQualifire%v\n", p)
				}
			}
			ty.Pointer = 2
			ty.Tq = append(ty.Tq, TqPointer)
			if p.Pointer.Case == 1 {
				log.Panicf("T1416: %v\n", p)
			}
		}
	}
	return ty
}

func getParam(p *cc.ParameterDeclaration, f Callee) (param Param) {
	switch p.Case {
	case 0:
		paramNameToken := getDirectDeclaratorToken(p.Declarator.DirectDeclarator)
		param = Param{f, false, paramNameToken, p}
	default:
		log.Panicf("T1209: %v, %v\n", p, f)
	}
	return param
}

func (mp *MethodDecl) CalleeName() string {
	return mp.sdecl.CefName() + "::" + mp.Funcname.Name()
}

func (mp *MethodDecl) ClassBaseName() string {
	return mp.sdecl.BaseName()
}

func (m MethodDecl) ReturnType() (retType Type) {
	ts := getTypeSpecifier(m.sd)
	// log.Printf("T811: %s\n", getDeclarator(m.sd).Type)
	retType = getTsType(ts)

	if m.sd.Case == 0 {
		pointer := m.sd.StructDeclaratorList.StructDeclarator.Declarator.PointerOpt
		if pointer != nil {
			retType.Pointer = 1
			retType.Tq = append(retType.Tq, TqPointer)
			if pointer.Pointer.Case != 0 {
				log.Panicf("T991: %v\n", pointer)
			}
		}
	} else {
		log.Panicf("T985: %v\n", m.sd)
	}
	return retType
}

func (m MethodDecl) IsBoolValueMethod() (boolMethod bool) {
	_, notBoolMethod := notBoolValueMethod[m.CalleeName()]
	if m.ReturnType().Ty == TyInt && !notBoolMethod {
		boolMethod = true
	}
	return boolMethod
}

func (m MethodDecl) ReturnGoType() string {
	if m.IsBoolValueMethod() {
		return "bool"
	}
	retType := m.ReturnType()
	if retType.Ty == TyStringUserfreeT {
		return "string"
	}
	return retType.GoType()
}

func (m MethodDecl) Type() string {
	d := m.sd.StructDeclaratorList.StructDeclarator.Declarator
	return d.Type.String()
}

func (t Type) Name() string {
	return t.Token.Name()
}

func (t Type) CType() (ret string) {
	switch t.Ty {
	case TyVoid:
		ret = "void"
	case TyInt:
		ret = "int"
	case TyInt32:
		ret = "int32"
	case TyInt64:
		ret = "int64"
	case TyUint16:
		ret = "uint16"
	case TyUint32:
		ret = "uint32"
	case TyUint64:
		ret = "uint64"
	case TySizeT:
		ret = "size_t"
	case TyStringT:
		ret = "cef_string_t"
	case TyStringUserfreeT:
		ret = "cef_string_userfree_t"
	case TyStructRefCounted, TyStructScoped:
		ret = "struct " + t.Name()
	case TyStructSimple:
		if t.Name()[0] == '_' {
			ret = "struct " + t.Name()
		} else {
			ret = t.Name()
		}
	case TyEnum:
		ret = t.Name()
	case TySimple:
		ret = t.Name()
	case TyTimeT:
		ret = "time_t"
	case TyFloat:
		ret = "float"
	case TyDouble:
		ret = "double"
	case TyLongLong:
		ret = "long long"
	case TyHWND:
		ret = "cef_window_handle_t"
	case TyMSG:
		if t.Pointer == 1 {
			t.Tq = []TypeQualifier{}
			ret = "cef_event_handle_t"
		} else {
			log.Panicf("T1556: %v\n", t)
		}
	case TyHCURSOR:
		ret = "cef_cursor_handle_t"
	case TyHINSTANCE:
		ret = "HINSTANCE"
	case TyDWORD:
		ret = "DWORD"
	case TyHMENU:
		ret = "HMENU"
	default:
		log.Panicf("T1561: %v\n", t)
	}
	for _, tq := range t.Tq {
		switch tq {
		case TqPointer:
			ret += "*"
		case TqConst:
			ret += " const "
		default:
			log.Panicf("T1535: Wrong Type Qualifire: %v", t)
		}
	}
	if t.Const {
		ret = "const " + ret
	}
	return ret
}

func (t Type) GoType() (ret string) {
	pointerOffset := 0
	switch t.Ty {
	case TyVoid:
		if t.Pointer == 1 {
			pointerOffset += 1
			ret = "unsafe.Pointer"
		} else {
			ret = ""
		}
	case TyInt:
		ret = "int"
	case TyInt32:
		ret = "int32"
	case TyInt64:
		ret = "int64"
	case TyUint16:
		ret = "uint16"
	case TyUint32:
		ret = "uint32"
	case TyUint64:
		ret = "uint64"
	case TySizeT:
		ret = "int64"
	case TyFloat:
		ret = "float32"
	case TyDouble:
		ret = "float64"
	case TyStructRefCounted, TyStructScoped, TyStructSimple, TySimple:
		ret = t.Token.GoName()
	case TyStringT:
		pointerOffset += 1
		ret = "string"
		if t.Pointer == 0 {
			log.Panicf("T1748: Not cef_string_t *\n")
		}
	case TyEnum:
		ret = t.Token.GoName()
	case TyStructNotDefined:
		if decl, ok := Defs[tagToTypdefName(t.Token.Name())]; ok {
			if s, ok := decl.(*CefClassDecl); ok {
				switch s.St {
				case StRefCounted:
					t.Ty = TyStructRefCounted
				case StScoped:
					t.Ty = TyStructScoped
				default:
					log.Panicf("T1164: %s: %v\n", t.Token.Name(), t)
				}
				return t.GoType()
			} else {
				log.Panicf("T1168: %v\n", t)
			}
		} else {
			log.Panicf("T1171: %v\n", t)
		}
	case TyTimeT:
		ret = "time.Time"
	case TyLongLong:
		ret = "int64"
	case TyHWND:
		ret = "CWindowHandleT"
	case TyMSG:
		pointerOffset += 1
		ret = "CEventHandleT"
	case TyHCURSOR:
		ret = "CCursorHandleT"
	case TyHINSTANCE:
		ret = "WinHinstance"
	case TyDWORD:
		ret = "WinDword"
	case TyHMENU:
		ret = "WinHmenu"
	default:
		log.Panicf("T841: %s, %s\n", t.Ty, t.Token)
	}

	for i := pointerOffset; i < t.Pointer; i++ {
		ret = "*" + ret
	}

	return ret
}

func (t Type) Deref() (t0 Type) {
	t0 = t
	if t0.Pointer > 0 {
		t0.Pointer--
	} else {
		log.Panicf("T1461: %v\n", t)
	}
	return t0
}

func (t Type) BaseName() (ret string) {
	return t.Token.BaseName()
}

func (t Type) TitleCase() (ret string) {
	return t.Token.TitleCase()
}

func (t Type) GoCType() (ct string) {
	pointer := t.Pointer
	if t.Ty == TyVoid {
		switch t.Pointer {
		case 0:
			return ""
		case 1, 2:
			pointer--
			ct = "VOIDP"
		default:
			log.Panicf("T1263: %v\n", t)
		}
	} else {
		switch t.Ty {
		case TyStructRefCounted, TyStructScoped, TyStructSimple:
			ct = "cef_" + t.Token.BaseName() + "_t"
		case TyLongLong:
			ct = "LONGLONG"
		case TyULongLong:
			ct = "ULONGLONG"
		default:
			ct = t.Name()
		}
	}
	ct = "C." + ct
	for i := 0; i < pointer; i++ {
		ct = "*" + ct
	}
	return ct
}

func (t Type) IsRefCountedClass() bool {
	if t.Pointer == 1 {
		if t.Ty == TyStructRefCounted {
			return true
		}
	}
	return false
}

func (t Type) IsScopedClass() bool {
	if t.Ty == TyStructScoped {
		return true
	}
	return false
}

func isHandlerClass(name string) (c bool) {
	_, c = handlerClasses[name]
	return c
}

func (t Type) IsHandlerClass() bool {
	return isHandlerClass(t.Name())
}

func (t Type) IsSimple() bool {
	n := t.Name()
	return isSimpleDefName(n)
}

func (t Type) IsStruct() bool {
	n := t.Name()
	return isStructDefName(n)
}
