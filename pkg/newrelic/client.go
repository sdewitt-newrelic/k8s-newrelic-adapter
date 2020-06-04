package newrelic

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
	"k8s.io/klog"
)

type newRelicClient struct {
	client *newrelic.NewRelic
}

// NewRelicClient creates a new NewRelic client.
func NewRelicClient() Client {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	if apiKey == "" {
		log.Fatal("an API key is required, please set the NEW_RELIC_ADMIN_API_KEY environment variable")
	}
	nr, err := newrelic.New(newrelic.ConfigPersonalAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to create a New Relic client with error %v", err)
	}
	return &newRelicClient{nr}
}

func (c *newRelicClient) Query() (float64, error) {
	query := `
	query($accountId: Int!, $nrqlQuery: Nrql!) {
		actor {
			account(id: $accountId) {
				nrql(query: $nrqlQuery, timeout: 5) {
					results
				}
			}
		}
  }`

	variables := map[string]interface{}{
		"accountId": 2364728,
		"nrqlQuery": "SELECT latest(lior.k8s.num) FROM Metric where metricName='lior.k8s.num'",
	}
	resp, err := c.client.NerdGraph.Query(query, variables)
	if err != nil {
		log.Fatal("error running NerdGraph query: ", err)
	}

	queryResp := resp.(nerdgraph.QueryResponse)
	actor := queryResp.Actor.(map[string]interface{})
	account := actor["account"].(map[string]interface{})
	nrql := account["nrql"].(map[string]interface{})
	results := nrql["results"].([]interface{})

	var durations float64

	klog.V(2).Info("#############")
	klog.V(2).Info(queryResp)
	klog.V(2).Info("#############")

	for _, r := range results {
		data := r.(map[string]interface{})
		durations = data["latest.lior.k8s.num"].(float64)
		return durations, nil
	}
	return 0, nil
}
