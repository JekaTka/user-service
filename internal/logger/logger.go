package logger

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ZeroLogger is an Fx event logger that logs events to Zero.
type ZeroLogger struct {
	Logger *zerolog.Logger
}

// FuncName returns a funcs formatted name
func FuncName(fn any) string {
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

// GrpcLogger middleware for grpc requests logging
func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

// HttpLogger middleware for http requests logging
func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}

var Module = fx.Provide(NewLogger, NewPtrLogger)
