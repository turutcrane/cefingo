#ifdef CEFINGO_BUILD
#define _WINDOWS_
typedef void *HINSTANCE;
typedef unsigned int DWORD;
typedef void *HWND;
typedef void *HMENU;
typedef unsigned int UINT;
typedef unsigned int UINT_PTR,*PUINT_PTR;
typedef UINT_PTR WPARAM;
typedef long LONG_PTR,*PLONG_PTR;
typedef LONG_PTR LPARAM;
typedef long LONG;
typedef void *PVOID;
typedef PVOID HANDLE;
typedef HANDLE HICON;
typedef HICON HCURSOR;
typedef struct tagPOINT {
    LONG x;
    LONG y;
} POINT;
typedef struct tagMSG {
    HWND hwnd;
    UINT message;
    WPARAM wParam;
    LPARAM lParam;
    DWORD time;
    POINT pt;
} MSG;
#endif
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_client_capi.h"
#include "include/capi/cef_urlrequest_capi.h"
#include "include/capi/cef_request_context_handler_capi.h"
#include "include/capi/views/cef_menu_button_capi.h"
#include "include/capi/views/cef_window_capi.h"
#include "include/capi/views/cef_browser_view_capi.h"
#include "include/capi/views/cef_fill_layout_capi.h"
#include "include/capi/views/cef_box_layout_capi.h"
#include "include/capi/views/cef_scroll_view_capi.h"
#include "include/capi/views/cef_textfield_capi.h"

#include "include/capi/cef_crash_util_capi.h"
#include "include/capi/cef_file_util_capi.h"
#include "include/capi/cef_origin_whitelist_capi.h"
/*
#include "include/capi/cef_parser_capi.h"
#include "include/capi/cef_path_util_capi.h"
#include "include/capi/cef_process_util_capi.h"
#include "include/capi/cef_server_capi.h"
#include "include/capi/cef_thread_capi.h"
#include "include/capi/cef_trace_capi.h"
#include "include/capi/cef_waitable_event_capi.h"
#include "include/capi/cef_xml_reader_capi.h"
#include "include/capi/cef_zip_reader_capi.h"
*/