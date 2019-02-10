#ifndef CEFINGO_VALUES_H_
#define CEFINGO_VALUES_H_

extern int cefingo_list_value_is_valid(struct _cef_list_value_t* self);
extern int cefingo_list_value_is_owned(struct _cef_list_value_t* self);
extern int cefingo_list_value_is_read_only(struct _cef_list_value_t* self);
extern int cefingo_list_value_is_same(struct _cef_list_value_t* self,
                                      struct _cef_list_value_t* that);
extern int cefingo_list_value_is_equal(struct _cef_list_value_t* self,
                                       struct _cef_list_value_t* that);
extern struct _cef_list_value_t* cefingo_list_value_copy(struct _cef_list_value_t* self);
extern int cefingo_list_value_set_size(struct _cef_list_value_t* self, size_t size);
extern size_t cefingo_list_value_get_size(struct _cef_list_value_t* self);
extern int cefingo_list_value_clear(struct _cef_list_value_t* self);
extern int cefingo_list_value_remove(struct _cef_list_value_t* self, size_t index);
extern cef_value_type_t cefingo_list_value_get_type(struct _cef_list_value_t* self,
        size_t index);
extern struct _cef_value_t* cefingo_list_value_get_value(struct _cef_list_value_t* self,
        size_t index);
extern int cefingo_list_value_get_int(struct _cef_list_value_t* self, size_t index);
extern int cefingo_list_value_get_bool(struct _cef_list_value_t* self, size_t index);
extern double cefingo_list_value_get_double(struct _cef_list_value_t* self, size_t index);
extern cef_string_userfree_t cefingo_list_value_get_string(struct _cef_list_value_t* self,
        size_t index);
extern struct _cef_binary_value_t* cefingo_list_value_get_binary(
    struct _cef_list_value_t* self, size_t index);
extern struct _cef_dictionary_value_t* cefingo_list_value_get_dictionary(
    struct _cef_list_value_t* self, size_t index);
extern struct _cef_list_value_t* cefingo_list_value_get_list(
    struct _cef_list_value_t* self, size_t index);
extern int cefingo_list_value_set_value(struct _cef_list_value_t* self,
                                        size_t index, struct _cef_value_t* value);
extern int cefingo_list_value_set_null(struct _cef_list_value_t* self, size_t index);
extern int cefingo_list_value_set_bool(struct _cef_list_value_t* self,
                                       size_t index, int value);
extern int cefingo_list_value_set_int(struct _cef_list_value_t* self,
                                      size_t index, int value);
extern int cefingo_list_value_set_double(struct _cef_list_value_t* self,
        size_t index, double value);
extern int cefingo_list_value_set_string(struct _cef_list_value_t* self,
        size_t index, cef_string_t* value);
extern int cefingo_list_value_set_binary(struct _cef_list_value_t* self,
        size_t index, struct _cef_binary_value_t* value);
extern int cefingo_list_value_set_dictionary(struct _cef_list_value_t* self,
        size_t index, struct _cef_dictionary_value_t* value);
extern int cefingo_list_value_set_list(struct _cef_list_value_t* self,
                                       size_t index, struct _cef_list_value_t* value);


#endif // CEFINGO_VALUES_H_