package cefingo

import (
	"C"
	"log"
	"runtime"
	"strings"
)

// enable Log Output
var LogOutput bool = false

func Logf(message string, v ...interface{}) {
	if LogOutput {
		fn := caller_name()
		log.Printf("("+fn+") "+message+"\n", v...)
	}
}

//export cefingo_cslog
func cefingo_cslog(fn *C.char, s *C.char) {
	if LogOutput {
		log.Println("(C."+C.GoString(fn)+")", strings.TrimRight(C.GoString(s), "\n"))
	}
}

//export cefingo_panic
func cefingo_panic(s *C.char) {
	fn := caller_name()
	log.Panicln("("+fn+")", strings.TrimRight(C.GoString(s), "\n"))
}

func caller_name() (fn string) {
	caller := []string{""}
	pt, _, _, ok := runtime.Caller(2)
	if ok {
		caller = strings.Split(runtime.FuncForPC(pt).Name(), "/")
	}
	fn = caller[len(caller)-1]
	if strings.Index(fn, "_cgo") >= 0 {
		fn = "C"
	}
	return
}
