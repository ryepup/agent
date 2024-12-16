// Copyright (c) F5, Inc.
//
// This source code is licensed under the Apache License, Version 2.0 license found in the
// LICENSE file in the root directory of this source tree.

package command

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nginx/agent/v3/internal/config"
	pkg "github.com/nginx/agent/v3/pkg/config"

	"github.com/nginx/agent/v3/internal/bus/busfakes"

	mpi "github.com/nginx/agent/v3/api/grpc/mpi/v1"
	"github.com/nginx/agent/v3/internal/bus"
	"github.com/nginx/agent/v3/internal/command/commandfakes"
	"github.com/nginx/agent/v3/internal/grpc/grpcfakes"
	"github.com/nginx/agent/v3/test/protos"
	"github.com/nginx/agent/v3/test/stub"
	"github.com/nginx/agent/v3/test/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandPlugin_Info(t *testing.T) {
	commandPlugin := NewCommandPlugin(types.AgentConfig(), &grpcfakes.FakeGrpcConnectionInterface{})
	info := commandPlugin.Info()

	assert.Equal(t, "command", info.Name)
}

func TestCommandPlugin_Subscriptions(t *testing.T) {
	commandPlugin := NewCommandPlugin(types.AgentConfig(), &grpcfakes.FakeGrpcConnectionInterface{})
	subscriptions := commandPlugin.Subscriptions()

	assert.Equal(
		t,
		[]string{
			bus.ResourceUpdateTopic,
			bus.InstanceHealthTopic,
			bus.DataPlaneHealthResponseTopic,
			bus.DataPlaneResponseTopic,
		},
		subscriptions,
	)
}

func TestCommandPlugin_Init(t *testing.T) {
	ctx := context.Background()
	messagePipe := busfakes.NewFakeMessagePipe()
	fakeCommandService := &commandfakes.FakeCommandService{}

	commandPlugin := NewCommandPlugin(types.AgentConfig(), &grpcfakes.FakeGrpcConnectionInterface{})
	err := commandPlugin.Init(ctx, messagePipe)
	require.NoError(t, err)

	require.NotNil(t, commandPlugin.messagePipe)
	require.NotNil(t, commandPlugin.commandService)

	commandPlugin.commandService = fakeCommandService

	closeError := commandPlugin.Close(ctx)
	require.NoError(t, closeError)
	require.Equal(t, 1, fakeCommandService.CancelSubscriptionCallCount())
}

func TestCommandPlugin_Process(t *testing.T) {
	ctx := context.Background()
	messagePipe := busfakes.NewFakeMessagePipe()
	fakeCommandService := &commandfakes.FakeCommandService{}

	commandPlugin := NewCommandPlugin(types.AgentConfig(), &grpcfakes.FakeGrpcConnectionInterface{})
	err := commandPlugin.Init(ctx, messagePipe)
	require.NoError(t, err)
	defer commandPlugin.Close(ctx)

	// Check CreateConnection
	fakeCommandService.IsConnectedReturnsOnCall(0, false)

	// Check UpdateDataPlaneStatus
	fakeCommandService.IsConnectedReturnsOnCall(1, true)
	fakeCommandService.IsConnectedReturnsOnCall(2, true)

	commandPlugin.commandService = fakeCommandService

	commandPlugin.Process(ctx, &bus.Message{Topic: bus.ResourceUpdateTopic, Data: protos.GetHostResource()})
	require.Equal(t, 1, fakeCommandService.CreateConnectionCallCount())

	commandPlugin.Process(ctx, &bus.Message{Topic: bus.ResourceUpdateTopic, Data: protos.GetHostResource()})
	require.Equal(t, 1, fakeCommandService.UpdateDataPlaneStatusCallCount())

	commandPlugin.Process(ctx, &bus.Message{Topic: bus.InstanceHealthTopic, Data: protos.GetInstanceHealths()})
	require.Equal(t, 1, fakeCommandService.UpdateDataPlaneHealthCallCount())

	commandPlugin.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: protos.OKDataPlaneResponse()})
	require.Equal(t, 1, fakeCommandService.SendDataPlaneResponseCallCount())

	commandPlugin.Process(ctx, &bus.Message{
		Topic: bus.DataPlaneHealthResponseTopic,
		Data:  protos.GetHealthyInstanceHealth(),
	})
	require.Equal(t, 1, fakeCommandService.UpdateDataPlaneHealthCallCount())
	require.Equal(t, 1, fakeCommandService.SendDataPlaneResponseCallCount())
}

