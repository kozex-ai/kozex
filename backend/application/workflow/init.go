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

package workflow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"gorm.io/gorm"

	"github.com/kozex-ai/kozex/backend/bizpkg/llm/modelbuilder"
	knowledge "github.com/kozex-ai/kozex/backend/domain/knowledge/service"
	dbservice "github.com/kozex-ai/kozex/backend/domain/memory/database/service"
	variables "github.com/kozex-ai/kozex/backend/domain/memory/variables/service"
	plugin "github.com/kozex-ai/kozex/backend/domain/plugin/service"
	search "github.com/kozex-ai/kozex/backend/domain/search/service"
	"github.com/kozex-ai/kozex/backend/domain/workflow"
	"github.com/kozex-ai/kozex/backend/domain/workflow/config"
	wfexecutor "github.com/kozex-ai/kozex/backend/domain/workflow/executor"
	wrapPlugin "github.com/kozex-ai/kozex/backend/domain/workflow/plugin"
	"github.com/kozex-ai/kozex/backend/domain/workflow/service"
	"github.com/kozex-ai/kozex/backend/infra/cache"
	"github.com/kozex-ai/kozex/backend/infra/coderunner"
	"github.com/kozex-ai/kozex/backend/infra/eventbus"
	"github.com/kozex-ai/kozex/backend/infra/idgen"
	"github.com/kozex-ai/kozex/backend/infra/imagex"
	"github.com/kozex-ai/kozex/backend/infra/storage"
	"github.com/kozex-ai/kozex/backend/types/consts"
)

type ServiceComponents struct {
	IDGen                    idgen.IDGenerator
	DB                       *gorm.DB
	Cache                    cache.Cmdable
	DatabaseDomainSVC        dbservice.Database
	VariablesDomainSVC       variables.Variables
	PluginDomainSVC          plugin.PluginService
	KnowledgeDomainSVC       knowledge.Knowledge
	DomainNotifier           search.ResourceEventBus
	Tos                      storage.Storage
	ImageX                   imagex.ImageX
	CPStore                  compose.CheckPointStore
	CodeRunner               coderunner.Runner
	WorkflowBuildInChatModel modelbuilder.BaseChatModel
	WorkflowJobProducer      eventbus.Producer
}

func initWorkflowConfig() (workflow.WorkflowConfig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configBs, err := os.ReadFile(filepath.Join(wd, "resources/conf/workflow/config.yaml"))
	if err != nil {
		return nil, err
	}
	var cfg *config.WorkflowConfig
	err = yaml.Unmarshal(configBs, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func InitService(_ context.Context, components *ServiceComponents) (*ApplicationService, error) {
	service.RegisterAllNodeAdaptors()

	cfg, err := initWorkflowConfig()
	if err != nil {
		return nil, err
	}

	workflowRepo, err := service.NewWorkflowRepository(components.IDGen, components.DB, components.Cache,
		components.Tos, components.CPStore, components.WorkflowBuildInChatModel, cfg)
	if err != nil {
		return nil, err
	}

	workflow.SetRepository(workflowRepo)

	workflowDomainSVC := service.NewWorkflowService(workflowRepo, components.WorkflowJobProducer)
	wrapPlugin.SetOSS(components.Tos)

	coderunner.SetCodeRunner(components.CodeRunner)
	callbacks.AppendGlobalHandlers(service.GetTokenCallbackHandler())

	setEventBus(components.DomainNotifier)

	if components.WorkflowJobProducer != nil {
		nameServer := os.Getenv(consts.MQServer)
		handler := wfexecutor.NewHandler(workflowDomainSVC)
		if err = eventbus.GetDefaultSVC().RegisterConsumer(
			nameServer,
			consts.RMQTopicWorkflowExecutor,
			consts.RMQConsumeGroupWorkflowExecutor,
			handler,
		); err != nil {
			return nil, fmt.Errorf("register workflow executor consumer failed: %w", err)
		}
	}

	SVC.DomainSVC = workflowDomainSVC
	SVC.ImageX = components.ImageX
	SVC.TosClient = components.Tos
	SVC.IDGenerator = components.IDGen

	err = SVC.InitNodeIconURLCache(context.Background())
	if err != nil {
		return nil, err
	}

	return SVC, nil
}
