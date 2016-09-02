package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"timezones_mc/datastore/elasticsearch"
	"timezones_mc/datastore/elasticsearch/configs"
	"timezones_mc/revel_app/app/models"

	"github.com/satori/go.uuid"
	"gopkg.in/olivere/elastic.v2"
)

// TODO: create a utils folder and move it there?
func panicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}

var (
	fileFlag        = flag.String("file", "", "file to parse")
	wg              sync.WaitGroup
	start           time.Time
	citiesRead      uint64 = 0
	citiesProcessed uint64 = 0
	bulkCounter     uint64 = 0
)

// go run scripts/cities/fan_in.go -file="cities/worldcities.txt"
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

	records := recordGenerator(csvReader)
	bulkRequest := esStore.Bulk()

	pipe := mergeCityChannels(
		getCityChan(records),
		getCityChan(records),
		getCityChan(records),
		getCityChan(records),
	)
	for city := range pipe {
		atomic.AddUint64(&citiesProcessed, 1)
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(esStore.IndexName).Type(esStore.TypeName).Id(uuid.NewV4().String()).Doc(city))
		atomic.AddUint64(&bulkCounter, 1)

		if counter := atomic.LoadUint64(&bulkCounter); counter%1000 == 0 {
			_, err := bulkRequest.Do()
			panicOnError(err)
			atomic.StoreUint64(&bulkCounter, 0)
		}

	}

	if counter := atomic.LoadUint64(&bulkCounter); counter != 0 {
		_, err := bulkRequest.Do()
		panicOnError(err)
	}

	citiesRead := atomic.LoadUint64(&citiesRead)
	citiesProcessed := atomic.LoadUint64(&citiesProcessed)
	elapsed := time.Since(start)

	fmt.Printf("Imported %d cities out of %d in %s \n", citiesProcessed, citiesRead, elapsed)
}

func recordGenerator(csvReader *csv.Reader) <-chan []string {
	records := make(chan []string, 1000)

	go func() {
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				close(records) // That's it, folks!
				break
			}
			panicOnError(err)
			records <- line
		}
	}()

	return records
}

func getCityChan(records <-chan []string) chan *models.City {
	cities := make(chan *models.City, 1000)

	go func() {
		for record := range records {
			latitude, _ := strconv.ParseFloat(record[5], 64)
			longitude, _ := strconv.ParseFloat(record[6], 64)

			city := &models.City{
				CountryCode: record[0],
				Name:        record[1], // TODO: All names are lowercase -- do something about it?
				AccentName:  record[2],
				Latitude:    latitude,
				Longitude:   longitude,
			}

			atomic.AddUint64(&citiesRead, 1)

			cities <- city
		}

		close(cities)
	}()

	return cities

}

func mergeCityChannels(cityChannels ...chan *models.City) chan *models.City {
	pipe := make(chan *models.City, 1000)

	output := func(cityChan <-chan *models.City) {
		for city := range cityChan {
			pipe <- city
		}
		wg.Done()
	}

	wg.Add(len(cityChannels))
	for _, cityChan := range cityChannels {
		go output(cityChan)
	}

	go func() {
		wg.Wait()
		close(pipe)
	}()

	return pipe
}