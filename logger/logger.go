package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	lvlDebug int = iota
	lvlInfo
	lvlWarning
	lvlError
)

var (
	severityName = []string{
		lvlDebug:   "DEBUG",
		lvlInfo:    "INFO",
		lvlWarning: "WARNING",
		lvlError:   "ERROR",
	}
	sprint   = fmt.Sprint
	sprintf  = fmt.Sprintf
	sprintln = fmt.Sprintln
)

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warning(args ...any)
	Error(args ...any)
	Debugln(args ...any)
	Infoln(args ...any)
	Warningln(args ...any)
	Errorln(args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warningf(format string, args ...any)
	Errorf(format string, args ...any)

	SetWriter(lvl int, w io.Writer) Logger
}

type loggerT struct {
	component string
	m         []*log.Logger
}

var _ Logger = (*loggerT)(nil)

func (l *loggerT) output(serverity int, s string) {
	sevStr := severityName[serverity]
	l.m[serverity].Output(4, fmt.Sprintf("[%s] [%s] %s", l.component, sevStr, s))
}

func (l *loggerT) print(serverity int, args ...any)   { l.output(serverity, sprint(args...)) }
func (l *loggerT) println(serverity int, args ...any) { l.output(serverity, sprintln(args...)) }
func (l *loggerT) printf(serverity int, format string, args ...any) {
	l.output(serverity, sprintf(format, args...))
}

func (l *loggerT) Debug(args ...any)                   { l.print(lvlDebug, args...) }
func (l *loggerT) Info(args ...any)                    { l.print(lvlInfo, args...) }
func (l *loggerT) Warning(args ...any)                 { l.print(lvlWarning, args...) }
func (l *loggerT) Error(args ...any)                   { l.print(lvlError, args...) }
func (l *loggerT) Debugln(args ...any)                 { l.println(lvlDebug, args...) }
func (l *loggerT) Infoln(args ...any)                  { l.println(lvlInfo, args...) }
func (l *loggerT) Warningln(args ...any)               { l.println(lvlWarning, args...) }
func (l *loggerT) Errorln(args ...any)                 { l.println(lvlError, args...) }
func (l *loggerT) Debugf(format string, args ...any)   { l.printf(lvlDebug, format, args...) }
func (l *loggerT) Infof(format string, args ...any)    { l.printf(lvlInfo, format, args...) }
func (l *loggerT) Warningf(format string, args ...any) { l.printf(lvlWarning, format, args...) }
func (l *loggerT) Errorf(format string, args ...any)   { l.printf(lvlError, format, args...) }

func (l *loggerT) SetWriter(lvl int, w io.Writer) Logger {
	l.m[lvl] = log.New(w, "", log.LstdFlags)
	return l
}

func NewLoggerW(component string, debugW, infoW, warningW, errorW io.Writer) Logger {
	prefix := ""
	flags := log.LstdFlags | log.Lshortfile

	m := []*log.Logger{
		log.New(debugW, prefix, flags),
		log.New(infoW, prefix, flags),
		log.New(warningW, prefix, flags),
		log.New(errorW, prefix, flags),
	}
	return &loggerT{component: component, m: m}
}

func NewLogger(component string) Logger {
	return NewLoggerW(component, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
}
