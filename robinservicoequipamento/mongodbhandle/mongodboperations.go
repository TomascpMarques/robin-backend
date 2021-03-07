package mongodbhandle

import "go.mongodb.org/mongo-driver/mongo"

// GetMongoDatabase -
func GetMongoDatabase(cl *mongo.Client, dbName string) *mongo.Database {
	return cl.Database(dbName)
}

// GetMongoCollection -
func GetMongoCollection(db *mongo.Database, collName string) *mongo.Collection {
	return db.Collection(collName)
}

