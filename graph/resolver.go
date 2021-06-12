package graph

import (
	"context"
	"github.com/victor-nach/time-tracker/db"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"github.com/victor-nach/time-tracker/lib/ulid"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver defines all the dependencies required by the resolver handlers
type Resolver struct {
	store        db.Datastore
	idGen        ulid.Idgenerator
	tokenHandler tokenhandler.TokenHandler
	logger       *zap.Logger
}

// NewResolver returns a new resolver
func NewResolver(store db.Datastore, tokenHandler tokenhandler.TokenHandler, logger *zap.Logger) *Resolver {
	return &Resolver{
		store: store,
		idGen: ulid.New(),
		tokenHandler: tokenHandler,
		logger: logger,
	}
}

func (R *Resolver) getClaimsFromCtx(ctx context.Context) (*tokenhandler.Claims, error) {
	return nil, nil
}
