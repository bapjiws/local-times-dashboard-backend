package configs

import (
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
)

var mapping = `{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "city": {
            "properties": {
                "name": {
                    "type": "string"
                },
                "accentName": {
                    "type": "string"
                },
                "countryCode": {
                    "type": "string"
                },
                "latitude": {
                    "type": "double"
                },
                "longitude": {
                    "type": "double"
                },
                "suggest": {
                    "type": "completion",
                    "analyzer": "simple",
                    "search_analyzer": "simple",
                    "payloads": true
                }
            }
        }
    }
}`

var CityStoreConfig *elasticsearch.ElasticConfig = &elasticsearch.ElasticConfig{
	IndexName: "timezones",
	TypeName:  "city",
	Mapping:   mapping,
}
