package infra

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"lfa.com/logs-master/application/domain"
	"log"
	"time"
)

type PersistenceMongo struct {
	clientMongo *mongo.Client
	//collection  *mongo.Collection
	database   string
	collection string
}

func NewPersistenceMongo(client *mongo.Client, database string, collection string) *PersistenceMongo {
	return &PersistenceMongo{
		clientMongo: client,
		//collection:  client.Database(database).Collection(collection),
		database:   database,
		collection: collection,
	}
}

func (p *PersistenceMongo) Save(log domain.Log) error {
	_, err := p.clientMongo.Database(p.database).Collection(p.collection).InsertOne(context.TODO(), log)

	return err
}

func (p *PersistenceMongo) SaveBatch(logs []interface{}) error {
	defer TimeTrack(time.Now(), "Save batch")
	coll := p.clientMongo.Database(p.database).Collection(p.collection)

	_, err := coll.InsertMany(context.TODO(), logs)
	if err != nil {
		log.Println(err)
	}

	return err
}
