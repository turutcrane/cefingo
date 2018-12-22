#include "cefingo.h"
#include "cefingo_v8.h"
#include "_cgo_export.h"

cef_v8value_t *cefingo_v8context_get_global(cef_v8context_t *self) {
    return self->get_global(self);
}

int cefingo_v8context_set_value_bykey(cef_v8value_t* self,
    cef_string_t* key,
    cef_v8value_t* value,
    cef_v8_propertyattribute_t attribute
) {
    return self->set_value_bykey(self, key, value, attribute);
    // return self->set_value_bykey(self, (const cef_string_t*) key, value, attribute);
}

int cefingo_v8context_has_value_bykey(cef_v8value_t* self,
    const cef_string_t* key) {
    return self->has_value_bykey(self, key);
}
cef_v8value_t* cefingo_v8context_get_value_bykey(
    struct _cef_v8value_t* self,
    const cef_string_t* key) {
    return self->get_value_bykey(self, key);

}

int cefingo_v8value_is_function(cef_v8value_t* self) {
    return self->is_function(self);
}

int cefingo_v8value_is_string(cef_v8value_t* self) {
    return self->is_string(self);
}
cef_string_userfree_t cefingo_v8value_get_string(cef_v8value_t* self) {
  return self->get_string_value(self);
}

  ///
  // Execute the function using the current V8 context. This function should
  // only be called from within the scope of a cef_v8handler_t or
  // cef_v8accessor_t callback, or in combination with calling enter() and
  // exit() on a stored cef_v8context_t reference. |object| is the receiver
  // ('this' object) of the function. If |object| is NULL the current context's
  // global object will be used. |arguments| is the list of arguments that will
  // be passed to the function. Returns the function return value on success.
  // Returns NULL if this function is called incorrectly or an exception is
  // thrown.
  ///
//   struct _cef_v8value_t*(CEF_CALLBACK* execute_function)(
//       struct _cef_v8value_t* self,
//       struct _cef_v8value_t* object,
//       size_t argumentsCount,
//       struct _cef_v8value_t* const* arguments);

void cefingo_construct_v8array_buffer_release_callback(cefingo_v8array_buffer_release_callback_wrapper_t *callback) {

    initialize_cefingo_base_ref_counted(
        offsetof(__typeof__(*callback), counter),
        (cef_base_ref_counted_t*) callback);

    callback->body.release_buffer = v8array_buffer_release_buffer;

}

typedef int(CEF_CALLBACK* cefingo_execute_t)(struct _cef_v8handler_t* self,
                             const cef_string_t* name,
                             struct _cef_v8value_t* object,
                             size_t argumentsCount,
                             struct _cef_v8value_t* const* arguments,
                             struct _cef_v8value_t** retval,
                             cef_string_t* exception);

void cefingo_construct_v8handler(cefingo_v8handler_wrapper_t *handler) {
    initialize_cefingo_base_ref_counted(
        offsetof(__typeof__(*handler), counter),
        (cef_base_ref_counted_t*) handler);

    handler->body.execute = (cefingo_execute_t) execute;
}

extern int cefingo_v8context_enter(cef_v8context_t* self) {
    return self->enter(self);
}

extern int cefingo_v8context_exit(cef_v8context_t* self) {
    return self->exit(self);
}
