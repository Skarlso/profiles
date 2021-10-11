module github.com/weaveworks/profiles

go 1.16

require (
	github.com/fluxcd/helm-controller/api v0.10.0
	github.com/fluxcd/kustomize-controller/api v0.15.2
	github.com/fluxcd/source-controller/api v0.12.1
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	knative.dev/pkg v0.0.0-20210412173742-b51994e3b312
	sigs.k8s.io/controller-runtime v0.10.1
)
