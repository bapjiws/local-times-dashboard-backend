package elasticsearch

import (
	"fmt"
	"log"
	"os"
	"timezones_mc/revel_app/app/models"

	"github.com/satori/go.uuid"
	"gopkg.in/olivere/elastic.v3"
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

//ElasticStore implements the DataStore interface
type ElasticStore struct {
	*ElasticConfig
	*elastic.Client
}

type ElasticConfig struct {
	IndexName string
	TypeName  string
	Mapping   string
}

func NewElasticStore(config *ElasticConfig) *ElasticStore {
	return &ElasticStore{config, connect()}
}

func (es *ElasticStore) AddDocument(doc models.Document) error {
	/*If city already exists, update the document and refresh the index.
	Refresh vs Flush: Changes to Lucene are only persisted to disk during a Lucene commit (flush), which is a relatively
	heavy operation and so cannot be performed after every index or delete operation. The refresh API allows to explicitly
	refresh one or more index, making all operations performed since the last refresh available for search.
	Also see: http://stackoverflow.com/questions/19963406/refresh-vs-flush.*/
	result, err := es.Client.
		Index().
		Index(es.IndexName).
		Type(es.TypeName).
		Id(uuid.NewV4().String()).
		BodyJson(doc).
		Do()
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	//fmt.Printf("Indexed city %s to index %s, type %s\n", result.Id, result.Index, result.Type)

	if result.Version == 1 && !result.Created {
		return fmt.Errorf("City %s has not been indexed.\n", doc.String())
	}

	return nil
}

// Delete and recreate the index if it exists, otherwise create a new index and an alias for it.
// TODO: reindex with zero downtime, see: https://www.elastic.co/blog/changing-mapping-with-zero-downtime
func (es *ElasticStore) Reindex() error {
	indexName := es.IndexName

	exists, err := es.IndexExists(indexName).Do()
	if err != nil {
		return err
	}
	if exists {
		deleteIndex, err := es.DeleteIndex(indexName).Do()
		if err != nil {
			return err
		}
		if !deleteIndex.Acknowledged {
			return fmt.Errorf("Deletion of index %s was not acknowledged.\n", indexName)
		}
	}

	createIndex, err := es.CreateIndex(indexName).Body(es.Mapping).Do()
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("Creation of index %s was not acknowledged.\n", indexName)
	}

	alias := fmt.Sprintf("%s-%s", indexName, "alias")
	aliasCreate, err := es.Alias().Add(indexName, alias).Do()
	if err != nil {
		return fmt.Errorf("Couldn't create alias for index %s: %s", indexName, err.Error())
	}
	if !aliasCreate.Acknowledged {
		return fmt.Errorf("Creation of alias for index %s was not acknowledged", indexName)
	}

	return nil
}
