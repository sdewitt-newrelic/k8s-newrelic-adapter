package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

func (p *newRelicProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	// Note:
	//		metric name and namespace is used to lookup for the CRD which contains configuration to
	//		call cloudwatch if not found then ignored and label selector is parsed for all the metrics
	klog.V(0).Infof("Received request for namespace: %s, metric name: %s, metric selectors: %s", namespace, info.Metric, metricSelector.String())

	_, selectable := metricSelector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}

	cwRequest, found := p.metricCache.GetCloudWatchRequest(namespace, info.Metric)
	if !found {
		klog.V(0).Info("$$$$$$$$$$")
	}
	klog.V(0).Info("@@@@@@@@@@")
	test := aws.StringValue(cwRequest)
	klog.V(0).Info("6. main  -- i  %T: &i=%p i=%v\n", test, &test, test)
	klog.V(0).Info("@@@@@@@@@@")
	metricValue, err := p.nrClient.Query()
	if err != nil {
		klog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	var quantity resource.Quantity
	if metricValue == 0 {
		quantity = *resource.NewMilliQuantity(0, resource.DecimalSI)
	} else {
		quantity = *resource.NewQuantity(int64(aws.Float64Value(&metricValue)), resource.DecimalSI)
	}
	externalmetric := external_metrics.ExternalMetricValue{
		MetricName: info.Metric,
		Value:      quantity,
		Timestamp:  metav1.Now(),
	}

	matchingMetrics := []external_metrics.ExternalMetricValue{}
	matchingMetrics = append(matchingMetrics, externalmetric)

	return &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}, nil
}

func (p *newRelicProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	p.valuesLock.RLock()
	defer p.valuesLock.RUnlock()

	// not implemented yet
	externalMetricsInfo := []provider.ExternalMetricInfo{}
	for _, name := range p.metricCache.ListMetricNames() {
		// only process if name is non-empty
		if name != "" {
			info := provider.ExternalMetricInfo{
				Metric: name,
			}
			externalMetricsInfo = append(externalMetricsInfo, info)
		}
	}
	return externalMetricsInfo
}
