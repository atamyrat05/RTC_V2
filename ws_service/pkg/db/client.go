package package_psql

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Close()
	Pool() *pgxpool.Pool
	StdDB() *sql.DB // optional, to access raw pool for trmpgx
}

type clientImpl struct {
	pool  *pgxpool.Pool
	stdDB *sql.DB
}

func (c *clientImpl) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, args...)
}

func (c *clientImpl) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

func (c *clientImpl) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

func (c *clientImpl) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

func (c *clientImpl) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.pool.SendBatch(ctx, b)
}

func (c *clientImpl) Close() {
	c.pool.Close()
	if c.stdDB != nil {
		_ = c.stdDB.Close()
	}
}

func (c *clientImpl) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *clientImpl) StdDB() *sql.DB {
	return c.stdDB
}