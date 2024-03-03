package vigor

// Log defines the logger to use in this library. This is by default a no-op
// logger to reduce binary size. This may be set to `log.Default()` from the Go
// stdlib to log out.
var Log Logger = &defaultLogger{}

// Logger abstracts the implementation of loggers,
type Logger interface {
	Print(v ...any)
	Fatal(v ...any)
}

type defaultLogger struct {
	Logger
}

// Fatal calls `panic` with the given args to stop execution.
func (l *defaultLogger) Fatal(v ...any) {
	panic(v)
}
