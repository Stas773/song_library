package logger

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	Logger := logrus.New()
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	return Logger
}
