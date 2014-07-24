package simplelog

import (
	"fmt"
	"testing"
)

type mockConsoleLogger struct {
	test *testing.T
	prefix string
	msg string
}

func newMockConsoleLogger(test *testing.T, prefix string) *mockConsoleLogger {
	return &mockConsoleLogger{test, prefix, ""}
}

func (mock *mockConsoleLogger) CheckMsg(msg string) {
	if mock.msg != msg {
		mock.test.Error("message was incorrect: \"%s\" != \"%s\"", mock.msg, msg)
	}
	mock.msg = ""
}

func (mock *mockConsoleLogger) Printf(format string, v ...interface{}) {
	mock.msg = fmt.Sprintf(format, v...)
}

func (mock *mockConsoleLogger) Prefix() string {
	return mock.prefix
}

type mockSyslogLogger struct {
	test *testing.T
	msg string
	level int
	closed bool
}

func newMockSyslogLogger(test *testing.T, prefix string) *mockSyslogLogger {
	return &mockSyslogLogger{test, "", -1, false}
}

func (mock *mockSyslogLogger) CheckMsg(level int, msg string) {
	if mock.msg != msg {
		mock.test.Error("message was incorrect: \"%s\" != \"%s\"", mock.msg, msg)
	}
	if mock.level != level {
		mock.test.Error("level was incorrect: %d != %d", mock.level, level)
	}
	mock.msg = ""
	mock.level = -1
}

func (mock *mockSyslogLogger) CheckClosed() {
	if !mock.closed {
		mock.test.Error("close was not called") 
	}
	mock.closed = false
}


func (mock *mockSyslogLogger) Debug(msg string) error {
	mock.level = DEBUG
	mock.msg = msg
	return nil
}

func (mock *mockSyslogLogger) Notice(msg string) error {
	mock.level = NOTICE
	mock.msg = msg
	return nil
}

func (mock *mockSyslogLogger) Info(msg string) error {
	mock.level = INFO
	mock.msg = msg
	return nil
}

func (mock *mockSyslogLogger) Warning(msg string) error {
	mock.level = WARN
	mock.msg = msg
	return nil
}

func (mock *mockSyslogLogger) Err(msg string) error {
	mock.level = ERROR
	mock.msg = msg
	return nil
}

func (mock *mockSyslogLogger) Crit(msg string) error {
	mock.level = FATAL
	mock.msg = msg
	return nil
}

func (mock *mockSyslogLogger) Close() error {
	mock.closed = true
	return nil
}

func TestNewEmptyLogger(t *testing.T) {
	if logger, err := NewLogger(0, "empty"); err == nil {
		if logger.Console() {
			t.Error("empty logger has console logging enabled!")
		}
		if logger.Syslog() {
			t.Error("empty logger has a syslog logging enabled!")
		}
		if logger.level != NOTICE {
			t.Error("empty logger level not set to NOTICE!")
		}
	} else {
		t.Errorf("empty logger creation failed: %s", err)
	}
}

func TestNewConsoleLogger(t *testing.T) {
	prefix := "console"
	if logger, err := NewLogger(CONSOLE, prefix); err == nil {
		if !logger.Console() {
			t.Error("console logger has console logging disabled!")
		} else if logger.console.Prefix() != prefix+" " {
			t.Error("console logger has an incorrect prefix!")
		}
		if logger.Syslog() {
			t.Error("console logger has syslog logging enabled!")
		}
		if logger.level != NOTICE {
			t.Error("console logger level not set to NOTICE!")
		}
		logger.Close()
	} else {
		t.Errorf("console logger creation failed: %s", err)
	}
}

func TestNewSyslogLogger(t *testing.T) {
	if logger, err := NewLogger(SYSLOG, "syslog"); err == nil {
		if logger.Console() {
			t.Error("syslog logger has console logging enabled!")
		}
		if !logger.Syslog() {
			t.Error("syslog logger has syslog logging disabled!")
		}
		if logger.level != NOTICE {
			t.Error("syslog logger level not set to NOTICE!")
		}
		logger.Close()
	} else {
		t.Errorf("syslog logger creation failed: %s", err)
	}
}

func TestNewAllLogger(t *testing.T) {
	if logger, err := NewLogger(CONSOLE|SYSLOG, "all"); err == nil {
		if !logger.Console() {
			t.Error("all logger has console logging disabled!")
		}
		if !logger.Syslog() {
			t.Error("all logger has syslog logging disabled!")
		}
		logger.Close()
	} else {
		t.Error("syslog logger creation failed:", err)
	}
}

