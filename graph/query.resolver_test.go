package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/victor-nach/time-tracker/graph/generated"
	types "github.com/victor-nach/time-tracker/graph/model"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"github.com/victor-nach/time-tracker/mocks"
	"github.com/victor-nach/time-tracker/models"
	"github.com/victor-nach/time-tracker/server/middlewares"
	"go.uber.org/zap/zaptest"
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

func TestQueryResolver_Me(t *testing.T) {
	const (
		success = iota
		invalidAuthError
		customerNotFoundError
	)

	var tests = []struct {
		name     string
		testType int
	}{
		{
			name:     "Successfully get user details",
			testType: success,
		},
		{
			name:     "Test invalid auth error",
			testType: invalidAuthError,
		},
		{
			name:     "Test customer not found error",
			testType: customerNotFoundError,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			storeMock, tokenHandlerMock := new(mocks.Datastore), new(mocks.TokenHandler)
			resolvers := NewResolver(storeMock, nil, zaptest.NewLogger(t))

			switch testCase.testType {
			case success:
				mockToken := "token"
				tokenHandlerMock.On("ValidateToken", mockToken).
					Return(&tokenhandler.Claims{UserId: "userId"}, nil)
				storeMock.On("GetUser", "userId").Return(&mockData.User, nil)

				srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))
				authMw := middlewares.NewAuthMiddleware(tokenHandlerMock, zaptest.NewLogger(t))
				gqlClient := client.New(authMw.HandleAuth(srv))

				query := `query {me { id name email Ts } }`
				var resp struct {
					Me types.User
				}
				gqlClient.MustPost(query, &resp, func(bd *client.Request) {
					bd.HTTP.Header.Add("Authorization", fmt.Sprintf("Bearer %v", mockToken))
				})
				assert.Equal(t, resp.Me.ID, mockData.User.ID)
				assert.Equal(t, *resp.Me.Name, mockData.User.Name)
				assert.Equal(t, resp.Me.Email, mockData.User.Email)

			case invalidAuthError:
				me, err := resolvers.Query().Me(context.Background())
				assert.IsType(t, &rerrors.Err{}, err)
				assert.Equal(t, rerrors.InvalidAuthErr, err.(*rerrors.Err).Code)
				fmt.Println(me, err)

			case customerNotFoundError:
				storeMock.On("GetUser", "userId").Return(nil, errors.New(""))

				ctx := context.WithValue(context.Background(), middlewares.AuthContextKey,
					tokenhandler.Claims{UserId: "userId"})
				me, err := resolvers.Query().Me(ctx)
				assert.IsType(t, &rerrors.Err{}, err)
				assert.Equal(t, rerrors.CustomerNotFoundErr, err.(*rerrors.Err).Code)
				fmt.Println(me, err)
			}
		})
	}
}

func TestQueryResolver_Session(t *testing.T) {
	const (
		success = iota
		invalidAuthError
		SessionNotFoundErr
	)

	var tests = []struct {
		name     string
		testType int
	}{
		{
			name:     "Successfully get user details",
			testType: success,
		},
		{
			name:     "Test invalid auth error",
			testType: invalidAuthError,
		},
		{
			name:     "Test customer not found error",
			testType: SessionNotFoundErr,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			storeMock, tokenHandlerMock := new(mocks.Datastore), new(mocks.TokenHandler)
			resolvers := NewResolver(storeMock, nil, zaptest.NewLogger(t))

			switch testCase.testType {
			case success:
				mockToken := "token"
				tokenHandlerMock.On("ValidateToken", mockToken).
					Return(&tokenhandler.Claims{UserId: "userId"}, nil)
				storeMock.On("GetSession", "id", "userId").Return(&mockData.Session, nil)

				srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))
				authMw := middlewares.NewAuthMiddleware(tokenHandlerMock, zaptest.NewLogger(t))
				gqlClient := client.New(authMw.HandleAuth(srv))

				query := `query { session { id owner title description start end duration Ts} }`
				var resp struct {
					Session types.Session
				}
				gqlClient.MustPost(query, &resp, func(bd *client.Request) {
					bd.HTTP.Header.Add("Authorization", fmt.Sprintf("Bearer %v", mockToken))
				})
				assert.Equal(t, resp.Session.ID, mockData.Session.ID)
				assert.Equal(t, *resp.Session.Title, mockData.Session.Title)
				assert.Equal(t, resp.Session.Duration, mockData.Session.Duration)
				assert.Equal(t, resp.Session.End, mockData.Session.End)

			case invalidAuthError:
				me, err := resolvers.Query().Me(context.Background())
				assert.IsType(t, &rerrors.Err{}, err)
				assert.Equal(t, rerrors.InvalidAuthErr, err.(*rerrors.Err).Code)
				fmt.Println(me, err)

			case SessionNotFoundErr:
				storeMock.On("GetSession", "id", "userId").Return(nil, errors.New(""))

				ctx := context.WithValue(context.Background(), middlewares.AuthContextKey,
					tokenhandler.Claims{UserId: "userId"})
				me, err := resolvers.Query().Session(ctx, "id")
				assert.IsType(t, &rerrors.Err{}, err)
				assert.Equal(t, rerrors.SessionNotFoundErr, err.(*rerrors.Err).Code)
				fmt.Println(me, err)
			}
		})
	}
}
