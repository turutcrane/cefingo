
#ifndef CEFINGO_H_
#define CEFINGO_H_
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_client_capi.h"
#include "include/capi/cef_v8_capi.h"
#include "include/cef_version.h"
#include "cefingo_base.h"

#define FUNCNAME_TO_GO ((char*)__func__)
extern void cefingo_cslogf(const char *fn, const char *format, ...);

CEFINGO_REF_COUNTER_WRAPPER(cef_app_t, cefingo_app_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_client_t, cefingo_client_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_browser_process_handler_t, cefingo_browser_process_handler_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_life_span_handler_t, cefingo_life_span_handler_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_render_process_handler_t, cefingo_render_process_handler_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_v8array_buffer_release_callback_t, cefingo_v8array_buffer_release_callback_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_v8handler_t, cefingo_v8handler_wrapper_t);
CEFINGO_REF_COUNTER_WRAPPER(cef_load_handler_t, cefingo_load_handler_wrapper_t);

extern void cefingo_construct_life_span_handler(cefingo_life_span_handler_wrapper_t *handler);
extern void cefingo_construct_browser_process_handler(cefingo_browser_process_handler_wrapper_t *handler);
extern void cefingo_construct_client(cefingo_client_wrapper_t* client);
extern void cefingo_construct_app(cefingo_app_wrapper_t* app);
extern void cefingo_construct_render_process_handler(cefingo_render_process_handler_wrapper_t* handler);
extern void cefingo_construct_load_handler(cefingo_load_handler_wrapper_t* handler);

extern cef_v8context_t *cefingo_frame_get_v8context(cef_frame_t *self);

#endif // CEFINGO_H_
