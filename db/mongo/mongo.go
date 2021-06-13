package mongo

import (
	"context"
	"github.com/victor-nach/time-tracker/db"
	"github.com/victor-nach/time-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	sessionCollection = "sessions"
	usersCollection   = "users"
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

func (m *mongoStore) col(collectionName string) *mongo.Collection {
	return m.client.Database(m.dbName).Collection(collectionName)
}

func (m mongoStore) CreateUser(user *models.User) (*models.User, error) {
	_, err := m.col(usersCollection).
		InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m mongoStore) GetUser(id string) (*models.User, error) {
	user := &models.User{}
	query := bson.M{
		"id": id,
	}
	err := m.col(usersCollection).FindOne(context.Background(), query).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m mongoStore) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := bson.M{
		"email": email,
	}
	err := m.col(usersCollection).FindOne(context.Background(), query).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m mongoStore) GetSession(id, owner string) (*models.Session, error) {
	session := &models.Session{}
	query := bson.M{
		"id":    id,
		"owner": owner,
	}
	err := m.col(sessionCollection).FindOne(context.Background(), query).Decode(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (m mongoStore) GetSessions(owner string, filter string) ([]*models.Session, error) {
	ctx := context.Background()
	query := bson.M{"owner": owner}

	if filter != "nil" {
		now := time.Now()
		var startTime time.Time

		switch filter {
		case "day":
			startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		case "week":
			startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).
				AddDate(0, 0, -int(now.Weekday()))
		case "month":
			startTime = time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local)
		}

		query["ts"] = bson.M{"$gt": startTime.Unix()}
	}

	// Sort by most recent
	findOptions := options.Find().SetSort(bson.M{"ts": -1})
	cursor, err := m.col(sessionCollection).Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}

	var sessions []*models.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (m mongoStore) CreateSession(session *models.Session) (*models.Session, error) {
	_, err := m.col(sessionCollection).
		InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (m mongoStore) UpdateSession(id string, info models.SessionInfo) error {
	filter := bson.M{
		"id": id,
	}
	setQuery := bson.M{}

	if info.Title != nil {
		setQuery["title"] = *info.Title
	}
	if info.Description != nil {
		setQuery["description"] = *info.Description
	}

	query := bson.M{
		"$set": setQuery,
	}

	_, err := m.col(sessionCollection).UpdateOne(context.Background(), filter, query)
	if err != nil {
		return err
	}
	return nil
}

func (m mongoStore) DeleteSession(id string) error {
	filter := bson.M{
		"id": id,
	}
	if _, err := m.col(sessionCollection).DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}
