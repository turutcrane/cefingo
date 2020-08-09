package cef

import (
	"os"
	"unsafe"

	"github.com/turutcrane/cefingo/capi"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) UNlock() {}

func ExecuteProcess(mainArgs *capi.CMainArgsT, app *capi.CAppT) {

	var sandBoxInfo unsafe.Pointer

	code := capi.ExecuteProcess(mainArgs, app, sandBoxInfo)
	if code >= 0 {
		os.Exit(int(code))
	}
}

func Initialize(mainArgs *capi.CMainArgsT, settings *capi.CSettingsT, app *capi.CAppT) {

	capi.Initialize(mainArgs, settings, app, nil)
}
