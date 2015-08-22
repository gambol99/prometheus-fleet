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

var (
	fleetMachines = []*Machine{
		&Machine{
			name: "10.50.20.100",
			metadata: map[string]string{
				"role": "kube",
			},
		},
		&Machine{
			name: "10.50.21.100",
			metadata: map[string]string{
				"role": "kube",
			},
		},
		&Machine{
			name: "10.50.22.100",
			metadata: map[string]string{
				"role": "kube",
			},
		},
		&Machine{
			name: "10.50.10.100",
			metadata: map[string]string{
				"role": "etcd",
			},
		},
		&Machine{
			name: "10.50.11.100",
			metadata: map[string]string{
				"role": "etcd",
			},
		},
		&Machine{
			name: "10.50.12.100",
			metadata: map[string]string{
				"role": "etcd",
			},
		},
	}
)

type fakeFleetClient struct{}

func newFakeFleetClient() (*fakeFleetClient, error) {
	return &fakeFleetClient{}, nil
}

func (r fakeFleetClient) GetMachines() ([]*Machine, error) {
	return fleetMachines, nil
}
