package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	LogTypeInfo  = "info"
	LogTypeWarn  = "warn"
	LogTypeError = "error"
	LogTypeFatal = "fatal"
)

// Logger interface defines all the logging methods to be implemented
type Logger interface {
	Info(payload ...interface{}) error
	Warn(payload ...interface{}) error
	Error(payload ...interface{}) error
	Fatal(payload ...interface{}) error
}

// LogHandler implements Logger
type LogHandler struct {
	Skipstack  int
	appName    string
	appVersion string
}

func (lh *LogHandler) defaultPayload(severity string) map[string]interface{} {
	_, file, line, _ := runtime.Caller(lh.Skipstack)
	return map[string]interface{}{
		"app":        lh.appName,
		"appVersion": lh.appVersion,
		"severity":   severity,
		"line":       fmt.Sprintf("%s:%d", file, line),
		"at":         time.Now(),
	}
}

func (lh *LogHandler) serialize(severity string, data ...interface{}) (string, error) {
	payload := lh.defaultPayload(severity)
	for idx, value := range data {
		payload[fmt.Sprintf("%d", idx)] = value
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (lh *LogHandler) log(severity string, payload ...interface{}) error {
	out, err := lh.serialize(severity, payload...)
	if err != nil {
		return err
	}

	switch severity {
	case "fatal":
		{
			fmt.Println(out)
			os.Exit(1)
		}
	}
	fmt.Println(out)

	return nil
}

func (lh *LogHandler) Info(payload ...interface{}) error {
	return lh.log(LogTypeInfo, payload...)
}

func (lh *LogHandler) Warn(payload ...interface{}) error {
	return lh.log(LogTypeWarn, payload...)
}

func (lh *LogHandler) Error(payload ...interface{}) error {
	return lh.log(LogTypeError, payload...)
}

func (lh *LogHandler) Fatal(payload ...interface{}) error {
	return lh.log(LogTypeFatal, payload...)
}

func New(appname string, appversion string, skipStack uint) *LogHandler {
	if skipStack <= 1 {
		skipStack = 4
	}

	return &LogHandler{
		Skipstack:  int(skipStack),
		appName:    appname,
		appVersion: appversion,
	}
}
