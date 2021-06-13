package mongo

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/victor-nach/time-tracker/lib/ulid"
	"github.com/victor-nach/time-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
	"time"
)

var (
	dbName = "tracker"
)
var mockData = struct {
	models.User
	models.Session
}{
	User: models.User{
		ID:       "userID",
		Email:    "victor@email.com",
		Name:     "firstname lastName",
		Password: "hashedPasscode",
		Ts:       time.Now().Unix(),
	},
	Session: models.Session{
		ID:          "sessionID",
		Owner:       "userID",
		Title:       "title",
		Description: "session description",
		Start:       time.Now().Unix(),
		End:         time.Now().Unix(),
		Duration:    24 * time.Hour.Milliseconds(),
		Ts:          time.Now().Unix(),
	},
}

var mongoDbPort = ""

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}

	resource, err := pool.Run("mongo", "4.2.9", []string{
		"MONGO_INITDB_DATABASE=tracker",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	mongoDbPort = resource.GetPort("27017/tcp")
	if err := pool.Retry(func() error {
		var err error
		connectUrl := fmt.Sprintf("mongodb://localhost:%s", mongoDbPort)
		_, _, err = New(connectUrl, "roava")
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	code := m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestMongoStore_GetUser(t *testing.T) {
	const (
		getByEmail = iota
		getById
		errorEmailNotFound
		errorUserIdNotFound
	)

	var tests = []struct {
		name     string
		arg      string
		testType int
	}{
		{
			name:     "Successfully get user by email",
			arg:      mockData.User.Email,
			testType: getByEmail,
		},
		{
			name:     "Successfully get user by ID",
			arg:      mockData.User.ID,
			testType: getById,
		},
		{
			name:     "Test email not found",
			arg:      "invalidEmail",
			testType: errorEmailNotFound,
		},
		{
			name:     "Test user Id not found",
			arg:      "invalidEmail",
			testType: errorUserIdNotFound,
		},
	}
	connectUri := "mongodb://localhost:" + mongoDbPort
	dataStore, client, err := New(connectUri, "tracker")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	_, err = client.Database(dbName).Collection(usersCollection).InsertOne(context.Background(), mockData.User)
	assert.Nil(t, err)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			switch testCase.testType {

			case getById:
				user, err := dataStore.GetUser(testCase.arg)
				assert.NotNil(t, user)
				assert.NoError(t, err)
			case getByEmail:
				user, err := dataStore.GetUserByEmail(testCase.arg)
				assert.NotNil(t, user)
				assert.NoError(t, err)
			case errorEmailNotFound:
				user, err := dataStore.GetUser(testCase.arg)
				assert.Nil(t, user)
				assert.Error(t, err)
			case errorUserIdNotFound:
				user, err := dataStore.GetUser(testCase.arg)
				assert.Nil(t, user)
				assert.Error(t, err)
			}
		})
	}

	//if err := client.Disconnect(context.Background()); err != nil {
	//	t.Fatal(err)
	//}
}

func TestMongoStore_CreateUser(t *testing.T) {
	connectUri := "mongodb://localhost:" + mongoDbPort
	dataStore, client, err := New(connectUri, "tracker")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	mockUser := mockData.User
	mockUser.ID = "userId2"
	mockUser.Email = "userEmail2"

	// assert user not found
	u, err := dataStore.GetUser(mockUser.ID)
	assert.Error(t, err)
	assert.Nil(t, u)

	// test create user
	usr, err := dataStore.CreateUser(&mockUser)
	assert.Nil(t, err)
	assert.NotNil(t, usr)

	user, err := dataStore.GetUser(mockUser.ID)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, mockUser, *user)
}

