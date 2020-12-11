package main

import (
	"flag"
	"os"
	"time"

	"github.com/pkg/errors"
	"k8s.io/component-base/logs"
	"k8s.io/klog"

	clientset "github.com/kidk/k8s-newrelic-adapter/pkg/client/clientset/versioned"
	informers "github.com/kidk/k8s-newrelic-adapter/pkg/client/informers/externalversions"
	"github.com/kidk/k8s-newrelic-adapter/pkg/controller"
	"github.com/kidk/k8s-newrelic-adapter/pkg/metriccache"
	"github.com/kidk/k8s-newrelic-adapter/pkg/newrelic"
	cwprov "github.com/kidk/k8s-newrelic-adapter/pkg/provider"
	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
)

// NewRelicAdapter represents a custom metrics BaseAdapter for NewRelic
type NewRelicAdapter struct {
	basecmd.AdapterBase
}

func (a *NewRelicAdapter) makeNewRelicClient() (newrelic.Client, error) {
	client := newrelic.NewRelicClient()
	return client, nil
}

func (a *NewRelicAdapter) newController(metriccache *metriccache.MetricCache) (*controller.Controller, informers.SharedInformerFactory) {
	clientConfig, err := a.ClientConfig()
	if err != nil {
		klog.Fatalf("unable to construct client config: %v", err)
	}
	adapterClientSet, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		klog.Fatalf("unable to construct lister client to initialize provider: %v", err)
	}

	adapterInformerFactory := informers.NewSharedInformerFactory(adapterClientSet, time.Second*30)
	handler := controller.NewHandler(
		adapterInformerFactory.Metrics().V1alpha1().ExternalMetrics().Lister(),
		metriccache)

	controller := controller.NewController(adapterInformerFactory.Metrics().V1alpha1().ExternalMetrics(), &handler)

	return controller, adapterInformerFactory
}

func (a *NewRelicAdapter) makeProvider(cwClient newrelic.Client, metriccache *metriccache.MetricCache) (provider.ExternalMetricsProvider, error) {
	client, err := a.DynamicClient()
	if err != nil {
		return nil, errors.Wrap(err, "unable to construct Kubernetes client")
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		return nil, errors.Wrap(err, "unable to construct RESTMapper")
	}

	cwProvider := cwprov.NewRelicProvider(client, mapper, cwClient, metriccache)
	return cwProvider, nil
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	// set up flags
	cmd := &NewRelicAdapter{}
	cmd.Name = "newrelic-metrics-adapter"
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the klog flags
	cmd.Flags().Parse(os.Args)

	stopCh := make(chan struct{})
	defer close(stopCh)

	metriccache := metriccache.NewMetricCache()

	// start and run contoller components
	controller, adapterInformerFactory := cmd.newController(metriccache)
	go adapterInformerFactory.Start(stopCh)
	go controller.Run(2, time.Second, stopCh)

	// create NewRelic client
	nrClient, err := cmd.makeNewRelicClient()
	if err != nil {
		klog.Fatalf("unable to construct NewRelic client: %v", err)
	}

	// construct the provider
	nrProvider, err := cmd.makeProvider(nrClient, metriccache)
	if err != nil {
		klog.Fatalf("unable to construct NewRelic metrics provider: %v", err)
	}

	cmd.WithExternalMetrics(nrProvider)

	klog.Info("NewRelic metrics adapter started")

	if err := cmd.Run(stopCh); err != nil {
		klog.Fatalf("unable to run NewRelic metrics adapter: %v", err)
	}
}
