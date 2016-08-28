package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"timezones_mc/datastore/elasticsearch"
	"timezones_mc/datastore/elasticsearch/configs"
	"timezones_mc/revel_app/app/models"
	//"io"
	"strconv"
)

// TODO: create a utils folder and move it there?
func panicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}

var fileFlag = flag.String("file", "", "file to parse")

// Unhandled characters: https://en.wikipedia.org/wiki/%C3%80
func main() {
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
	//csvReader.LazyQuotes = true

	citiesImported := 0

	headers, err := csvReader.Read()
	panicOnError(err)
	fmt.Printf("Headers: %v\n", headers) // [Country City AccentCity Region Population Latitude Longitude]

	testLine, err := csvReader.Read()
	panicOnError(err)

	latitude, _ := strconv.ParseFloat(testLine[5], 64)  // TODO: check for error?
	longitude, _ := strconv.ParseFloat(testLine[6], 64) // TODO: check for error?

	println(testLine[2])

	city := &models.City{
		CountryCode: testLine[0],
		Name:        testLine[1], // All names are lowercase -- do something about it?
		AccentName:  testLine[2], //TODO: handle exotic characters (see the comment above)
		Latitude:    latitude,
		Longitude:   longitude,
	}

	err = esStore.AddDocument(city)
	panicOnError(err)

	citiesImported++

	//for {
	//	line, err := csvReader.Read()
	//	if err == io.EOF {
	//		break
	//	}
	//	panicOnError(err) // TODO: use another util function?
	//	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	//	//fmt.Println(record)
	//
	//	city := &models.City{
	//		CountryCode: ,
	//		Name:        ,
	//		Latitude:    ,
	//		Longitude:   ,
	//	}
	//
	//	err = esStore.AddDocument(city)
	//	panicOnError(err)
	//
	//	citiesImported++
	//}

	fmt.Printf("Importing cities successfully completed!\n%d cities has been imported.\n", citiesImported)
}

//tempCity := models.City{
//CountryCode: "us",
//Name:        "NYC",
//Latitude:    123.45,
//Longitude:   -123.45,
//}
