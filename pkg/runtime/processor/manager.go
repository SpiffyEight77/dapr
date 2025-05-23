/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package processor

import (
	"context"
	"fmt"

	"github.com/dapr/durabletask-go/backend"

	"github.com/dapr/components-contrib/bindings"
	componentsapi "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	subapi "github.com/dapr/dapr/pkg/apis/subscriptions/v2alpha1"
	"github.com/dapr/dapr/pkg/runtime/meta"
	rtpubsub "github.com/dapr/dapr/pkg/runtime/pubsub"
)

// manager implements the life cycle events of a component category.
type manager interface {
	Init(context.Context, componentsapi.Component) error
	Close(componentsapi.Component) error
}

type StateManager interface {
	ActorStateStoreName() (string, bool)
	manager
}

type SecretManager interface {
	ProcessResource(context.Context, meta.Resource) (bool, string)
	manager
}

type SubscribeManager interface {
	InitProgramaticSubscriptions(context.Context) error
	StartAppSubscriptions() error
	StopAppSubscriptions()
	StopAllSubscriptionsForever()
	ReloadDeclaredAppSubscription(name, pubsubName string) error
	StartStreamerSubscription(sub *subapi.Subscription, connectionID rtpubsub.ConnectionID) error
	StopStreamerSubscription(sub *subapi.Subscription, connectionID rtpubsub.ConnectionID)
	ReloadPubSub(string) error
	StopPubSub(string)
}

type BindingManager interface {
	SendToOutputBinding(context.Context, string, *bindings.InvokeRequest) (*bindings.InvokeResponse, error)

	StartReadingFromBindings(context.Context) error
	StopReadingFromBindings(forever bool)
	manager
}

type WorkflowBackendManager interface {
	Backend() (backend.Backend, bool)
}

func (p *Processor) managerFromComp(comp componentsapi.Component) (manager, error) {
	category := p.category(comp)
	m, ok := p.managers[category]
	if !ok {
		return nil, fmt.Errorf("unknown component category: %q", category)
	}
	return m, nil
}

func (p *Processor) State() StateManager {
	return p.state
}

func (p *Processor) Secret() SecretManager {
	return p.secret
}

func (p *Processor) Binding() BindingManager {
	return p.binding
}

func (p *Processor) WorkflowBackend() WorkflowBackendManager {
	return p.workflowBackend
}

func (p *Processor) Subscriber() SubscribeManager {
	return p.subscriber
}
