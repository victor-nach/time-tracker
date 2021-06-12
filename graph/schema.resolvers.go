package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"github.com/victor-nach/time-tracker/models"
	"time"

	"github.com/victor-nach/time-tracker/graph/generated"
	types "github.com/victor-nach/time-tracker/graph/models"
)

func (r *mutationResolver) SignUp(ctx context.Context) (*types.AuthResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context) (*types.AuthResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (*types.AuthResponse, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	tokenExpiry := time.Now().Add(tokenhandler.AuthTokenDuration)
	authToken, err := r.tokenHandler.NewToken(claims.UserId, tokenExpiry)

	refreshExpiry := time.Now().Add(tokenhandler.RefreshTokenDuration)
	refreshToken, err := r.tokenHandler.NewToken(claims.UserId, refreshExpiry)
	if err != nil {
		return nil, err
	}

	resp := &types.AuthResponse{
		JwtToken: authToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func (r *mutationResolver) SaveSession(ctx context.Context, input *types.SessionInput) (*types.Response, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	session := models.Session{
		ID: r.idGen.Generate(),
		Owner: claims.UserId,
	}
}

func (r *mutationResolver) UpdateSessionInfo(ctx context.Context, id string, input *types.UpdateSessionInput) (*types.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteSession(ctx context.Context, id string) (*types.Response, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Session(ctx context.Context, id string) (*types.Session, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sessions(ctx context.Context, filter *types.FilterType) ([]*types.Session, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
