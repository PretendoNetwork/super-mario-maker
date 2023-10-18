package database

import (
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func initPostgres() {
	_, err := Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS pretendo_smm`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("pretendo_smm Postgres schema created")

	// * Super Mario Maker has non-standard DataID requirements.
	// * DataID 900000 is reserved for the event course metadata
	// * file, so to prevent a collision eventually, we need to
	// * start course IDs AFTER 900000. DataIDs are stored and
	// * processed as uint64, however Super Mario Maker can not
	// * use the full uint64 range. This is because course share
	// * codes are generated from the courses DataID. A course
	// * share code is an 8 byte hex string, where the upper 2
	// * bytes are the checksum of the lower 6 bytes. The lower
	// * 6 bytes are the courses DataID. Super Mario Maker is
	// * only capable of displaying codes up to 0xFFFFFFFFFFFF,
	// * essentially truncating DataIDs down to 48 bit integers
	// * instead of 64 bit. I doubt we will ever hit even the 32
	// * bit limit, let alone 48, but this is here just in case
	_, err = Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS pretendo_smm.courses_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 900001
		CACHE 1`,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_smm.courses (
		data_id int NOT NULL DEFAULT nextval('pretendo_smm.courses_seq') PRIMARY KEY,
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
		meta_binary bytea,
		flag int,
		extra_data text[],
		data_type smallint,
		period smallint
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_smm.ratings (
		data_id bigint PRIMARY KEY,
		stars serial,
		attempts serial,
		failures serial,
		completions serial
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_smm.buffer_queues (
		id uuid PRIMARY KEY,
		data_id bigint,
		slot int,
		buffer bytea
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_smm.user_play_info (
		pid int PRIMARY KEY,
		starred_courses bigint[],
		played_courses bigint[]
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres tables created")
}
