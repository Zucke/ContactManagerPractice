package data

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	data *Data
	once sync.Once
)

//Data have the db info an client
type Data struct {
	client  *mongo.Client
	context context.Context
}

//Close finish the db connection
func (d *Data) Close() error {
	return d.client.Disconnect(context.Background())
}

//DBCollection return a reference to a collection
func (d *Data) DBCollection(name string) *mongo.Collection {
	return d.client.Database(DBNAME).Collection(name)
}

//New return the db connection and init if not was created
func New() *Data {
	once.Do(initData)
	return data
}

func initData() {
	dbclient, err := createDBSession()
	if err != nil {
		log.Fatal(err)
	}

	data = &Data{
		client:  dbclient,
		context: context.Background(),
	}

}
func createDBSession() (*mongo.Client, error) {
	var connString = os.Getenv("DATABASE_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if connString == "" {
		connString = "mongodb://127.0.0.1:27017"
	}

	return mongo.Connect(ctx, options.Client().ApplyURI(connString))

}

//Close call te Close method of Data if data exist
func Close() error {
	if data == nil {
		return nil
	}

	return data.Close()

}
