package cef

import (
	"runtime"

	"github.com/turutcrane/cefingo/capi"
)

type mStringMap struct {
	cef capi.CStringMapT
}

type StringMap struct {
	m *mStringMap
}

func NewStringMap() StringMap {
	cefObject := capi.StringMapAlloc()
	goObject := mStringMap{cefObject}
	runtime.SetFinalizer(&goObject, func(o *mStringMap) {
		capi.Logf("L20:")
		capi.StringMapFree(o.cef)
	})
	return StringMap{&goObject}
}

func (o StringMap) CefObject() capi.CStringMapT {
	return o.m.cef
}

type mStringMultimap struct {
	cef capi.CStringMultimapT
}

type StringMultimap struct {
	m *mStringMultimap
}

func NewStringMultimap() StringMultimap {
	cefObject := capi.StringMultimapAlloc()
	goObject := mStringMultimap{cefObject}
	runtime.SetFinalizer(&goObject, func(o *mStringMultimap) {
		capi.Logf("L43:")
		capi.StringMultimapFree(o.cef)
	})
	return StringMultimap{&goObject}
}

func (o StringMultimap) CefObject() capi.CStringMultimapT {
	return o.m.cef
}
