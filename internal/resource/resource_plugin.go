// Copyright (c) F5, Inc.
//
// This source code is licensed under the Apache License, Version 2.0 license found in the
// LICENSE file in the root directory of this source tree.

package resource

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	mpi "github.com/nginx/agent/v3/api/grpc/mpi/v1"
	"github.com/nginx/agent/v3/internal/config"
	"github.com/nginx/agent/v3/internal/logger"
	"github.com/nginx/agent/v3/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nginx/agent/v3/internal/bus"
)

// The resource plugin listens for a writeConfigSuccessfulTopic from the file plugin after the config apply
// files have been written. The resource plugin then, validates the config, reloads the instance and monitors the logs.
// This is done in the resource plugin to make the file plugin usable for every type of instance.

type Resource struct {
	messagePipe     bus.MessagePipeInterface
	resourceService resourceServiceInterface
	agentConfig     *config.Config
}

var _ bus.Plugin = (*Resource)(nil)

func NewResource(agentConfig *config.Config) *Resource {
	return &Resource{
		agentConfig: agentConfig,
	}
}

func (r *Resource) Init(ctx context.Context, messagePipe bus.MessagePipeInterface) error {
	slog.DebugContext(ctx, "Starting resource plugin")

	r.messagePipe = messagePipe
	r.resourceService = NewResourceService(ctx, r.agentConfig)

	return nil
}

func (*Resource) Close(ctx context.Context) error {
	slog.DebugContext(ctx, "Closing resource plugin")
	return nil
}

func (*Resource) Info() *bus.Info {
	return &bus.Info{
		Name: "resource",
	}
}

// cyclomatic complexity 11 max is 10
// nolint: revive, cyclop
func (r *Resource) Process(ctx context.Context, msg *bus.Message) {
	switch msg.Topic {
	case bus.AddInstancesTopic:
		instanceList, ok := msg.Data.([]*mpi.Instance)
		if !ok {
			slog.ErrorContext(ctx, "Unable to cast message payload to []*mpi.Instance", "payload", msg.Data)

			return
		}

		resource := r.resourceService.AddInstances(instanceList)

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.ResourceUpdateTopic, Data: resource})

		return
	case bus.UpdatedInstancesTopic:
		instanceList, ok := msg.Data.([]*mpi.Instance)
		if !ok {
			slog.ErrorContext(ctx, "Unable to cast message payload to []*mpi.Instance", "payload", msg.Data)

			return
		}
		resource := r.resourceService.UpdateInstances(instanceList)

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.ResourceUpdateTopic, Data: resource})

		return

	case bus.DeletedInstancesTopic:
		instanceList, ok := msg.Data.([]*mpi.Instance)
		if !ok {
			slog.ErrorContext(ctx, "Unable to cast message payload to []*mpi.Instance", "payload", msg.Data)

			return
		}
		resource := r.resourceService.DeleteInstances(instanceList)

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.ResourceUpdateTopic, Data: resource})

		return
	case bus.WriteConfigSuccessfulTopic:
		r.handleWriteConfigSuccessful(ctx, msg)
	case bus.RollbackWriteTopic:
		r.handleRollbackWrite(ctx, msg)
	case bus.APIActionRequestTopic:
		r.handleAPIActionRequest(ctx, msg)
	default:
		slog.DebugContext(ctx, "Unknown topic", "topic", msg.Topic)
	}
}

func (*Resource) Subscriptions() []string {
	return []string{
		bus.AddInstancesTopic,
		bus.UpdatedInstancesTopic,
		bus.DeletedInstancesTopic,
		bus.WriteConfigSuccessfulTopic,
		bus.RollbackWriteTopic,
		bus.APIActionRequestTopic,
	}
}

func (r *Resource) handleAPIActionRequest(ctx context.Context, msg *bus.Message) {
	managementPlaneRequest, ok := msg.Data.(*mpi.ManagementPlaneRequest)

	if !ok {
		slog.ErrorContext(ctx, "Unable to cast message payload to *mpi.ManagementPlaneRequest", "payload",
			msg.Data)

		return
	}

	request, requestOk := managementPlaneRequest.GetRequest().(*mpi.ManagementPlaneRequest_ActionRequest)
	if !requestOk {
		slog.ErrorContext(ctx, "Unable to cast message payload to *mpi.ManagementPlaneRequest_ActionRequest",
			"payload", msg.Data)
	}

	instanceID := request.ActionRequest.GetInstanceId()

	switch request.ActionRequest.GetAction().(type) {
	case *mpi.APIActionRequest_NginxPlusAction:
		r.handleNginxPlusActionRequest(ctx, request.ActionRequest.GetNginxPlusAction(), instanceID)
	default:
		slog.DebugContext(ctx, "API action request not implemented yet")
	}
}

