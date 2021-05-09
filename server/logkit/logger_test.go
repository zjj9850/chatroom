package logkit_test

import (
	"chatroom/logkit"
	"testing"
)

func TestInfo(t *testing.T) {
	logkit.Info("Test Info")
}

func TestInfof(t *testing.T) {
	logkit.Infof("%d Test Info", 123)
}

func TestDebug(t *testing.T) {
	logkit.Debug("Test Debug")
}

func TestDebugf(t *testing.T) {
	logkit.Debugf("%d Test Debug", 123)
}

func TestWarning(t *testing.T) {
	logkit.Warning("Test Warnning")
}

func TestWarningf(t *testing.T) {
	logkit.Warningf("%d Test Warnning", 123)
}

func TestError(t *testing.T) {
	logkit.Error("Test Error")
}

func TestErrorf(t *testing.T) {
	logkit.Errorf("%d Test Error", 123)
}

func TestCritical(t *testing.T) {
	logkit.Critical("Test Critical")
}

func TestCriticalf(t *testing.T) {
	logkit.Criticalf("%d Test Critical", 123)
}

// func TestFatalf(t *testing.T) {
// 	logkit.Fatalf("%d Test Fatal", 123)
// }

// func TestFatal(t *testing.T) {
// 	logkit.Fatal("Test Fatal")
// }

// func TestPanicf(t *testing.T) {
// 	logkit.Panicf("%d Test Panic", 123)
// }

// func TestPanic(t *testing.T) {
// 	logkit.Panic("Test Panic")
// }
