package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GooseLogger struct{}

func (g GooseLogger) Fatalf(format string, v ...interface{}) {
	StdFatal(format, v...)
}

func (g GooseLogger) Printf(format string, v ...interface{}) {
	StdInfo(format, v...)
}

var Default = New(os.Stderr, zapcore.InfoLevel, "production")
var Goose GooseLogger

var (
	Info    = Default.Info
	Infow   = Default.Infow
	Warn    = Default.Warn
	Warnw   = Default.Warnw
	Error   = Default.Error
	Errorw  = Default.Errorw
	DPanic  = Default.DPanic
	DPanicw = Default.DPanicw
	Panic   = Default.Panic
	Panicw  = Default.Panicw
	Fatal   = Default.Fatal
	Fatalw  = Default.Fatalw
	Debug   = Default.Debug
	Debugw  = Default.Debugw
)

func New(writer io.Writer, level zapcore.Level, env string, extraOpts ...zap.Option) *zap.SugaredLogger {
	var cfg zapcore.EncoderConfig
	var encoder zapcore.Encoder
	opts := make([]zap.Option, 0, len(extraOpts))

	switch env {
	case "production":
		cfg = zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(cfg)
	case "development":
		cfg = zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
		opts = append(opts, zap.WithCaller(true))
	}
	opts = append(opts, extraOpts...)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		level,
	)

	return zap.New(core, opts...).Sugar()
}

func ResetDefault(l *zap.SugaredLogger) {
	Default = l
	Info = Default.Info
	Infow = Default.Infow
	Warn = Default.Warn
	Warnw = Default.Warnw
	Error = Default.Error
	Errorw = Default.Errorw
	DPanic = Default.DPanic
	DPanicw = Default.DPanicw
	Panic = Default.Panic
	Panicw = Default.Panicw
	Fatal = Default.Fatal
	Fatalw = Default.Fatalw
	Debug = Default.Debug
	Debugw = Default.Debugw
}

func Init() *zap.SugaredLogger {
	var logLevel zapcore.Level
	var env string
	devMode := os.Getenv("DEVMODE")

	if strings.EqualFold(devMode, "true") {
		logLevel = zapcore.DebugLevel
		env = "development"
	} else {
		logLevel = zapcore.InfoLevel
		env = "production"
	}

	Default = New(
		os.Stderr,
		logLevel,
		env,
		// for potentially adding default fields to logger
		// zap.Fields(zap.Field{Key: "log_type", Type: zapcore.StringType, String: "default"}),
	)

	ResetDefault(Default)

	Goose = GooseLogger{}
	return Default
}

func StdInfo(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

func StdFatal(format string, v ...any) {
	Fatal(fmt.Sprintf(format, v...))
}
