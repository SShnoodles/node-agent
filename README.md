# node-agent

Provides lifecycle management and heartbeat reporting capabilities for Go applications.

Designed to replace legacy Java services with consideration for limited host memory.

## Features

- Generate control files (PID, port, stop script) in the working directory on startup
- Report heartbeats to the central server every second
- Capture panic / gin.Error and report error information
- Auto cleanup PID files on graceful shutdown
- Cross-platform support (Linux / macOS / Windows)

## Installation

```bash
go get github.com/SShnoodles/node-agent
```

## Quick Start

```go
import node "github.com/SShnoodles/node-agent"

// Start registration
n := node.New(node.DefaultConfig())
if err := n.Start(8080); err != nil {
    log.Fatal(err)
}

// Exception reporting (optional)
// For example, use in gin Recovery middleware
n.reporter.ReportError(...);
```

## Configuration

Use `node.DefaultConfig()` to get the default configuration, or manually construct a `Config`:

```go
cfg := node.Config{
    Enabled:        true,
    WorkDir:        "target",                                      // Control file output directory
    ReportURL:      "http://127.0.0.1:28080/report_beatheart",    // Heartbeat report URL
    ReportErrorURL: "http://127.0.0.1:28080/report_error",        // Error report URL
    PrintHTTP:      false,                                         // Enable HTTP debug logging with slog
}
```

## Generated Control Files

After the application starts, Node generates the following files in the `WorkDir` (default `target/`) directory:

| File        | Content                                    |
|-------------|--------------------------------------------|
| `pid`       | Process ID (numeric only)                  |
| `pid.txt`   | Process ID (format: `#1234#`) for backward compatibility |
| `port.txt`  | Listening port number                      |
| `stop.sh`   | Stop script for Linux/macOS                |
| `stop.bat`  | Stop script for Windows                    |

## Heartbeat and Error Reporting

**Heartbeat data** (POST to `ReportURL` every second):

```json
{
  "jar":       "myapp-1.0.0",
  "port":      8080,
  "app":       "myapp",
  "version":   "1.0.0",
  "startTime": 1700000000000,
  "path":      "/home/user/myapp",
  "pid":       "12345"
}
```

**Error data** (POST to `ReportErrorURL` on panic):

```json
{
  "app":   "myapp",
  "error": "runtime error: index out of range",
  "timestamp": 1700000000000
}
```

> Report request timeout is 3 seconds. Failures are silently ignored and do not affect the main business logic.

