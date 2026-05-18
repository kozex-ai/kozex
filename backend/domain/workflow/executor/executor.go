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

	workflowModel "github.com/kozex-ai/kozex/backend/crossdomain/workflow/model"
	"github.com/kozex-ai/kozex/backend/domain/workflow"
	"github.com/kozex-ai/kozex/backend/infra/eventbus"
	"github.com/kozex-ai/kozex/backend/pkg/logs"
	"github.com/kozex-ai/kozex/backend/pkg/sonic"
)

// Handler implements eventbus.ConsumerHandler for workflow execution jobs.
type Handler struct {
	svc workflow.Service
}

func NewHandler(svc workflow.Service) eventbus.ConsumerHandler {
	return &Handler{svc: svc}
}

func (h *Handler) HandleMessage(ctx context.Context, msg *eventbus.Message) error {
	var job workflowModel.WorkflowJob
	if err := sonic.UnmarshalString(string(msg.Body), &job); err != nil {
		return fmt.Errorf("executor: failed to unmarshal job: %w", err)
	}

	logs.CtxInfof(ctx, "executor: processing workflow job execute_id=%d workflow_id=%d",
		job.ExecuteID, job.Config.ID)

	if err := h.svc.ExecuteJob(ctx, job); err != nil {
		logs.CtxErrorf(ctx, "executor: job failed execute_id=%d err=%v", job.ExecuteID, err)
		return err
	}

	return nil
}