func (r *Resource) handleNginxPlusActionRequest(ctx context.Context, action *mpi.NGINXPlusAction, instanceID string) {
	instance := r.resourceService.Instance(instanceID)
	if instance == nil {
		slog.Info("Unable to find instance with ID", "id", instanceID)
	}
	correlationID := logger.GetCorrelationID(ctx)

	if instance.GetInstanceMeta().GetInstanceType() != mpi.InstanceMeta_INSTANCE_TYPE_NGINX_PLUS {
		slog.ErrorContext(ctx, "", "err", errors.New("failed to preform API action, instance is not NGINX Plus"))
		resp := r.createDataPlaneResponse(correlationID, mpi.CommandResponse_COMMAND_STATUS_FAILURE,
			"", instanceID, "failed to preform API action, instance is not NGINX Plus")

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: resp})

		return
	}

	switch action.GetAction().(type) {
	case *mpi.NGINXPlusAction_UpdateHttpUpstreamServers:
		add, update, del, err := r.resourceService.UpdateHTTPUpstreams(ctx, instance,
			action.GetUpdateHttpUpstreamServers().GetHttpUpstreamName(),
			action.GetUpdateHttpUpstreamServers().GetServers())
		if err != nil {
			resp := r.createDataPlaneResponse(correlationID, mpi.CommandResponse_COMMAND_STATUS_FAILURE,
				"", instanceID, err.Error())
			r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: resp})

			return
		}

		slog.DebugContext(ctx, "successfully updated http upstreams", "add", len(add),
			"update", len(update), "delete", len(del))
		resp := r.createDataPlaneResponse(correlationID, mpi.CommandResponse_COMMAND_STATUS_OK,
			"Successfully updated HTTP Upstreams", instanceID, "")

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: resp})

		slog.DebugContext(ctx, "Updating http upstream servers", "", err)
	case *mpi.NGINXPlusAction_GetHttpUpstreamServers:
		slog.DebugContext(ctx, "Get http upstream servers")
		upstreams, err := r.resourceService.GetUpstreams(ctx, instance,
			action.GetGetHttpUpstreamServers().GetHttpUpstreamName())
		if err != nil {
			resp := r.createDataPlaneResponse(correlationID, mpi.CommandResponse_COMMAND_STATUS_FAILURE,
				"", instanceID, err.Error())
			r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: resp})

			return
		}

		upstreamsJSON, err := json.Marshal(upstreams)
		if err != nil {
			slog.ErrorContext(ctx, "Unable to marshal http upstreams", "err", err)
		}
		resp := r.createDataPlaneResponse(correlationID, mpi.CommandResponse_COMMAND_STATUS_OK,
			string(upstreamsJSON), instanceID, "")

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: resp})
	default:
		slog.DebugContext(ctx, "NGINX Plus action not implemented yet")
	}
}

func (r *Resource) handleWriteConfigSuccessful(ctx context.Context, msg *bus.Message) {
	data, ok := msg.Data.(*model.ConfigApplyMessage)
	if !ok {
		slog.ErrorContext(ctx, "Unable to cast message payload to *model.ConfigApplyMessage", "payload", msg.Data)

		return
	}
	err := r.resourceService.ApplyConfig(ctx, data.InstanceID)
	if err != nil {
		data.Error = err
		slog.Error("errors found during config apply, sending error status, rolling back config", "err", err)
		response := r.createDataPlaneResponse(data.CorrelationID, mpi.CommandResponse_COMMAND_STATUS_ERROR,
			"Config apply failed, rolling back config", data.InstanceID, err.Error())
		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: response})

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.ConfigApplyFailedTopic, Data: data})

		return
	}

	response := r.createDataPlaneResponse(data.CorrelationID, mpi.CommandResponse_COMMAND_STATUS_OK,
		"Config apply successful", data.InstanceID, "")

	r.messagePipe.Process(
		ctx,
		&bus.Message{
			Topic: bus.ConfigApplySuccessfulTopic,
			Data:  response,
		},
	)
}

func (r *Resource) handleRollbackWrite(ctx context.Context, msg *bus.Message) {
	data, ok := msg.Data.(*model.ConfigApplyMessage)
	if !ok {
		slog.ErrorContext(ctx, "Unable to cast message payload to *model.ConfigApplyMessage", "payload", msg.Data)

		return
	}
	err := r.resourceService.ApplyConfig(ctx, data.InstanceID)
	if err != nil {
		slog.Error("errors found during rollback, sending failure status", "err", err)

		rollbackResponse := r.createDataPlaneResponse(data.CorrelationID,
			mpi.CommandResponse_COMMAND_STATUS_ERROR, "Rollback failed", data.InstanceID, err.Error())

		applyResponse := r.createDataPlaneResponse(data.CorrelationID,
			mpi.CommandResponse_COMMAND_STATUS_FAILURE, "Config apply failed, rollback failed",
			data.InstanceID, data.Error.Error())

		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.DataPlaneResponseTopic, Data: rollbackResponse})
		r.messagePipe.Process(ctx, &bus.Message{Topic: bus.ConfigApplyCompleteTopic, Data: applyResponse})

		return
	}

	applyResponse := r.createDataPlaneResponse(data.CorrelationID,
		mpi.CommandResponse_COMMAND_STATUS_FAILURE,
		"Config apply failed, rollback successful", data.InstanceID, data.Error.Error())

	r.messagePipe.Process(ctx, &bus.Message{Topic: bus.ConfigApplyCompleteTopic, Data: applyResponse})
}

func (*Resource) createDataPlaneResponse(correlationID string, status mpi.CommandResponse_CommandStatus,
	message, instanceID, err string,
) *mpi.DataPlaneResponse {
	return &mpi.DataPlaneResponse{
		MessageMeta: &mpi.MessageMeta{
			MessageId:     uuid.NewString(),
			CorrelationId: correlationID,
			Timestamp:     timestamppb.Now(),
		},
		CommandResponse: &mpi.CommandResponse{
			Status:  status,
			Message: message,
			Error:   err,
		},
		InstanceId: instanceID,
	}
}
