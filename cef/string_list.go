package cef

import (
	"runtime"

	"github.com/turutcrane/cefingo/capi"
)

type mStringList struct {
	cef capi.CStringListT
}

type StringList struct {
	m *mStringList
}

func makeStringList(cefObject capi.CStringListT) StringList {
	goObject := mStringList{cefObject}
	runtime.SetFinalizer(&goObject, func(o *mStringList) {
		capi.Logf("L20:")
		capi.StringListFree(o.cef)
	})
	return StringList{&goObject}

}

func NewStringList() StringList {
	cefObject := capi.StringListAlloc()
	return makeStringList(cefObject)
}

func (o StringList) CefObject() capi.CStringListT {
	return o.m.cef
}

func (o StringList) Size() (ret int64) {
	return capi.StringListSize(o.m.cef)
}

func (o StringList) Value(index int64) (value string, ok bool) {
	ok, value = capi.StringListValue(o.m.cef, index)
	return value, ok
}

func (o StringList) Append(value string) {
	capi.StringListAppend(o.m.cef, value)
}

func (o StringList) Clear() {
	capi.StringListClear(o.m.cef)
}

func (o StringList) Copy() StringList {
	cefObject := capi.StringListCopy(o.m.cef)
	return makeStringList(cefObject)
}
