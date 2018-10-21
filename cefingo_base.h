
#ifndef CEFINGO_BASE_H_
#define CEFINGO_BASE_H_
#include <unistd.h>
#include "include/capi/cef_base_capi.h"
#include <stddef.h>
// size_t offsetof(type， member);

typedef struct _cefingo_ref_counter {
    int64 ref_count;
} cefingo_ref_counter;

#define CEFINGO_REF_COUNTER_WRAPPER(body_t, wapper_name) \
typedef struct { \
    body_t body; \
    cefingo_ref_counter counter;\
} wapper_name

extern int REF_COUNT_LOG_OUTPUT;

extern void initialize_cefingo_base_ref_counted(size_t size, cef_base_ref_counted_t* base);
extern void CEF_CALLBACK cefingo_add_ref(cef_base_ref_counted_t* self);
extern int CEF_CALLBACK cefingo_release(cef_base_ref_counted_t* self);
extern int CEF_CALLBACK cefingo_has_one_ref(cef_base_ref_counted_t* self);

extern void cefingo_base_add_ref(cef_base_ref_counted_t *self);
extern int cefingo_base_release(cef_base_ref_counted_t *self);
extern int cefingo_base_has_one_ref(cef_base_ref_counted_t *self);

#endif // CEFINGO_BASE_H_
