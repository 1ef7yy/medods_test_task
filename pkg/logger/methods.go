package logger

import (
	"fmt"
	"os"
)

func (l logger) Debug(msg string) {
	l.log.Debug(msg)
}
func (l logger) Info(msg string) {
	l.log.Info(msg)
}
func (l logger) Warn(msg string) {
	l.log.Warn(msg)
}
func (l logger) Error(msg string) {
	l.log.Error(msg)
}
func (l logger) Fatal(msg string) {
	l.log.Error(msg)
	os.Exit(1)
}

func (l logger) Debugf(msg string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(msg, args...))
}
func (l logger) Infof(msg string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, args...))
}
func (l logger) Warnf(msg string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, args...))
}
func (l logger) Errorf(msg string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, args...))
}
func (l logger) Fatalf(msg string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}
