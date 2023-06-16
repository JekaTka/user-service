package logger

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// ZeroLogger is an Fx event logger that logs events to Zero.
type ZeroLogger struct {
	Logger *zerolog.Logger
}

// FuncName returns a funcs formatted name
func FuncName(fn interface{}) string {
	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		return fmt.Sprint(fn)
	}

	function := runtime.FuncForPC(fnV.Pointer()).Name()
	return fmt.Sprintf("%s()", function)
}

// LogEvent logs the given event to the provided logger.
func (l *ZeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.Supplied:
		l.Logger.Info().
			Str("type", e.TypeName).
			Msg("supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Info().
				Str("constructor", FuncName(e.ConstructorName)).
				Str("type", rtype).
				Msg("provided")
		}
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("error encountered while applying options")
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Str("trace", e.Trace).
				Str("function", FuncName(e.FunctionName)).
				Msg("invoke failed")

		} else {
			l.Logger.Info().
				Str("function", FuncName(e.FunctionName)).
				Msg("invoked")
		}
	case *fxevent.Stopping:
		l.Logger.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.Logger.Error().
			Err(e.StartErr).
			Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("start failed")
		} else {
			l.Logger.Info().
				Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("custom logger installation failed")
		} else {
			l.Logger.Info().
				Str("constructor", FuncName(e.ConstructorName)).
				Msg("installed custom fxevent.Logger")
		}
	}
}

const timeFormat = time.RFC3339

// NewLogger initialize and return zerolog.Logger
func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = timeFormat
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat}

	var writers []io.Writer
	writers = append(writers, consoleWriter)

	log.Logger = log.Output(zerolog.MultiLevelWriter(writers...))

	return log.Logger
}

// NewPtrLogger return logger pointer
func NewPtrLogger(logger zerolog.Logger) *zerolog.Logger {
	return &logger
}

var Module = fx.Provide(NewLogger, NewPtrLogger)