func TestCommandPlugin_monitorSubscribeChannel(t *testing.T) {
	tests := []struct {
		managementPlaneRequest *mpi.ManagementPlaneRequest
		expectedTopic          *bus.Message
		name                   string
		request                string
	}{
		{
			name: "Test 1: Config Upload Request",
			managementPlaneRequest: &mpi.ManagementPlaneRequest{
				Request: &mpi.ManagementPlaneRequest_ConfigUploadRequest{
					ConfigUploadRequest: &mpi.ConfigUploadRequest{},
				},
			},
			expectedTopic: &bus.Message{Topic: bus.ConfigUploadRequestTopic},
			request:       "UploadRequest",
		},
		{
			name: "Test 2: Config Apply Request",
			managementPlaneRequest: &mpi.ManagementPlaneRequest{
				Request: &mpi.ManagementPlaneRequest_ConfigApplyRequest{
					ConfigApplyRequest: &mpi.ConfigApplyRequest{},
				},
			},
			expectedTopic: &bus.Message{Topic: bus.ConfigApplyRequestTopic},
			request:       "ApplyRequest",
		},
		{
			name: "Test 3: Health Request",
			managementPlaneRequest: &mpi.ManagementPlaneRequest{
				Request: &mpi.ManagementPlaneRequest_HealthRequest{
					HealthRequest: &mpi.HealthRequest{},
				},
			},
			expectedTopic: &bus.Message{Topic: bus.DataPlaneHealthRequestTopic},
		},
		{
			name: "Test 4: API Action Request Request",
			managementPlaneRequest: &mpi.ManagementPlaneRequest{
				Request: &mpi.ManagementPlaneRequest_ActionRequest{
					ActionRequest: &mpi.APIActionRequest{
						Action: &mpi.APIActionRequest_NginxPlusAction{},
					},
				},
			},
			expectedTopic: &bus.Message{Topic: bus.APIActionRequestTopic},
			request:       "APIActionRequest",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			messagePipe := busfakes.NewFakeMessagePipe()

			agentConfig := types.AgentConfig()

			if test.request == "APIActionRequest" {
				agentConfig.Features = append(config.DefaultFeatures(), pkg.FeatureAPIAction)
			}
			commandPlugin := NewCommandPlugin(agentConfig, &grpcfakes.FakeGrpcConnectionInterface{})
			err := commandPlugin.Init(ctx, messagePipe)
			require.NoError(tt, err)
			defer commandPlugin.Close(ctx)

			go commandPlugin.monitorSubscribeChannel(ctx)

			commandPlugin.subscribeChannel <- test.managementPlaneRequest

			assert.Eventually(
				t,
				func() bool { return len(messagePipe.GetMessages()) == 1 },
				2*time.Second,
				10*time.Millisecond,
			)

			messages := messagePipe.GetMessages()
			assert.Len(tt, messages, 1)
			assert.Equal(tt, test.expectedTopic.Topic, messages[0].Topic)

			_, ok := messages[0].Data.(*mpi.ManagementPlaneRequest)

			switch test.request {
			case "UploadRequest":
				assert.True(tt, ok)
				require.NotNil(tt, test.managementPlaneRequest.GetConfigUploadRequest())
			case "ApplyRequest":
				assert.True(tt, ok)
				require.NotNil(tt, test.managementPlaneRequest.GetConfigApplyRequest())
			case "APIActionRequest":
				assert.True(tt, ok)
				require.NotNil(tt, test.managementPlaneRequest.GetActionRequest())
			}
		})
	}
}

func TestMonitorSubscribeChannel(t *testing.T) {
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	logBuf := &bytes.Buffer{}
	stub.StubLoggerWith(logBuf)

	cp := NewCommandPlugin(types.AgentConfig(), &grpcfakes.FakeGrpcConnectionInterface{})

	message := protos.CreateManagementPlaneRequest()

	// Run in a separate goroutine
	go cp.monitorSubscribeChannel(ctx)

	// Give some time to exit the goroutine
	time.Sleep(100 * time.Millisecond)

	cp.subscribeChannel <- message

	// Give some time to process the message
	time.Sleep(100 * time.Millisecond)

	cncl()

	time.Sleep(100 * time.Millisecond)

	// Verify the logger was called
	if s := logBuf.String(); !strings.Contains(s, "Received management plane request") {
		t.Errorf("Unexpected log %s", s)
	}

	// Clear the log buffer
	logBuf.Reset()
}
