module github.com/kidk/k8s-newrelic-adapter

go 1.13

require (
	github.com/NYTimes/gziphandler v1.0.1 // indirect
	github.com/aws/aws-sdk-go v1.25.11
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170208215640-dcef7f557305 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20190918110929-3d9be26a50eb
	github.com/newrelic/newrelic-client-go v0.28.1
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.12.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.0.0-20191004115701-31ade1b30762
	k8s.io/apiserver v0.0.0-20191109104256-50c872e90e34 // indirect
	k8s.io/client-go v0.0.0-20191029021442-5f2132fc4383
	k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base v0.0.0-20191004121406-d5138742ad72
	k8s.io/klog v1.0.0
	k8s.io/metrics v0.0.0-20191004123503-ae3d6ea895be
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

replace k8s.io/client-go => k8s.io/client-go v12.0.0+incompatible
