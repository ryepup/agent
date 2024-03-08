// Copyright (c) F5, Inc.
//
// This source code is licensed under the Apache License, Version 2.0 license found in the
// LICENSE file in the root directory of this source tree.

package plugin

import (
	"log/slog"
	"testing"

	"github.com/nginx/agent/v3/internal/bus"
	"github.com/nginx/agent/v3/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadPLugins(t *testing.T) {
	tests := []struct {
		name     string
		input    *config.Config
		expected []bus.Plugin
	}{
		{
			name: "Only process manager plugin enabled",
			input: &config.Config{
				ProcessMonitor: &config.ProcessMonitor{
					MonitoringFrequency: 500,
				},
			},
			expected: []bus.Plugin{
				&ProcessMonitor{},
				&Instance{},
				&Config{},
			},
		}, {
			name: "DataPlane API plugin enabled",
			input: &config.Config{
				DataPlaneAPI: &config.DataPlaneAPI{
					Host: "localhost",
					Port: 8080,
				},
			},
			expected: []bus.Plugin{
				&Instance{},
				&Config{},
				&DataPlaneServer{},
			},
		}, {
			name: "Metrics plugin enabled",
			input: &config.Config{
				DataPlaneAPI: &config.DataPlaneAPI{
					Host: "localhost",
					Port: 8080,
				},
				Metrics: &config.Metrics{},
			},
			expected: []bus.Plugin{
				&Instance{},
				&Metrics{},
				&Config{},
				&DataPlaneServer{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			t.Logf("running test %s", test.name)
			result := LoadPlugins(test.input, slog.New(&slog.TextHandler{}))
			assert.Equal(tt, len(test.expected), len(result))
			for i, expectedPlugin := range test.expected {
				assert.IsType(tt, expectedPlugin, result[i])
			}
		})
	}
}
