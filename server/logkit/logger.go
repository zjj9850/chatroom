package logkit

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const PREFIX_COLOR = ""

var LOG_COLOR map[string]string = map[string]string{
	"INFO":     "\x1b[95;1m%s 0 :INFO %s:%d \x1b[1m\x1b[97;1m%s\x1b[1m",     // white
	"ERROR":    "\x1b[95;1m%s 0 :ERROR %s:%d \x1b[1m\x1b[91;1m%s\x1b[1m",    // red
	"TLOG":     "\x1b[95;1m%s 0 :TLOG none:0 \x1b[1m\x1b[92;1m%s\x1b[1m",    // green
	"WARNING":  "\x1b[95;1m%s 0 :WARNING %s:%d \x1b[1m\x1b[93;1m%s\x1b[1m",  // yellow
	"DEBUG":    "\x1b[95;1m%s 0 :DEBUG %s:%d \x1b[1m\x1b[94;1m%s\x1b[1m",    // blue
	"CRITICAL": "\x1b[95;1m%s 0 :CRITICAL %s:%d \x1b[1m\x1b[91;1m%s\x1b[1m", // bold-red
}

func get_args_format(n int) string {
	s := strings.Repeat("%v ", n-1)
	s += "%v"
	return s
}

func get_caller(calldepth int) (string, int) {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		return "none", 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file, line
}

func get_date() string {
	tm := time.Unix(0, time.Now().UnixNano())
	return tm.Format("2006-01-02 15:04:05.000")
}

func Info(v ...interface{}) {
	logFmt, _ := LOG_COLOR["INFO"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Infof(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["INFO"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Debug(v ...interface{}) {
	logFmt, _ := LOG_COLOR["DEBUG"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Debugf(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["DEBUG"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Warning(v ...interface{}) {
	logFmt, _ := LOG_COLOR["WARNING"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Warningf(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["WARNING"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Error(v ...interface{}) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Errorf(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Critical(v ...interface{}) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Criticalf(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
}

func Fatalf(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	fmt.Println(fmt.Sprintf(logFmt, get_date(), file, line, msg))
	os.Exit(1)
}

func Panicf(format string, v ...interface{}) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	s := fmt.Sprintf(logFmt, get_date(), file, line, msg)
	fmt.Println(s)
	panic(s)
}

func Panic(v ...interface{}) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	s := fmt.Sprintf(logFmt, get_date(), file, line, msg)
	fmt.Println(s)
	panic(s)
}
