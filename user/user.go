package user

import "context"

// IDFetcher is an interface for getting a user id using a token.
type IDFetcher interface {
	GetUserID(ctx context.Context, token string) (string, error)
}

// IDFetcherTokenRemover is an interface for getting a user id using a token.
type IDFetcherTokenRemover interface {
	GetUserIDRemoveToken(ctx context.Context, token string) (string, error)
}

type PersonalData struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

// PersonalDataFetcher is an interface for fetching personal data about a
// user from the user id.
type PersonalDataFetcher interface {
	FetchPersonalData(ctx context.Context, userID string) (*PersonalData, error)
}
