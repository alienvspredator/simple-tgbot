// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logging sets up and configures logging.
package logging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logger is stored.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

type config struct {
	debug bool
	ws    zapcore.WriteSyncer
}

type Option func(cfg config) config

func WithDebugEnabled(enabled bool) Option {
	return func(cfg config) config {
		cfg.debug = enabled
		return cfg
	}
}

func WithWriteSyncer(ws zapcore.WriteSyncer) Option {
	return func(cfg config) config {
		cfg.ws = zapcore.NewMultiWriteSyncer(cfg.ws, ws)
		return cfg
	}
}

// NewLogger creates a new logger with the given configuration.
func NewLogger(opts ...Option) *zap.SugaredLogger {
	sink, closeOut, err := zap.Open(outputStderr...)
	if err != nil {
		return zap.NewNop().Sugar()
	}
	closeOut()

	zapOpts := []zap.Option{zap.ErrorOutput(sink)}

	cfg := config{
		debug: false,
		ws:    sink,
	}
	for _, apply := range opts {
		cfg = apply(cfg)
	}

	lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	if cfg.debug {
		zapOpts = append(zapOpts, zap.Development())
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		zapOpts = append(zapOpts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if samplingConfig.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(samplingConfig.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				samplingConfig.Initial,
				samplingConfig.Thereafter,
				samplerOpts...,
			)
		}))
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), cfg.ws, lvl)
	logger := zap.New(core, zapOpts...)

	return logger.Sugar()
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}

const (
	timestamp  = "timestamp"
	severity   = "severity"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	stacktrace = "stacktrace"

	levelDebug     = "DEBUG"
	levelInfo      = "INFO"
	levelWarning   = "WARNING"
	levelError     = "ERROR"
	levelCritical  = "CRITICAL"
	levelAlert     = "ALERT"
	levelEmergency = "EMERGENCY"

	encodingJSON = "json"
)

var outputStderr = []string{"stderr"}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        timestamp,
	LevelKey:       severity,
	NameKey:        logger,
	CallerKey:      caller,
	MessageKey:     message,
	StacktraceKey:  stacktrace,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    levelEncoder(),
	EncodeTime:     timeEncoder(),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var samplingConfig = &zap.SamplingConfig{
	Initial:    250,
	Thereafter: 250,
}

// levelEncoder transforms a zap level to the associated stackdriver level.
func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(levelDebug)
		case zapcore.InfoLevel:
			enc.AppendString(levelInfo)
		case zapcore.WarnLevel:
			enc.AppendString(levelWarning)
		case zapcore.ErrorLevel:
			enc.AppendString(levelError)
		case zapcore.DPanicLevel:
			enc.AppendString(levelCritical)
		case zapcore.PanicLevel:
			enc.AppendString(levelAlert)
		case zapcore.FatalLevel:
			enc.AppendString(levelEmergency)
		}
	}
}

// TraceFromContext adds the correct Stackdriver trace fields.
//
// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func TraceFromContext(ctx context.Context) []zap.Field {
	span := trace.FromContext(ctx)

	if span == nil {
		return nil
	}

	sc := span.SpanContext()

	return []zap.Field{
		zap.String("trace", fmt.Sprintf("traces/%s", sc.TraceID)),
		zap.String("spanId", sc.SpanID.String()),
		zap.Bool("traceSampled", sc.IsSampled()),
	}
}

// timeEncoder encodes the time as RFC3339 nano
func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}
