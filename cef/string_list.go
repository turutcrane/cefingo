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

func NewStringList() StringList {
	cefObject := capi.StringListAlloc()
	goObject := mStringList{cefObject}
	runtime.SetFinalizer(&goObject, func(o *mStringList) {
		capi.Logf("L20:")
		capi.StringListFree(o.cef)
	})
	return StringList{&goObject}
}

func (o StringList) CefObject() capi.CStringListT {
	return o.m.cef
}
