package configs

import (
	"timezones_mc/datastore/elasticsearch"
)

var mapping = `{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "mappings":{
        "city":{
            "properties":{
                "accentName":{
                    "type":"string"
                },
                "countryCode":{
                    "type":"string"
                },
                "latitude":{
                    "type":"double"
                },
                "longitude":{
                    "type":"double"
                },
                "name":{
                    "type":"string"
                }
            }
        }
    }
}`

var CityStoreConfig *elasticsearch.ElasticConfig = &elasticsearch.ElasticConfig{
	IndexName: "timezones",
	TypeName:   "city",
	Mapping:   mapping,
}
