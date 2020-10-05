package observability

import (
	"context"
	"fmt"
	"io"
)

type Exporter interface {
	io.Closer
	StartExporter(func() error) error
}

func NewFromEnv(ctx context.Context, config *Config) (Exporter, error) {
	switch config.ExporterType {
	case ExporterNoop:
		return NewNoop(ctx)
	case ExporterPrometheus, ExporterOCAgent:
		return NewOpenCensus(ctx, config.OpenCensus)
	default:
		return nil, fmt.Errorf("unknown observability exporter type %v", config.ExporterType)
	}
}
