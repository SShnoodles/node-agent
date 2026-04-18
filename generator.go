package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// NodeFileGenerator generates control files (pid, port, stop scripts) in the work directory.
type NodeFileGenerator struct {
	cfg     Config
	pidFile string
}

// NewNodeFileGenerator creates a new NodeFileGenerator.
func NewNodeFileGenerator(cfg Config) *NodeFileGenerator {
	return &NodeFileGenerator{cfg: cfg}
}

// Start creates the work directory and generates all control files.
func (g *NodeFileGenerator) Start(port int) error {
	dir := g.cfg.WorkDir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create work dir %q: %w", dir, err)
	}
	if err := g.createCompatJar(dir); err != nil {
		return fmt.Errorf("failed to create compat jar file: %w", err)
	}
	if err := g.createPID(dir); err != nil {
		return err
	}
	if err := g.createPort(dir, port); err != nil {
		return err
	}
	if err := g.createStopBat(dir); err != nil {
		return err
	}
	return g.createStopSh(dir)
}

// Stop removes the PID file to signal that the application has stopped.
func (g *NodeFileGenerator) Stop() {
	if g.pidFile != "" {
		os.Remove(g.pidFile)
	}
}

func (g *NodeFileGenerator) createPID(dir string) error {
	pid := strconv.Itoa(os.Getpid())

	// Plain PID file (used by stop scripts)
	if err := os.WriteFile(filepath.Join(dir, "pid"), []byte(pid), 0644); err != nil {
		return err
	}

	// PID file, for old version
	bracketedPID := fmt.Sprintf("#%s#", pid)
	pidFilePath := filepath.Join(dir, PIDFileName)
	if err := os.WriteFile(pidFilePath, []byte(bracketedPID), 0644); err != nil {
		return err
	}
	g.pidFile = pidFilePath
	return nil
}

func (g *NodeFileGenerator) createCompatJar(dir string) error {
	execPath, err := os.Executable()
	if err != nil || execPath == "" {
		execPath = os.Args[0]
	}

	name := compatJarBaseName(execPath)
	if name == "" {
		return fmt.Errorf("empty executable name")
	}

	jarPath := filepath.Join(dir, name+".jar")
	f, err := os.OpenFile(jarPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func compatJarBaseName(execPath string) string {
	name := filepath.Base(execPath)
	if strings.EqualFold(filepath.Ext(name), ".exe") {
		return strings.TrimSuffix(name, filepath.Ext(name))
	}
	return name
}

func (g *NodeFileGenerator) createPort(dir string, port int) error {
	return os.WriteFile(filepath.Join(dir, PortFileName), []byte(strconv.Itoa(port)), 0644)
}

func (g *NodeFileGenerator) createStopBat(dir string) error {
	script := "@echo off\r\nset /p pid=<pid\r\ntaskkill /PID %pid% /F\r\n"
	return os.WriteFile(filepath.Join(dir, WindowsStopFileName), []byte(script), 0644)
}

func (g *NodeFileGenerator) createStopSh(dir string) error {
	script := "#!/bin/bash\npid=$(cat pid)\nkill -9 $pid\n"
	return os.WriteFile(filepath.Join(dir, LinuxStopFileName), []byte(script), 0755)
}
