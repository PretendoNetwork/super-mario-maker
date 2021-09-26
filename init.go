package main

import (
	"io/ioutil"
	"runtime"
)

/*
type Config struct {
	Mongo struct {
	}
	Cassandra struct{}
}
*/

var hmacSecret []byte
var dataStoreIDGenerators []*DataStoreIDGenerator

func init() {
	var err error

	hmacSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}

	// Connect to and setup databases
	connectMongo()
	connectCassandra()
	createDataStoreIDGenerators()
	//insertCourseDataRow(dataID, param.Size, param.Name, param.Flag, param.ExtraData)
	//updateCourseMetaBinary(dataID, param.MetaBinary)
	//fmt.Printf("%+v", getCourseMetadataByDataID(288230376151711744))
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
