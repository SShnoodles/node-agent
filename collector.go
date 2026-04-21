package node

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

// AppInfoCollector collects application runtime information.
type AppInfoCollector struct {
	port      int
	cfg       Config
	startTime int64
}

// NewAppInfoCollector creates a new AppInfoCollector.
func NewAppInfoCollector(cfg Config) *AppInfoCollector {
	return &AppInfoCollector{
		cfg:       cfg,
		startTime: time.Now().UnixMilli(),
	}
}

// SetPort sets the application listening port used during collection.
func (c *AppInfoCollector) SetPort(port int) {
	c.port = port
}

// Collect gathers all application runtime information.
func (c *AppInfoCollector) Collect() (*AppInfo, error) {
	jarInfo, err := GetJarInfo()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(c.cfg.WorkDir)
	if err != nil {
		return nil, err
	}

	return &AppInfo{
		Path:      filepath.Dir(absPath),
		PID:       strconv.Itoa(os.Getpid()),
		Port:      c.port,
		Jar:       jarInfo.Jar,
		App:       jarInfo.App,
		Version:   jarInfo.Version,
		StartTime: c.startTime,
	}, nil
}

// IsWindows returns true when the current OS is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// GetJarInfo parses the start script (start.sh / start.bat) to extract
// the binary name, application name, and version.
func GetJarInfo() (*JarInfo, error) {
	startFile := LinuxStartFileName
	if IsWindows() {
		startFile = WindowsStartFileName
	}

	content, err := os.ReadFile(startFile)
	if err != nil {
		return nil, fmt.Errorf("cannot find start file: %s", startFile)
	}

	re := regexp.MustCompile(BinaryFileRegex)
	matches := re.FindSubmatch(content)
	if len(matches) < 3 {
		return nil, fmt.Errorf("cannot get binary name from start file: %s", startFile)
	}

	return &JarInfo{
		Jar:     string(matches[0]),
		App:     string(matches[1]),
		Version: string(matches[2]),
	}, nil
}
