package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type UserDetailGen struct {
	db *dbx.DB
}

func NewUserDetailGen(db *dbx.DB) *UserDetailGen {
	return &UserDetailGen{db: db}
}

func (d *UserDetailGen) Insert(ctx context.Context, m do.UserDetail) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO user_detail (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
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

func (d *UserDetailGen) FindByIds(ctx context.Context, ids []int32) ([]*model.UserDetail, error) {
	var list []*model.UserDetail
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

func (d *UserDetailGen) UpdateById(ctx context.Context, m do.UserDetail, uid int32) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, uid)
	query := "UPDATE user_detail SET " + strings.Join(sets, ",") + " WHERE uid=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserDetailGen) UpdateByIds(ctx context.Context, m do.UserDetail, ids []int32) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
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

func (d *UserDetailGen) nonNilFields(m do.UserDetail) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Uid != nil { cols = append(cols, "uid"); vals = append(vals, "?"); args = append(args, m.Uid) }
	if m.Address != nil { cols = append(cols, "address"); vals = append(vals, "?"); args = append(args, m.Address) }
	return cols, vals, args
}
