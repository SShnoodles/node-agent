// Package node provides application lifecycle management and heartbeat reporting.
package node

import (
	"context"
	"fmt"
)

type Node struct {
	cfg       Config
	collector *AppInfoCollector
	reporter  *AppInfoReporter
	generator *NodeFileGenerator
	cancel    context.CancelFunc
}

// New creates a new Node instance with the given configuration.
func New(cfg Config) *Node {
	collector := NewAppInfoCollector(cfg)
	reporter := NewAppInfoReporter(cfg, collector)
	generator := NewNodeFileGenerator(cfg)
	return &Node{
		cfg:       cfg,
		collector: collector,
		reporter:  reporter,
		generator: generator,
	}
}

// Start initialises Node: generates control files and begins heartbeat reporting.
// Call this after the HTTP server has started listening on port.
func (a *Node) Start(port int) error {
	if !a.cfg.Enabled {
		return nil
	}
	a.collector.SetPort(port)
	if err := a.generator.Start(port); err != nil {
		return fmt.Errorf("node: generator failed: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	go a.reporter.StartHeartbeat(ctx)
	return nil
}

// Stop shuts down Node: cancels heartbeat reporting and removes the PID file.
func (a *Node) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
	a.generator.Stop()
}
