package main

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// DatabaseAccess defines data access functionality
type DatabaseAccess interface {
	Get(ctx context.Context, id uuid.UUID) error
}

// comptime check that database implements the DatabaseAccess interface
var _ DatabaseAccess = new(database)

// NewDatabaseAccess returns a new implementation of the database access interface
func NewDatabaseAccess(repo *sql.DB) DatabaseAccess { return &database{repo} }

type database struct {
	dbConn *sql.DB
}

// Get fetches a user from the database using an ID
func (d *database) Get(ctx context.Context, id uuid.UUID) error {
	query := "select * from users where user_id = $1 limit 1"

	tx, err := d.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.Query(query, id); err != nil {
		return sql.ErrNoRows
	}

	return nil
}
