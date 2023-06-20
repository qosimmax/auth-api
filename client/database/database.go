package database

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/route-kz/auth-api/config"
)

// Client holds the database client and prepared statements.
type Client struct {
	DB                         *sqlx.DB
	RecordUserIDToObjectIDStmt *sqlx.Stmt
	GetUserIDByObjectIDStmt    *sqlx.Stmt
	RecordTokenToUserIDStmt    *sqlx.Stmt
	GetUserIDByTokenStmt       *sqlx.Stmt
	FetchPersonalDataStmt      *sqlx.Stmt
	GetUserIDRemoveTokenStmt   *sqlx.Stmt
}

// Init sets up a new database client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseURL,
		config.DatabasePort,
		config.DatabaseDB,
		config.DatabaseOptions,
	)

	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(config.DatabaseMaxConnections)
	db.SetMaxIdleConns(config.DatabaseMaxIdleConnections)

	c.DB = db

	if err := c.prepareRecordUserIDToObjectIDStmt(); err != nil {
		return err
	}

	if err := c.prepareRecordTokenToUserIDStmt(); err != nil {
		return err
	}

	if err := c.prepareGetUserIDByObjectIDStmt(); err != nil {
		return err
	}

	if err := c.prepareGetUserIDByTokenStmt(); err != nil {
		return err
	}

	if err := c.prepareFetchPersonalDataStmt(); err != nil {
		return err
	}

	if err := c.prepareGetUserIDRemoveTokenStmt(); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection and statements.
func (c *Client) Close() error {

	if err := c.RecordUserIDToObjectIDStmt.Close(); err != nil {
		return fmt.Errorf("error on closing record user id statement: %w", err)
	}

	if err := c.RecordTokenToUserIDStmt.Close(); err != nil {
		return fmt.Errorf("error on closing record token statement: %w", err)
	}

	if err := c.GetUserIDByTokenStmt.Close(); err != nil {
		return fmt.Errorf("error on closing get user id by token statement: %w", err)
	}

	if err := c.GetUserIDByObjectIDStmt.Close(); err != nil {
		return fmt.Errorf("error on closing get user id by object id statement: %w", err)
	}

	if err := c.FetchPersonalDataStmt.Close(); err != nil {
		return fmt.Errorf("error on closing get personal data statement: %w", err)
	}

	if err := c.GetUserIDRemoveTokenStmt.Close(); err != nil {
		return fmt.Errorf("error on closing get user id and remove token statement: %w", err)
	}

	err := c.DB.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	return nil
}
