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
	"time"
)

var config struct {
	// the interval to check with fleet
	interval time.Duration
	// the fleet socket
	fleetSocket string
	// the file to save the configuration to
	prometheusFile string
	// a dry-run flag
	dryRun bool
	// include all nodes, even those not defined in a job
	includeAll bool
	// the job name of the default group
	defaultJobName string
	// the default port to use for default jobs
	defaultPort int
	// a array of jobs to perform
	jobs *Jobs
}

func init() {
	config.jobs = new(Jobs)
	config.jobs.items = make([]Job, 0)
	flag.DurationVar(&config.interval, "interval", time.Duration(10)*time.Second, "the interval to check with fleet for machines")
	flag.StringVar(&config.fleetSocket, "socket", "unix://var/run/fleet.sock", "the path to the fleet api socket")
	flag.StringVar(&config.prometheusFile, "config", "/etc/prometheus/targets.d/nodes.yaml", "the location to write the nodes configuration")
	flag.StringVar(&config.defaultJobName, "group", "nodes", "the job name of the default group, i.e. those hosts not matched by a tag")
	flag.IntVar(&config.defaultPort, "port", 9100, "the port to use for machines which have been placed into the default group")
	flag.BoolVar(&config.includeAll, "all", false, "include all nodes, even those not matched by a job spec; these will be placed into the default group")
	flag.BoolVar(&config.dryRun, "dryrun", false, "perform a dry run and display the output to screen")
	flag.Var(config.jobs, "job", "add a job to group the machines (i.e. 'name;tag=value;port[;labels]')")
}
