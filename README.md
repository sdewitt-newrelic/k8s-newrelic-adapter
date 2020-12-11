This work is based on [kuperiu/k8s-newrelic-adapter](https://github.com/kuperiu/k8s-newrelic-adapter).

# WIP: Kubernetes Custom Metrics Adapter for NewRelic

This is an implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for New Relic events, metrics and logs.

This adapter allows you to scale your Kubernetes deployment using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with data coming from NewRelic and NRQL (New Relic Query Language).

## Installation instructions

### Create `newrelic-custom-metrics` namespace

`kubectl create namespace newrelic-custom-metrics`

### Create account ID and Personal API Token secrets

This adapter requires the following to access data from within New Relic:

- Account ID - Change **NEW_RELIC_ACCOUNT_ID** in deploy/adapter.yaml to your NewRelic Account ID

- Personal API Token - Create a secret called `newrelic` with the key `personal_api_key`. (The key should be encode to base64)

You can get the instructions on getting a personal API token on [the New Relic documentation](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key).

```
kubectl create secret generic newrelic \
  -n newrelic-custom-metrics \
  --from-literal=account_id=**NEW_RELIC_ACCOUNT_ID** \
  --from-literal=personal_api_key=**NEWRELIC_API_KEY**
```

## Deploy

Now deploy the adapter to your Kubernetes cluster.

```bash
$ kubectl apply -f https://raw.githubusercontent.com/kidk/k8s-newrelic-adapter/master/deploy/adapter.yaml
namespace/newrelic-custom-metrics created
clusterrolebinding.rbac.authorization.k8s.io/k8s-newrelic-adapter:system:auth-delegator created
rolebinding.rbac.authorization.k8s.io/k8s-newrelic-adapter-auth-reader created
deployment.apps/k8s-newrelic-adapter created
clusterrolebinding.rbac.authorization.k8s.io/k8s-newrelic-adapter-resource-reader created
serviceaccount/k8s-newrelic-adapter created
service/k8s-newrelic-adapter created
apiservice.apiregistration.k8s.io/v1beta1.external.metrics.k8s.io created
clusterrole.rbac.authorization.k8s.io/k8s-newrelic-adapter:external-metrics-reader created
clusterrole.rbac.authorization.k8s.io/k8s-newrelic-adapter-resource-reader created
clusterrolebinding.rbac.authorization.k8s.io/k8s-newrelic-adapter:external-metrics-reader created
customresourcedefinition.apiextensions.k8s.io/externalmetrics.metrics.newrelic created
clusterrole.rbac.authorization.k8s.io/k8s-newrelic-adapter:crd-metrics-reader created
clusterrolebinding.rbac.authorization.k8s.io/k8s-newrelic-adapter:crd-metrics-reader created
```

This creates a new namespace `newrelic-custom-metrics` and deploys the necessary ClusterRole, Service Account,
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

## Deploying The sample application
There is a sample application provided in this repository for you to test how the adapter works.
Refer to this [guide](sample/README.md)

## License
This library is licensed under the Apache 2.0 License.

## Issues
Report any issues in the [Github Issues](https://github.com/kidk/k8s-newrelic-adapter/issues)
