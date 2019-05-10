#include "cefingo.h"
#include "_cgo_export.h"

int cefingo_v8context_is_valid(cef_v8context_t *self)
{
    return self->is_valid(self);
}

cef_browser_t *cefingo_v8context_get_browser(cef_v8context_t* self)
{
    return self->get_browser(self);
}

cef_frame_t *cefingo_v8context_get_frame(cef_v8context_t *self)
{
    return self->get_frame(self);
}

cef_v8value_t *cefingo_v8context_get_global(cef_v8context_t *self)
{
    return self->get_global(self);
}

int cefingo_v8value_is_valid(cef_v8value_t* self)
{
    return self->is_valid(self);
}

int cefingo_v8value_is_undefined(cef_v8value_t* self)
{
    return self->is_undefined(self);
}

int cefingo_v8value_is_null(cef_v8value_t* self)
{
    return self->is_null(self);
}

int cefingo_v8value_is_bool(cef_v8value_t* self)
{
    return self->is_bool(self);
}

int cefingo_v8value_is_int(cef_v8value_t* self)
{
    return self->is_int(self);
}

int cefingo_v8value_is_uint(cef_v8value_t* self)
{
    return self->is_uint(self);
}

int cefingo_v8value_is_double(cef_v8value_t* self)
{
    return self->is_double(self);
}

int cefingo_v8value_is_date(cef_v8value_t* self)
{
    return self->is_date(self);
}

int cefingo_v8value_is_string(cef_v8value_t* self)
{
    return self->is_string(self);
}

int cefingo_v8value_is_object(cef_v8value_t* self)
{
    return self->is_object(self);
};

int cefingo_v8value_is_array(cef_v8value_t* self)
{
    return self->is_array(self);
};

int cefingo_v8value_is_array_buffer(cef_v8value_t* self)
{
    return self->is_array_buffer(self);
};

int cefingo_v8value_is_function(cef_v8value_t* self)
{
    return self->is_function(self);
}

int cefingo_v8value_is_same(cef_v8value_t* self, cef_v8value_t* that)
{
    return self->is_same(self, that);
}

int cefingo_v8value_get_bool_value(cef_v8value_t* self)
{
    return self->get_bool_value(self);
}

int cefingo_v8value_get_int_value(cef_v8value_t* self)
{
    return self->get_int_value(self);
}

uint32 cefingo_v8value_get_uint_value(cef_v8value_t* self)
{
    return self->get_uint_value(self);
}

double cefingo_v8value_get_double_value(cef_v8value_t* self)
{
    return self->get_double_value(self);
}

cef_time_t cefingo_v8value_get_date_value(cef_v8value_t* self)
{
    return self->get_date_value(self);
}

cef_string_userfree_t cefingo_v8value_get_string_value(cef_v8value_t* self)
{
    return self->get_string_value(self);
}

int cefingo_v8value_has_exception(cef_v8value_t* self)
{
    return self->has_exception(self);
}

cef_v8exception_t *cefingo_v8value_get_exception(cef_v8value_t *self)
{
    return self->get_exception(self);
}

int cefingo_v8value_clear_exception(cef_v8value_t *self)
{
    return self->clear_exception(self);
}

int cefingo_v8value_has_value_bykey(cef_v8value_t* self,
                                    const cef_string_t* key)
{
    return self->has_value_bykey(self, key);
}

int cefingo_v8value_delete_value_bykey(cef_v8value_t* self,
                                       const cef_string_t* key)
{
    return self->delete_value_bykey(self, key);
}

cef_v8value_t* cefingo_v8value_get_value_bykey(
    struct _cef_v8value_t* self,
    const cef_string_t* key)
{
    return self->get_value_bykey(self, key);

}

int cefingo_v8context_set_value_bykey(cef_v8value_t* self,
                                      cef_string_t* key,
                                      cef_v8value_t* value,
                                      cef_v8_propertyattribute_t attribute
                                     )
{
    return self->set_value_bykey(self, key, value, attribute);
    // return self->set_value_bykey(self, (const cef_string_t*) key, value, attribute);
}

cef_string_userfree_t cefingo_v8value_get_function_name(cef_v8value_t* self)
{
    return self->get_function_name(self);
}

cef_v8value_t* cefingo_v8value_execute_function(
    cef_v8value_t* self,
    cef_v8value_t* object,
    size_t argumentsCount,
    cef_v8value_t ** arguments)
{
    return self->execute_function(self, object, argumentsCount, arguments);
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

cef_v8array_buffer_release_callback_t *cefingo_construct_v8array_buffer_release_callback(cefingo_v8array_buffer_release_callback_wrapper_t *callback)
{

    initialize_cefingo_base_ref_counted(
        offsetof(__typeof__(*callback), counter),
        (cef_base_ref_counted_t*) callback);

    callback->body.release_buffer = cefingo_v8array_buffer_release_callback_release_buffer;

    return (cef_v8array_buffer_release_callback_t*) callback;
}

typedef int(CEF_CALLBACK* cefingo_execute_t)(struct _cef_v8handler_t* self,
        const cef_string_t* name,
        struct _cef_v8value_t* object,
        size_t argumentsCount,
        struct _cef_v8value_t* const* arguments,
        struct _cef_v8value_t** retval,
        cef_string_t* exception);

cef_v8handler_t *cefingo_construct_v8handler(cefingo_v8handler_wrapper_t *handler)
{
    initialize_cefingo_base_ref_counted(
        offsetof(__typeof__(*handler), counter),
        (cef_base_ref_counted_t*) handler);

    handler->body.execute = (cefingo_execute_t) cefingo_v8handler_execute;

    return (cef_v8handler_t *)handler;
}

int cefingo_v8context_enter(cef_v8context_t* self)
{
    return self->enter(self);
}

int cefingo_v8context_exit(cef_v8context_t* self)
{
    return self->exit(self);
}

int cefingo_v8context_is_same(cef_v8context_t* self, cef_v8context_t *that)
{
    return self->is_same(self, that);
}

int cefingo_v8context_eval(cef_v8context_t* self,
                           const cef_string_t* code,
                           const cef_string_t* script_url,
                           int start_line,
                           struct _cef_v8value_t** retval,
                           struct _cef_v8exception_t** exception)
{
    self->eval(self, code, script_url, start_line, retval, exception);
}

cef_string_userfree_t cefingo_v8exception_get_message(cef_v8exception_t *self)
{
    return self->get_message(self);
}
