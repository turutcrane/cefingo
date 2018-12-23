#ifndef CEFINGO_V8_H_
#define CEFINGO_V8_H_

extern cef_v8value_t *cefingo_v8context_get_global(cef_v8context_t *self);
extern int cefingo_v8context_set_value_bykey(cef_v8value_t* self, cef_string_t* key,
    cef_v8value_t* value, cef_v8_propertyattribute_t attribute);
extern void cefingo_construct_v8array_buffer_release_callback(cefingo_v8array_buffer_release_callback_wrapper_t *callback);
extern void cefingo_construct_v8handler(cefingo_v8handler_wrapper_t *handler);


extern int cefingo_v8value_is_valid(cef_v8value_t* self);
extern int cefingo_v8value_is_function(cef_v8value_t* self);
extern int cefingo_v8value_is_string(cef_v8value_t* self);
extern int cefingo_v8value_is_object(cef_v8value_t* self);

extern cef_string_userfree_t cefingo_v8value_get_string(cef_v8value_t* self);

extern int cefingo_v8value_has_value_bykey(cef_v8value_t* self, const cef_string_t* key);
extern cef_v8value_t* cefingo_v8value_get_value_bykey(cef_v8value_t* self, const cef_string_t* key);

extern cef_string_userfree_t cefingo_v8value_get_function_name(cef_v8value_t* self);
extern cef_v8value_t* cefingo_v8value_execute_function(cef_v8value_t* self,
 cef_v8value_t* object, size_t argumentsCount, cef_v8value_t ** arguments);

extern int cefingo_v8context_enter(cef_v8context_t* self);
extern int cefingo_v8context_exit(cef_v8context_t* self);
extern int cefingo_v8context_is_same(cef_v8context_t* self, cef_v8context_t *that);
extern int cefingo_v8context_eval(cef_v8context_t* self, const cef_string_t* code,
    const cef_string_t* script_url, int start_line,
    struct _cef_v8value_t** retval, struct _cef_v8exception_t** exception);

#endif //CEFINGO_V8_H_
