package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/victor-nach/time-tracker/graph/generated"
	types "github.com/victor-nach/time-tracker/graph/model"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"go.uber.org/zap"
)

func (r *queryResolver) Me(ctx context.Context) (*types.User, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		err = rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	user, err := r.store.GetUser(claims.UserId)
	if err != nil {
		err := rerrors.Format(rerrors.CustomerNotFoundErr, err)
		r.logger.Error("sign up", zap.Error(err))
		return nil, err
	}

	resp := mapUser(user)

	return resp, nil
}

func (r *queryResolver) Session(ctx context.Context, id string) (*types.Session, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		err = rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	session, err := r.store.GetSession(id, claims.UserId)
	if err != nil {
		err = rerrors.Format(rerrors.SessionNotFoundErr, err)
		r.logger.Error("delete session", zap.Error(err))
		return nil, err
	}

	resp := mapSession(session)

	return resp, nil
}

func (r *queryResolver) Sessions(ctx context.Context, filter *types.FilterType) ([]*types.Session, error) {
	fmt.Println("sessions query ...")
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	fil := ""
	if filter != nil {
		fil = filter.String()
	}
	sessions, err := r.store.GetSessions(claims.UserId, fil)
	if err != nil {
		err = rerrors.Format(rerrors.DatabaseErr, err)
		r.logger.Error("delete session", zap.Error(err))
		return nil, err
	}

	sessionsResp := make([]*types.Session, len(sessions))
	for i, s := range sessions {
		sessionsResp[i] = mapSession(s)
	}

	return sessionsResp, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
