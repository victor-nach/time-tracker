package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/victor-nach/time-tracker/graph/generated"
	types "github.com/victor-nach/time-tracker/graph/model"
)

func (r *queryResolver) Me(ctx context.Context) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Session(ctx context.Context, id string) (*types.Session, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sessions(ctx context.Context, filter *types.FilterType) ([]*types.Session, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
