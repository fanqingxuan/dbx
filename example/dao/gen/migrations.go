package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type MigrationsGen struct {
	db *dbx.DB
}

func NewMigrationsGen(db *dbx.DB) *MigrationsGen {
	return &MigrationsGen{db: db}
}

func (d *MigrationsGen) Insert(ctx context.Context, m do.Migrations) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO migrations (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
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

func (d *MigrationsGen) FindByIds(ctx context.Context, ids []int64) ([]*model.Migrations, error) {
	var list []*model.Migrations
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM migrations WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *MigrationsGen) DeleteByIds(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM migrations WHERE id IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *MigrationsGen) UpdateById(ctx context.Context, m do.Migrations, id int64) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, id)
	query := "UPDATE migrations SET " + strings.Join(sets, ",") + " WHERE id=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *MigrationsGen) UpdateByIds(ctx context.Context, m do.Migrations, ids []int64) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	query := "UPDATE migrations SET " + strings.Join(sets, ",") + " WHERE id IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
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

func (d *MigrationsGen) nonNilFields(m do.Migrations) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Id != nil { cols = append(cols, "id"); vals = append(vals, "?"); args = append(args, m.Id) }
	if m.Migration != nil { cols = append(cols, "migration"); vals = append(vals, "?"); args = append(args, m.Migration) }
	if m.Batch != nil { cols = append(cols, "batch"); vals = append(vals, "?"); args = append(args, m.Batch) }
	if m.CreatedAt != nil { cols = append(cols, "created_at"); vals = append(vals, "?"); args = append(args, m.CreatedAt) }
	return cols, vals, args
}
