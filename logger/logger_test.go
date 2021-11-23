package logger

import (
	"os"
	"os/exec"
	"testing"
)

func TestLevel_validate(t *testing.T) {
	tests := []struct {
		name    string
		l       Level
		wantErr bool
	}{
		{
			name:    "Valid - Debug",
			l:       LevelDebug,
			wantErr: false,
		},
		{
			name:    "Valid - Info",
			l:       LevelInfo,
			wantErr: false,
		},
		{
			name:    "Valid - Warn",
			l:       LevelWarn,
			wantErr: false,
		},
		{
			name:    "Valid - Error",
			l:       LevelError,
			wantErr: false,
		},
		{
			name:    "Valid - Fatal",
			l:       LevelFatal,
			wantErr: false,
		},
		{
			name:    "Invalid level",
			l:       Level(-1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Level.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name    string
		l       string
		want    Level
		wantErr bool
	}{
		{
			name:    "Parse debug",
			l:       "DEBUG",
			want:    LevelDebug,
			wantErr: false,
		},
		{
			name:    "Parse info",
			l:       "info",
			want:    LevelInfo,
			wantErr: false,
		},
		{
			name:    "Parse warn",
			l:       "warn",
			want:    LevelWarn,
			wantErr: false,
		},
		{
			name:    "Parse error",
			l:       "error",
			want:    LevelError,
			wantErr: false,
		},
		{
			name:    "Parse fatal",
			l:       "fatal",
			want:    LevelFatal,
			wantErr: false,
		},
		{
			name:    "Invalid level",
			l:       "Foo",
			want:    Level(-1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockWriter struct {
	calls [][]byte
}

func (w *mockWriter) Write(b []byte) (int, error) {
	w.calls = append(w.calls, b)
	return len(b), nil
}

func TestNewLevelLoggerValidatesLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   Level
		wantErr bool
	}{
		{
			name:    "Valid level",
			level:   LevelDebug,
			wantErr: false,
		},
		{
			name:    "Invalid level",
			level:   Level(-1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewLevelLogger(tt.level, nil); (err != nil) != tt.wantErr {
				t.Errorf("NewLevelLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewLevelLoggerDefaultsToStdout(t *testing.T) {
	l, err := NewLevelLogger(LevelDebug, nil)
	if err != nil {
		t.Fatalf("NewLevelLogger() unexpected error = %v", err)
	}

	if l.w != os.Stdout {
		t.Errorf("NewLevelLogger() did not use os.Stdout")
	}
}

func TestLevelLogger_Debug(t *testing.T) {
	type fields struct {
		level Level
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWrites int
	}{
		{
			name:       "Logs when LevelDebug - no newline",
			fields:     fields{level: LevelDebug},
			args:       args{format: "%s", args: []interface{}{"foo"}},
			wantWrites: 3,
		},
		{
			name:       "Logs when LevelDebug - newline",
			fields:     fields{level: LevelDebug},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Does not log when higher level",
			fields:     fields{level: LevelInfo},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockWriter{calls: make([][]byte, 0)}
			l := &LevelLogger{
				w:     w,
				level: tt.fields.level,
			}
			l.Debug(tt.args.format, tt.args.args...)

			if len(w.calls) != tt.wantWrites {
				t.Errorf("Expected Debug to write %d times, but got %d writes", tt.wantWrites, (w.calls))
			}
		})
	}
}

func TestLevelLogger_Info(t *testing.T) {
	type fields struct {
		level Level
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWrites int
	}{
		{
			name:       "Logs when LevelInfo - no newline",
			fields:     fields{level: LevelInfo},
			args:       args{format: "%s", args: []interface{}{"foo"}},
			wantWrites: 3,
		},
		{
			name:       "Logs when LevelInfo - newline",
			fields:     fields{level: LevelInfo},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Logs when lower level",
			fields:     fields{level: LevelDebug},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Does not log when higher level",
			fields:     fields{level: LevelFatal},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockWriter{calls: make([][]byte, 0)}
			l := &LevelLogger{
				w:     w,
				level: tt.fields.level,
			}
			l.Info(tt.args.format, tt.args.args...)

			if len(w.calls) != tt.wantWrites {
				t.Errorf("Expected Info to write %d times, but got %d writes", tt.wantWrites, (w.calls))
			}
		})
	}
}

func TestLevelLogger_Warn(t *testing.T) {
	type fields struct {
		level Level
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWrites int
	}{
		{
			name:       "Logs when LevelWarn - no newline",
			fields:     fields{level: LevelWarn},
			args:       args{format: "%s", args: []interface{}{"foo"}},
			wantWrites: 3,
		},
		{
			name:       "Logs when LevelWarn - newline",
			fields:     fields{level: LevelWarn},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Logs when lower level",
			fields:     fields{level: LevelDebug},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Does not log when higher level",
			fields:     fields{level: LevelFatal},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockWriter{calls: make([][]byte, 0)}
			l := &LevelLogger{
				w:     w,
				level: tt.fields.level,
			}
			l.Warn(tt.args.format, tt.args.args...)

			if len(w.calls) != tt.wantWrites {
				t.Errorf("Expected Warn to write %d times, but got %d writes", tt.wantWrites, (w.calls))
			}
		})
	}
}

func TestLevelLogger_Error(t *testing.T) {
	type fields struct {
		level Level
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWrites int
	}{
		{
			name:       "Logs when LevelError - no newline",
			fields:     fields{level: LevelError},
			args:       args{format: "%s", args: []interface{}{"foo"}},
			wantWrites: 3,
		},
		{
			name:       "Logs when LevelError - newline",
			fields:     fields{level: LevelError},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Logs when lower level",
			fields:     fields{level: LevelDebug},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 2,
		},
		{
			name:       "Does not log when higher level",
			fields:     fields{level: LevelFatal},
			args:       args{format: "%s\n", args: []interface{}{"foo"}},
			wantWrites: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockWriter{calls: make([][]byte, 0)}
			l := &LevelLogger{
				w:     w,
				level: tt.fields.level,
			}
			l.Error(tt.args.format, tt.args.args...)

			if len(w.calls) != tt.wantWrites {
				t.Errorf("Expected Error to write %d times, but got %d writes", tt.wantWrites, (w.calls))
			}
		})
	}
}

func TestLevelLogger_Fatal(t *testing.T) {
	if os.Getenv("TESTLEVELLOGGER_FATAL") == "1" {
		(&LevelLogger{
			w:     &mockWriter{calls: make([][]byte, 0)},
			level: LevelFatal,
		}).Fatal("")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLevelLogger_Fatal")
	cmd.Env = append(os.Environ(), "TESTLEVELLOGGER_FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
