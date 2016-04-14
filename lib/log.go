// Package lib provides ...
package lib

import (
	"os"
	"fmt"
	"time"
	"runtime"
	"strings"
	"path/filepath"
	logging "log"
)

const (
	LOG_DEBUG = 1
	LOG_INFO = 2
	LOG_WARN = 3
	LOG_ERROR = 4
	LOG_FATAL = 5
)

var log *Logger
var level map[int]string = map[int]string{
	LOG_DEBUG: "DEBUG",
	LOG_INFO:  "INFO ",
	LOG_WARN:  "WARN ",
	LOG_ERROR: "ERROR",
	LOG_FATAL: "FATAL",
}
var logLevel int = 1
var logEnable bool = true
var logOut = os.Stdout

type Logger struct {
	Log *logging.Logger
	IsEnabled bool
	Level int
	Path string
	BasePath string
}

func NewLogger(p string) *Logger {
	var err error
	if log != nil {
		return log
	}
	if CliParams.LogPath != "" {
		cpath := CliParams.LogPath
		d := time.Now().Format("2006-01-02")
		var absPath string
		if strings.HasPrefix(cpath, "/") {
			absPath, err = filepath.Abs(cpath)
		} else {
			realPath := filepath.Join(filepath.Base("."),  cpath)
			absPath, err = filepath.Abs(filepath.Join(realPath))
		}
		if err != nil {
			log.Fatalf("Path to log wrong: %v\n", err)
		}
		path, err := filepath.Abs(filepath.Join(absPath, d + "_geep-server.log"))
		if err == nil {
			os.Remove(path)
		}
		f, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("Unable open log file: %v\n", err)
		}
		log = &Logger{ Log: logging.New(f, "", logging.Ltime) }
		// log.Log.SetOutput(f) // Not working on Linux (but working on OS X!)
	} else {
		buf := *os.Stdout
		log = &Logger{ Log: logging.New(&buf, "", logging.Ltime) }
	}

	log.BasePath, err = filepath.Abs(".")
	if err != nil {
		fmt.Printf("Unable open log file: %v\n", err)
	}
	log.BasePath = strings.ToLower(log.BasePath) + "/"
	log.IsEnabled = CliParams.LogEnable
	log.Level = CliParams.LogLevel
	return log
}

func getLogPath(s string) *os.File {
	d := time.Now().Format("2006-01-02")
	path := filepath.Join(s, d + "_geep-server.log")
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

// Debug
func (l *Logger) Debug(s string) {
	l.print(LOG_DEBUG, s)
}
func (l *Logger) Debugf(s string, args ...interface{}) {
	l.printf(LOG_DEBUG, s, args...)
}

// Info
func (l *Logger) Info(s string) {
	l.print(LOG_INFO, s)
}
func (l *Logger) Infof(s string, args ...interface{}) {
	l.printf(LOG_INFO, s, args...)
}

// Warn
func (l *Logger) Warn(s string) {
	l.printf(LOG_WARN, s)
}
func (l *Logger) Warnf(s string, args ...interface{}) {
	l.printf(LOG_WARN, s, args...)
}

// Error
func (l *Logger) Error(s string) {
	l.printf(LOG_ERROR, s)
}
func (l *Logger) Errorf(s string, args ...interface{}) {
	l.printf(LOG_ERROR, s, args...)
}

// Fatal
func (l *Logger) Fatal(s string) {
	l.printf(LOG_FATAL, s)
	os.Exit(1)
}
func (l *Logger) Fatalf(s string, args ...interface{}) {
	l.printf(LOG_FATAL, s, args...)
	os.Exit(1)
}

// Internal
func (l *Logger) print(t int, s string) {
	if l.IsEnabled && t >= l.Level {
		l._print(level[t], s)
	}
}
func (l *Logger) printf(t int, s string, a ...interface{}) {
	if l.IsEnabled && t >= l.Level {
		msg := fmt.Sprintf(s, a...)
		l._print(level[t], msg)
	}
}

func (l *Logger) _print(lvl, msg string) {
	_, file, line, _ := runtime.Caller(3)
	bfile := strings.TrimPrefix(strings.ToLower(file), l.BasePath)
	l.Log.Printf("| %s | %s:%d | %s", lvl, bfile, line, msg)
}
