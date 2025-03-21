package infrastructures

import (
	"github.com/olivere/elastic/v7"
	"go-clean/src/app/config"
	"log"
	"time"
)

func CreateElasticsearchClient() *elastic.Client {
	var (
		client     *elastic.Client
		err        error
		maxRetries = 5
	)

	for i := 1; i <= maxRetries; i++ {
		newClient, err := elastic.NewClient(
			elastic.SetURL(config.GetConfig().OutboundURL.ElasticsearchOutbound),
			elastic.SetBasicAuth(config.GetConfig().Elasticsearch.Username, config.GetConfig().Elasticsearch.Password),
		)

		if err != nil {
			log.Printf("(Attempt %d/%d)\n", i, maxRetries)
			log.Printf("Failed to connect to the elasticsearch. Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		} else {
			client = newClient
		}
	}

	if err != nil {
		panic(err)
	}

	return client
}