func TestLevel(t *testing.T) {
	msg := "hello world"
	mock := newMockSyslogLogger(t, "syslog")
	logger, err := NewLogger(SYSLOG, "syslog")
	if err != nil {
		t.Errorf("syslog logger creation failed: %s", err)
	}
	logger.syslog = mock

	logger.Debug(msg)
	mock.CheckMsg(-1, "")
	logger.Info(msg)
	mock.CheckMsg(-1, "")
	logger.Notice(msg)
	mock.CheckMsg(NOTICE, msg)
	logger.Warn(msg)
	mock.CheckMsg(WARN, msg)
	logger.Error(msg)
	mock.CheckMsg(ERROR, msg)

	logger.SetLevel(ERROR)
	if logger.level != ERROR {
		t.Errorf("level not set to ERROR")
	}

	logger.Debug(msg)
	mock.CheckMsg(-1, "")
	logger.Info(msg)
	mock.CheckMsg(-1, "")
	logger.Notice(msg)
	mock.CheckMsg(-1, "")
	logger.Warn(msg)
	mock.CheckMsg(-1, "")
	logger.Error(msg)
	mock.CheckMsg(ERROR, msg)
}

func TestConsole(t *testing.T) {
	var fullMsg string
	prefix := "console"
	format := "test: %s"
	value := "some value"
	msg := "test: some value\n"
	mock := newMockConsoleLogger(t, prefix)
	logger, err := NewLogger(CONSOLE, prefix)
	logger.SetLevel(DEBUG)
	if err != nil {
		t.Errorf("console logger creation failed: %s", err)
	}
	logger.console = mock

	// test debug
	fullMsg = fmt.Sprintf("[DEBUG]  %s", msg)
	logger.Log(DEBUG, format, value)
	mock.CheckMsg(fullMsg)
	logger.Debug(format, value)
	mock.CheckMsg(fullMsg)

	// test notice
	fullMsg = fmt.Sprintf("[NOTICE] %s", msg)
	logger.Log(NOTICE, format, value)
	mock.CheckMsg(fullMsg)
	logger.Notice(format, value)
	mock.CheckMsg(fullMsg)

	// test info
	fullMsg = fmt.Sprintf("[INFO]   %s", msg)
	logger.Log(INFO, format, value)
	mock.CheckMsg(fullMsg)
	logger.Info(format, value)
	mock.CheckMsg(fullMsg)

	// test warn
	fullMsg = fmt.Sprintf("[WARN]   %s", msg)
	logger.Log(WARN, format, value)
	mock.CheckMsg(fullMsg)
	logger.Warn(format, value)
	mock.CheckMsg(fullMsg)

	// test error
	fullMsg = fmt.Sprintf("[ERROR]  %s", msg)
	logger.Log(ERROR, format, value)
	mock.CheckMsg(fullMsg)
	logger.Error(format, value)
	mock.CheckMsg(fullMsg)

	// test close
	if err := logger.Close(); err != nil {
		t.Errorf("close returned an error: %s", err)
	}
}

func TestSyslog(t *testing.T) {
	prefix := "syslog"
	format := "test: %s"
	value := "some value"
	msg := "test: some value"
	mock := newMockSyslogLogger(t, prefix)
	logger, err := NewLogger(SYSLOG, prefix)
	logger.SetLevel(DEBUG)
	if err != nil {
		t.Errorf("syslog logger creation failed: %s", err)
	}
	logger.syslog = mock

	// test debug
	logger.Log(DEBUG, format, value)
	mock.CheckMsg(DEBUG, msg)
	logger.Debug(format, value)
	mock.CheckMsg(DEBUG, msg)

	// test notice
	logger.Log(NOTICE, format, value)
	mock.CheckMsg(NOTICE, msg)
	logger.Notice(format, value)
	mock.CheckMsg(NOTICE, msg)

	// test info
	logger.Log(INFO, format, value)
	mock.CheckMsg(INFO, msg)
	logger.Info(format, value)
	mock.CheckMsg(INFO, msg)

	// test warn
	logger.Log(WARN, format, value)
	mock.CheckMsg(WARN, msg)
	logger.Warn(format, value)
	mock.CheckMsg(WARN, msg)

	// test error
	logger.Log(ERROR, format, value)
	mock.CheckMsg(ERROR, msg)
	logger.Error(format, value)
	mock.CheckMsg(ERROR, msg)

	// test close
	if err := logger.Close(); err == nil {
		mock.CheckClosed()
	} else {
		t.Errorf("close returned an error: %s", err)
	}
}
