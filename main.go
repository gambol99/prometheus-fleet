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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

func main() {
	flag.Parse()

	// step: ensure we have some jobs to process
	if config.jobs.Size() <= 0 && !config.includeAll {
		glog.Infof("Zero job specifications have been defined and the includeAll option is false, nothing to do!")
		os.Exit(0)
	}

	glog.Infof("Starting the %s service, version: %s, git+sha: %s", Prog, Version, GitSha)

	// step: we grab a client to fleet api
	client, err := newFleetClient()
	if err != nil {
		glog.Errorf("Failed to create a fleet client on socket: %s, error: %s", config.fleetSocket, err)
		os.Exit(1)
	}

	// step: setup the termination signal
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		var targets TargetGroups

		// step: grab a list of machines from fleet
		machines, err := client.GetMachines()
		if err != nil {
			glog.Errorf("Failed to retrieve a list of machines from fleet, error: %s", err)
			goto NEXT_LOOP
		}

		// step: if we have zero jobs; we don't need to filter, just place everything into the default group
		switch config.jobs.Size() {
		case 0:
			targets = produceDefaultTargets(machines)
		default:
			targets = produceFilteredTargets(config.jobs, machines)
		}

		// step: write the targets to file
		if err := persistTargets(targets); err != nil {
			glog.Errorf("Failed to persist the targets to file: %s, error: %s", config.prometheusFile, err)
		} else {
			glog.V(4).Infof("Successfully wrote the targets (%d) to file: %s", len(targets.targets), config.prometheusFile)
		}

	NEXT_LOOP:
		<-time.After(config.interval)
	}
}

func produceDefaultTargets(machines []*Machine) TargetGroups {
	var group TargetGroups
	target := group.AddTarget(config.defaultJobName)
	// step: iterate and place all the machine inside it
	for _, machine := range machines {
		target.targets = append(target.targets, fmt.Sprintf("%s:%d", machine.name, config.defaultPort))
	}
	return group
}

func produceFilteredTargets(jobs *Jobs, machines []*Machine) TargetGroups {
	var groups TargetGroups
	// we use a map to keep track of the machines which have been added to a group
	added := make(map[string]bool, 0)

	// step: for each of the jobs we need to produce the targets
	for _, job := range jobs.items {
		target := groups.AddTarget(job.name)
		glog.V(6).Infof("Processing the %d machines against the job: %s", len(machines), job)
		// step: iterate the machines and find any one that matches our metadata
		for _, machine := range machines {
			// step: does it have a matching tag
			if value, found := machine.metadata[job.tagName]; !found {
				glog.V(6).Infof("The machine: %s does not hae matching tag: %s in it's metadata", machine, job.tagName)
				continue
			} else if value == job.tagValue {
				// we can append this machine to the list of targets
				target.targets = append(target.targets, fmt.Sprintf("%s:%d", machine.name, job.port))
				// keep the tracker updates
				added[machine.name] = true
			}
		}
	}

	// step: are we adding 'all' the machines, even those not matching.
	// We check if we have left overs, then then iterate the machines and find any machine which has not been
	// added already and add to the default group
	if config.includeAll {
		if len(added) != len(machines) {
			// we have machines that haven't been added
			target := groups.AddTarget(config.defaultJobName)
			for _, machine := range machines {
				if _, found := added[machine.name]; !found {
					target.targets = append(target.targets, fmt.Sprintf("%s:%d", machine.name, config.defaultPort))
				}
			}
		}
	}

	return groups
}

func persistTargets(targets TargetGroups) error {
	// step: first we encode the structure
	content, err := encodeTargets(targets)
	if err != nil {
		glog.Errorf("Failed to encode the target stucture into yaml, error: %s", err)
		return err
	}

	// step: attempt to write the file
	if config.dryRun {
		fmt.Printf("%s", content)
		return nil
	}

	err = ioutil.WriteFile(config.prometheusFile, []byte(content), os.FileMode(0444))
	if err != nil {
		glog.Errorf("Failed to write to file: '%s', error: %s", config.prometheusFile, err)
		return err
	}

	return nil
}

func encodeTargets(group TargetGroups) (string, error) {
	output, err := yaml.Marshal(group.targets)
	if err != nil {
		glog.Errorf("Failed to marshall the structure to json, %s, error: %s", group, err)
		return "", fmt.Errorf("marshalling failure, data: %V, error: %s", group, err)
	}

	return string(output), nil
}
