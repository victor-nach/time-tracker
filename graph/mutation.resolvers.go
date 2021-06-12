package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/victor-nach/time-tracker/graph/generated"
	types "github.com/victor-nach/time-tracker/graph/model"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"github.com/victor-nach/time-tracker/models"
	"go.uber.org/zap"
)

func (r *mutationResolver) SignUp(ctx context.Context, email string, passcode string, name string) (*types.AuthResponse, error) {
	if _, err := r.store.GetUserByEmail(email); err != nil {
		err := rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("sign up", zap.Error(err))
		return nil, err
	}

	hashPasscode, err := r.encryptor.HashPassword(passcode)
	if err != nil {
		err := rerrors.Format(rerrors.InternalErr, nil)
		r.logger.Error("sign up", zap.Error(err))
		return nil, err
	}

	userId := r.idGen.Generate()
	user := models.User{
		ID:       userId,
		Name:     name,
		Email:    email,
		Password: hashPasscode,
		Ts:       time.Now().Unix(),
	}

	if _, err := r.store.CreateUser(&user); err != nil {
		err = rerrors.Format(rerrors.DatabaseErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	authToken, refreshToken, err := r.genAuthTokens(userId)
	if err != nil {
		return nil, err
	}

	resp := &types.AuthResponse{
		Success:      true,
		Message:      "Sign up Successful",
		JwtToken:     authToken,
		RefreshToken: refreshToken,
		//User: user,
	}

	return resp, nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, passcode string) (*types.AuthResponse, error) {
	user, err := r.store.GetUserByEmail(email)
	if err != nil {
		err := rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("sign up", zap.Error(err))
		return nil, err
	}

	if ok := r.encryptor.ComparePasscode(passcode, user.Password); !ok {
		err := rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("sign up", zap.Error(err))
		return nil, err
	}

	authToken, refreshToken, err := r.genAuthTokens(user.ID)

	resp := &types.AuthResponse{
		Success:      true,
		Message:      "Sign up Successful",
		JwtToken:     authToken,
		RefreshToken: refreshToken,
		User:         mapUser(user),
	}

	return resp, nil
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
		JwtToken:     authToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func (r *mutationResolver) SaveSession(ctx context.Context, input *types.SessionInput) (*types.Response, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		err = rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	sessionId := r.idGen.Generate()
	session := models.Session{
		ID:    sessionId,
		Owner: claims.UserId,
		Start: int64(input.Start),
		End:   int64(input.End),
		Ts:    time.Now().Unix(),
	}
	if input.Title != nil {
		session.Title = *input.Title
	}
	if input.Description != nil {
		session.Title = *input.Description
	}

	if _, err := r.store.CreateSession(&session); err != nil {
		err = rerrors.Format(rerrors.DatabaseErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	resp := &types.Response{
		Success: true,
		Message: "Successfully created session!",
		Token:   &sessionId,
	}

	return resp, nil
}

func (r *mutationResolver) UpdateSessionInfo(ctx context.Context, id string, input *types.UpdateSessionInput) (*types.Response, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		err = rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	if _, err := r.store.GetSession(id, claims.UserId); err != nil {
		err = rerrors.Format(rerrors.SessionNotFoundErr, err)
		r.logger.Error("delete session", zap.Error(err))
		return nil, err
	}

	sessionInfo := models.SessionInfo{
		Title:       input.Title,
		Description: input.Description,
	}
	if err := r.store.UpdateSession(id, sessionInfo); err != nil {
		err = rerrors.Format(rerrors.DatabaseErr, err)
		r.logger.Error("delete session", zap.Error(err))
		return nil, err
	}

	return &types.Response{
		Success: true,
		Message: "Successfully updated session",
	}, nil
}

func (r *mutationResolver) DeleteSession(ctx context.Context, id string) (*types.Response, error) {
	claims, err := r.getClaimsFromCtx(ctx)
	if err != nil {
		err = rerrors.Format(rerrors.InvalidAuthErr, err)
		r.logger.Error("save session", zap.Error(err))
		return nil, err
	}

	if _, err := r.store.GetSession(id, claims.UserId); err != nil {
		err = rerrors.Format(rerrors.SessionNotFoundErr, err)
		r.logger.Error("delete session", zap.Error(err))
		return nil, err
	}

	if err := r.store.DeleteSession(id); err != nil {
		err = rerrors.Format(rerrors.DatabaseErr, err)
		r.logger.Error("delete session", zap.Error(err))
		return nil, err
	}

	return &types.Response{
		Success: true,
		Message: "Successfully deleted session",
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
