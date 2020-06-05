package controller

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	api "github.com/kuperiu/k8s-newrelic-adapter/pkg/apis/metrics/v1alpha1"
	"github.com/kuperiu/k8s-newrelic-adapter/pkg/metriccache"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kuperiu/k8s-newrelic-adapter/pkg/client/clientset/versioned/fake"
	informers "github.com/kuperiu/k8s-newrelic-adapter/pkg/client/informers/externalversions"
)

func getExternalKey(externalMetric *api.ExternalMetric) namespacedQueueItem {
	return namespacedQueueItem{
		namespaceKey: fmt.Sprintf("%s/%s", externalMetric.Namespace, externalMetric.Name),
		kind:         externalMetric.TypeMeta.Kind,
	}
}

func TestExternalMetricValueIsStored(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	queueItem := getExternalKey(externalMetric)
	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	metricRequest, exists := metriccache.GetNewRelicQuery(externalMetric.Namespace, externalMetric.Name)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateExternalMetricResult(metricRequest, externalMetric, t)
}

func TestShouldBeAbleToStoreCustomAndExternalWithSameNameAndNamespace(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	externalItem := getExternalKey(externalMetric)
	err := handler.Process(externalItem)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	externalRequest, exists := metriccache.GetNewRelicQuery(externalMetric.Namespace, externalMetric.Name)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateExternalMetricResult(externalRequest, externalMetric, t)
}

func TestShouldFailOnInvalidCacheKey(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	queueItem := namespacedQueueItem{
		namespaceKey: "invalidkey/with/extrainfo",
		kind:         "somethingwrong",
	}
	err := handler.Process(queueItem)

	if err == nil {
		t.Errorf("error after processing nil, want non nil")
	}

	_, exists := metriccache.GetNewRelicQuery(externalMetric.Namespace, externalMetric.Name)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func TestWhenExternalItemHasBeenDeleted(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")

	// don't put anything in the stores
	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	// add the item to the cache then test if it gets deleted
	queueItem := getExternalKey(externalMetric)
	metriccache.Update(queueItem.Key(), "test", cloudwatch.GetMetricDataInput{})

	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error == %v, want nil", err)
	}

	_, exists := metriccache.GetNewRelicQuery(externalMetric.Namespace, externalMetric.Name)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func TestWhenItemKindIsUnknown(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	// don't put anything in the stores, as we are not looking anything up
	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	// add the item to the cache then test if it gets deleted
	queueItem := namespacedQueueItem{
		namespaceKey: "default/unknown",
		kind:         "Unknown",
	}

	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error == %v, want nil", err)
	}

	_, exists := metriccache.GetNewRelicQuery("default", "unknown")

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func newHandler(storeObjects []runtime.Object, externalMetricsListerCache []*api.ExternalMetric) (Handler, *metriccache.MetricCache) {
	fakeClient := fake.NewSimpleClientset(storeObjects...)
	i := informers.NewSharedInformerFactory(fakeClient, 0)

	externalMetricLister := i.Metrics().V1alpha1().ExternalMetrics().Lister()

	for _, em := range externalMetricsListerCache {
		i.Metrics().V1alpha1().ExternalMetrics().Informer().GetIndexer().Add(em)
	}

	metriccache := metriccache.NewMetricCache()
	handler := NewHandler(externalMetricLister, metriccache)

	return handler, metriccache
}
