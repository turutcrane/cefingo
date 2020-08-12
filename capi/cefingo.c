#include "cefingo.h"
#include "_cgo_export.h"
#include <assert.h>

void cefingo_init()
{
	assert(sizeof(long) <= 8);
	assert(sizeof(long long) == 8);
}

void cefingo_post_data_get_elements(
	struct _cef_post_data_t* self,
	size_t* elementsCount,
	struct _cef_post_data_element_t** elements
)
{
	self->get_elements(
		self,
		elementsCount,
		elements
	);
}
