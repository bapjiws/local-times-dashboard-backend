package elasticsearch

import (
	"os"
	"log"
	"gopkg.in/olivere/elastic.v3"
	"fmt"
)

func connect() (client *elastic.Client) {
	var err error
	var debugMode = os.Getenv("DEBUG_MODE")
	if debugMode == "" {
		debugMode = "none"
	}

	var elasticSearchDomain = "http://localhost:9200" // TODO: eliminate?

	// TODO: consult with https://github.com/olivere/elastic/wiki/Logging
	switch debugMode {
	case "none":
		if client, err = elastic.NewClient(elastic.SetURL(elasticSearchDomain), elastic.SetSniff(false)); err != nil {
			panic("ElasticDomain: " + elasticSearchDomain + " Error: " + err.Error())
		}
	case "console":
		if client, err = elastic.NewClient(elastic.SetURL(elasticSearchDomain), elastic.SetTraceLog(log.New(os.Stdout, "AA_", 1)), elastic.SetSniff(false)); err != nil {
			panic("ElasticDomain: " + elasticSearchDomain + " Error: " + err.Error())
		}
	case "file":
		file, err := os.OpenFile("elastic.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}
		if client, err = elastic.NewClient(elastic.SetURL(elasticSearchDomain), elastic.SetTraceLog(log.New(file, "ELASTIC ", log.LstdFlags)), elastic.SetSniff(false)); err != nil {
			panic("ElasticDomain: " + elasticSearchDomain + " Error: " + err.Error())
		}
	}

	return client
}

type ElasticIndex struct {
	*ElasticConfig
	*elastic.Client
	Error error
}

type ElasticConfig struct {
	Index   string
	Type    string
	Mapping string
}

func NewElasticIndex(config *ElasticConfig, client *elastic.Client) *ElasticIndex {
	return &ElasticIndex{config, client, nil}
}

// Delete and recreate the index if it exists, otherwise create a new index.
// TODO: reindex with zero downtime, see: https://www.elastic.co/blog/changing-mapping-with-zero-downtime
func (i *ElasticIndex) Reindex() error {
	indexName := i.Index

	exists, err := i.IndexExists(indexName).Do()
	if err != nil {
		return err
	}
	if exists {
		deleteIndex, err := i.DeleteIndex(indexName).Do()
		if err != nil {
			return err
		}
		if !deleteIndex.Acknowledged {
			return fmt.Errorf("Deletion of index %s was not acknowledged.\n", indexName)
		}
	}

	createIndex, err := i.CreateIndex(indexName).Body(i.Mapping).Do()
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("Creation of index %s was not acknowledged.\n", indexName)
	}

	alias := fmt.Sprintf("%s-%s", indexName, "alias")
	aliasCreate, err := i.Alias().Add(indexName, alias).Do()
	if err != nil {
		return fmt.Errorf("Couldn't create alias for index %s: %s", indexName, err.Error())
	}
	if !aliasCreate.Acknowledged {
		return fmt.Errorf("Creation of alias for index %s was not acknowledged", indexName)
	}

	return nil
}
