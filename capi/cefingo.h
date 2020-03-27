
#ifndef CEFINGO_H_
#define CEFINGO_H_
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_client_capi.h"
#include "include/capi/cef_request_context_handler_capi.h"
#include "include/capi/cef_urlrequest_capi.h"
#include "include/cef_version.h"
#include "include/capi/views/cef_menu_button_capi.h"
#include "include/capi/views/cef_window_capi.h"
#include "include/capi/views/cef_browser_view_capi.h"
#include "include/capi/views/cef_fill_layout_capi.h"
#include "include/capi/views/cef_box_layout_capi.h"
#include "include/capi/views/cef_scroll_view_capi.h"
#include "include/capi/views/cef_textfield_capi.h"
#include "cefingo_base.h"
#include "cefingo_gen.h"

typedef void *VOIDP;
typedef long long LONGLONG;
typedef unsigned long long ULONGLONG;

extern void cefingo_init();

#endif // CEFINGO_H_
