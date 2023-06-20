package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"gitlab.com/route-kz/auth-api/user"
)

func (c *Client) prepareGetUserIDByTokenStmt() error {
	stmt, err := c.DB.Preparex(`
		SELECT
			user_id
		FROM tokens
		WHERE token = $1;
	`)
	if err != nil {
		return fmt.Errorf("error preparing get user id by token statement: %w", err)
	}
	c.GetUserIDByTokenStmt = stmt
	return nil
}

func (c *Client) GetUserID(ctx context.Context, token string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetUserID")
	defer span.Finish()

	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	r := c.GetUserIDByTokenStmt.QueryRowContext(cctx, token)

	var userID string
	err := r.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("error scanning for user id from token: %w", err)
	}

	return userID, nil
}

func (c *Client) prepareGetUserIDRemoveTokenStmt() error {
	stmt, err := c.DB.Preparex(`select get_user_id_token_remover($1);`)
	if err != nil {
		return fmt.Errorf("error preparing get user id and token remover statement: %w", err)
	}

	c.GetUserIDRemoveTokenStmt = stmt
	return nil
}

func (c *Client) GetUserIDRemoveToken(ctx context.Context, token string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetUserIDRemoveToken")
	defer span.Finish()

	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var userID *string
	err := c.GetUserIDRemoveTokenStmt.GetContext(cctx, &userID, token)
	if err != nil || userID == nil {
		return "", err
	}

	return *userID, nil
}

func (c *Client) prepareFetchPersonalDataStmt() error {
	stmt, err := c.DB.Preparex(`
		select user_id, login, auth_user_type, auth_method from user_ids where user_id=$1;`)

	if err != nil {
		return fmt.Errorf("error preparing get personal data: %w", err)
	}

	c.FetchPersonalDataStmt = stmt
	return nil
}

func (c *Client) FetchPersonalData(ctx context.Context, userID string) (*user.PersonalData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FetchPersonalData")
	defer span.Finish()

	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var row personalDataRow

	err := c.FetchPersonalDataStmt.GetContext(cctx, &row, userID)
	if err != nil {
		return nil, err
	}

	personalData := &user.PersonalData{
		UserID:      userID,
		PhoneNumber: row.Login, //@TODO check for object type phoneNumber, email, etc...
		Email:       "",
	}

	return personalData, nil
}

type personalDataRow struct {
	UserID       string `db:"user_id"`
	Login        string `db:"login"`
	AuthUserType string `db:"auth_user_type"`
	AuthMethod   string `db:"auth_method"`
}
