package gen

import (
	"context"
	"reflect"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type UserGen struct {
	db *dbx.DB
}

func NewUserGen(db *dbx.DB) *UserGen {
	return &UserGen{db: db}
}

func (d *UserGen) Insert(ctx context.Context, m *model.User) error {
	query := "INSERT INTO user (id,name,status,created_at,updated_at,deleted_at,age) VALUES (:id,:name,:status,:created_at,:updated_at,:deleted_at,:age)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserGen) InsertSelective(ctx context.Context, m *model.User) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO user (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *UserGen) Update(ctx context.Context, m *model.User) error {
	query := "UPDATE user SET name=:name,status=:status,created_at=:created_at,updated_at=:updated_at,deleted_at=:deleted_at,age=:age WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserGen) UpdateSelective(ctx context.Context, m *model.User) error {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, m.Id)
	query := "UPDATE user SET " + strings.Join(sets, ",") + " WHERE id=?"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *UserGen) Delete(ctx context.Context, id int64) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM user WHERE id=?", id)
	return err
}

func (d *UserGen) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var m model.User
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM user WHERE id=?", id)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *UserGen) FindByIds(ctx context.Context, ids []int64) ([]model.User, error) {
	var list []model.User
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM user WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *UserGen) DeleteByIds(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM user WHERE id IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserGen) UpdateByIds(ctx context.Context, ids []int64, fields map[string]any) (int64, error) {
	if len(ids) == 0 || len(fields) == 0 { return 0, nil }
	var sets []string
	var args []any
	for k, v := range fields { sets = append(sets, k+"=?"); args = append(args, v) }
	query := "UPDATE user SET " + strings.Join(sets, ",") + " WHERE id IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *UserGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *UserGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *UserGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *UserGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}

func (d *UserGen) nonNilFields(m *model.User) ([]string, []string, []any) {
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
