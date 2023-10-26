package infrastructures

import (
	"github.com/olivere/elastic/v7"
	"go-clean/src/app/config"
)

func CreateElasticsearchClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(config.GetConfig().OutboundURL.ElasticsearchOutbound),
		elastic.SetBasicAuth(config.GetConfig().Elasticsearch.Username, config.GetConfig().Elasticsearch.Password),
	)
	if err != nil {
		panic(err)
	}
	return client
}
