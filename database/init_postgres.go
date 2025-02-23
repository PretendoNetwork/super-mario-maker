package database

import (
	"os"

	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func initPostgres() {
	_, err := Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS datastore`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("datastore Postgres schema created")

	// * Super Mario Maker has non-standard DataID requirements.
	// * DataID 900000 is reserved for the event course metadata
	// * file, and official event courses begin at DataID 930010
	// * and end at DataID 930050. To prevent a collision
	// * eventually, we need to start course IDs AFTER 930050
	// *
	// * DataIDs are stored and processed as uint64, however
	// * Super Mario Maker can not use the full uint64 range.
	// * This is because course share codes are generated from the
	// * courses DataID. A course share code is an 8 byte hex
	// * string, where the upper 2 bytes are the checksum of the
	// * lower 6 bytes. The lower 6 bytes are the courses DataID
	// *
	// * Super Mario Maker is only capable of displaying codes up
	// * to 0xFFFFFFFFFFFF, essentially truncating DataIDs down to
	// * 48 bit integers instead of 64 bit. I doubt we will ever
	// * hit even the 32 bit limit, let alone 48, but this is here
	// * just in case
	_, err = Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.object_data_id_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 940000
		CACHE 1`, // * Honestly I don't know what CACHE does but I saw it recommended so here it is
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	// * "deletion_reason" and "under_review" are specific to SMM.
	// * Everything else is stock
	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.objects (
		data_id bigint NOT NULL DEFAULT nextval('datastore.object_data_id_seq') PRIMARY KEY,
		upload_completed boolean NOT NULL DEFAULT FALSE,
		deleted boolean NOT NULL DEFAULT FALSE,
		deletion_reason int NOT NULL DEFAULT 0,
		under_review boolean NOT NULL DEFAULT FALSE,
		owner int,
		size int,
		name text,
		data_type int,
		meta_binary bytea,
		permission int,
		permission_recipients int[],
		delete_permission int,
		delete_permission_recipients int[],
		flag int,
		period int,
		refer_data_id bigint,
		tags text[],
		persistence_slot_id int,
		extra_data text[],
		access_password bigint NOT NULL DEFAULT 0,
		update_password bigint NOT NULL DEFAULT 0,
		creation_date timestamp,
		update_date timestamp
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	// * Unsure what like half of this is but the client sends it so we saves it
	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.object_ratings (
		data_id bigint,
		slot smallint,
		flag smallint,
		internal_flag smallint,
		lock_type smallint,
		initial_value bigint,
		range_min int,
		range_max int,
		period_hour smallint,
		period_duration int,
		total_value bigint,
		count int NOT NULL DEFAULT 0,
		PRIMARY KEY(data_id, slot)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	// * Custom rankings are specific to SMM
	// TODO - Store the period? What even is the period of custom rankings?
	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.object_custom_rankings (
		data_id bigint,
		application_id bigint,
		value bigint,
		PRIMARY KEY(data_id, application_id)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	// * BufferQueues are specific to SMM
	// * Real server does not allow duplicate buffers in a given slot for an object,
	// * even if uploaded by different users. We could change this, but I don't see
	// * much point
	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.buffer_queues (
		data_id bigint,
		slot int,
		creation_date timestamp,
		buffer bytea,
		PRIMARY KEY(data_id, slot, buffer)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	// * Course records are specific to SMM
	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.course_records (
		data_id bigint,
		slot int,
		first_pid int,
		best_pid int,
		best_score int,
		creation_date timestamp,
		update_date timestamp,
		PRIMARY KEY(data_id, slot)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("Postgres tables created")

	ensureEventCourseMetaDataFileExists()
}

func ensureEventCourseMetaDataFileExists() {
	// * The event course metadata file is REQUIRED for
	// * Super Mario Maker to open Course World. This
	// * ensures that it exists both in S3 and in the
	// * database before booting. It has the reserved
	// * DataID 900000
	// bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	// key := "900000.bin"

	// objectSizeS3, err := globals.S3ObjectSize(bucket, key)
	// if err != nil {
	// 	globals.Logger.Errorf("Failed to stat event course metadata file. Ensure your S3 credentials are correct and the 900000.bin file is uploaded to your bucket. S3 error: %s", err.Error())
	// 	os.Exit(0)
	// }

	// globals.Logger.Success("Event course metadata file found. Verifying database")

	// var exists bool
	// err = Postgres.QueryRow(`SELECT EXISTS(SELECT 1 FROM datastore.objects WHERE data_id=900000) AS "exists"`).Scan(&exists)
	// if err != nil {
	// 	globals.Logger.Errorf("Error querying for event course metadata object in Postgres: %s", err.Error())
	// 	os.Exit(0)
	// }

	// now := time.Now()

	// if !exists {
	// 	var dataID uint64

	// 	globals.Logger.Info("Event course metadata object not found in Postgres. Creating")

	// 	err := Postgres.QueryRow(`INSERT INTO datastore.objects (
	// 		data_id,
	// 		upload_completed,
	// 		owner,
	// 		size,
	// 		name,
	// 		data_type,
	// 		meta_binary,
	// 		permission,
	// 		permission_recipients,
	// 		delete_permission,
	// 		delete_permission_recipients,
	// 		flag,
	// 		period,
	// 		refer_data_id,
	// 		tags,
	// 		persistence_slot_id,
	// 		extra_data,
	// 		creation_date,
	// 		update_date
	// 	) VALUES (
	// 		$1,
	// 		$2,
	// 		$3,
	// 		$4,
	// 		$5,
	// 		$6,
	// 		$7,
	// 		$8,
	// 		$9,
	// 		$10,
	// 		$11,
	// 		$12,
	// 		$13,
	// 		$14,
	// 		$15,
	// 		$16,
	// 		$17,
	// 		$18,
	// 		$19
	// 	) RETURNING data_id`,
	// 		900000,
	// 		true,
	// 		2, // * "Quazal Rendez-Vous" special account
	// 		objectSizeS3,
	// 		"",       // * Has no name
	// 		50,       // * Metadata file has DataType 50. Event courses have DataType 51
	// 		[]byte{}, // * No MetaBinary
	// 		0,        // * Accessible by everyone
	// 		pq.Array([]uint32{}),
	// 		3, // * THE REAL SERVER HAS THIS SET TO 0, FOR EVERYONE. THAT'S SUPER INSECURE.
	// 		pq.Array([]uint32{}),
	// 		0,
	// 		64306,
	// 		0,
	// 		pq.Array([]string{}),
	// 		0, // * Unsure what the slot ID actually is
	// 		pq.Array([]string{}),
	// 		now,
	// 		now,
	// 	).Scan(&dataID)
	// 	if err != nil {
	// 		globals.Logger.Errorf("Error creating event course metadata object: %s", err.Error())
	// 		os.Exit(0)
	// 	}
	// } else {
	// 	var objectSizeDB uint32

	// 	err := Postgres.QueryRow(`SELECT size FROM datastore.objects WHERE data_id=900000`).Scan(&objectSizeDB)
	// 	if err != nil {
	// 		globals.Logger.Errorf("Error querying event course metadata object size: %s", err.Error())
	// 		os.Exit(0)
	// 	}

	// 	if objectSizeS3 != uint64(objectSizeDB) {
	// 		globals.Logger.Success("Event course metadata object found in Postgres! Updating size")

	// 		_, err := Postgres.Exec(`UPDATE datastore.objects SET size=$1, update_date=$2 WHERE data_id=900000`, objectSizeS3, now)
	// 		if err != nil {
	// 			globals.Logger.Errorf("Error updating event course metadata object size: %s", err.Error())
	// 			os.Exit(0)
	// 		}
	// 	} else {
	// 		globals.Logger.Success("Event course metadata object found in Postgres!")
	// 	}

	// 	globals.Logger.Info("Ensuring event course metadata object has correct delete permission")

	// 	_, err = Postgres.Exec(`UPDATE datastore.objects SET delete_permission=3 WHERE data_id=900000`)
	// 	if err != nil {
	// 		globals.Logger.Errorf("Error updating event course metadata object delete permission: %s", err.Error())
	// 		os.Exit(0)
	// 	}
	// }

	globals.Logger.Success("Event course metadata object found and is up to date!")
}
