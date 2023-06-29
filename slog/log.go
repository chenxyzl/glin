package slog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

// processName process name
// pid process id
// rollingFileWriter rolling file writer
var processName string
var pid = os.Getpid()
var rollingFileWriter io.Writer

// consoleLogger represents output to stdout
// rollingFileLogger represents output rolling file
// DefaultLogger default use
var consoleLogger *zap.SugaredLogger
var rollingFileLogger *zap.SugaredLogger
var DefaultLogger *zap.SugaredLogger

func init() {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executableName := filepath.Base(executablePath)
	processName = executableName[:len(executableName)-len(filepath.Ext(executableName))]

	//rollingFileWriter = &lumberjack.Logger{Filename: fmt.Sprintf("%v_%v_%v.log", processName, time.Now().Format("20060102150405"), pid), MaxSize: 100, MaxAge: 30, MaxBackups: 100, LocalTime: true, Compress: false}
	rollingFileWriter = &lumberjack.Logger{Filename: fmt.Sprintf("%v.log", processName), MaxSize: 100, MaxAge: 7, MaxBackups: 10, LocalTime: true, Compress: false}
	consoleLogger = New(DebugLevel).Sugar()
	rollingFileLogger = New(InfoLevel, rollingFileWriter).Sugar()
	DefaultLogger = rollingFileLogger
}

// Debug logs to INFO level.
func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

// Debugf logs to INFO level
func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

// Info logs to INFO level.
func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

// Infof logs to INFO level
func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Warning logs to the WARNING level.
func Warning(v ...interface{}) {
	DefaultLogger.Warn(v...)
}

// Warningf logs to the WARNING level.
func Warningf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

// Error logs to the ERROR level.
func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

// Errorf logs to the ERROR level.
func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Fatal logs to the FATAL level followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

// Fatalf logs to the FATAL level followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

// Panic logs to the PANIC level followed by a call to panic().
func Panic(v ...interface{}) {
	DefaultLogger.Panic(v...)
}

// Panicf logs to the PANIC level followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	DefaultLogger.Panicf(format, v...)
}

// Sync sync log
func Sync() error {
	return DefaultLogger.Sync()
}

// New creates an instance of Log
func New(level Level, writers ...io.Writer) *Log {
	// create the zap Log configuration
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.FunctionKey = "F"
	// create the zap log core
	var core zapcore.Core
	//sync
	//sync := zapcore.AddSync(writer)
	var writeSyncers = []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)}
	for _, wt := range writers {
		writeSyncers = append(writeSyncers, zapcore.AddSync(wt))
	}
	sync := zapcore.NewMultiWriteSyncer(writeSyncers...)
	//config := zapcore.NewJSONEncoder(cfg.EncoderConfig),
	config := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	// set the log level
	switch level {
	case DebugLevel:
		core = zapcore.NewCore(config, sync, zapcore.DebugLevel)
	case InfoLevel:
		core = zapcore.NewCore(config, sync, zapcore.InfoLevel)
	case WarningLevel:
		core = zapcore.NewCore(config, sync, zapcore.WarnLevel)
	case ErrorLevel:
		core = zapcore.NewCore(config, sync, zapcore.ErrorLevel)
	case PanicLevel:
		core = zapcore.NewCore(config, sync, zapcore.PanicLevel)
	case FatalLevel:
		core = zapcore.NewCore(config, sync, zapcore.FatalLevel)
	default:
		core = zapcore.NewCore(config, sync, zapcore.DebugLevel)
	}
	// get the zap Log
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	// create the instance of Log and returns it
	return &Log{zapLogger}
}

func NewWith(args ...interface{}) Logger {
	return DefaultLogger.Desugar().WithOptions(zap.AddCallerSkip(-1)).Sugar().With(args...)
	//return DefaultLogger.Desugar().WithOptions().Sugar().With(args...)
}

// Log implements Logger interface with the underlying zap as
// the underlying logging library
type Log struct {
	*zap.Logger
}

// Debug starts a message with debug level
func (l *Log) Debug(v ...any) {
	l.Logger.Debug(fmt.Sprint(v...))
}

// Debugf starts a message with debug level
func (l *Log) Debugf(format string, v ...any) {
	l.Logger.Debug(fmt.Sprintf(format, v...))
}

// Panic starts a new message with panic level. The panic() function
// is called which stops the ordinary flow of a goroutine.
func (l *Log) Panic(v ...any) {
	l.Logger.Panic(fmt.Sprint(v...))
}

// Panicf starts a new message with panic level. The panic() function
// is called which stops the ordinary flow of a goroutine.
func (l *Log) Panicf(format string, v ...any) {
	l.Logger.Panic(fmt.Sprintf(format, v...))
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called which terminates the program immediately.
func (l *Log) Fatal(v ...any) {
	l.Logger.Fatal(fmt.Sprint(v...))
}

// Fatalf starts a new message with fatal level. The os.Exit(1) function
// is called which terminates the program immediately.
func (l *Log) Fatalf(format string, v ...any) {
	l.Logger.Fatal(fmt.Sprintf(format, v...))
}

// Error starts a new message with error level.
func (l *Log) Error(v ...any) {
	l.Logger.Error(fmt.Sprint(v...))
}

// Errorf starts a new message with error level.
func (l *Log) Errorf(format string, v ...any) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// Warn starts a new message with warn level
func (l *Log) Warn(v ...any) {
	l.Logger.Warn(fmt.Sprint(v...))
}

// Warnf starts a new message with warn level
func (l *Log) Warnf(format string, v ...any) {
	l.Logger.Warn(fmt.Sprintf(format, v...))
}

// Info starts a message with info level
func (l *Log) Info(v ...any) {
	l.Logger.Info(fmt.Sprint(v...))
}

// Infof starts a message with info level
func (l *Log) Infof(format string, v ...any) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l *Log) Sync() {
	l.Logger.Sync()
}

// LogLevel returns the log level that is used
func (l *Log) LogLevel() Level {
	switch l.Level() {
	case zapcore.FatalLevel:
		return FatalLevel
	case zapcore.PanicLevel:
		return PanicLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.WarnLevel:
		return WarningLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.DebugLevel:
		return DebugLevel
	default:
		return InvalidLevel
	}
}
