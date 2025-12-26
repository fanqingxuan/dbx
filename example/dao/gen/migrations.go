package gen

import (
	"context"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type MigrationsGen struct {
	db *dbx.DB
}

func NewMigrationsGen(db *dbx.DB) *MigrationsGen {
	return &MigrationsGen{db: db}
}

func (d *MigrationsGen) Insert(ctx context.Context, m *model.Migrations) error {
	query := "INSERT INTO migrations (id,migration,batch,created_at) VALUES (:id,:migration,:batch,:created_at)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *MigrationsGen) Update(ctx context.Context, m *model.Migrations) error {
	query := "UPDATE migrations SET migration=:migration,batch=:batch,created_at=:created_at WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *MigrationsGen) Delete(ctx context.Context, id int64) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM migrations WHERE id=?", id)
	return err
}

func (d *MigrationsGen) FindByID(ctx context.Context, id int64) (*model.Migrations, error) {
	var m model.Migrations
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM migrations WHERE id=?", id)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *MigrationsGen) FindByIds(ctx context.Context, ids []int64) ([]model.Migrations, error) {
	var list []model.Migrations
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM migrations WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *MigrationsGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *MigrationsGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *MigrationsGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *MigrationsGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *MigrationsGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}
