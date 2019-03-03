#include <stdio.h>
#include <stdarg.h>
#include "cefingo_base.h"
#include "_cgo_export.h"

int REF_COUNT_LOG_OUTPUT = 0;

#define MAXLOGBUF 1000
void cefingo_cslogf(cef_base_ref_counted_t* self, const char *fn, const char *format, ...)
{
    static char buf[MAXLOGBUF + 1];

    va_list ap;
    va_start(ap, format);
    int n = vsnprintf(buf, MAXLOGBUF, format, ap);
    va_end(ap);
    buf[MAXLOGBUF] = '\0';

    cefingo_cslog(self, (char *) fn, buf);
}

void cefingo_panicf(const char *fn, const char *format, ...)
{
    static char buf[MAXLOGBUF + 1];

    va_list ap;
    va_start(ap, format);
    int n = vsnprintf(buf, MAXLOGBUF, format, ap);
    va_end(ap);
    buf[MAXLOGBUF] = '\0';

    cefingo_panic((char *) fn, buf);
}

void cefingo_base_add_ref(cef_base_ref_counted_t *self)
{
    self->add_ref(self);
}

int cefingo_base_release(cef_base_ref_counted_t *self)
{
    int status = 0;

    status = self->release(self);

    return status;
}

int cefingo_base_has_one_ref(cef_base_ref_counted_t *self)
{
    return self->has_one_ref(self);
}

///
// Increment the reference count.
///
// void CEF_CALLBACK cefingo_add_ref(cef_base_ref_counted_t* self) {
void CEF_CALLBACK cefingo_add_ref(cef_base_ref_counted_t* self)
{
    cefingo_ref_counter *counter = (cefingo_ref_counter *)(((void *)self) + self->size);

    // counter->ref_count++;
    int64 count = __atomic_add_fetch(&counter->ref_count, 1, __ATOMIC_SEQ_CST);
    if (REF_COUNT_LOG_OUTPUT) cefingo_cslogf(self, __func__, "L64: 0x%llx +count: %d", self, count);
}

///
// Called to decrement the reference count for the object. If the reference
// count falls to 0 the object should self-delete. Returns true (1) if the
// resulting reference count is 0.
///
extern int CEF_CALLBACK cefingo_release(cef_base_ref_counted_t* self)
{
    cefingo_ref_counter *counter = (cefingo_ref_counter *)(((void *)self) + self->size);
    // counter->ref_count--;
    int64 count = __atomic_sub_fetch(&counter->ref_count, 1, __ATOMIC_SEQ_CST);

    if (count >= 0) {
        if (REF_COUNT_LOG_OUTPUT) cefingo_cslogf(self, __func__, "L79: 0x%llx -count: %d", self, count);
        if (count == 0) {
            cefingo_base_deassoc(self);
            // free(self);
        }
    } else {
        cefingo_panicf(__func__, "L85: 0x%llx -count:%d", self, count);
    }
    return (count == 0 ? 1 : 0);
}

///
// Returns true (1) if the current reference count is 1.
///
int CEF_CALLBACK cefingo_has_one_ref(cef_base_ref_counted_t* self)
{
    cefingo_ref_counter *counter = (cefingo_ref_counter *)(((void *)self) + self->size);
    int64 count = __atomic_load_n(&counter->ref_count, __ATOMIC_SEQ_CST);

    if (REF_COUNT_LOG_OUTPUT) cefingo_cslogf(self, __func__, "L98: 0x%llx has-one: %d", self, count);
    return (count == 1 ? 1 : 0);
}

void initialize_cefingo_base_ref_counted(size_t size, cef_base_ref_counted_t* base)
{
    if (REF_COUNT_LOG_OUTPUT) cefingo_cslogf(base, __func__, "L104: size: %d base: 0x%llx", size, base);
    base->size = size; // size_t size = base->size;

    if (size <= sizeof(cef_base_ref_counted_t)) {
        cefingo_cslogf(NULL, __func__, "FATAL: initialize_cef_base failed, size member not set");
        _exit(1);
    }
    base->add_ref = cefingo_add_ref;
    base->release = cefingo_release;
    base->has_one_ref = cefingo_has_one_ref;
}
