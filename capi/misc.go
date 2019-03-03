package cefingo

// #include "cefingo.h"
import "C"

func (app *CAppT) HasOneRef() bool {
	return BaseHasOneRef(app.p_app)
}

func (v *CV8valueT) HasOneRef() bool {
	return BaseHasOneRef(v.p_v8value)
}

func (h *CV8handlerT) HasOneRef() bool {
	return BaseHasOneRef(h.p_v8handler)
}
