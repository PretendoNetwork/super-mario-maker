package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
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

func initializeCourseData(courseID uint64, ownerPID uint32, size uint32, name string, flag uint32, extraData []string, dataType uint16, period uint16) {
	datetime := nex.NewDateTime(0)
	now := datetime.Now()

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
		world_record_first_pid,
		world_record_pid,
		world_record_creation_date,
		world_record_update_date,
		world_record,
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
		?,
		?,
		0,
		0,
		0,
		0,
		0,
		?,
		?
	) IF NOT EXISTS`,
		courseID, ownerPID, size, name, flag, extraData, now, now, dataType, period).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET stars=stars+0, attempts=attempts+0, failures=failures+0, completions=completions+0 WHERE data_id=?`, courseID).Exec(); err != nil {
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

	if sliceMap, err = cassandraClusterSession.Query(`SELECT data_id, owner_pid, size, name, meta_binary, flag, creation_date, update_date, data_type, period FROM pretendo_smm.courses LIMIT ?`, limit).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	courseMetadatas := make([]*CourseMetadata, 0)

	for i := 0; i < len(sliceMap); i++ {
		dataID := uint64(sliceMap[i]["data_id"].(int64))

		var stars uint32
		var attempts uint32
		var failures uint32
		var completions uint32

		_ = cassandraClusterSession.Query(`SELECT stars, attempts, failures, completions FROM pretendo_smm.ratings WHERE data_id=?`, dataID).Scan(&stars, &attempts, &failures, &completions)

		courseMetadata := &CourseMetadata{
			DataID:      dataID,
			OwnerPID:    uint32(sliceMap[i]["owner_pid"].(int)),
			Size:        uint32(sliceMap[i]["size"].(int)),
			CreatedTime: nex.NewDateTime(uint64(sliceMap[i]["creation_date"].(int64))),
			UpdatedTime: nex.NewDateTime(uint64(sliceMap[i]["update_date"].(int64))),
			Name:        sliceMap[i]["name"].(string),
			MetaBinary:  sliceMap[i]["meta_binary"].([]byte),
			Stars:       stars,
			Attempts:    attempts,
			Failures:    failures,
			Completions: completions,
			Flag:        uint32(sliceMap[i]["flag"].(int)),
			DataType:    uint16(sliceMap[i]["data_type"].(int16)),
			Period:      uint16(sliceMap[i]["period"].(int16)),
		}

		courseMetadatas = append(courseMetadatas, courseMetadata)
	}

	return courseMetadatas
}

func getCourseMetadataByDataID(dataID uint64) *CourseMetadata {
	var ownerPID uint32
	var size uint32
	var name string
	var metaBinary []byte
	var flag uint32
	var createdTime uint64
	var updatedTime uint64
	var dataType uint16
	var period uint16

	_ = cassandraClusterSession.Query(`SELECT owner_pid, size, name, meta_binary, flag, creation_date, update_date, data_type, period FROM pretendo_smm.courses WHERE data_id=?`, dataID).Scan(&ownerPID, &size, &name, &metaBinary, &flag, &createdTime, &updatedTime, &dataType, &period)

	var stars uint32
	var attempts uint32
	var failures uint32
	var completions uint32

	_ = cassandraClusterSession.Query(`SELECT stars, attempts, failures, completions FROM pretendo_smm.ratings WHERE data_id=?`, dataID).Scan(&stars, &attempts, &failures, &completions)

	courseMetadata := &CourseMetadata{
		DataID:      dataID,
		OwnerPID:    ownerPID,
		Size:        size,
		CreatedTime: nex.NewDateTime(createdTime),
		UpdatedTime: nex.NewDateTime(updatedTime),
		Name:        name,
		MetaBinary:  metaBinary,
		Stars:       stars,
		Attempts:    attempts,
		Failures:    failures,
		Completions: completions,
		Flag:        flag,
		DataType:    dataType,
		Period:      period,
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

func getCourseMetadatasByPID(pid uint32) []*CourseMetadata {
	courseMetadatas := make([]*CourseMetadata, 0)

	// TODO: Fix this query? Seems like a weird way of doing this...
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT data_id FROM pretendo_smm.courses WHERE owner_pid=? ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	for _, course := range sliceMap {
		dataID := uint64(course["data_id"].(int64))
		courseMetadatas = append(courseMetadatas, getCourseMetadataByDataID(dataID))
	}

	return courseMetadatas
}

func insertBufferQueueData(dataID uint64, slot uint32, buffer []byte) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.buffer_queues( id, data_id, slot, buffer ) VALUES ( now(), ?, ?, ? ) IF NOT EXISTS`, dataID, slot, buffer).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getBufferQueueDeathData(dataID uint64) [][]byte {
	pBufferQueue := make([][]byte, 0)

	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT buffer FROM pretendo_smm.buffer_queues WHERE data_id=? AND slot=3 ALLOW FILTERING`, dataID).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(sliceMap); i++ {
		pBufferQueue = append(pBufferQueue, sliceMap[i]["buffer"].([]byte))
	}

	return pBufferQueue
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
	datetime := nex.NewDateTime(0)
	now := datetime.Now()

	if getCourseWorldRecord(courseID) == nil {
		if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET world_record_first_pid=?, world_record_creation_date=? WHERE data_id=?`, ownerPID, now, courseID).Exec(); err != nil {
			log.Fatal(err)
		}
	}

	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET world_record_pid=?, world_record_update_date=?, world_record=? WHERE data_id=?`, ownerPID, now, score, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func incrementCourseClearCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET completions=completions+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func incrementCourseStarCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET stars=stars+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func incrementCourseFailCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET failures=failures+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func incrementCourseAttemptCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET attempts=attempts+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getUserStarredCourses(pid uint32) []*CourseMetadata {
	var dataIDs []uint64
	_ = cassandraClusterSession.Query(`SELECT starred_courses FROM pretendo_smm.user_play_info WHERE pid=?`, pid).Scan(&dataIDs)

	return getCourseMetadataByDataIDs(dataIDs)
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
