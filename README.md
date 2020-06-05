[![Build Status](https://travis-ci.org/kuperiu/k8s-newrelic-adapter.svg?branch=master)](https://travis-ci.org/kuperiu/k8s-newrelic-adapter)
[![GitHub
release](https://img.shields.io/github/release/kuperiu/k8s-newrelic-adapter/all.svg)](https://github.com/kuperiu/k8s-newrelic-adapter/releases)

[![docker image
size](https://shields.beevelop.com/docker/image/image-size/kuperiu/k8s-newrelic-adapter/latest.svg)](https://hub.docker.com/r/kuperiu/k8s-newrelic-adapter)
[![image
layers](https://shields.beevelop.com/docker/image/layers/kuperiu/k8s-newrelic-adapter/latest.svg)](https://hub.docker.com/r/kuperiu/k8s-newrelic-adapter)
[![image
pulls](https://shields.beevelop.com/docker/pulls/kuperiu/k8s-newrelic-adapter.svg)](https://hub.docker.com/r/kuperiu/k8s-newrelic-adapter)

# Kubernetes Custom Metrics Adapter for Kubernetes


An implementation of the Kubernetes [Custom Metrics API and External Metrics
API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis)
for NewRelic metrics.

This adapter allows you to scale your Kubernetes deployment using the [Horizontal Pod
Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with
metrics from NewRelic.

## Prerequsites
This adapter requires the following to access metric data from Amazon CloudWatch.
- Account ID - Change **NEW_RELIC_ACCOUNT_ID** in deploy/adapter.yaml to your NewRelic Account ID
- Personal API Token - Create a secret called newrelic with The key api_key (The key should be encode to base64)
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: newrelic
  namespace: custom-metrics
type: Opaque
data:
  api_key: 1234=
```


You can get The personal API toke from
[Here](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key). 


## Deploy
Requires a Kubernetes cluster with Metric Server deployed, Amazon EKS cluster is fine too.

Now deploy the adapter to your Kubernetes cluster.

```bash
$ kubectl apply -f https://raw.githubusercontent.com/kuperiu/k8s-newrelic-adapter/master/deploy/adapter.yaml
namespace/custom-metrics created
clusterrolebinding.rbac.authorization.k8s.io/k8s-cloudwatch-adapter:system:auth-delegator created
rolebinding.rbac.authorization.k8s.io/k8s-cloudwatch-adapter-auth-reader created
deployment.apps/k8s-cloudwatch-adapter created
clusterrolebinding.rbac.authorization.k8s.io/k8s-cloudwatch-adapter-resource-reader created
serviceaccount/k8s-cloudwatch-adapter created
service/k8s-cloudwatch-adapter created
apiservice.apiregistration.k8s.io/v1beta1.external.metrics.k8s.io created
clusterrole.rbac.authorization.k8s.io/k8s-cloudwatch-adapter:external-metrics-reader created
clusterrole.rbac.authorization.k8s.io/k8s-cloudwatch-adapter-resource-reader created
clusterrolebinding.rbac.authorization.k8s.io/k8s-cloudwatch-adapter:external-metrics-reader created
customresourcedefinition.apiextensions.k8s.io/externalmetrics.metrics.aws created
clusterrole.rbac.authorization.k8s.io/k8s-cloudwatch-adapter:crd-metrics-reader created
clusterrolebinding.rbac.authorization.k8s.io/k8s-cloudwatch-adapter:crd-metrics-reader created
```

This creates a new namespace `custom-metrics` and deploys the necessary ClusterRole, Service Account,
Role Binding, along with the deployment of the adapter.

### Verifying the deployment
Next you can query the APIs to see if the adapter is deployed correctly by running:

```bash
$ kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "external.metrics.k8s.io/v1beta1",
  "resources": [
  ]
}
```

## Deploying the sample application
There is a sample application provided in this repository for you to test how the adapter works.
Refer to this [guide](sample/README.md)

## License

This library is licensed under the Apache 2.0 License. 

## Issues
Report any issues in the [Github Issues](https://github.com/kuperiu/k8s-newrelic-adapter/issues)
