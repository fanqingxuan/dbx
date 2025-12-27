package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type UsersGen struct {
	db *dbx.DB
}

func NewUsersGen(db *dbx.DB) *UsersGen {
	return &UsersGen{db: db}
}

func (d *UsersGen) Insert(ctx context.Context, m do.Users) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO users (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *UsersGen) Delete(ctx context.Context, id int64) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM users WHERE id=?", id)
	return err
}

func (d *UsersGen) FindByID(ctx context.Context, id int64) (*model.Users, error) {
	var m model.Users
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM users WHERE id=?", id)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *UsersGen) FindByIds(ctx context.Context, ids []int64) ([]*model.Users, error) {
	var list []*model.Users
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM users WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *UsersGen) DeleteByIds(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM users WHERE id IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UsersGen) UpdateById(ctx context.Context, m do.Users, id int64) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, id)
	query := "UPDATE users SET " + strings.Join(sets, ",") + " WHERE id=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *UsersGen) UpdateByIds(ctx context.Context, m do.Users, ids []int64) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	query := "UPDATE users SET " + strings.Join(sets, ",") + " WHERE id IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UsersGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *UsersGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *UsersGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *UsersGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *UsersGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}

func (d *UsersGen) nonNilFields(m do.Users) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Id != nil { cols = append(cols, "id"); vals = append(vals, "?"); args = append(args, m.Id) }
	if m.Name != nil { cols = append(cols, "name"); vals = append(vals, "?"); args = append(args, m.Name) }
	if m.Email != nil { cols = append(cols, "email"); vals = append(vals, "?"); args = append(args, m.Email) }
	if m.EmailVerifiedAt != nil { cols = append(cols, "email_verified_at"); vals = append(vals, "?"); args = append(args, m.EmailVerifiedAt) }
	if m.Password != nil { cols = append(cols, "password"); vals = append(vals, "?"); args = append(args, m.Password) }
	if m.RememberToken != nil { cols = append(cols, "remember_token"); vals = append(vals, "?"); args = append(args, m.RememberToken) }
	if m.CreatedAt != nil { cols = append(cols, "created_at"); vals = append(vals, "?"); args = append(args, m.CreatedAt) }
	if m.UpdatedAt != nil { cols = append(cols, "updated_at"); vals = append(vals, "?"); args = append(args, m.UpdatedAt) }
	return cols, vals, args
}
