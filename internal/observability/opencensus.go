package observability

import (
	"context"
	"fmt"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type opencensusExporter struct {
	exporter *ocagent.Exporter
	config   *OpenCensusConfig
}

var _ Exporter = (*opencensusExporter)(nil)

// NewOpenCensus creates a new metrics and trace exporter for OpenCensus.
func NewOpenCensus(_ context.Context, config *OpenCensusConfig) (Exporter, error) {
	var opts []ocagent.ExporterOption
	if config.Insecure {
		opts = append(opts, ocagent.WithInsecure())
	}
	if config.Endpoint != "" {
		opts = append(opts, ocagent.WithAddress(config.Endpoint))
	}

	oc, err := ocagent.NewExporter(opts...)
	if err != nil {
		return nil, fmt.Errorf("create opencensus exporter: %w", err)
	}
	return &opencensusExporter{
		exporter: oc,
		config:   config,
	}, nil
}

func (e *opencensusExporter) StartExporter(register func() error) error {
	trace.ApplyConfig(
		trace.Config{
			DefaultSampler: trace.ProbabilitySampler(e.config.SampleRate),
		},
	)
	trace.RegisterExporter(e.exporter)
	view.RegisterExporter(e.exporter)

	if err := register(); err != nil {
		return fmt.Errorf("start opencensus exporter: %w", err)
	}

	return nil
}

func (e *opencensusExporter) Close() error {
	if err := e.exporter.Stop(); err != nil {
		return fmt.Errorf("stop opencensus exporter: %w", err)
	}

	trace.UnregisterExporter(e.exporter)
	view.UnregisterExporter(e.exporter)

	return nil
}
