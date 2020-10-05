package observability

import "context"

// noopExporter is an observability exporter that does nothing.
type noopExporter struct{}

var _ Exporter = (*noopExporter)(nil)

func NewNoop(_ context.Context) (Exporter, error) {
	return &noopExporter{}, nil
}

func (e *noopExporter) StartExporter(_ func() error) error {
	return nil
}

func (e *noopExporter) Close() error {
	return nil
}
