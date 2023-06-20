package user

import "context"

type CreateTokenPayload struct {
	Login        string `json:"login"`
	AuthUserType string `json:"auth_user_type"`
	AuthCode     string `json:"auth_code"`
	AuthMethod   string `json:"auth_method"`
}

// IDFetcherCreator is an interface for getting a user id if that already
// exists, or creating it if it does not. Stores user id -> object id.
type IDFetcherCreator interface {
	GetOrCreateUserID(ctx context.Context, payload CreateTokenPayload) (string, error)
}

// TokenCreator is an interface for creating a token given a user id.
type TokenCreator interface {
	CreateToken(ctx context.Context, userID string) (string, error)
}

// IDFetcherTokenCreator is an interface for getting or creating a user id
// and creating a token for the user.
type IDFetcherTokenCreator interface {
	IDFetcherCreator
	TokenCreator
}

// TokenRefresher is an interface for refresh token
type TokenRefresher interface {
	IDFetcherTokenRemover
	TokenCreator
}
