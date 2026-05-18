/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Copyright 2025 kozex-ai Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	workflowModel "github.com/kozex-ai/kozex/backend/crossdomain/workflow/model"
	"github.com/kozex-ai/kozex/backend/domain/workflow"
	"github.com/kozex-ai/kozex/backend/domain/workflow/entity"
	"github.com/kozex-ai/kozex/backend/domain/workflow/entity/vo"
	mock_workflow "github.com/kozex-ai/kozex/backend/internal/mock/domain/workflow"
	mockeventbus "github.com/kozex-ai/kozex/backend/internal/mock/infra/eventbus"
)

func TestQueueExecute_PublishesJobAndReturnsID(t *testing.T) {
	mockey.PatchConvey("QueueExecute publishes job to MQ and returns queued execute ID", t, func() {
		ctrl := gomock.NewController(t, gomock.WithOverridableExpectations())
		defer ctrl.Finish()

		mockRepo := mock_workflow.NewMockRepository(ctrl)
		mockProducer := mockeventbus.NewMockProducer(ctrl)

		// Wire mock repo as the global repository
		mockey.Mock(workflow.GetRepository).Return(mockRepo).Build()

		svc := NewWorkflowService(mockRepo, mockProducer).(*impl)

		cfg := workflowModel.ExecuteConfig{
			ID:          42,
			From:        workflowModel.FromDraft,
			Operator:    100,
			Mode:        workflowModel.ExecuteModeDebug,
			SyncPattern: workflowModel.SyncPatternAsync,
		}
		input := map[string]any{"key": "value"}

		// Stub: Get (MetaOnly) returns basic workflow info
		mockRepo.EXPECT().GetEntity(gomock.Any(), &vo.GetPolicy{
			ID:       cfg.ID,
			QType:    cfg.From,
			MetaOnly: true,
			Version:  cfg.Version,
		}).Return(&entity.Workflow{
			ID:   42,
			Meta: &vo.Meta{SpaceID: 99},
		}, nil)

		// Stub: GenID returns a known execute ID
		const wantExecuteID = int64(1001)
		mockRepo.EXPECT().GenID(gomock.Any()).Return(wantExecuteID, nil)

		// Stub: CreateWorkflowExecution must receive Status=Queued
		mockRepo.EXPECT().CreateWorkflowExecution(gomock.Any(), gomock.AssignableToTypeOf(&entity.WorkflowExecution{})).
			DoAndReturn(func(_ context.Context, exe *entity.WorkflowExecution) error {
				assert.Equal(t, wantExecuteID, exe.ID)
				assert.Equal(t, entity.WorkflowQueued, exe.Status)
				assert.Equal(t, int64(42), exe.WorkflowID)
				return nil
			})

		// Stub: SetTestRunLatestExeID (debug mode)
		mockRepo.EXPECT().SetTestRunLatestExeID(gomock.Any(), cfg.ID, cfg.Operator, wantExecuteID).
			Return(nil)

		// Stub: Producer.Send — verify it's called (body check omitted; JSON content tested separately)
		mockProducer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)

		gotID, err := svc.QueueExecute(context.Background(), cfg, input)
		assert.NoError(t, err)
		assert.Equal(t, wantExecuteID, gotID)
	})
}

func TestQueueExecute_ReturnsErrorWhenProducerNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_workflow.NewMockRepository(ctrl)
	svc := NewWorkflowService(mockRepo, nil).(*impl)

	_, err := svc.QueueExecute(context.Background(), workflowModel.ExecuteConfig{ID: 1}, nil)
	assert.ErrorContains(t, err, "producer is not configured")
}

func TestQueueExecute_ReturnsErrorWhenPublishFails(t *testing.T) {
	mockey.PatchConvey("QueueExecute returns error when MQ publish fails", t, func() {
		ctrl := gomock.NewController(t, gomock.WithOverridableExpectations())
		defer ctrl.Finish()

		mockRepo := mock_workflow.NewMockRepository(ctrl)
		mockProducer := mockeventbus.NewMockProducer(ctrl)
		mockey.Mock(workflow.GetRepository).Return(mockRepo).Build()

		svc := NewWorkflowService(mockRepo, mockProducer).(*impl)

		cfg := workflowModel.ExecuteConfig{ID: 42, From: workflowModel.FromDraft}

		mockRepo.EXPECT().GetEntity(gomock.Any(), gomock.Any()).Return(&entity.Workflow{
			ID:   42,
			Meta: &vo.Meta{SpaceID: 1},
		}, nil)
		mockRepo.EXPECT().GenID(gomock.Any()).Return(int64(999), nil)
		mockRepo.EXPECT().CreateWorkflowExecution(gomock.Any(), gomock.Any()).Return(nil)
		mockProducer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(assert.AnError)

		_, err := svc.QueueExecute(context.Background(), cfg, nil)
		assert.ErrorContains(t, err, "failed to publish workflow job")
	})
}

func TestQueueExecute_StatusIsQueuedNotRunning(t *testing.T) {
	mockey.PatchConvey("status written to DB must be WorkflowQueued (6), not WorkflowRunning (1)", t, func() {
		ctrl := gomock.NewController(t, gomock.WithOverridableExpectations())
		defer ctrl.Finish()

		mockRepo := mock_workflow.NewMockRepository(ctrl)
		mockProducer := mockeventbus.NewMockProducer(ctrl)
		mockey.Mock(workflow.GetRepository).Return(mockRepo).Build()

		svc := NewWorkflowService(mockRepo, mockProducer).(*impl)

		var capturedStatus entity.WorkflowExecuteStatus
		mockRepo.EXPECT().GetEntity(gomock.Any(), gomock.Any()).Return(&entity.Workflow{
			ID: 1, Meta: &vo.Meta{}, CommitID: "abc",
		}, nil)
		mockRepo.EXPECT().GenID(gomock.Any()).Return(int64(100), nil)
		mockRepo.EXPECT().CreateWorkflowExecution(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, exe *entity.WorkflowExecution) error {
				capturedStatus = exe.Status
				return nil
			})
		mockProducer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)

		_, err := svc.QueueExecute(context.Background(), workflowModel.ExecuteConfig{
			ID:   1,
			From: workflowModel.FromDraft,
		}, nil)
		assert.NoError(t, err)
		assert.Equal(t, entity.WorkflowQueued, capturedStatus, "status must be WorkflowQueued(6), not WorkflowRunning(1)")
		assert.NotEqual(t, entity.WorkflowRunning, capturedStatus)

		// Sanity check: the constant values are what we expect
		assert.Equal(t, entity.WorkflowExecuteStatus(6), entity.WorkflowQueued)
		assert.Equal(t, entity.WorkflowExecuteStatus(1), entity.WorkflowRunning)
	})
}
