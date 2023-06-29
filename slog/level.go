package slog

// Level specifies the log level
type Level int

const (
	// DebugLevel indicates Debug log level
	DebugLevel = -1
	// InfoLevel indicates Info log level.
	InfoLevel Level = iota
	// WarningLevel indicates Warning log level.
	WarningLevel
	// ErrorLevel indicates Error log level.
	ErrorLevel
	// FatalLevel indicates Fatal log level.
	FatalLevel
	// PanicLevel indicates Panic log level
	PanicLevel
	InvalidLevel
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warn"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case ErrorLevel:
		return "error"
	default:
		return "" // FIXME: Surely we should blast here
	}
}
