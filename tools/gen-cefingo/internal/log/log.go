package log
import(
	golog "log"
)

var traceon bool

func Trace(on bool) {
	traceon = on
}

func Tracef(message string, v ...interface{}) {
	if traceon {
		golog.Printf(message, v...)
	}
}

func Panicf(message string, v ...interface{}) {
	golog.Panicf(message, v...)
}

func Panicln(v ...interface{}) {
	golog.Panicln(v...)
}

func Fatalln(v ...interface{}) {
	golog.Fatalln(v...)
}