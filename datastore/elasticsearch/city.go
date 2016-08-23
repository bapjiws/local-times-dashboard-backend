package elasticsearch

import (
	"fmt"

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
//	Mapping string
//}

func CreateIndex() error { // TODO: make it a method a with pointer receiver?
	client, err := elastic.NewClient()
	if err != nil {
		return fmt.Errorf("Couldn't create a client: %s.\n", err.Error())
	}

	createIndex, err := client.CreateIndex("city").BodyString(mapping).Do()
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// TODO: Not acknowledged
	}

	return nil
}
