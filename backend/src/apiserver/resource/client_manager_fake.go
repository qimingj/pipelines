// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	workflowclient "github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	"github.com/golang/glog"
	"github.com/googleprivate/ml/backend/src/apiserver/storage"
	"github.com/googleprivate/ml/backend/src/common/util"
	scheduledworkflowclient "github.com/googleprivate/ml/backend/src/crd/pkg/client/clientset/versioned/typed/scheduledworkflow/v1alpha1"
)

const (
	DefaultFakeUUID = "123e4567-e89b-12d3-a456-426655440000"
)

type FakeClientManager struct {
	db                          *storage.DB
	experimentStore             storage.ExperimentStoreInterface
	pipelineStore               storage.PipelineStoreInterface
	jobStore                    storage.JobStoreInterface
	runStore                    storage.RunStoreInterface
	resourceReferenceStore      storage.ResourceReferenceStoreInterface
	objectStore                 storage.ObjectStoreInterface
	workflowClientFake          *storage.FakeWorkflowClient
	scheduledWorkflowClientFake *FakeScheduledWorkflowClient
	time                        util.TimeInterface
	uuid                        util.UUIDGeneratorInterface
}

func NewFakeClientManager(time util.TimeInterface, uuid util.UUIDGeneratorInterface) (
	*FakeClientManager, error) {

	if time == nil {
		glog.Fatalf("The time parameter must not be null.") // Must never happen
	}

	if uuid == nil {
		glog.Fatalf("The UUID generator must not be null.") // Must never happen
	}

	// Initialize GORM
	db, err := storage.NewFakeDb()
	if err != nil {
		return nil, err
	}

	return &FakeClientManager{
		db:                          db,
		experimentStore:             storage.NewExperimentStore(db, time, uuid),
		pipelineStore:               storage.NewPipelineStore(db, time, uuid),
		jobStore:                    storage.NewJobStore(db, time),
		runStore:                    storage.NewRunStore(db, time),
		workflowClientFake:          storage.NewWorkflowClientFake(),
		resourceReferenceStore:      storage.NewResourceReferenceStore(db),
		objectStore:                 storage.NewFakeObjectStore(),
		scheduledWorkflowClientFake: NewScheduledWorkflowClientFake(),
		time:                        time,
		uuid:                        uuid,
	}, nil
}

func NewFakeClientManagerOrFatal(time util.TimeInterface) *FakeClientManager {
	uuid := util.NewFakeUUIDGeneratorOrFatal(DefaultFakeUUID, nil)
	fakeStore, err := NewFakeClientManager(time, uuid)
	if err != nil {
		glog.Fatalf("The fake store doesn't create successfully. Fail fast.")
	}
	return fakeStore
}

func (f *FakeClientManager) ExperimentStore() storage.ExperimentStoreInterface {
	return f.experimentStore
}

func (f *FakeClientManager) PipelineStore() storage.PipelineStoreInterface {
	return f.pipelineStore
}

func (f *FakeClientManager) ObjectStore() storage.ObjectStoreInterface {
	return f.objectStore
}

func (f *FakeClientManager) Time() util.TimeInterface {
	return f.time
}

func (f *FakeClientManager) UUID() util.UUIDGeneratorInterface {
	return f.uuid
}

func (f *FakeClientManager) DB() *storage.DB {
	return f.db
}

func (f *FakeClientManager) Workflow() workflowclient.WorkflowInterface {
	return f.workflowClientFake
}

func (f *FakeClientManager) JobStore() storage.JobStoreInterface {
	return f.jobStore
}

func (f *FakeClientManager) RunStore() storage.RunStoreInterface {
	return f.runStore
}

func (f *FakeClientManager) ResourceReferenceStore() storage.ResourceReferenceStoreInterface {
	return f.resourceReferenceStore
}

func (f *FakeClientManager) ScheduledWorkflow() scheduledworkflowclient.ScheduledWorkflowInterface {
	return f.scheduledWorkflowClientFake
}

func (f *FakeClientManager) Close() error {
	return f.db.Close()
}