package main

import (
	"log"
	"os"
	"runtime"

	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Endpoint := os.Getenv("PN_SMM_CONFIG_S3_ENDPOINT")
	//s3Region := os.Getenv("PN_SMM_CONFIG_S3_REGION")
	s3AccessKey := os.Getenv("PN_SMM_CONFIG_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv("PN_SMM_CONFIG_S3_ACCESS_SECRET")

	staticCredentials := credentials.NewStaticV4(s3AccessKey, s3AccessSecret, "")

	minIOClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  staticCredentials,
		Secure: true,
	})
	if err != nil {
		panic(err)
	}

	globals.MinIOClient = minIOClient
	globals.Presigner = globals.NewS3Presigner(globals.MinIOClient)

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
