package database

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var cluster *gocql.ClusterConfig
var cassandraClusterSession *gocql.Session

func connectCassandra() {
	// Connect to Cassandra

	var err error

	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Timeout = 30 * time.Second

	createKeyspace("pretendo_smm")

	cluster.Keyspace = "pretendo_smm"

	cassandraClusterSession, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}

	// Create tables if missing

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.courses (
			data_id bigint PRIMARY KEY,
			playable boolean,
			owner_pid int,
			name text,
			size int,
			creation_date bigint,
			update_date bigint,
			world_record_first_pid int,
			world_record_pid int,
			world_record_creation_date bigint,
			world_record_update_date bigint,
			world_record int,
			meta_binary blob,
			flag int,
			extra_data list<text>,
			data_type smallint,
			period smallint
		)`).Exec(); err != nil {
		fmt.Println("pretendo_smm.courses")
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.ratings (
			data_id bigint PRIMARY KEY,
			stars counter,
			attempts counter,
			failures counter,
			completions counter
		)`).Exec(); err != nil {
		fmt.Println("pretendo_smm.ratings")
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.buffer_queues (
			id uuid PRIMARY KEY,
			data_id bigint,
			slot int,
			buffer blob
		)`).Exec(); err != nil {
		fmt.Println("pretendo_smm.buffer_queues")
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.generator_last_id (
		node_id int PRIMARY KEY,
		last_id int
	)`).Exec(); err != nil {
		fmt.Println("pretendo_smm.generator_last_id")
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.user_play_info (
		pid int PRIMARY KEY,
		starred_courses set<bigint>,
		played_courses set<bigint>
	)`).Exec(); err != nil {
		fmt.Println("pretendo_smm.user_play_info")
		log.Fatal(err)
	}

	fmt.Println("Connected to Cassandra")
}

// Adapted from gocql common_test.go
func createKeyspace(keyspace string) {
	flagRF := flag.Int("rf", 1, "replication factor for pretendo_smm keyspace")

	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 30 * time.Second

	s, err := c.CreateSession()

	if err != nil {
		panic(err)
	}

	defer s.Close()

	if err := s.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, *flagRF)).Exec(); err != nil {
		log.Fatal(err)
	}
}
