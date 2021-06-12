package graph

import (
	"context"
	"github.com/victor-nach/time-tracker/db"
	types "github.com/victor-nach/time-tracker/graph/model"
	"github.com/victor-nach/time-tracker/lib/encryptor"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"github.com/victor-nach/time-tracker/lib/ulid"
	"github.com/victor-nach/time-tracker/models"
	"go.uber.org/zap"
	"time"
)

// Resolver defines all the dependencies required by the resolver handlers
type Resolver struct {
	store        db.Datastore
	idGen        ulid.Idgenerator
	encryptor    encryptor.Encryptor
	tokenHandler tokenhandler.TokenHandler
	logger       *zap.Logger
}

// NewResolver returns a new resolver
func NewResolver(store db.Datastore, tokenHandler tokenhandler.TokenHandler, logger *zap.Logger) *Resolver {
	return &Resolver{
		store:        store,
		idGen:        ulid.New(),
		encryptor:    encryptor.NewEncryptor(),
		tokenHandler: tokenHandler,
		logger:       logger,
	}
}

func (r *Resolver) getClaimsFromCtx(ctx context.Context) (*tokenhandler.Claims, error) {
	return nil, nil
}

func (r *mutationResolver) genAuthTokens(userId string) (authToken string, refreshToken string, err error) {
	tokenExpiry := time.Now().Add(tokenhandler.AuthTokenDuration)
	authToken, err = r.tokenHandler.NewToken(userId, tokenExpiry)
	if err != nil {
		err := rerrors.Format(rerrors.InternalErr, nil)
		r.logger.Error("generate token", zap.Error(err))
		return "", "", err
	}
	refreshExpiry := time.Now().Add(tokenhandler.RefreshTokenDuration)
	refreshToken, err = r.tokenHandler.NewToken(userId, refreshExpiry)
	if err != nil {
		err := rerrors.Format(rerrors.InternalErr, nil)
		r.logger.Error("generate token", zap.Error(err))
		return "", "", err
	}
	return authToken, refreshToken, nil
}

// mapSession converts models.Session the corresponding graphql type
func mapSession(data *models.Session) *types.Session {
	return &types.Session{
		ID:          data.ID,
		Owner:       data.Owner,
		Title:       &data.Title,
		Description: &data.Description,
		Start:       int(data.Start),
		End:         int(data.End),
		Ts:          int(data.Ts),
	}
}

// mapUser converts models.Session the corresponding graphql type
func mapUser(data *models.User) *types.User {
	return &types.User{
		ID:    data.ID,
		Name:  &data.Name,
		Email: data.Email,
		Ts:    int(data.Ts),
	}
}
