package capi

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"unsafe"
)

// #include "cefingo.h"
import "C"

var Logger *log.Logger

var ref_count_log struct {
	output   bool
	trace    bool
	traceSet map[unsafe.Pointer]bool
}

func init() {
	ref_count_log.traceSet = map[unsafe.Pointer]bool{}
}

func traceuf(up int, p unsafe.Pointer, message string, v ...interface{}) {
	if Logger != nil && ref_count_log.output {
		_, trace := ref_count_log.traceSet[p]
		if p == nil || trace {
			fn := caller_name(up)
			msg := fmt.Sprintf(message, v...)
			Logger.Printf("(%s) %s %p\n", fn, msg, p)
		}
	}
}

func traceOn(rc refCounted) {
	p := rc.cast_to_p_base_ref_counted_t()
	ref_count_log.traceSet[unsafe.Pointer(p)] = true
}

func Tracef(p unsafe.Pointer, message string, v ...interface{}) {
	traceuf(1, p, message, v...)
}

func Logf(message string, v ...interface{}) {
	if Logger != nil {
		fn := caller_name(0)
		Logger.Printf("("+fn+") "+message+"\n", v...)
	}
}

func Panicf(message string, v ...interface{}) {
	fn := caller_name(0)
	if Logger != nil {
		Logger.Panicf("("+fn+") "+message+"\n", v...)
	} else {
		log.Panicf("("+fn+") "+message+"\n", v...)
	}
}

func RefCountLogOutput(enable bool) {
	ref_count_log.output = enable
	if enable {
		C.REF_COUNT_LOG_OUTPUT = C.TRUE
	} else {
		C.REF_COUNT_LOG_OUTPUT = C.FALSE
	}
}

func RefCountLogTrace(on bool) {
	ref_count_log.trace = on
}

//export cefingo_cslog
func cefingo_cslog(p unsafe.Pointer, fn *C.char, s *C.char) {
	if Logger != nil {
		_, trace := ref_count_log.traceSet[p]
		if p == nil || trace {
			Logger.Println("(C."+C.GoString(fn)+")", strings.TrimRight(C.GoString(s), "\n"))
		}
	}
}

//export cefingo_panic
func cefingo_panic(fn *C.char, s *C.char) {
	if Logger != nil {
		Logger.Panicln("(C."+C.GoString(fn)+")", strings.TrimRight(C.GoString(s), "\n"))
	}
}

func caller_name(upper int) (fn string) {
	caller := []string{""}
	pt, _, _, ok := runtime.Caller(upper + 2)
	if ok {
		caller = strings.Split(runtime.FuncForPC(pt).Name(), "/")
	}
	fn = caller[len(caller)-1]
	if strings.Index(fn, "_cgo") >= 0 {
		fn = "C"
	}
	return fn
}
