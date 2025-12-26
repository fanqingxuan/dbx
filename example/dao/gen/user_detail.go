package gen

import (
	"context"
	"reflect"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type UserDetailGen struct {
	db *dbx.DB
}

func NewUserDetailGen(db *dbx.DB) *UserDetailGen {
	return &UserDetailGen{db: db}
}

func (d *UserDetailGen) Insert(ctx context.Context, m *model.UserDetail) error {
	query := "INSERT INTO user_detail (uid,address) VALUES (:uid,:address)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserDetailGen) InsertSelective(ctx context.Context, m *model.UserDetail) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO user_detail (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *UserDetailGen) Update(ctx context.Context, m *model.UserDetail) error {
	query := "UPDATE user_detail SET address=:address WHERE uid=:uid"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserDetailGen) UpdateSelective(ctx context.Context, m *model.UserDetail) error {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, m.Uid)
	query := "UPDATE user_detail SET " + strings.Join(sets, ",") + " WHERE uid=?"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *UserDetailGen) Delete(ctx context.Context, uid int32) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM user_detail WHERE uid=?", uid)
	return err
}

func (d *UserDetailGen) FindByID(ctx context.Context, uid int32) (*model.UserDetail, error) {
	var m model.UserDetail
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM user_detail WHERE uid=?", uid)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *UserDetailGen) FindByIds(ctx context.Context, ids []int32) ([]model.UserDetail, error) {
	var list []model.UserDetail
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM user_detail WHERE uid IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *UserDetailGen) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM user_detail WHERE uid IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserDetailGen) UpdateByIds(ctx context.Context, ids []int32, fields map[string]any) (int64, error) {
	if len(ids) == 0 || len(fields) == 0 { return 0, nil }
	var sets []string
	var args []any
	for k, v := range fields { sets = append(sets, k+"=?"); args = append(args, v) }
	query := "UPDATE user_detail SET " + strings.Join(sets, ",") + " WHERE uid IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserDetailGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *UserDetailGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *UserDetailGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *UserDetailGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *UserDetailGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}

func (d *UserDetailGen) nonNilFields(m *model.UserDetail) ([]string, []string, []any) {
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
