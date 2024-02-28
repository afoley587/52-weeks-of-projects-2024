## 
minikube start --cpus=2 --memory=2048
## IDEA: make a controller which sends pings to certain places


https://sdk.operatorframework.io/docs/installation/

$ brew install operator-sdk
$ operator-sdk version
operator-sdk version: "v1.31.0", commit: "e67da35ef4fff3e471a208904b2a142b27ae32b1", kubernetes version: "v1.26.0", go version: "go1.20.6", GOOS: "darwin", GOARCH: "arm64"
mkdir ping-operator && cd ping-operator


operator-sdk init --domain=engineeringwithalex.io --repo=github.com/afoley587/52-weeks-of-projects-2023/08-kubernetes-controllers/ping-operator
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
go: downloading k8s.io/api v0.26.0
go: downloading k8s.io/apiextensions-apiserver v0.26.0
go: downloading github.com/prometheus/client_model v0.3.0
go: downloading sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2
go: downloading golang.org/x/net v0.3.1-0.20221206200815-1e63c2f08a10
go: downloading github.com/imdario/mergo v0.3.6
go: downloading github.com/evanphx/json-patch v4.12.0+incompatible
go: downloading golang.org/x/term v0.3.0
go: downloading gopkg.in/inf.v0 v0.9.1
go: downloading sigs.k8s.io/yaml v1.3.0
go: downloading github.com/google/gnostic v0.5.7-v3refs
go: downloading k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280
go: downloading github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
go: downloading github.com/google/uuid v1.1.2
go: downloading github.com/cespare/xxhash/v2 v2.1.2
go: downloading golang.org/x/sys v0.3.0
go: downloading google.golang.org/protobuf v1.28.1
go: downloading github.com/fsnotify/fsnotify v1.6.0
go: downloading github.com/matttproud/golang_protobuf_extensions v1.0.2
go: downloading golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b
go: downloading github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
go: downloading golang.org/x/text v0.5.0
go: downloading github.com/emicklei/go-restful/v3 v3.9.0
go: downloading github.com/go-openapi/swag v0.19.14
go: downloading github.com/go-openapi/jsonreference v0.20.0
go: downloading google.golang.org/appengine v1.6.7
go: downloading github.com/go-openapi/jsonpointer v0.19.5
go: downloading github.com/mailru/easyjson v0.7.6
go: downloading github.com/josharian/intern v1.0.0
Update dependencies:
$ go mod tidy
go: downloading github.com/go-logr/zapr v1.2.3
go: downloading go.uber.org/zap v1.24.0
go: downloading github.com/stretchr/testify v1.8.0
go: downloading github.com/onsi/ginkgo/v2 v2.6.0
go: downloading github.com/onsi/gomega v1.24.1
go: downloading gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f
go: downloading go.uber.org/atomic v1.7.0
go: downloading go.uber.org/multierr v1.6.0
go: downloading github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e
go: downloading github.com/benbjohnson/clock v1.1.0
Next: define a resource with:
$ operator-sdk create api

operator-sdk create api --group monitors --version v1beta1 --kind Ping --namespaced --controller --resource --make

## make cahnges in ping_types.go

## make changes to reconcile

make manifests
make install
make run

