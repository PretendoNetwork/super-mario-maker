package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/PretendoNetwork/nex-go"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cluster *gocql.ClusterConfig
var cassandraClusterSession *gocql.Session

var mongoClient *mongo.Client
var mongoContext context.Context
var mongoDatabase *mongo.Database
var mongoCollection *mongo.Collection

func connectMongo() {
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	mongoDatabase = mongoClient.Database("pretendo")
	mongoCollection = mongoDatabase.Collection("pnids")
}

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
			stars int,
			world_record_first_pid int,
			world_record_pid int,
			world_record_creation_date bigint,
			world_record_update_date bigint,
			world_record int,
			attempts int,
			failures int,
			completions int,
			meta_binary blob,
			flag int,
			extra_data list<text>,
			data_type smallint,
			period smallint
		)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.buffer_queues (
			id uuid PRIMARY KEY,
			data_id bigint,
			slot int,
			buffer blob
		)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_smm.generator_last_id (
		node_id int PRIMARY KEY,
		last_id int
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	/*
		if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.courses(data_id, meta_binary) VALUES (?, ?) IF NOT EXISTS`, 0, []byte{0, 1, 2, 3}).Exec(); err != nil {
			log.Fatal(err)
		}
	*/

	/*
		type MetaData struct {
			meta_binary []byte
		}
		meta := &MetaData{}
		_ = cassandraClusterSession.Query(`SELECT meta_binary FROM pretendo_smm.courses WHERE data_id=? LIMIT 1`, 0).Scan(&meta)

		fmt.Println(meta)
	*/

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

////////////////////////////////
//                            //
// Cassandra database methods //
//                            //
////////////////////////////////

func createDataStoreIDGeneratorRow(nodeID int) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.generator_last_id(node_id, last_id) VALUES (?, ?) IF NOT EXISTS`, nodeID, 0).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getDataStoreIDGeneratorLastID(nodeID int) uint32 {
	var lastID uint32
	_ = cassandraClusterSession.Query(`SELECT last_id FROM pretendo_smm.generator_last_id WHERE node_id=?`, nodeID).Scan(&lastID)

	return lastID
}

func setDataStoreIDGeneratorLastID(nodeID int, value uint32) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.generator_last_id SET last_id=? WHERE node_id=?`, value, nodeID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func insertCourseDataRow(courseID uint64, ownerPID uint32, size uint32, name string, flag uint32, extraData []string, dataType uint16, period uint16) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.courses(
		data_id,
		owner_pid,
		size,
		name,
		flag,
		extra_data,
		playable,
		creation_date,
		update_date,
		stars,
		world_record_pid,
		world_record_creation_date,
		world_record_update_date,
		attempts,
		failures,
		completions,
		data_type,
		period
	)
	VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		false,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		?,
		?
	) IF NOT EXISTS`,
		courseID, ownerPID, size, name, flag, extraData, dataType, period).Exec(); err != nil {
		log.Fatal(err)
	}
}

func updateCourseMetaBinary(courseID uint64, metaBinary []byte) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET meta_binary=? WHERE data_id=?`, metaBinary, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func setCoursePlayable(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET playable=true WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getCourseMetadatasByLimit(limit uint32) []*CourseMetadata {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT stars, data_id, owner_pid, size, name, meta_binary, flag, data_type, period FROM pretendo_smm.courses LIMIT ?`, limit).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	courseMetadatas := make([]*CourseMetadata, 0)

	for i := 0; i < len(sliceMap); i++ {
		courseMetadata := &CourseMetadata{
			Stars:      uint32(sliceMap[i]["stars"].(int)),
			DataID:     uint64(sliceMap[i]["data_id"].(int64)),
			OwnerPID:   uint32(sliceMap[i]["owner_pid"].(int)),
			Size:       uint32(sliceMap[i]["size"].(int)),
			Name:       sliceMap[i]["name"].(string),
			MetaBinary: sliceMap[i]["meta_binary"].([]byte),
			Flag:       uint32(sliceMap[i]["flag"].(int)),
			DataType:   uint16(sliceMap[i]["data_type"].(int16)),
			Period:     uint16(sliceMap[i]["period"].(int16)),
		}

		courseMetadatas = append(courseMetadatas, courseMetadata)
	}

	return courseMetadatas
}

func getCourseMetadataByDataID(dataID uint64) *CourseMetadata {
	var stars uint32
	var ownerPID uint32
	var size uint32
	var name string
	var metaBinary []byte
	var flag uint32
	var dataType uint16
	var period uint16

	_ = cassandraClusterSession.Query(`SELECT stars, owner_pid, size, name, meta_binary, flag, data_type, period FROM pretendo_smm.courses WHERE data_id=?`, dataID).Scan(&stars, &ownerPID, &size, &name, &metaBinary, &flag, &dataType, &period)

	courseMetadata := &CourseMetadata{
		Stars:      stars,
		DataID:     dataID,
		OwnerPID:   ownerPID,
		Size:       size,
		Name:       name,
		MetaBinary: metaBinary,
		Flag:       flag,
		DataType:   dataType,
		Period:     period,
	}

	return courseMetadata
}

func getCourseMetadataByDataIDs(dataIDs []uint64) []*CourseMetadata {
	// TODO: Do this in one query?
	courseMetadatas := make([]*CourseMetadata, 0)

	for i := 0; i < len(dataIDs); i++ {
		courseMetadatas = append(courseMetadatas, getCourseMetadataByDataID(dataIDs[i]))
	}

	return courseMetadatas
}

func insertBufferQueueData(dataID uint64, slot uint32, buffer []byte) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.buffer_queues( id, data_id, slot, buffer ) VALUES ( now(), ?, ?, ? ) IF NOT EXISTS`, dataID, slot, buffer).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getCourseWorldRecord(dataID uint64) *CourseWorldRecord {
	var worldRecordFirstPID uint32
	var worldRecordPID uint32
	var worldRecordCreatedTime uint64
	var worldRecordUpdatedTime uint64
	var worldRecord int32

	_ = cassandraClusterSession.Query(`SELECT world_record_first_pid, world_record_pid, world_record_creation_date, world_record_update_date, world_record FROM pretendo_smm.courses WHERE data_id=?`, dataID).Scan(&worldRecordFirstPID, &worldRecordPID, &worldRecordCreatedTime, &worldRecordUpdatedTime, &worldRecord)

	if worldRecordFirstPID == 0 {
		return nil
	}

	return &CourseWorldRecord{
		FirstPID:    worldRecordFirstPID,
		BestPID:     worldRecordPID,
		CreatedTime: nex.NewDateTime(worldRecordCreatedTime),
		UpdatedTime: nex.NewDateTime(worldRecordUpdatedTime),
		Score:       worldRecord,
	}
}

func updateCourseWorldRecord(courseID uint64, ownerPID uint32, score int32) {
	now := uint64(time.Now().Unix())

	if getCourseWorldRecord(courseID) == nil {
		if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET world_record_first_pid=?, world_record_creation_date=? WHERE data_id=?`, ownerPID, now, courseID).Exec(); err != nil {
			log.Fatal(err)
		}
	}

	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET world_record_pid=?, world_record_update_date=?, world_record=? WHERE data_id=?`, ownerPID, now, score, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

//////////////////////////////
//                          //
// MongoDB database methods //
//                          //
//////////////////////////////

func getUserMiiInfoByPID(pid uint32) bson.M {
	var result bson.M

	err := mongoCollection.FindOne(context.TODO(), bson.D{{Key: "pid", Value: pid}}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		panic(err)
	}

	return result["mii"].(bson.M)
}
