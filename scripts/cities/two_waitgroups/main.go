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
	"time"
	//"log"
	"sync/atomic"
)

// TODO: create a utils folder and move it there?
func panicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}

var (
	fileFlag = flag.String("file", "", "file to parse")

	jobWG sync.WaitGroup
	docWG sync.WaitGroup

	start time.Time
	//counter uint64

	citiesRead      uint64 = 0
	citiesProcessed uint64 = 0
)



// go run scripts/cities/two_waitgroups/main.go -file="cities/worldcities.txt"
func main() {
	start = time.Now()

	flag.Parse()

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

	headers, err := csvReader.Read()
	panicOnError(err)
	fmt.Printf("Headers: %v\n", headers) // [Country City AccentCity Region Population Latitude Longitude]

	jobs := make(chan []string, 1000)
	docs := make(chan models.Document, 1000)

	go func() {
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				close(jobs) // That's it, folks!
				break
			}
			panicOnError(err)
			jobs <- line
		}
	}()

	for w := 1; w <= 4; w++ {
		jobWG.Add(1)
		go processLine(&jobWG, jobs, docs)
	}

	for w := 1; w <= 4; w++ {
		docWG.Add(1)
		go processDoc(&docWG, docs, esStore)
	}

	//citiesImported++
	//fmt.Printf("Importing cities successfully completed!\n%d cities has been imported.\n", citiesImported)

	//go func() {
	//	jobWG.Wait()
	//	println("close docs")
	//	close(docs)
	//}()

	jobWG.Wait()
	close(docs)

	bulkRequest := esStore.Bulk()
	for doc := range docs {
		println(doc.String())
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(esStore.IndexName).Type(esStore.TypeName).Id(uuid.NewV4().String()).Doc(doc))
	}
	_, err = bulkRequest.Do()
	panicOnError(err)

	docWG.Wait()

	// TODO: remove all below
	//go func() {
	//	docWG.Wait()
	//	println("HERE")
	//}()

	//if bulkCounter != 0 {
	//	_, err := bulkRequest.Do() // TODO: extract response, too
	//	panicOnError(err)
	//	bulkCounter = 0
	//}

	citiesRead := atomic.LoadUint64(&citiesRead)
	citiesProcessed := atomic.LoadUint64(&citiesProcessed)
	elapsed := time.Since(start)

	fmt.Printf("Imported %d cities out of %d in %s \n", citiesProcessed, citiesRead, elapsed)

}

func processLine(jobWG *sync.WaitGroup, jobs <-chan []string, docs chan<- models.Document) {
	defer jobWG.Done()

	for job := range jobs {
		latitude, _ := strconv.ParseFloat(job[5], 64)
		longitude, _ := strconv.ParseFloat(job[6], 64)

		city := &models.City{
			CountryCode: job[0],
			Name:        job[1], // TODO: All names are lowercase -- do something about it?
			AccentName:  job[2],
			Latitude:    latitude,
			Longitude:   longitude,
		}

		atomic.AddUint64(&citiesRead, 1)
		docs <- city
	}
}

func processDoc (docWG *sync.WaitGroup, docs <-chan models.Document, es *elasticsearch.ElasticStore) {
	defer docWG.Done()

	bulkCounter := 0
	bulkRequest := es.Bulk()

	for doc := range docs {
		atomic.AddUint64(&citiesProcessed, 1)

		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(es.IndexName).Type(es.TypeName).Id(uuid.NewV4().String()).Doc(doc))
		bulkCounter++

		if bulkCounter%1000 == 0 {
			_, err := bulkRequest.Do()
			panicOnError(err)
			bulkCounter = 0
		}


		//atomic.AddUint64(&counter, 1)
		//check := atomic.LoadUint64(&counter)
		//if check == 3173958 {
		//	elapsed := time.Since(start)
		//	log.Printf("Took %s", elapsed)
		//	os.Exit(0)
		//}

	}
}