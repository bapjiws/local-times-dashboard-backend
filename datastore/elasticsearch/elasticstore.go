package elasticsearch

import (
	"fmt"
	"log"
	"os"

	"github.com/bapjiws/timezones_mc/models/document"
	"github.com/bapjiws/timezones_mc/models/suggest"
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
		client, err = elastic.NewClient() // or, if using the default behavior of elasticsearch container, elastic.SetURL("http://172.17.0.2:9200")
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

func (es *ElasticStore) AddDocument(doc document.Document) error {
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

func (es *ElasticStore) FindDocumentById(id string) (document.Document, error) {
	searchResult, err := es.Get().Index(es.IndexName).Type(es.TypeName).Id(id).Do()
	if err != nil {
		return nil, err
	}

	return searchResult.Source, nil
}

func (es *ElasticStore) SuggestDocuments(s suggest.Suggest) ([]document.Document, error) {
	suggestResult, err := es.Search(es.IndexName).
		Suggester(elastic.NewCompletionSuggester(s.SuggesterName).Text(s.Text).Field(s.Field)).Do()
	if err != nil {
		return nil, err
	}

	suggestions := make([]document.Document, 0, len(suggestResult.Suggest[s.SuggesterName][0].Options))
	for _, option := range suggestResult.Suggest[s.SuggesterName][0].Options {
		suggestion := struct {
			Text string `json:"text"`
			Id   string `json:"id"`
		}{
			Text: option.Text,
			Id:   option.Payload.(map[string]interface{})[s.PayloadKeys["city_id"]].(string),
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
