package simplelog

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
)

const (
	SYSLOG = 1 << iota
	CONSOLE
)

const (
	DEBUG = iota
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
)

// Get a level for a string value.
func StringToLevel(level string) int {
	level = strings.ToUpper(strings.TrimSpace(level))
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	}
	return NOTICE
}

// Get a string for a level value.
func LevelToString(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "NOTICE"
}

type consoleLogger interface {
	Printf(string, ...interface{})
	Prefix() string
}

type syslogLogger interface {
	Debug(string) error
	Notice(string) error
	Info(string) error
	Warning(string) error
	Err(string) error
	Crit(string) error
	Close() error
}

// Simple logger instance.
type Logger struct {
	outputs int
	console consoleLogger
	syslog  syslogLogger
	level   int
}

// Create a new logger with the given outputs and log prefix. Level is set up INFO.
func NewLogger(outputs int, prefix string) (l *Logger, err error) {
	var outConsole *log.Logger
	var outSyslog *syslog.Writer
	if outputs&CONSOLE == CONSOLE {
		outConsole = log.New(os.Stdout, prefix+" ", log.LstdFlags|log.Lmicroseconds)
	}
	if outputs&SYSLOG == SYSLOG {
		if outSyslog, err = syslog.New(syslog.LOG_DAEMON|syslog.LOG_NOTICE, prefix); err != nil {
			return
		}
	}
	l = &Logger{outputs, outConsole, outSyslog, NOTICE}
	return
}

// Set the logging level.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// Check if console logging is enabled.
func (l *Logger) Console() bool {
	return l.outputs&CONSOLE == CONSOLE
}

// Check if syslog logging is enabled.
func (l *Logger) Syslog() bool {
	return l.outputs&SYSLOG == SYSLOG
}

// Log to the console.
func (l *Logger) logConsole(level int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.console.Printf("%-8s %s\n", fmt.Sprintf("[%s]", LevelToString(level)), msg)
}

// Log to syslog.
func (l *Logger) logSyslog(level int, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	switch level {
	case DEBUG:
		return l.syslog.Debug(msg)
	case INFO:
		return l.syslog.Info(msg)
	case WARN:
		return l.syslog.Warning(msg)
	case ERROR:
		return l.syslog.Err(msg)
	case FATAL:
		return l.syslog.Crit(msg)
	}
	return l.syslog.Notice(msg)
}

// Log to all configured outputs.
func (l *Logger) Log(level int, format string, args ...interface{}) (err error) {
	if level < l.level{
		return
	}
	if l.Console() {
		l.logConsole(level, format, args...)
	}
	if l.Syslog() {
		err = l.logSyslog(level, format, args...)
	}
	if level == FATAL {
		os.Exit(1)
	}
	return
}

// Log to DEBUG level.
func (l *Logger) Debug(format string, args ...interface{}) error {
	return l.Log(DEBUG, format, args...)
}

// Log to NOTICE level.
func (l *Logger) Notice(format string, args ...interface{}) error {
	return l.Log(NOTICE, format, args...)
}

// Log to INFO level.
func (l *Logger) Info(format string, args ...interface{}) error {
	return l.Log(INFO, format, args...)
}

// Log to WARN level.
func (l *Logger) Warn(format string, args ...interface{}) error {
	return l.Log(WARN, format, args...)
}

// Log to ERROR level.
func (l *Logger) Error(format string, args ...interface{}) error {
	return l.Log(ERROR, format, args...)
}

// Log to FATAL level.
func (l *Logger) Fatal(format string, args ...interface{}) error {
	return l.Log(FATAL, format, args...)
}

// Close the logger.
func (l *Logger) Close() (err error) {
	if l.Syslog() {
		err = l.syslog.Close()
	}
	return
}
