module github.com/kuperiu/k8s-newrelic-adapter

go 1.12

require (
	github.com/aws/aws-sdk-go v1.28.5
	github.com/awslabs/k8s-cloudwatch-adapter v0.8.0 // indirect
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20190918110929-3d9be26a50eb
	github.com/newrelic/newrelic-client-go v0.28.1
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.0.0-20191004115701-31ade1b30762
	k8s.io/client-go v0.0.0-20191029021442-5f2132fc4383
	k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base v0.0.0-20191004121406-d5138742ad72
	k8s.io/klog v1.0.0
	k8s.io/metrics v0.0.0-20191004123503-ae3d6ea895be
)
