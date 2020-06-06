module github.com/kuperiu/k8s-newrelic-adapter

go 1.12

require (
	github.com/aws/aws-sdk-go v1.25.11
	github.com/emicklei/go-restful v2.12.0+incompatible // indirect
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6 // indirect
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20190918110929-3d9be26a50eb
	github.com/newrelic/newrelic-client-go v0.28.1
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.3
	k8s.io/apiserver v0.18.3 // indirect
	k8s.io/client-go v0.18.3
	k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base v0.18.3
	k8s.io/klog v1.0.0
	k8s.io/metrics v0.0.0-20191004123503-ae3d6ea895be
)
