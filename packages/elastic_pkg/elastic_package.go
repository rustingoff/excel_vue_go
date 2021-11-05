package elastic_pkg

import (
	"github.com/olivere/elastic/v7"
	"log"
)

func NewElasticSearchConnection() *elastic.Client {

	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "amazon_campaign"),
	)

	if err != nil {
		log.Println("failed connection to Elastic Search")
		panic(err.Error())
	}

	log.Println("Successfully connected to Elastic Search.")

	return client
}
