[![Build Status](https://travis-ci.org/gambol99/prometheus-fleet.svg?branch=master)](https://travis-ci.org/gambol99/prometheus-fleet)
[![GoDoc](http://godoc.org/github.com/gambol99/prometheus-fleet?status.png)](http://godoc.org/github.com/gambol99/prometheus-fleet)

### **Prometheus Fleet**
---
The service is used to generate one or more prometheus endpoints from machines in a fleet (CoreOS) cluster. You can slice and dice the machines up by fleet metadata and generate various prometheus jobs from them.

```shell
    Usage of bin/prometheus-fleet:
      -all=false: include all nodes, even those not matched by a job spec; these will be placed into the default group
      -alsologtostderr=false: log to standard error as well as files
      -config="/etc/prometheus/targets.d/nodes.yaml": the location to write the nodes configuration
      -dryrun=false: perform a dry run and display the output to screen
      -group="nodes": the job name of the default group, i.e. those hosts not matched by a tag
      -interval=10s: the interval to check with fleet for machines
      -job=jobs: 0: add a job to group the machines (i.e. 'name;tag=value;port[;labels]')
      -json=false: produce the targets file in json rather than default yaml
      -log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
      -log_dir="": If non-empty, write log files in this directory
      -logtostderr=false: log to standard error instead of files
      -port=9100: the port to use for machines which have been placed into the default group
      -socket="unix://var/run/fleet.sock": the path to the fleet api socket
      -stderrthreshold=0: logs at or above this threshold go to stderr
      -v=0: log level for V logs
      -vmodule=: comma-separated list of pattern=N settings for file-filtered logging
```

#### **Jobs Usage**
----

A job in this context provides a means to slice up the machines in fleet into prometheus jobs i.e. we have a bunch of compute boxes which have the fleet metadata role=compute which are running a metrics endpoint on port 9100. The syntax for a job is as such; *-job  <JOB_NAME>:TAG:PORT[;LABELS] [-job ...]* . Note: multiple jobs can be defined by repeating the -job option on the command line.

```shell
 <JOB_NAME>:TAG:PORT[;LABELS] ...
 compute;role=compute;9100;role=compute,region=eu-west-1
```

The above jobs would produce something like;

```YAML
 - targets: ['NODE:9100', 'NODE:9100']
   labels:
     job: 'compute'
	 role: 'compute'
     region: 'eu-west-1'
```

Or taking the example metadata below;

``` shell
[jest@starfury ~]$ bin/prometheus-fleet -job='compute;role=compute;9100' -job=''ceph_store;role=ceph_store;9100' -all
```

```shell
[jest@starfury ~]$ fleetctl --endpoint=https://127.0.0.1:2379 list-machines
MACHINE		IP		METADATA
10137a9a...	10.50.1.79	env=prod,private_ipv4=10.50.1.79,region=eu-west-1,role=kubernetes
1ca43013...	10.50.12.200	env=prod,private_ipv4=10.50.12.200,region=eu-west-1,role=etcd
5d13b10c...	10.50.0.248	env=prod,private_ipv4=10.50.0.248,region=eu-west-1,role=kubernetes
6d9ae038...	10.50.21.100	env=prod,private_ipv4=10.50.21.100,region=eu-west-1,role=ceph_store
71e1befe...	10.50.11.200	env=prod,private_ipv4=10.50.11.200,region=eu-west-1,role=etcd
866f610f...	10.50.2.213	env=prod,private_ipv4=10.50.2.213,region=eu-west-1,role=kubernetes
c9994868...	10.50.22.100	env=prod,private_ipv4=10.50.22.100,region=eu-west-1,role=ceph_store
d5270139...	10.50.20.100	env=prod,private_ipv4=10.50.20.100,region=eu-west-1,role=ceph_store
faca46da...	10.50.10.200	env=prod,private_ipv4=10.50.10.200,region=eu-west-1,role=etcd
```

Would produce the following prometheus targets;

```YAML
 - targets: ['10.50.2.213:9100', '10.50.0.248:9100', '10.50.1.79:9100']
   labels:
     job: 'compute'

 - targets: ['10.50.20.100:9100', '10.50.21.100:9100', '10.50.22.100:9100']
   labels:
     job: 'ceph_store'

 - targets: ['10.50.10.200:9100', '10.50.11.200:9100', '10.50.12.200:9100']
   labels:
     job: 'nodes'

```

#### **Example Usage**:
----

Vairous, but I've deployed  prometheus server within [kubernetes](https://github.com/kubernetes/kubernetes), thus the deployment pattern for me is a pod containing prometheus and a collection of endpoint discovery containers which share a filesystem via the empthPath volume. The discovery containers simply write their *.yaml files to /etc/prometheus/targets.d which are picked up on intervals by the services [file discovery](http://prometheus.io/blog/2015/06/01/advanced-service-discovery/) plugin.

```YAML
#
#   Date: 2015-07-20 16:46:35 +0100 (Mon, 20 Jul 2015)
#
#  vim:ts=2:sw=2:et
#
---
apiVersion: v1
kind: ReplicationController
metadata:
  name: prometheus
spec:
  replicas: 1
  selector:
    name: prometheus
  template:
    metadata:
      labels:
        name: prometheus
    spec:
      containers:
      - name: prometheus-k8s
        image: gambol99/prometheus-k8s
        args:
          - -config=/etc/prometheus/targets.d
          - -bearer-token-file=/etc/tokens/node-register.token
          - -api=10.101.0.1
          - -api-protocol=https
          - -insecure=true
          - -logtostderr=true
          - -v=3
          - -nodes=false
        volumeMounts:
        - name: targets
          mountPath: /etc/prometheus/targets.d
        - name: tokens
          mountPath: /etc/tokens/node-register.token
      - name: prometheus-fleet
        image: gambol99/prometheus-fleet:0.0.1
        args:
          - -config=/etc/prometheus/targets.d/fleet-nodes.yml
          - -job=compute;role=kubernetes;9100
          - -job=etcd;role=etcd;9100
          - -job=ceph_store;role=ceph_store;9100
          - -job=ceph_monitor;role=ceph_monitor;9100
          - -all
          - -logtostderr=true
          - -v=3
        volumeMounts:
        - name: fleet
          mountPath: /var/run/fleet.sock
        - name: targets
          mountPath: /etc/prometheus/targets.d
      - name: prometheus
        image: gambol99/prometheus
        ports:
        - containerPort: 9090
        args:
          - -config.file=/etc/prometheus/prometheus.yml
          - -storage.local.path=/prometheus
          - -web.console.libraries=/etc/prometheus/console_libraries
          - -web.console.templates=/etc/prometheus/console
        volumeMounts:
        - name: targets
          mountPath: /etc/prometheus/targets.d
      imagePullPolicy: Always
      volumes:
      - name: tokens
        hostPath:
          path: /run/kube-kubelet/node-register.token
      - name: targets
        source:
          emptyDir: {}
      - name: fleet
        hostPath:
          path: /var/run/fleet.sock

```

#### **Status / Todo List**
----

>- Need to add the additional labels support into a job spec
- Need to add the ability to filter by multiple tags

#### **Contributing**
---

>  - Fork it
 - Create your feature branch (git checkout -b my-new-feature)
 - Commit your changes (git commit -am 'Add some feature')
 - Push to the branch (git push origin my-new-feature)
 - Create new Pull Request
 - If applicable, update the README.md
