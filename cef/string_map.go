package cef

import (
	"runtime"

	"github.com/turutcrane/cefingo/capi"
)

type StringMap struct {
	noCopy noCopy
	cef    capi.CStringMapT
}

func NewStringMap() *StringMap {
	cefObject := capi.StringMapAlloc()
	goObject := &StringMap{noCopy{}, cefObject}
	runtime.SetFinalizer(goObject, func(o *StringMap) {
		capi.Logf("L20:")
		capi.StringMapFree(o.cef)
	})
	return goObject
}

func (o *StringMap) CefObject() capi.CStringMapT {
	return o.cef
}

type StringMultimap struct {
	noCopy noCopy
	cef    capi.CStringMultimapT
}

func NewStringMultimap() *StringMultimap {
	cefObject := capi.StringMultimapAlloc()
	goObject := &StringMultimap{noCopy{}, cefObject}
	runtime.SetFinalizer(goObject, func(o *StringMultimap) {
		capi.Logf("L43:")
		capi.StringMultimapFree(o.cef)
	})
	return goObject
}

func (o *StringMultimap) CefObject() capi.CStringMultimapT {
	return o.cef
}

func (o *StringMultimap) Append(key, value string) bool {
	return capi.StringMultimapAppend(o.CefObject(), key, value)
}
