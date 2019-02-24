package cefingo

// #include "cefingo.h"
import "C"

func (app *CAppT) HasOneRef() bool {
	return BaseHasOneRef(app.p_app)
}
