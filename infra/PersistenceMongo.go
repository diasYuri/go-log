package infra

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"lfa.com/logs-master/application/domain"
)

type PersistenceMongo struct {
	clientMongo *mongo.Client
	collection  *mongo.Collection
}

func NewPersistenceMongo(client *mongo.Client, database string, collection string) *PersistenceMongo {
	return &PersistenceMongo{
		clientMongo: client,
		collection:  client.Database(database).Collection(collection),
	}
}

func (p *PersistenceMongo) Save(log domain.Log) error {
	_, err := p.collection.InsertOne(context.TODO(), log)

	return err
}

func (p *PersistenceMongo) SaveBatch(logs []interface{}) error {
	_, err := p.collection.InsertMany(context.TODO(), logs)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
