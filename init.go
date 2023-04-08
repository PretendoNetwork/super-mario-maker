package main

import (
	"io/ioutil"
	"log"
	"runtime"

	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/joho/godotenv"
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	globals.HMACSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}

	// Connect to and setup databases
	database.ConnectAll()
	createDataStoreIDGenerators()
}

func createDataStoreIDGenerators() {
	globals.DataStoreIDGenerators = make([]*database.DataStoreIDGenerator, 0)
	regionID := 0 // USA

	for corenum := 0; corenum < runtime.NumCPU(); corenum++ {
		database.CreateDataStoreIDGeneratorRow(corenum)

		lastID := database.GetDataStoreIDGeneratorLastID(corenum)

		generator := database.NewDataStoreIDGenerator(uint8(regionID), uint8(corenum), lastID)
		globals.DataStoreIDGenerators = append(globals.DataStoreIDGenerators, generator)
	}
}
