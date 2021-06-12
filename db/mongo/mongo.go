package mongo

import (
	"context"
	"github.com/victor-nach/time-tracker/db"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"github.com/victor-nach/time-tracker/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	rateCollection = "rates"
)

type mongoStore struct {
	client *mongo.Client
	dbName string
}

// ensure mongostore implements the datastore interface
var _ db.Datastore = &mongoStore{}

//New returns an instance of mongo store
func New(dbUrl, dbName string) (db.Datastore, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	return &mongoStore{
		client: client,
		dbName: dbName,
	}, client, nil
}

func (m *mongoStore) SaveRate(rate *models.Rate) (*models.Rate, error) {
	_, err := m.client.Database(m.dbName).Collection(rateCollection).InsertOne(context.Background(), rate)
	if err != nil {
		return nil, rerrors.LogFormat(rerrors.DatabaseErr, err)
	}
	return rate, nil
}
