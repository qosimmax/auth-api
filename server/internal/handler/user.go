package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/route-kz/auth-api/user"
)

// Identity is a handler that gets the user's identity (user id) from
// a token.
//
//	GET /api/v1/tokens/
//	Responds: 200, 500
//	Query Parameters:
//		token: The token to exchange to get the user id
func Identity(
	db user.IDFetcher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.URL.Query().Get("token")

		// Get the user ID
		userID, err := db.GetUserID(ctx, token)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error getting user id in identity handler: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}

		// Marshal data and respond
		response, err := json.Marshal(struct {
			UserID string `json:"user_id"`
		}{
			UserID: userID,
		})
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error marshalling user id in identity handler: %w", err),
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

func PersonalData(
	db user.PersonalDataFetcher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := r.URL.Query().Get("userId")

		personalData, err := db.FetchPersonalData(ctx, userID)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error getting personal data: %w", err),
				http.StatusInternalServerError,
				true,
			)

			return
		}

		response, _ := json.Marshal(personalData)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)

	}
}
