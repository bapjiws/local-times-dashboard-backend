package elasticsearch

import (
	"timezones_mc/datastore/elasticsearch"
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

var CityStoreConfig *elasticsearch.ElasticConfig = &elasticsearch.ElasticConfig{
	IndexName: "timezones",
	TypeName:   "city", //TODO: "citieS"?
	Mapping:   `{"mappings": ` + mapping + `}`,
}

//type CityStore struct {
//	*ElasticStore
//}
//
//func NewCityStore(config *ElasticConfig) *CityStore {
//	return &CityStore{NewElasticStore(config)}
//}
