module github.com/weaveworks/profiles

go 1.16

require (
	github.com/fluxcd/helm-controller/api v0.10.0
	github.com/fluxcd/kustomize-controller/api v0.11.1
	github.com/fluxcd/source-controller/api v0.12.1
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.1
	github.com/onsi/gomega v1.11.0
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	knative.dev/pkg v0.0.0-20210412173742-b51994e3b312
	sigs.k8s.io/controller-runtime v0.8.3
)
