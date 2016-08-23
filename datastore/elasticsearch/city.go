package elasticsearch

import (
	"fmt"

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
