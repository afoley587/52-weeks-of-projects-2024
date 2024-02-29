# Building A Custom Kubernetes Operator
## How The Operator-SDK Makes It Wicked Easy

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
## 4. Define our reconciliation logic
## 5. Use the Operator SDK to build and deploy our CRDs and operator
## make cahnges in ping_types.go

## make changes to reconcile

make manifests
make install
make run

