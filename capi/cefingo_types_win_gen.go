// Code generated by "gen-cefingo.go" DO NOT EDIT.
// +build windows

package capi

// #include "cefingo.h"
import "C"

// cef_types_win.h, include/internal/cef_types_win.h:56:57,

///
// Structure representing CefExecuteProcess arguments.
///
type CMainArgsT C.cef_main_args_t

func NewCMainArgsT() *CMainArgsT {
	s := &CMainArgsT{}
	return s
}

func (st *CMainArgsT) Instance() WinHinstance {
	return WinHinstance(st.instance)
}

func (st *CMainArgsT) SetInstance(v WinHinstance) {
	st.instance = (C.HINSTANCE)(v)
}

///
// Structure representing window information.
///
type CWindowInfoT C.cef_window_info_t

func NewCWindowInfoT() *CWindowInfoT {
	s := &CWindowInfoT{}
	return s
}

func (st *CWindowInfoT) ExStyle() WinDword {
	return WinDword(st.ex_style)
}

func (st *CWindowInfoT) SetExStyle(v WinDword) {
	st.ex_style = (C.DWORD)(v)
}

func (st *CWindowInfoT) WindowName() string {
	return string_from_cef_string(&st.window_name)
}

func (st *CWindowInfoT) SetWindowName(v string) {
	set_cef_string(&st.window_name, v)
}

func (st *CWindowInfoT) Style() WinDword {
	return WinDword(st.style)
}

func (st *CWindowInfoT) SetStyle(v WinDword) {
	st.style = (C.DWORD)(v)
}

func (st *CWindowInfoT) X() int {
	return int(st.x)
}

func (st *CWindowInfoT) SetX(v int) {
	st.x = (C.int)(v)
}

func (st *CWindowInfoT) Y() int {
	return int(st.y)
}

func (st *CWindowInfoT) SetY(v int) {
	st.y = (C.int)(v)
}

func (st *CWindowInfoT) Width() int {
	return int(st.width)
}

func (st *CWindowInfoT) SetWidth(v int) {
	st.width = (C.int)(v)
}

func (st *CWindowInfoT) Height() int {
	return int(st.height)
}

func (st *CWindowInfoT) SetHeight(v int) {
	st.height = (C.int)(v)
}

func (st *CWindowInfoT) ParentWindow() CWindowHandleT {
	return CWindowHandleT(st.parent_window)
}

func (st *CWindowInfoT) SetParentWindow(v CWindowHandleT) {
	st.parent_window = (C.HWND)(v)
}

func (st *CWindowInfoT) Menu() WinHmenu {
	return WinHmenu(st.menu)
}

func (st *CWindowInfoT) SetMenu(v WinHmenu) {
	st.menu = (C.HMENU)(v)
}

func (st *CWindowInfoT) WindowlessRenderingEnabled() bool {
	return st.windowless_rendering_enabled != 0
}

func (st *CWindowInfoT) SetWindowlessRenderingEnabled(v bool) {
	if v {
		st.windowless_rendering_enabled = 1
	} else {
		st.windowless_rendering_enabled = 0
	}
}

func (st *CWindowInfoT) SharedTextureEnabled() bool {
	return st.shared_texture_enabled != 0
}

func (st *CWindowInfoT) SetSharedTextureEnabled(v bool) {
	if v {
		st.shared_texture_enabled = 1
	} else {
		st.shared_texture_enabled = 0
	}
}

func (st *CWindowInfoT) ExternalBeginFrameEnabled() bool {
	return st.external_begin_frame_enabled != 0
}

func (st *CWindowInfoT) SetExternalBeginFrameEnabled(v bool) {
	if v {
		st.external_begin_frame_enabled = 1
	} else {
		st.external_begin_frame_enabled = 0
	}
}

func (st *CWindowInfoT) Window() CWindowHandleT {
	return CWindowHandleT(st.window)
}

func (st *CWindowInfoT) SetWindow(v CWindowHandleT) {
	st.window = (C.HWND)(v)
}
