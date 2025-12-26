package dbx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
	logger Logger
}

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

type nopLogger struct{}

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}

func New(db *sqlx.DB) *DB {
	return &DB{DB: db, logger: nopLogger{}}
}

func Open(driverName, dsn string) (*DB, error) {
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return New(db), nil
}

func (d *DB) SetLogger(l Logger) {
	d.logger = l
}

func (d *DB) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	d.logger.Debug("QueryRows: %s, args: %v", query, args)
	return d.SelectContext(ctx, dest, query, args...)
}

func (d *DB) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	d.logger.Debug("QueryRow: %s, args: %v", query, args)
	return d.GetContext(ctx, dest, query, args...)
}

func (d *DB) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	d.logger.Debug("QueryValue: %s, args: %v", query, args)
	return d.GetContext(ctx, dest, query, args...)
}

func (d *DB) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	d.logger.Debug("QueryColumn: %s, args: %v", query, args)
	return d.SelectContext(ctx, dest, query, args...)
}

func (d *DB) ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error) {
	d.logger.Debug("Exec: %s, args: %v", query, args)
	return d.ExecContext(ctx, query, args...)
}

func (d *DB) In(query string, args ...any) (string, []any, error) {
	return sqlx.In(query, args...)
}
