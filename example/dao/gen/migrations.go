package gen

import (
	"context"
	"reflect"
	"strings"

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

func (d *MigrationsGen) InsertSelective(ctx context.Context, m *model.Migrations) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO migrations (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *MigrationsGen) Update(ctx context.Context, m *model.Migrations) error {
	query := "UPDATE migrations SET migration=:migration,batch=:batch,created_at=:created_at WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *MigrationsGen) UpdateSelective(ctx context.Context, m *model.Migrations) error {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, m.Id)
	query := "UPDATE migrations SET " + strings.Join(sets, ",") + " WHERE id=?"
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

func (d *MigrationsGen) FindByIds(ctx context.Context, ids []int64) ([]model.Migrations, error) {
	var list []model.Migrations
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

func (d *MigrationsGen) UpdateByIds(ctx context.Context, ids []int64, fields map[string]any) (int64, error) {
	if len(ids) == 0 || len(fields) == 0 { return 0, nil }
	var sets []string
	var args []any
	for k, v := range fields { sets = append(sets, k+"=?"); args = append(args, v) }
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

func (d *MigrationsGen) nonNilFields(m *model.Migrations) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	v := reflect.ValueOf(m).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Ptr && f.IsNil() { continue }
		tag := t.Field(i).Tag.Get("db")
		if tag == "" { continue }
		cols = append(cols, tag)
		vals = append(vals, "?")
		if f.Kind() == reflect.Ptr { args = append(args, f.Elem().Interface()) } else { args = append(args, f.Interface()) }
	}
	return cols, vals, args
}
