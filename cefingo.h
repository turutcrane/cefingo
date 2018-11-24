
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

extern void construct_cefingo_life_span_handler(cefingo_life_span_handler_wrapper_t *handler);
extern void construct_cefingo_browser_process_handler(cefingo_browser_process_handler_wrapper_t *handler);
extern void construct_cefingo_client(cefingo_client_wrapper_t* client);
extern void construct_cefingo_app(cefingo_app_wrapper_t* app);
extern void construct_cefingo_render_process_handler(cefingo_render_process_handler_wrapper_t* handler);

extern cef_v8value_t *v8context_get_global(cef_v8context_t *self);
extern int v8context_set_value_bykey(cef_v8value_t* self, cef_string_t* key,
    cef_v8value_t* value, cef_v8_propertyattribute_t attribute);
extern void construct_cefingo_v8array_buffer_release_callback(cefingo_v8array_buffer_release_callback_wrapper_t *callback);
extern void construct_cefingo_v8handler(cefingo_v8handler_wrapper_t *handler);

#endif // CEFINGO_H_
