package main

import (
	"io/ioutil"
	"log"
	"runtime"

	"github.com/joho/godotenv"
)

var hmacSecret []byte
var dataStoreIDGenerators []*DataStoreIDGenerator

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hmacSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}

	// Connect to and setup databases
	connectMongo()
	connectCassandra()
	createDataStoreIDGenerators()
}

func createDataStoreIDGenerators() {
	dataStoreIDGenerators = make([]*DataStoreIDGenerator, 0)
	regionID := 0 // USA

	for corenum := 0; corenum < runtime.NumCPU(); corenum++ {
		createDataStoreIDGeneratorRow(corenum)

		lastID := getDataStoreIDGeneratorLastID(corenum)

		generator := NewDataStoreIDGenerator(uint8(regionID), uint8(corenum), lastID)
		dataStoreIDGenerators = append(dataStoreIDGenerators, generator)
	}
}
