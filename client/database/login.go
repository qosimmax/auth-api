package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/google/uuid"
	"gitlab.com/route-kz/auth-api/user"
)

func (c *Client) prepareRecordUserIDToObjectIDStmt() error {
	stmt, err := c.DB.Preparex(`
		INSERT INTO
			user_ids (user_id, login, auth_method, auth_user_type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (login, auth_user_type) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("error preparing record user id to object id statement: %w", err)
	}
	c.RecordUserIDToObjectIDStmt = stmt
	return nil
}

func (c *Client) prepareGetUserIDByObjectIDStmt() error {
	stmt, err := c.DB.Preparex(`
		SELECT
			user_id
		FROM user_ids
		WHERE login = $1 and auth_user_type=$2;
	`)
	if err != nil {
		return fmt.Errorf("error preparing get user id by object id statement: %w", err)
	}
	c.GetUserIDByObjectIDStmt = stmt
	return nil
}

func (c *Client) prepareRecordTokenToUserIDStmt() error {
	stmt, err := c.DB.Preparex(`
		INSERT INTO
			tokens as t (token, user_id, created_at)
		VALUES ($1, $2, now())
		ON CONFLICT (token) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("error preparing record token to user id statement: %w", err)
	}
	c.RecordTokenToUserIDStmt = stmt
	return nil
}

func (c *Client) GetOrCreateUserID(ctx context.Context, payload user.CreateTokenPayload) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetOrCreateUserID")
	defer span.Finish()

	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	userID, err := c.getUserIDFromObjectID(ctx, payload)
	if err != nil {
		return "", fmt.Errorf(
			"error getting user id from object id: %w", err)
	}

	if userID != "" {
		return userID, nil
	}

	userID, err = generateUUID()
	if err != nil {
		return "", fmt.Errorf("error generating new user id: %w", err)
	}

	_, err = c.RecordUserIDToObjectIDStmt.ExecContext(
		cctx,
		userID,
		payload.Login,
		payload.AuthMethod,
		payload.AuthUserType,
	)
	if err != nil {
		return "", fmt.Errorf("error recording token to user id: %w", err)
	}

	return userID, nil

}

func (c *Client) getUserIDFromObjectID(ctx context.Context, payload user.CreateTokenPayload) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getUserIDFromObjectID")
	defer span.Finish()

	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	r := c.GetUserIDByObjectIDStmt.QueryRowContext(cctx, payload.Login, payload.AuthUserType)

	var userID string
	err := r.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("error scanning for user id from object id: %w", err)
	}

	return userID, nil
}

func (c *Client) CreateToken(ctx context.Context, userID string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateToken")
	defer span.Finish()

	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	token, err := generateUUID()
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	_, err = c.RecordTokenToUserIDStmt.ExecContext(
		cctx,
		token,
		userID,
	)
	if err != nil {
		return "", fmt.Errorf("error recording token to user id: %w", err)
	}

	return token, nil
}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating uuid: %w", err)
	}
	return id.String(), nil
}
