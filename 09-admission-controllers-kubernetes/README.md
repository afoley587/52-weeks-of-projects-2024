go get -u github.com/gin-gonic/gin
go get -u k8s.io/api/core/v1
go get -u k8s.io/api/admission/v1
go get -u k8s.io/apimachinery/pkg/types


# what is it and how it works

# high level overview of what has to be done

# source code / dependencies

# building / generating certs

certsdir="$PWD/certs"
mkdir -p $certsdir
_pwd=$PWD
cd $certsdir
openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 365 -key ca.key \
  -subj "/C=AU/CN=example-admission-webhook"\
  -out ca.crt

openssl req -newkey rsa:2048 -nodes -keyout server.key \
  -subj "/C=AU/CN=example-admission-webhook" \
  -out server.csr

openssl x509 -req \
  -extfile <(printf "subjectAltName=DNS:example-admission-webhook.com") \
  -days 365 \
  -in server.csr \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt

cd $_pwd

# start minikube
minikube start --cpus=2 --memory=4096

# deploy manifests

# deploy pods

