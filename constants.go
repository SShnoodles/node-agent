package node

const (
	PIDFileName          = "pid.txt"
	PortFileName         = "port.txt"
	WindowsStopFileName  = "stop.bat"
	LinuxStopFileName    = "stop.sh"
	WindowsStartFileName = "start.bat"
	LinuxStartFileName   = "start.sh"

	// BinaryFileRegex matches both JAR filenames (myapp-1.0.0.jar) and
	// Go binary names (myapp-1.0.0) inside start scripts.
	BinaryFileRegex = `([\w-]+)-([\d]+\.[\d]+\.[\d]+)(?:\.jar)?`
)
