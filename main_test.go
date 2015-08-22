/*
Copyright 2014 Rohith All rights reserved.

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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFakeTargets(t *testing.T) TargetGroups {
	client, _ := newFakeFleetClient()
	machines, _ := client.GetMachines()
	targets := produceDefaultTargets(machines)
	if targets.Size() == 0 {
		assert.Fail(t, "unable to produce a fake target group")
	}
	return targets
}

func TestNewFakeFleetClient(t *testing.T) {
	client, _ := newFakeFleetClient()
	assert.NotNil(t, client)
}

func TestGetMachines(t *testing.T) {
	client, _ := newFakeFleetClient()
	machines, err := client.GetMachines()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(machines))
}

func TestDefaultGroup(t *testing.T) {
	client, _ := newFakeFleetClient()
	machines, _ := client.GetMachines()
	targets := produceDefaultTargets(machines)
	assert.NotNil(t, targets)
	assert.Equal(t, 1, targets.Size())
	assert.Equal(t, "nodes", targets.targets[0].labels["job"])
	assert.Equal(t, len(machines), len(targets.targets[0].targets))
}

func TestFilterGroup(t *testing.T) {
	var jobs Jobs
	client, _ := newFakeFleetClient()
	machines, _ := client.GetMachines()
	jobs.Set("compute;role=kube;9001")
	jobs.Set("etcd;role=etcd;9001")
	targets := produceFilteredTargets(&jobs, machines)
	assert.NotNil(t, targets)
	assert.Equal(t, 2, targets.Size())
	assert.Equal(t, 3, len(targets.targets[0].targets))
	assert.Equal(t, 3, len(targets.targets[1].targets))
	assert.Equal(t, "compute", targets.targets[0].labels["job"])
	assert.Equal(t, "etcd", targets.targets[1].labels["job"])
}

func TestFilterWithDefaults(t *testing.T) {
	var jobs Jobs
	client, _ := newFakeFleetClient()
	machines, _ := client.GetMachines()
	config.includeAll = true
	jobs.Set("compute;role=kube;9001")
	targets := produceFilteredTargets(&jobs, machines)
	assert.NotNil(t, targets)
	assert.Equal(t, 2, targets.Size())
	assert.Equal(t, 3, len(targets.targets[0].targets))
	assert.Equal(t, 3, len(targets.targets[1].targets))
	assert.Equal(t, "compute", targets.targets[0].labels["job"])
	assert.Equal(t, config.defaultJobName, targets.targets[1].labels["job"])
}

func TestEncodeTargets(t *testing.T) {
	content, err := encodeTargets(createFakeTargets(t))
	assert.Nil(t, err)
	assert.NotEmpty(t, content)
	assert.NotEqual(t, 0, len(content))
}
