package database

func ConnectAll() {
	connectMongo()
	connectCassandra()
}
