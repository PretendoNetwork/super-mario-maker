package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Endpoint := os.Getenv("PN_SMM_CONFIG_S3_ENDPOINT")
	s3Region := os.Getenv("PN_SMM_CONFIG_S3_REGION")
	s3AccessKey := os.Getenv("PN_SMM_CONFIG_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv("PN_SMM_CONFIG_S3_ACCESS_SECRET")

	staticCredentials := credentials.NewStaticCredentialsProvider(s3AccessKey, s3AccessSecret, "")

	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s3Endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(staticCredentials),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)

	if err != nil {
		panic(err)
	}

	globals.S3Client = s3.NewFromConfig(cfg)

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
