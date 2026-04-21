package node

// Config holds the configuration for Node lifecycle management.
type Config struct {
	// Enabled controls whether Node is active. Default: true.
	Enabled bool

	// WorkDir is the directory where control files (pid, port, stop scripts) are generated.
	// Default: "target"
	WorkDir string

	// ReportURL is the heartbeat report endpoint.
	// Default: "http://127.0.0.1:28080/report_beatheart"
	ReportURL string

	// ReportErrorURL is the error report endpoint.
	// Default: "http://127.0.0.1:28080/report_error"
	ReportErrorURL string

	// PrintHTTP controls whether to print HTTP debug information. Default: false.
	PrintHTTP bool
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		Enabled:        true,
		WorkDir:        "target",
		ReportURL:      "http://127.0.0.1:28080/report_beatheart",
		ReportErrorURL: "http://127.0.0.1:28080/report_error",
		PrintHTTP:      false,
	}
}
