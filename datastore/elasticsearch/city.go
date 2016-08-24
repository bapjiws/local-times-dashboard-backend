package elasticsearch

import (
	"fmt"
	"log"
	"os"

	"../../revel_app/app/models"
	"gopkg.in/olivere/elastic.v3"
)

// TODO: complete
var mapping = `{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "mappings":{
        "city":{
            "properties":{

            }
        }
    }
}`

//type CityIndex struct {
//	Client *elastic.Client
//	Mapping string
//}

// TODO: implement NewCityIndex function

func CreateIndex() error { // TODO: make it a method a with pointer receiver?
	client, err := elastic.NewClient()
	if err != nil {
		return fmt.Errorf("Couldn't create a client: %s.\n", err.Error())
	}

	createIndex, err := client.CreateIndex("city").BodyString(mapping).Do() // TODO: abstract this away via index interface?
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// TODO: Not acknowledged
	}

	return nil
}

func AddDocument(city models.City) error {
	// TODO: creating a client for the second time, make this a method?
	client, err := elastic.NewClient()
	if err != nil {
		return fmt.Errorf("Couldn't create a client: %s.\n", err.Error())
	}

	put, err := client.Index().
		Index("timezones").
		Type("city").
		Id("1"). // TODO: work this out
		BodyJson(city).
		Do()
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put.Id, put.Index, put.Type)

	return nil
}

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
	IndexName  string
	TypeName   string
	Mapping    string
	IndexAlias string
}

func NewElasticIndex(config *ElasticConfig, client *elastic.Client) *ElasticIndex {
	return &ElasticIndex{config, client, nil}
}

type CityIndex struct {
	*ElasticIndex
}

// TODO: check createFreelancerInSearchStore's logic
func (i *CityIndex) IndexCity(city *models.City) (*models.City, *elastic.IndexResponse, error) {
	var result *elastic.IndexResponse
	var err error

	/*If city already exists, update the document and refresh the index.
	Refresh vs Flush: Changes to Lucene are only persisted to disk during a Lucene commit (flush), which is a relatively
	heavy operation and so cannot be performed after every index or delete operation. The refresh API allows to explicitly
	refresh one or more index, making all operations performed since the last refresh available for search.
	Also see: http://stackoverflow.com/questions/19963406/refresh-vs-flush.*/
	// TODO: result, err = ...BodyJson(elFreel).Refresh(true).Do()

	if err != nil {
		return nil, nil, err
	}

	if result.Version == 1 && !result.Created {
		return nil, nil, fmt.Errorf("City has not been indexed.\n")
	}

	return city, result, nil
}

// Delete and recreate the index if it exists, otherwise create a new index.
// TODO: reindex with zero downtime, see: https://www.elastic.co/blog/changing-mapping-with-zero-downtime
func (i *ElasticIndex) Reindex() error {
	indexName := i.IndexName

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
