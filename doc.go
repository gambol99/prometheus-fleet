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

import "fmt"

// Job ...A job in this context provides a means to slice up the machines in fleet into prometheus jobs
// i.e. we have a bunch of compute boxes which have the fleet metadata role=compute which are running
// a metrics endpoint on port 9100
//
// <JOB_NAME>:TAG:PORT[;LABELS] ...
// compute;role=compute;9100;role=compute,region=eu-west-1
//
// - targets: ['NODE:9100', 'NODE:9100']
//	 labels:
//     job: 'compute'
//	   role: 'compute'
//     region: 'eu-west-1'
//
type Job struct {
	// the name of the job
	name string
	// the tag we are looking for
	tagName string
	// the value of the tag
	tagValue string
	// the labels associated to the job
	labels map[string]string
	// the port the target should be running on
	port int
	// a list of targets linked to the job
	targets []string
}

// Machine ... the structure of a machine from fleet
type Machine struct {
	// the name of the machine - normally the ip address
	name string
	// the metadata associated to the machine
	metadata map[string]string
}

// TargetGroups ... is just a array of target and use more an a help wrapper than anything else
type TargetGroups struct {
	// an array of targets
	targets []*Target
}

// Target ... is the structure for prometheus endpoints
type Target struct {
	// the array of hosts within this target
	Targets []string `yaml:"targets""`
	// the labels associated to these targets
	Labels map[string]string `yaml:"labels"`
}

// String ... returns a string representation of a Machine
func (r Machine) String() string {
	return fmt.Sprintf("machine: %s, metadata: %s", r.name, r.metadata)
}

// Size ... the number of targets in the Target Group
func (r TargetGroups) Size() int {
	return len(r.targets)
}

// AddTarget ... create a new target for you, add to the group and return it
func (r *TargetGroups) AddTarget(name string) *Target {
	target := &Target{
		Targets: make([]string, 0),
		Labels: map[string]string{
			"job": name,
		},
	}
	r.targets = append(r.targets, target)
	return target
}
