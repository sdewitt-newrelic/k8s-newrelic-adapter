package controller

import (
	"fmt"
	"testing"

	api "github.com/kuperiu/k8s-newrelic-adapter/pkg/apis/metrics/v1alpha1"
	"github.com/kuperiu/k8s-newrelic-adapter/pkg/metriccache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kuperiu/k8s-newrelic-adapter/pkg/client/clientset/versioned/fake"
	informers "github.com/kuperiu/k8s-newrelic-adapter/pkg/client/informers/externalversions"
)

var query = "SELECT latest(test.k8s.num) FROM Metric WHERE metricName='test.k8s.num'"

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

	if metricRequest != query {
		t.Errorf("externalRequest = %s, want %v", metricRequest, query)
	}
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

	if externalRequest != query {
		t.Errorf("externalRequest = %s, want %v", externalRequest, query)
	}
}

func TestShouldFailOnInvalidCacheKey(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	queueItem := namespacedQueueItem{
		namespaceKey: "default/with/extrainfo",
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
	metriccache.Update(queueItem.Key(), "test", query)

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

	_, exists := metriccache.GetNewRelicQuery("default", "test")

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

func newFullExternalMetric(name string) *api.ExternalMetric {
	return &api.ExternalMetric{
		TypeMeta: metav1.TypeMeta{APIVersion: api.SchemeGroupVersion.String(), Kind: "ExternalMetric"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: api.MetricSeriesSpec{
			Name: "Name",
			Queries: []api.MetricDataQuery{
				{
					ID:    "query1",
					Query: query,
				},
			},
		},
	}
}
