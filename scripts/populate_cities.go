package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"timezones_mc/datastore/elasticsearch"
	"timezones_mc/datastore/elasticsearch/configs"
	"timezones_mc/revel_app/app/models"
	"gopkg.in/olivere/elastic.v2"
	"github.com/satori/go.uuid"
)

// TODO: create a utils folder and move it there?
func panicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}

var fileFlag = flag.String("file", "", "file to parse")

var wg sync.WaitGroup

// Unhandled characters: https://en.wikipedia.org/wiki/%C3%80
func main() {
	// TODO: wg := new(sync.WaitGroup)?

	flag.Parse()

	// go run scripts/populate_cities.go -file="cities/worldcities.txt"
	println(*fileFlag)

	if *fileFlag == "" {
		fmt.Fprintf(os.Stderr, "CSV file has not been specified. Use the 'file' flag:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Printf("Importing cities from %s\n", *fileFlag)

	esStore := elasticsearch.NewElasticStore(configs.CityStoreConfig)
	err := esStore.Reindex()
	panicOnError(err)

	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening the file: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = true // panic: line 19970, column 7: bare " in non-quoted-field

	//citiesImported := 0 // TODO: use an atomic counter or something like that + print every,say, 1000.

	headers, err := csvReader.Read()
	panicOnError(err)
	fmt.Printf("Headers: %v\n", headers) // [Country City AccentCity Region Population Latitude Longitude]

	jobs := make(chan []string)
	docs := make(chan models.Document)

	go func() {
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			panicOnError(err)
			jobs <- line
		}
	}()

	for w := 1; w <= 200; w++ {
		wg.Add(1)
		go processLine( &wg, jobs, docs)
	}

	for w := 1; w <= 200; w++ {
		wg.Add(1)
		go processDoc(&wg, docs, esStore)
	}

	//citiesImported++
	//fmt.Printf("Importing cities successfully completed!\n%d cities has been imported.\n", citiesImported)

	wg.Wait()
	// TODO: check if there are still some jobs to do and clean the mess up.
	//if bulkCounter != 0 {
	//	_, err := bulkRequest.Do() // TODO: extract response, too
	//	panicOnError(err)
	//	bulkCounter = 0
	//}

}

//TODO: rename if perform bulk updates.
func processLine(wg *sync.WaitGroup, jobs <-chan []string, docs chan<- models.Document) {
	defer wg.Done()

	for job := range jobs {
		//TODO: wait until get enough jobs to perform a bulk update request.

		latitude, _ := strconv.ParseFloat(job[5], 64)  // TODO: check for error?
		longitude, _ := strconv.ParseFloat(job[6], 64) // TODO: check for error?

		city := &models.City{
			CountryCode: job[0],
			Name:        job[1], // All names are lowercase -- do something about it?
			AccentName:  job[2], //TODO: handle exotic characters (see the comment above)
			Latitude:    latitude,
			Longitude:   longitude,
		}

		//err := es.AddDocument(city)
		//panicOnError(err)

		docs <- city
	}
}

func processDoc (wg *sync.WaitGroup, docs <-chan models.Document, es *elasticsearch.ElasticStore) {
	defer wg.Done()

	bulkCounter := 0
	bulkRequest := es.Bulk()

	for doc := range docs {
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(es.IndexName).Type(es.TypeName).Id(uuid.NewV4().String()).Doc(doc))
		bulkCounter++

		if bulkCounter%1000 == 0 {
			_, err := bulkRequest.Do() // TODO: extract response, too
			panicOnError(err)
			bulkCounter = 0
		}
	}
}