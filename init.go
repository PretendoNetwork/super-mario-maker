package main

import "io/ioutil"

/*
type Config struct {
	Mongo struct {
	}
	Cassandra struct{}
}
*/

var hmacSecret []byte

func init() {
	var err error

	hmacSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}

	// Connect to and setup databases
	connectMongo()
	connectCassandra()
}
