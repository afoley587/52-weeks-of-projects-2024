# Building A Custom Kubernetes Operator
## How The Operator-SDK Makes It Wicked Easy

![Thumbnail](./images/thumbnail.png)

If you're in the kubernetes space, either professionally or personally, 
you've probably heard of the term `operator` or `controller`. You've 
probably heard people say "Install <XYZ> with the operator" or 
"I just built a custom controller to do <ABC>". What does that really mean?
Well, let's first start with some basics.

## Operators Vs Controllers
There's a lot of confusion vs what is a controller and what is an operator.
To be a bit brief about it, in my opinion, operators are a specialized 
controller. They try to embody more workload-specific knowledge into
the controller. 

So, next question, what's a controller? Well, when you create some kubernetes
object, you say "This is my desired state". If we make a Deployment object and
set the number of replicas to 3, we are telling the replica set controller to
watch this object and ensure there are always 3 pods running to service the
deployment's requested state. The process that controllers use to do this
is something called reconciliation (keep this in mind for later). In short, 
controllers use a reconciliation process to make sure our desired state is 
always met.

So, what graduates a controller to an operator? Well, there's a good bit
of debate and discussion on that:

1. Here is one of many [GitHub discussions](https://github.com/kubeflow/training-operator/issues/300)
2. A blog on [controller vs operators](https://joshrosso.com/docs/2019/2019-10-13-controllers-and-operators/)
3. It's even made its way to [StackOverflow](https://stackoverflow.com/questions/47848258/what-is-the-difference-between-a-kubernetes-controller-and-a-kubernetes-operator)

In our case, we will say that a controller becomes an operator when it 
satisfies the criteria from number 2 above (Credit to https://joshrosso.com):

1. Contains workload-specific knowledge
2. Manages workload lifecycle
3. Offers a CRD

I believe that what we're building today matches that criteria...
so let's say we're building an operator today!

## Operator-SDK

To build out our operator, we will be using the 
[operator-sdk](https://sdk.operatorframework.io/). This SDK does a lot
of the heavy lifting for us and includes a good bit of skeleton code
for our operators/controllers. It provides utilities to then generate
your CRD's from your source code, install resources into your cluster,
and run/test the operators. Installation is easy. If you're like me (and
on a mac), you can install with [homebrew](https://brew.sh/):

```shell
$ brew install operator-sdk
$ operator-sdk version
operator-sdk version: "v1.31.0", commit: "e67da35ef4fff3e471a208904b2a142b27ae32b1", kubernetes version: "v1.26.0", go version: "go1.20.6", GOOS: "darwin", GOARCH: "arm64"
```

However, they support a bunch of different installation methods described
[here](https://sdk.operatorframework.io/docs/installation/) in their documentation.

## Setting Up Our Environment
We will be using [minikube](https://minikube.sigs.k8s.io/docs/start/) 
as our local development environment. If you don't have it installed, 
follow the link above to download/install it. Let's now start our
minikube cluster:

```shell
$ minikube start --cpus=2 --memory=2048
```



## Our Project
So, what are we really building today? Well, we're going to build an operator
which responds to certain custom resource definitions (CRDs). Our CRDs will
define ping checks to be sent out to a hostname with a certain number of attempts.
In short, we will write an operator to:

1. Create a Kubernetes job when a new object of type `Ping` is requested
2. This job will run `ping -c <number_of_attempts> <hostname>`

The steps in making this happen are as following:

1. Use the Operator SDK to create a new project for us
2. Use the Operator SDK to create new CRD and operator skeletons for us
3. Define our Kubernetes CRD Schema
4. Define our reconciliation logic
5. Use the Operator SDK to build and deploy our CRDs and operator

## 1. Use the Operator SDK to create a new project for us

Let's first start by creating a new directory for our operators workspace 
and then use `operator-sdk init` to initialize our go workspace and download
any dependencies:

```shell
$ mkdir ping-operator && cd ping-operator
$ # Change your domain / repo accordingly!
$ operator-sdk init \
    --domain=engineeringwithalex.io \
    --repo=github.com/afoley587/52-weeks-of-projects-2023/08-kubernetes-controllers/ping-operator
WARN[0000] the platform of this environment (darwin/arm64) is not suppported by kustomize v3 (v3.8.7) which is used in this scaffold. You will be unable to download a binary for the kustomize version supported and used by this plugin. The currently supported platforms are: ["linux/amd64" "linux/arm64" "darwin/amd64"] 
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.14.1
go: downloading sigs.k8s.io/controller-runtime v0.14.1
go: downloading k8s.io/apimachinery v0.26.0
go: downloading k8s.io/client-go v0.26.0
go: downloading github.com/go-logr/logr v1.2.3
go: downloading k8s.io/utils v0.0.0-20221128185143-99ec85e7a448
go: downloading k8s.io/klog/v2 v2.80.1
go: downloading k8s.io/component-base v0.26.0
go: downloading sigs.k8s.io/structured-merge-diff/v4 v4.2.3
go: downloading github.com/google/gofuzz v1.1.0
go: downloading github.com/evanphx/json-patch/v5 v5.6.0
go: downloading gomodules.xyz/jsonpatch/v2 v2.2.0
.
.
.
.
Next: define a resource with:
$ operator-sdk create api
```

## Use the Operator SDK to create new CRD and operator skeletons for us

We will be creating a new API for kubernetes to use. To do that, we need a few
pieces of information:

* API Group
* API Version
* Kind

If we look at the stateful set example, the group is `apps`, the
version is `v1`, and the kind is `StatefulSet`. This is how Kubernetes
will organize the API that we are adding.

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: ...
  namespace: ...
```

So, let's create our API:

```shell
$ operator-sdk create api \
    --group monitors \
    --version v1beta1 \
    --kind Ping \
    --namespaced \
    --controller \
    --resource \
    --make
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1beta1/ping_types.go
controllers/ping_controller.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
mkdir -p ...
test -s .../controller-gen && .../controller-gen --version | grep -q v0.11.1 || \
        GOBIN=... go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
.../controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

With the above command, we asked the operator SDK to make us a new API with
the following flags:

* `--group monitors` - Set the API group to `monitors`
* `--version v1beta1` - Set the API version to `v1beta1`
* `--kind monitors` - Create a new object of kind `Ping`
* `--namespaced` - This resource will be namespaced
* `--controller` - Generate the controller without prompting us
* `--resource` - Generate the resource without prompting us
* `--make` - Run make generate after generating files

Let's do an `ls` on a few directories to get a better handle on what happened:

```shell
$ ls api/v1beta1
groupversion_info.go  ping_types.go  zz_generated.deepcopy.go
$ ls controllers 
ping_controller.go  suite_test.go
```

Wow - the `operator-sdk` went and created a ton of files for us. Of particular
importance are `api/v1beta1/ping_types.go` where we will specifiy our CRD schema
and `controllers/ping_controller.go` where we will define our reconciliation logic!

## 3. Define our Kubernetes CRD Schema

As mentioned above, our CRD schema will go into `api/v1beta1/ping_types.go`.

If you look in the default file, you will see something that looks like this:

```golang
// PingSpec defines the desired state of Ping
type PingSpec struct {
        // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
        // Important: Run "make" to regenerate code after modifying this file

        // Foo is an example field of Ping. Edit ping_types.go to remove/update
        Foo string `json:"foo,omitempty"`
}

// PingStatus defines the observed state of Ping
type PingStatus struct {
        // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
        // Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Ping is the Schema for the pings API
type Ping struct {
        metav1.TypeMeta   `json:",inline"`
        metav1.ObjectMeta `json:"metadata,omitempty"`

        Spec   PingSpec   `json:"spec,omitempty"`
        Status PingStatus `json:"status,omitempty"`
}
```

This is the default setup for new resources. It's bare! For the sake
of this post, we will only be updating the `PingSpec` and not the `PingStatus`.

So, from a high-level, what is `PingSpec`, `PingStatus`, and `Ping`? Well, 
`PingSpec` is the specification that users need to define when requesting 
a `Ping` object. Let's look at an example for a pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec: ### THIS STARTS THE SPEC!
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

We can see that the pod spec has `containers` which, in itself, has another
schema. We won't be nesting any objects in our CRD. Again, we just want the user
to be able to define the hostname to send a ping to and the number of attempts
that the ping should perform. So, a user should be able to do something like:

```yaml
apiVersion: monitors.engineeringwithalex.io/v1beta1
kind: Ping
metadata:
  labels:
    app.kubernetes.io/name: ping
    app.kubernetes.io/instance: ping-sample
    app.kubernetes.io/part-of: ping-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: ping-operator
  name: ping-sample
spec: ### THIS STARTS THE SPEC!
  hostname: "www.google.com"
  attempts: 1
```

To accomplish that, we need to update our Go code:

```golang
// PingSpec defines the desired state of Ping
type PingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	/**
	Need to add the Hostname to the spec
	**/
	Hostname string `json:"hostname,omitempty"`
	Attempts int    `json:"attempts,omitempty"`
}

// PingStatus defines the observed state of Ping
type PingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Ping is the Schema for the pings API
type Ping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PingSpec   `json:"spec,omitempty"`
	Status PingStatus `json:"status,omitempty"`
}
```

This adds a string for `hostname` and an integer for `attempts`.

The `PingStatus` would report the status of the object. Again, we won't
be updating that here, but it might include some data about your ping such
as:

* Was the ping successful?
* Is the job finished?
* etc.

Finally, the `Ping` object combines everything. It combines the basic type
metadata (such as group, API version, and kind), the object metadata (such as name, 
namespace, etc.), the `PingSpec`, and the `PingStatus`!

At this point, we can run `make manifests` command:

```shell
$ make manifests
test -s ...08-kubernetes-controllers/ping-operator/bin/controller-gen && ...08-kubernetes-controllers/ping-operator/bin/controller-gen --version | grep -q v0.11.1 || \
        GOBIN=...08-kubernetes-controllers/ping-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
...08-kubernetes-controllers/ping-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

We should now see some manifests in `config/crd/bases`:

```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.11.1
  creationTimestamp: null
  name: pings.monitors.engineeringwithalex.io
spec:
  group: monitors.engineeringwithalex.io
  names:
    kind: Ping
    listKind: PingList
    plural: pings
    singular: ping
  scope: Namespaced
  versions:
  - name: v1beta1
    schema:
      openAPIV3Schema:
        description: Ping is the Schema for the pings API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: PingSpec defines the desired state of Ping
            properties:
              attempts:
                type: integer
              hostname:
                description: '* Need to add the Hostname to the spec *'
                type: string
            type: object
          status:
            description: PingStatus defines the observed state of Ping
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
```

## 4. Define our reconciliation logic

Our schemas and CRDs are now out of the way. But, what should
the operator do when we request a new one of these? For example, 
if I applied the yaml for a `Ping` resource - what should the operator
do? Enter the `Reconcile` function. 

The code for this lives in the `controllers/ping_controller.go` file. The
`operator-sdk` made us a nice skeleton:

```golang
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ping object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
        _ = log.FromContext(ctx)

        // TODO(user): your logic here

        return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PingReconciler) SetupWithManager(mgr ctrl.Manager) error {
        return ctrl.NewControllerManagedBy(mgr).
                For(&monitorsv1beta1.Ping{}).
                Complete(r)
}
```

But that won't do much for us. This `Reconcile` method is called 
whenever a `Ping` resource is created, updated, or deleted. 
A bunch of user data, such as name and namespace, are
then passed to this function as a `ctrl.Request`. Now, as previously noted,
when a user requests a `Ping` resource, we want our controller to spin up a new
kubernetes job to handle that request. Our function becomes:

```golang
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ping object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var ping monitorsv1beta1.Ping

	if err := r.Get(ctx, req.NamespacedName, &ping); err != nil {
		log.FromContext(ctx).Error(err, "Unable to fetch Ping")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	job, err := r.BuildJob(ping)

	if err != nil {
		log.FromContext(ctx).Error(err, "Unable to get job definition")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Create(ctx, &job); err != nil {
		log.FromContext(ctx).Error(err, "Unable to create job")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
}
```

We first try to get the `Ping` object from the kubernetes API. If, for example, 
resource type didn't exist (i.e. CRDs were never applied), this would throw an
error. We wouldn't be able to handle the request in that case. Next, we call
this custom `BuildJob` function. From a high level, `BuildJob` will build the 
kuberetes job definition for us (see below). The definition is then returned to
the `Reconcile` function, who applies and deploys this job. If the job deployment
fails, we again return an error. 

So, what's `BuildJob`? It's definition is below:

```golang
func (r *PingReconciler) BuildJob(ping monitorsv1beta1.Ping) (batchv1.Job, error) {
	attempts := "-c" + strconv.Itoa(ping.Spec.Attempts)
	host := ping.Spec.Hostname
	j := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: batchv1.SchemeGroupVersion.String(),
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ping.Name + "-job",
			Namespace: ping.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "ping",
							Image:   "bash",
							Command: []string{"/bin/ping"},
							Args:    []string{attempts, host},
						},
					},
				},
			},
		},
	}
	return j, nil
}
```

First, we pull the spec arguments (hostname and attempts) from the `Ping`
resource that the `Reconcile` function got from the kubernetes API. We then
create a new job. We will give it the name of `<name of ping>-job` and we will
add a new container to it. The container will run the `bash` docker image with 
the command of `/bin/ping` and arguments `-c<number of attempts> <hostname>`.
So, when the `Reconcile` function builds this, we would expect to see a new
job get created. That job should create one pod which issues the ping command.

Note: We aren't implementing any [finalizers](https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/).
So, if you deploy the same `Ping` resource twice, you'll see some errors about
duplicate resources!

## 5. Use the Operator SDK to build and deploy our CRDs and operator
All of the heavy listing is done now! Let's use the `Makefile` provided
by the `operator-sdk` to build and deploy our manifests and operator!

```shell
$ make manifests                     
test -s ...08-kubernetes-controllers/ping-operator/bin/controller-gen && ...08-kubernetes-controllers/ping-operator/bin/controller-gen --version | grep -q v0.11.1 || \
        GOBIN=...08-kubernetes-controllers/ping-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
...08-kubernetes-controllers/ping-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
$ make install
test -s ...08-kubernetes-controllers/ping-operator/bin/controller-gen && ...08-kubernetes-controllers/ping-operator/bin/controller-gen --version | grep -q v0.11.1 || \
        GOBIN=...08-kubernetes-controllers/ping-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
...08-kubernetes-controllers/ping-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
...08-kubernetes-controllers/ping-operator/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/pings.monitors.engineeringwithalex.io unchanged
$ make run
test -s ...08-kubernetes-controllers/ping-operator/bin/controller-gen && ...08-kubernetes-controllers/ping-operator/bin/controller-gen --version | grep -q v0.11.1 || \
        GOBIN=...08-kubernetes-controllers/ping-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
...08-kubernetes-controllers/ping-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
...08-kubernetes-controllers/ping-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go run ./main.go
2024-02-29T15:15:56-07:00       INFO    controller-runtime.metrics      Metrics server is starting to listen    {"addr": ":8080"}
2024-02-29T15:15:56-07:00       INFO    setup   starting manager
2024-02-29T15:15:56-07:00       INFO    Starting server {"path": "/metrics", "kind": "metrics", "addr": "[::]:8080"}
2024-02-29T15:15:56-07:00       INFO    Starting server {"kind": "health probe", "addr": "[::]:8081"}
2024-02-29T15:15:56-07:00       INFO    Starting EventSource    {"controller": "ping", "controllerGroup": "monitors.engineeringwithalex.io", "controllerKind": "Ping", "source": "kind source: *v1beta1.Ping"}
2024-02-29T15:15:56-07:00       INFO    Starting Controller     {"controller": "ping", "controllerGroup": "monitors.engineeringwithalex.io", "controllerKind": "Ping"}
2024-02-29T15:15:56-07:00       INFO    Starting workers        {"controller": "ping", "controllerGroup": "monitors.engineeringwithalex.io", "controllerKind": "Ping", "worker count": 1}
```

In another terminal, let's apply a kubernetes manifest to see how our operator
reacts:

```shell
$ kubectl apply -f - <<EOF
apiVersion: monitors.engineeringwithalex.io/v1beta1
kind: Ping
metadata:
  labels:
    app.kubernetes.io/name: ping
    app.kubernetes.io/instance: ping-sample
    app.kubernetes.io/part-of: ping-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: ping-operator
  name: ping-sample
spec:
  hostname: "www.google.com"
  attempts: 1
EOF
ping.monitors.engineeringwithalex.io/ping-sample created
```

And now we can check on that object!

```shell
# First, let's look at the Ping object
$ kubectl get Ping
NAME          AGE
ping-sample   26s
# Let's look at the job it made
$ kubectl get job
NAME              COMPLETIONS   DURATION   AGE
ping-sample-job   1/1           6s         33s
# Let's look at the pod that the job made
$ kubectl get pod
NAME                    READY   STATUS      RESTARTS   AGE
ping-sample-job-hk6kq   0/1     Completed   0          42s
# Let's read the pod logs
$ kubectl logs ping-sample-job-hk6kq
PING www.google.com (142.250.72.4): 56 data bytes
64 bytes from 142.250.72.4: seq=0 ttl=62 time=18.363 ms

--- www.google.com ping statistics ---
1 packets transmitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 18.363/18.363/18.363 ms
```

Here's an image for a condensed format!

![Demo](./images/demo.png)


I hope you liked following along! As always, all code can be found
[here](https://github.com/afoley587/52-weeks-of-projects-2024/tree/main/08-kubernetes-controllers) on GitHub.

## References

* [OperatorSDK](https://sdk.operatorframework.io/)
* [Controllers Vs Operators](https://joshrosso.com/docs/2019/2019-10-13-controllers-and-operators/)
* [Building A Controller](https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/)