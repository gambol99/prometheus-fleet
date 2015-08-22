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

func TestJobsSet(t *testing.T) {
	jobs := &Jobs{}
	assert.Nil(t, jobs.Set("come;role=value;8000"))
	assert.NotNil(t, jobs.Set("come;8009"))
	assert.NotNil(t, jobs.Set("come;role=kube;"))
	assert.NotNil(t, jobs.Set("come;role=kube;jki"))
	assert.Nil(t, jobs.Set("\"come;role=value;8000\""))
}

func TestJobsSize(t *testing.T) {
	jobs := &Jobs{}
	assert.Nil(t, jobs.Set("come;role=value;8000"))
	assert.Nil(t, jobs.Set("come;role=value;8000"))
	assert.Nil(t, jobs.Set("come;role=value;8000"))
	assert.Nil(t, jobs.Set("come;role=value;8000"))
	assert.Equal(t, 4, jobs.Size())
}

func TestJobsPort(t *testing.T) {
	jobs := &Jobs{}
	assert.NotNil(t, jobs.Set("come;role=value;80000"))
}

func TestJobs(t *testing.T) {
	jobs := &Jobs{}
	assert.Nil(t, jobs.Set("come;role=value;8000"))
	assert.Equal(t, 1, jobs.Size())
	assert.Equal(t, "come", jobs.items[0].name)
	assert.Equal(t, 8000, jobs.items[0].port)
	assert.Equal(t, "role", jobs.items[0].tagName)
	assert.Equal(t, "value", jobs.items[0].tagValue)
}
