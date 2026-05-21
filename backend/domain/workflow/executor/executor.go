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

// Package executor contains the workflow job consumer.
// It is intentionally separate from the HTTP layer so it can later be split
// into an independent process or pod with minimal refactoring.
package executor

import (
	"context"
	"fmt"
	"runtime/debug"

	workflowModel "github.com/kozex-ai/kozex/backend/crossdomain/workflow/model"
	"github.com/kozex-ai/kozex/backend/domain/workflow"
	"github.com/kozex-ai/kozex/backend/domain/workflow/entity"
	"github.com/kozex-ai/kozex/backend/infra/eventbus"
	"github.com/kozex-ai/kozex/backend/pkg/lang/ptr"
	"github.com/kozex-ai/kozex/backend/pkg/logs"
	"github.com/kozex-ai/kozex/backend/pkg/sonic"
	"github.com/kozex-ai/kozex/backend/types/consts"
)

// Handler implements eventbus.ConsumerHandler for workflow execution jobs.
type Handler struct {
	svc workflow.Service
}

func NewHandler(svc workflow.Service) eventbus.ConsumerHandler {
	return &Handler{svc: svc}
}

func (h *Handler) HandleMessage(ctx context.Context, msg *eventbus.Message) (retErr error) {
	var job workflowModel.WorkflowJob
	if err := sonic.UnmarshalString(string(msg.Body), &job); err != nil {
		return fmt.Errorf("executor: failed to unmarshal job: %w", err)
	}

	// Restore the original log ID from the execution record so all async logs
	// share the same log ID as the HTTP request that enqueued the job.
	if job.ExecuteID != 0 {
		repo := workflow.GetRepository()
		if exe, found, err := repo.GetWorkflowExecution(ctx, job.ExecuteID); err == nil && found {
			if exe.LogID != "" {
				logs.CtxInfof(ctx, "executor: restored log ID from execution record execute_id=%d log_id=%s", job.ExecuteID, exe.LogID)
				ctx = context.WithValue(ctx, consts.CtxLogIDKey, exe.LogID)
			}
			// Job was canceled while waiting in queue — skip execution.
			if exe.Status == entity.WorkflowCancel {
				logs.CtxInfof(ctx, "executor: skipping canceled job execute_id=%d", job.ExecuteID)
				return nil
			}
		}
	}

	defer func() {
		if r := recover(); r != nil {
			logs.CtxErrorf(ctx, "executor: panic recovered execute_id=%d panic=%v\n%s",
				job.ExecuteID, r, debug.Stack())
			if job.ExecuteID != 0 {
				failReason := fmt.Sprintf("panic: %v", r)
				repo := workflow.GetRepository()
				if _, _, err := repo.UpdateWorkflowExecution(ctx, &entity.WorkflowExecution{
					ID:         job.ExecuteID,
					Status:     entity.WorkflowFailed,
					FailReason: ptr.Of(failReason),
				}, []entity.WorkflowExecuteStatus{entity.WorkflowQueued, entity.WorkflowRunning}); err != nil {
					logs.CtxErrorf(ctx, "executor: failed to mark execution as failed after panic: %v", err)
				}
			}
			retErr = nil // do not ask MQ to retry on panic — it would just panic again
		}
	}()

	logs.CtxInfof(ctx, "executor: processing workflow job execute_id=%d workflow_id=%d", job.ExecuteID, job.Config.ID)

	if err := h.svc.ExecuteJob(ctx, job); err != nil {
		logs.CtxErrorf(ctx, "executor: job failed execute_id=%d err=%v", job.ExecuteID, err)
		return err
	}

	return nil
}
