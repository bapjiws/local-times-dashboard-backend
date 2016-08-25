package elasticsearch

import (
	"fmt"
	"timezones_mc/revel_app/app/models"

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

var CityStorageConfig *ElasticConfig = &ElasticConfig{
	IndexName: "timezones",
	TypeName:   "city", //TODO: "citieS"?
	Mapping:   `{"mappings": ` + mapping + `}`,
}

type CityStorage struct {
	*ElasticStorage
}

func CreateIndex() error { // TODO: make it a method a with pointer receiver?
	// TODO: use connect() and make the latter return an error instead of panicking
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

func NewCityStorage(config *ElasticConfig, client *elastic.Client) *CityStorage {
	return &CityStorage{NewElasticStorage(config)}
}

// TODO: move into
func (cs *CityStorage) IndexCity(city *models.City) (*models.City, *elastic.IndexResponse, error) {
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
