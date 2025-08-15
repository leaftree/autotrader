package logger

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	flags  = log.LstdFlags // | log.Llongfile
	stdout = log.New(os.Stdout, "", flags)
)

func TestNewLoggerW(t *testing.T) {
	type args struct {
		component string
	}
	tests := []struct {
		name         string
		args         args
		want         Logger
		wantDebugW   string
		wantInfoW    string
		wantWarningW string
		wantErrorW   string
	}{
		// TODO: Add test cases.
		//	{name: "new logger test", args: args{component: "component"}, want: NewLogger("component", stdout, stdout, stdout, stdout), wantDebugW: "wantDebugw", wantInfoW: "info debug", wantWarningW: "warning", wantErrorW: "error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugW := &bytes.Buffer{}
			infoW := &bytes.Buffer{}
			warningW := &bytes.Buffer{}
			errorW := &bytes.Buffer{}
			if got := NewLoggerW(tt.args.component, debugW, infoW, warningW, errorW); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
			if gotDebugW := debugW.String(); gotDebugW != tt.wantDebugW {
				t.Errorf("NewLogger() = %v, want %v", gotDebugW, tt.wantDebugW)
			}
			if gotInfoW := infoW.String(); gotInfoW != tt.wantInfoW {
				t.Errorf("NewLogger() = %v, want %v", gotInfoW, tt.wantInfoW)
			}
			if gotWarningW := warningW.String(); gotWarningW != tt.wantWarningW {
				t.Errorf("NewLogger() = %v, want %v", gotWarningW, tt.wantWarningW)
			}
			if gotErrorW := errorW.String(); gotErrorW != tt.wantErrorW {
				t.Errorf("NewLogger() = %v, want %v", gotErrorW, tt.wantErrorW)
			}
		})
	}
}

func Test_loggerT_Debug(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		l    *loggerT
		args args
	}{
		// TODO: Add test cases.
		{
			name: "debug test",
			l:    &loggerT{component: "unit test component", m: []*log.Logger{stdout, stdout, stdout, stdout}},
			args: args{args: []any{"eth", "btc", "down trend"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debug(tt.args.args...)
		})
	}
}

func Test_loggerT_Info(t *testing.T) {
	var flags = log.LstdFlags | log.Lshortfile
	stdout := log.New(os.Stdout, "", flags)
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		l    *loggerT
		args args
	}{
		// TODO: Add test cases.
		{
			name: "info test",
			l:    &loggerT{component: "unit test component", m: []*log.Logger{stdout, stdout, stdout, stdout}},
			args: args{args: []any{"eth", "btc", "down trend"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Info(tt.args.args...)
		})
	}
}

func Test_loggerT_Warningln(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		l    *loggerT
		args args
	}{
		// TODO: Add test cases.
		{
			name: "warningln test",
			l:    &loggerT{component: "unit test component", m: []*log.Logger{stdout, stdout, stdout, stdout}},
			args: args{args: []any{"eth", "btc", "down trend"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Warningln(tt.args.args...)
		})
	}
}

func Test_loggerT_Errorln(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		l    *loggerT
		args args
	}{
		// TODO: Add test cases.
		{
			name: "errorln test",
			l:    &loggerT{component: "unit test component", m: []*log.Logger{stdout, stdout, stdout, stdout}},
			args: args{args: []any{"eth", "btc", "down trend"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Errorln(tt.args.args...)
		})
	}
}

func Test_loggerT_Debugf(t *testing.T) {
	type args struct {
		format string
		args   []any
	}
	tests := []struct {
		name string
		l    *loggerT
		args args
	}{
		// TODO: Add test cases.
		{
			name: "debugf test",
			l:    &loggerT{component: "unit test component", m: []*log.Logger{stdout, stdout, stdout, stdout}},
			args: args{format: "%s and %s is %s", args: []any{"eth", "btc", "down trend"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debugf(tt.args.format, tt.args.args...)
		})
	}
}

func Test_loggerT_Warningf(t *testing.T) {
	type args struct {
		format string
		args   []any
	}
	tests := []struct {
		name string
		l    *loggerT
		args args
	}{
		// TODO: Add test cases.
		{
			name: "warningf test",
			l:    &loggerT{component: "unit test component", m: []*log.Logger{stdout, stdout, stdout, stdout}},
			args: args{format: "test %s and %s, is %s", args: []any{"eth", "btc", "down trend"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Warningf(tt.args.format, tt.args.args...)
		})
	}
}
