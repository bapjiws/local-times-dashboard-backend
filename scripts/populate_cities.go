package main

import (
	"timezones_mc/datastore/elasticsearch"
	"timezones_mc/datastore/elasticsearch/configs"
	"timezones_mc/revel_app/app/models"
	"encoding/csv"
	"flag"
	"os"
	"fmt"
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

//TODO: check http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file
//TODO: also check https://gobyexample.com/reading-files
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
	defer file.Close()
	csvReader := csv.NewReader(file)
	//csvReader.LazyQuotes = true

	citiesImported := 0

	headers, _ := csvReader.Read()
	fmt.Printf("Headers: %v\n", headers)

	testLine, _ := csvReader.Read()

	latitude, _ := strconv.ParseFloat(testLine[5], 64) // TODO: check for error?
	longitude, _ := strconv.ParseFloat(testLine[6], 64) // TODO: check for error?

	println(testLine[2])

	city := &models.City{
		CountryCode: testLine[0],
		Name:        testLine[2], //TODO: handle exotic characters!
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
