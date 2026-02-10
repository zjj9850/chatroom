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

type LogLevel int

const (
	logLevel_Info LogLevel = iota
	logLevel_Warning
	logLevel_Error
	logLevel_Debug
	logLevel_Critical
	logLevel_Fatal
	logLevel_Panic
)

type logdata struct {
	data  string
	level LogLevel
}

var logChan chan logdata

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

func Info(v ...any) {
	logFmt, _ := LOG_COLOR["INFO"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Info, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Infof(format string, v ...any) {
	logFmt, _ := LOG_COLOR["INFO"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	logChan <- logdata{level: logLevel_Info, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Debug(v ...any) {
	logFmt, _ := LOG_COLOR["DEBUG"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Debug, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Debugf(format string, v ...any) {
	logFmt, _ := LOG_COLOR["DEBUG"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	logChan <- logdata{level: logLevel_Debug, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Warning(v ...any) {
	logFmt, _ := LOG_COLOR["WARNING"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Warning, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Warningf(format string, v ...any) {
	logFmt, _ := LOG_COLOR["WARNING"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	logChan <- logdata{level: logLevel_Warning, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Error(v ...any) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Error, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Errorf(format string, v ...any) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	logChan <- logdata{level: logLevel_Error, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Critical(v ...any) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Critical, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Criticalf(format string, v ...any) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	logChan <- logdata{level: logLevel_Critical, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Fatalf(format string, v ...any) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(format, v...)
	logChan <- logdata{level: logLevel_Fatal, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Fatal(v ...any) {
	logFmt, _ := LOG_COLOR["CRITICAL"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Fatal, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func Panicf(format string, v ...any) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(logFmt, get_date(), file, line, fmt.Sprintf(format, v...))
	logChan <- logdata{level: logLevel_Panic, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}

}

func Panic(v ...any) {
	logFmt, _ := LOG_COLOR["ERROR"]
	file, line := get_caller(2)
	msg := fmt.Sprintf(get_args_format(len(v)), v...)
	logChan <- logdata{level: logLevel_Panic, data: fmt.Sprintf(logFmt, get_date(), file, line, msg)}
}

func init() {
	logChan = make(chan logdata, 1024)
	go func() {
		for msgData := range logChan {
			fmt.Println(msgData.data)
			switch msgData.level {
			case logLevel_Fatal:
				time.Sleep(10 * time.Millisecond)
				os.Exit(1)
			case logLevel_Panic:
				time.Sleep(10 * time.Millisecond)
				panic(msgData.data)
			}
		}
	}()
}
