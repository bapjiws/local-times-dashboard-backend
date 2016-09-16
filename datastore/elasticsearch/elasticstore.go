package elasticsearch

import (
	"fmt"
	"log"
	"os"
	"github.com/bapjiws/timezones_mc/models"
	"github.com/bapjiws/timezones_mc/utils"

	"github.com/satori/go.uuid"
	"gopkg.in/olivere/elastic.v3"
)

func connect() (client *elastic.Client) {
	var err error
	var debugMode = os.Getenv("DEBUG_MODE")
	if debugMode == "" {
		debugMode = "NONE"
	}

	switch debugMode {
	case "NONE":
		client, err = elastic.NewClient()
		utils.PanicOnError(err)
	case "CONSOLE":
		client, err = elastic.NewClient(
			elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
			elastic.SetTraceLog(log.New(os.Stderr, "[[ELASTIC]]", 0)),
		)
		utils.PanicOnError(err)
	case "FILE": // Stuff adopted from: https://github.com/olivere/elastic/wiki/Logging
		file, err := os.OpenFile("elastic.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		utils.PanicOnError(err)

		client, err = elastic.NewClient(
			elastic.SetInfoLog(log.New(file, "ELASTIC ", log.LstdFlags)),
			elastic.SetErrorLog(log.New(file, "ELASTIC ", log.LstdFlags)),
			elastic.SetTraceLog(log.New(file, "ELASTIC ", log.LstdFlags)),
		)
		utils.PanicOnError(err)
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
	/*If the document already exists, update it and refresh the index.
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
		return fmt.Errorf("Error while attempting to index the document: %s.\n", err.Error())
	}

	if result.Version == 1 && !result.Created {
		return fmt.Errorf("Document %s has not been indexed.\n", doc)
	}

	return nil
}

/*{
	"query": {},
	"suggest": {
		"city_suggest": {
			"text": "ain",
			"completion": {
				"field": "suggest"
			}
		}
	}
}*/

func (es *ElasticStore) FindDocumentById(id string) (models.Document, error) {
	searchResult, err := es.Get().Index(es.IndexName).Type(es.TypeName).Id(id).Do()
	if err != nil {
		return nil, err
	}

	return searchResult.Source, nil
}

func (es *ElasticStore) SuggestDocuments(suggesterName string, text string, field string, payloadKey string) ([]models.Document, error) {
	suggestResult, err := es.Search(es.IndexName).
		Query(elastic.NewBoolQuery()).Size(0).
		Suggester(elastic.NewCompletionSuggester(suggesterName).Text(text).Field(field)).Do()
	if err != nil {
		return nil, err
	}

	// type SearchSuggest map[string][]SearchSuggestion

	//type SearchSuggestion struct {
	//	Text    string                   `json:"text"`
	//	Offset  int                      `json:"offset"`
	//	Length  int                      `json:"length"`
	//	Options []SearchSuggestionOption `json:"options"`
	//}

	// "payload": { "city_id": "0e2997f0-36c4-4995-8115-4f433b693775"}

	suggestions := make([]models.Document, 0, len(suggestResult.Suggest[suggesterName][0].Options))
	for _, option := range suggestResult.Suggest[suggesterName][0].Options {
		suggestion := struct {
			Text string `json:"text"`
			Id   string `json:"id"`
		}{
			Text: option.Text,
			Id:   option.Payload.(map[string]interface{})[payloadKey].(string),
		}

		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// Delete the index if it exists, otherwise create a new index and an alias for it.
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
