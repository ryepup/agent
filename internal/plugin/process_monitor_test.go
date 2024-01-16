/**
 * Copyright (c) F5, Inc.
 *
 * This source code is licensed under the Apache License, Version 2.0 license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/nginx/agent/v3/internal/bus"
	"github.com/nginx/agent/v3/internal/model/os"
	"github.com/stretchr/testify/assert"
)

func TestProcessMonitor_Init(t *testing.T) {
	testProcesses := []*os.Process{{Pid: 123, Name: "nginx"}}

	processMonitor := NewProcessMonitor(&ProcessMonitorParameters{
		getProcessesFunc: func() ([]*os.Process, error) {
			return testProcesses, nil
		},
		MonitoringFrequency: time.Millisecond,
	})

	messagePipe := bus.NewMessagePipe(context.TODO(), 100)
	err := messagePipe.Register(100, []bus.Plugin{processMonitor})
	assert.NoError(t, err)
	go messagePipe.Run()

	time.Sleep(10 * time.Millisecond)

	assert.NotNil(t, processMonitor.messagePipe)
	assert.Equal(t, testProcesses, processMonitor.processes)
}

func TestProcessMonitor_Info(t *testing.T) {
	processMonitor := NewProcessMonitor(&ProcessMonitorParameters{})
	info := processMonitor.Info()
	assert.Equal(t, "process-monitor", info.Name)
}

func TestProcessMonitor_Subscriptions(t *testing.T) {
	processMonitor := NewProcessMonitor(&ProcessMonitorParameters{})
	subscriptions := processMonitor.Subscriptions()
	assert.Equal(t, []string{}, subscriptions)
}