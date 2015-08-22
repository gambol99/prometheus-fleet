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
	"fmt"
	"strconv"
	"regexp"
)

var (
	jobRegex = regexp.MustCompile("^([a-zA-Z0-9]*);([a-zA-Z0-9]+)=([a-zA-Z0-9]+);([0-9]{1,5})")
)

func (r Job) String() string {
	return fmt.Sprintf("name: %s, labels: %V, targets: %v", r.name, r.labels, r.targets)
}

// Jobs ... is an array of job items to perform
type Jobs struct {
	// an array of jobs
	items []Job
}

// Size ... return the number of jobs
func (r Jobs) Size() int {
	return len(r.items)
}

func (r Jobs) String() string {
	return fmt.Sprintf("jobs: %d", len(r.items))
}

func (r *Jobs) Set(value string) error {
	var err error
	// check: ensure it matches the regex for a job spec
	if !jobRegex.MatchString(value) {
		return fmt.Errorf("invalid job specification, should be <NAME>;TAG=<VALUE>;<PORT> (e.g compute;role=kube;9010)")
	}

	matches := jobRegex.FindAllStringSubmatch(value, -1)
	job := Job{
		name:    matches[0][1],
		tagName: matches[0][2],
		tagValue: matches[0][3],
		labels:  make(map[string]string, 0),
	}

	// step: parse the port and check it
	job.port, err = strconv.Atoi(matches[0][4])
	if err != nil {
		return fmt.Errorf("invalid port specification in job: %s", value)
	}

	if job.port <= 0 || job.port >= 65536 {
		return fmt.Errorf("invalid port specification in job: %s, should be between 0-65536", value)
	}
	r.items = append(r.items, job)

	return nil
}

