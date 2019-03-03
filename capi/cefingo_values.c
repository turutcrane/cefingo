#include "cefingo.h"
#include "_cgo_export.h"

int cefingo_list_value_is_valid(struct _cef_list_value_t* self)
{
    return self->is_valid(self);
}

int cefingo_list_value_is_owned(struct _cef_list_value_t* self)
{
    return self->is_owned(self);
};

int cefingo_list_value_is_read_only(struct _cef_list_value_t* self)
{
    return self->is_read_only(self);
}

int cefingo_list_value_is_same(struct _cef_list_value_t* self,
                               struct _cef_list_value_t* that)
{
    return self->is_same(self, that);
}

int cefingo_list_value_is_equal(struct _cef_list_value_t* self,
                                struct _cef_list_value_t* that)
{
    return self->is_equal(self, that);
}

struct _cef_list_value_t* cefingo_list_value_copy(struct _cef_list_value_t* self)
{
    return self->copy(self);
}

int cefingo_list_value_set_size(struct _cef_list_value_t* self, size_t size)
{
    return self->set_size(self, size);
}

size_t cefingo_list_value_get_size(struct _cef_list_value_t* self)
{
    return self->get_size(self);
}

int cefingo_list_value_clear(struct _cef_list_value_t* self)
{
    return self->clear(self);
}

int cefingo_list_value_remove(struct _cef_list_value_t* self, size_t index)
{
    return self->remove(self, index);
};

cef_value_type_t cefingo_list_value_get_type(struct _cef_list_value_t* self,
        size_t index)
{
    return self->get_type(self, index);
}

struct _cef_value_t* cefingo_list_value_get_value(struct _cef_list_value_t* self,
        size_t index)
{
    return self->get_value(self, index);
}

int cefingo_list_value_get_int(struct _cef_list_value_t* self, size_t index)
{
    return self->get_int(self, index);
}

int cefingo_list_value_get_bool(struct _cef_list_value_t* self, size_t index)
{
    return self->get_bool(self, index);
}

double cefingo_list_value_get_double(struct _cef_list_value_t* self, size_t index)
{
    return self->get_double(self, index);
}

cef_string_userfree_t
cefingo_list_value_get_string(struct _cef_list_value_t* self, size_t index)
{
    self->get_string(self, index);
}

struct _cef_binary_value_t* cefingo_list_value_get_binary(
    struct _cef_list_value_t* self, size_t index)
{
    return self->get_binary(self, index);
}

struct _cef_dictionary_value_t* cefingo_list_value_get_dictionary(
    struct _cef_list_value_t* self, size_t index)
{
    return self->get_dictionary(self, index);
}

struct _cef_list_value_t* cefingo_list_value_get_list(
    struct _cef_list_value_t* self, size_t index)
{
    return self->get_list(self, index);
}

int cefingo_list_value_set_value(struct _cef_list_value_t* self,
                                 size_t index,
                                 struct _cef_value_t* value)
{
    return self->set_value(self, index, value);
}

int cefingo_list_value_set_null(struct _cef_list_value_t* self, size_t index)
{
    return self->set_null(self, index);
}

int cefingo_list_value_set_bool(struct _cef_list_value_t* self,
                                size_t index,
                                int value)
{
    return self->set_bool(self, index, value);
}

int cefingo_list_value_set_int(struct _cef_list_value_t* self,
                               size_t index, int value)
{
    return self->set_int(self, index, value);
}

int cefingo_list_value_set_double(struct _cef_list_value_t* self,
                                  size_t index, double value)
{
    return self->set_double(self, index, value);
}

int cefingo_list_value_set_string(struct _cef_list_value_t* self,
                                  size_t index, cef_string_t* value)
{
    return self->set_string(self, index, value);
}

int cefingo_list_value_set_binary(struct _cef_list_value_t* self,
                                  size_t index, struct _cef_binary_value_t* value)
{
    return self->set_binary(self, index, value);
}

int cefingo_list_value_set_dictionary(struct _cef_list_value_t* self,
                                      size_t index, struct _cef_dictionary_value_t* value)
{
    return self->set_dictionary(self, index, value);
}

int cefingo_list_value_set_list(struct _cef_list_value_t* self,
                                size_t index, struct _cef_list_value_t* value)
{
    return self->set_list(self, index,value);
}
