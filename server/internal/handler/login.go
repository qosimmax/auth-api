package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/route-kz/auth-api/user"
)

// CreateToken is a handler that creates tokens identifying the user.
//
//	POST /api/v1/tokens
//	Responds: 200, 400, 500
//	Body:
//		type createTokenPayload struct {
//			AuthUserType    string `json:"userType"`
//			AuthCode    string `json:"authCode"`
//			Login string `json:"objectId"`
//		}
//
// The handler will get the data about the user and create a token
// that can be exchanged to get data about the user. It will also
// create a user id for the user if it does not exist.
func CreateToken(
	db user.IDFetcherTokenCreator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Decode request body
		var payload user.CreateTokenPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error decoding create token payload: %w", err),
				http.StatusBadRequest,
				true,
			)
			return
		}

		// Get (or create, if this is a new user) the  user id for the user
		userID, err := db.GetOrCreateUserID(ctx, payload)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error getting or creating a user id in create token handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}

		// Create token
		token, err := db.CreateToken(ctx, userID)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error creating token in create token handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}

		// Marshal data and respond
		response, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error marshalling token in create token handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}
}

// RefreshToken is a handler that refresh tokens
func RefreshToken(
	db user.TokenRefresher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.URL.Query().Get("token")
		// Get the user ID and remove token
		userID, err := db.GetUserIDRemoveToken(ctx, token)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error getting user id in refresh-token handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}

		if userID == "" {
			handleError(
				w,
				fmt.Errorf("invalid token"),
				http.StatusBadRequest,
				true,
			)
			return
		}

		// Create new token
		newToken, err := db.CreateToken(ctx, userID)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error creating token in create refresh-token handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}

		// Marshal data and respond
		response, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			Token: newToken,
		})
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error marshalling token in refresh-token handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}
}
