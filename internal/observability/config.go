package observability

type ExporterType string

const (
	ExporterPrometheus ExporterType = "PROMETHEUS"
	ExporterOCAgent    ExporterType = "OCAGENT"
	ExporterNoop       ExporterType = "NOOP"
)

type Config struct {
	ExporterType ExporterType `env:"OBSERVABILITY_EXPORTER, default=NOOP"`

	OpenCensus *OpenCensusConfig
}

// OpenCensusConfig holds the configuration options for the open census exporter
type OpenCensusConfig struct {
	SampleRate float64 `env:"TRACE_PROBABILITY, default=0.40"`

	Insecure bool   `env:"OCAGENT_INSECURE"`
	Endpoint string `env:"OCAGENT_TRACE_EXPORTER_ENDPOINT"`
}
