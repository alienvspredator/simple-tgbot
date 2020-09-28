package zapfluentd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"go.uber.org/zap/zapcore"
)

type writeSyncer struct {
	fl *fluent.Fluent
}

var _ zapcore.WriteSyncer = (*writeSyncer)(nil)

func NewWriteSyncer(cfg *Config) (zapcore.WriteSyncer, error) {
	fl, err := fluent.New(
		fluent.Config{
			FluentPort: cfg.Port,
			FluentHost: cfg.Host,
			TagPrefix:  "tgbot",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating fluent: %w", err)
	}

	return zapcore.Lock(&writeSyncer{fl: fl}), nil
}

func (s *writeSyncer) Write(p []byte) (n int, err error) {
	var ent map[string]interface{}
	if err := json.Unmarshal(p, &ent); err != nil {
		return 0, fmt.Errorf("unmarshalling bytes: %w", err)
	}
	ts := time.Now()
	if logTs, ok := ent["timestamp"].(time.Time); ok {
		ts = logTs
	}
	tag := "default"
	if loggerName, ok := ent["logger"].(string); ok {
		tag = loggerName
	}
	if err := s.fl.PostWithTime(tag, ts, ent); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (s *writeSyncer) Sync() error {
	return nil
}
