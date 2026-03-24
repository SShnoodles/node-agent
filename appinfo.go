package node

type AppInfo struct {
	Jar       string `json:"jar"` // backward compatible
	Port      int    `json:"port"`
	App       string `json:"app"`
	Version   string `json:"version"`
	StartTime int64  `json:"startTime"`
	Path      string `json:"path"`
	PID       string `json:"pid"`
}

type JarInfo struct {
	// Jar is the full matched filename, e.g. "myapp-1.0.0.jar" or "myapp-1.0.0"
	Jar string

	// App is the application name extracted from the filename, e.g. "myapp"
	App string

	// Version is the semantic version extracted from the filename, e.g. "1.0.0"
	Version string
}

type AppError struct {
	App       string `json:"app"`
	Error     string `json:"error"`
	Timestamp int64  `json:"timestamp"`
}
