module github.com/weaveworks/profiles

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/fluxcd/helm-controller/api v0.12.0
	github.com/fluxcd/kustomize-controller/api v0.31.0
	github.com/fluxcd/pkg/apis/meta v0.18.0
	github.com/fluxcd/pkg/version v0.1.0
	github.com/fluxcd/source-controller v0.16.0
	github.com/fluxcd/source-controller/api v0.17.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-logr/logr v1.2.3
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.20.1
	github.com/weaveworks/schemer v0.0.0-20210802122110-338b258ad2ca
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21
	google.golang.org/grpc v1.47.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.28.0
	k8s.io/api v0.25.4
	k8s.io/apimachinery v0.25.4
	k8s.io/client-go v0.25.4
	sigs.k8s.io/controller-runtime v0.13.1
)
