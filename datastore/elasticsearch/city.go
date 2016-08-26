package elasticsearch

import (
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

type CityStore struct {
	*ElasticStore
}

func NewCityStore(config *ElasticConfig, client *elastic.Client) *CityStore {
	return &CityStore{NewElasticStore(config)}
}