func TestMongoStore_CreateSession(t *testing.T) {
	connectUri := "mongodb://localhost:" + mongoDbPort
	dataStore, client, err := New(connectUri, "tracker")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	session, err := dataStore.CreateSession(&mockData.Session)
	assert.Nil(t, err)
	assert.NotNil(t, session)

	s := &models.Session{}
	err = client.Database(dbName).Collection(sessionCollection).
		FindOne(context.Background(), bson.M{"id": mockData.Session.ID}).Decode(s)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, session, s)

	if err := client.Disconnect(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestMongoStore_GetSession(t *testing.T) {
	const (
		getSingleSession = iota
		errorSessionNotFound
	)

	var tests = []struct {
		name     string
		testType int
	}{
		{
			name:     "Get single session",
			testType: getSingleSession,
		},
		{
			name:     "Test user session not found",
			testType: errorSessionNotFound,
		},
	}
	connectUri := "mongodb://localhost:" + mongoDbPort
	dataStore, client, err := New(connectUri, "tracker")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	// seed owner
	mockUser := mockData.User
	mockUser.ID = ulid.New().Generate()
	mockUser.Email = "testt@mail2.com"

	mockSession := mockData.Session
	mockSession.ID = ulid.New().Generate()
	mockSession.Owner = mockUser.ID

	_, err = client.Database(dbName).Collection(usersCollection).InsertOne(context.Background(), mockUser)
	assert.Nil(t, err)

	_, err = client.Database(dbName).Collection(sessionCollection).InsertOne(context.Background(), mockSession)
	assert.Nil(t, err)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			switch testCase.testType {
			case errorSessionNotFound:
				user, err := dataStore.GetSession("invalidId", mockUser.ID)
				assert.Nil(t, user)
				assert.Error(t, err)
			case getSingleSession:
				session, err := dataStore.GetSession(mockSession.ID, mockSession.Owner)
				assert.NotNil(t, session)
				assert.Nil(t, err)
			}
		})
	}
}

func TestMongoStore_GetSessions(t *testing.T) {
	const (
		getAllSessions = iota
		getSessionsByDay
		getSessionsByWeek
		getSessionsByMonth
	)

	var tests = []struct {
		name     string
		testType int
		filter string
		expectedLength int
	}{
		{
			name:     "Get all sessions",
			testType: getAllSessions,
			filter: "",
			expectedLength: 2,
		},
		{
			name:     "Get sessions by day",
			testType: getSessionsByDay,
			filter: "day",
			expectedLength: 1,
		},
		{
			name:     "Get sessions by week",
			testType: getSessionsByWeek,
			filter: "week",
			expectedLength: 1,
		},
		{
			name:     "Get sessions by month",
			testType: getSessionsByMonth,
			filter: "month",
			expectedLength: 2,
		},
	}


	connectUri := "mongodb://localhost:" + mongoDbPort
	dataStore, client, err := New(connectUri, "tracker")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	// seed owner
	mockUser := mockData.User
	mockUser.ID = ulid.New().Generate()
	mockUser.Email = "testt@mail2.com"

	mockSession := mockData.Session
	mockSession.ID = ulid.New().Generate()
	mockSession.Owner = mockUser.ID

	mockSession2 := mockSession
	mockSession2.Ts = time.Now().AddDate(0,0,-10).Unix() // 10 days ago

	_, err = client.Database(dbName).Collection(usersCollection).InsertOne(context.Background(), mockUser)
	assert.Nil(t, err)

	_, err = client.Database(dbName).Collection(sessionCollection).InsertOne(context.Background(), mockSession)
	assert.Nil(t, err)

	_, err = client.Database(dbName).Collection(sessionCollection).InsertOne(context.Background(), mockSession2)
	assert.Nil(t, err)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			sessions, err := dataStore.GetSessions(mockSession.Owner, testCase.filter)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedLength, len(sessions))
		})
	}
}

func TestMongoStore_ManageSession(t *testing.T) {
	connectUri := "mongodb://localhost:" + mongoDbPort
	dataStore, client, err := New(connectUri, "tracker")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	mockUser := mockData.User
	mockUser.ID = ulid.New().Generate()
	mockUser.Email = "victor@mail2.com"

	mockSession := mockData.Session
	mockSession.ID = ulid.New().Generate()
	mockSession.Owner = mockUser.ID

	_, err = client.Database(dbName).Collection(usersCollection).InsertOne(context.Background(), mockUser)
	assert.Nil(t, err)
	_, err = client.Database(dbName).Collection(sessionCollection).InsertOne(context.Background(), mockSession)
	assert.Nil(t, err)

	// test update session
	title := "new title"
	err = dataStore.UpdateSession(mockSession.ID, models.SessionInfo{Title: &title})
	assert.NoError(t, err)

	s, err := dataStore.GetSession(mockSession.ID, mockSession.Owner)
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.Equal(t, title ,s.Title)

	//	 test delete session
	err = dataStore.DeleteSession(mockSession.ID)
	assert.NoError(t, err)

	ss, err := dataStore.GetSession(mockSession.ID, mockSession.Owner)
	assert.Nil(t, ss)
	assert.Error(t, err)
}
