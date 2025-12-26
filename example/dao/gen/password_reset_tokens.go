package gen

import (
	"context"
	"reflect"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type PasswordResetTokensGen struct {
	db *dbx.DB
}

func NewPasswordResetTokensGen(db *dbx.DB) *PasswordResetTokensGen {
	return &PasswordResetTokensGen{db: db}
}

func (d *PasswordResetTokensGen) Insert(ctx context.Context, m *model.PasswordResetTokens) error {
	query := "INSERT INTO password_reset_tokens (email,token,created_at) VALUES (:email,:token,:created_at)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *PasswordResetTokensGen) InsertSelective(ctx context.Context, m *model.PasswordResetTokens) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO password_reset_tokens (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *PasswordResetTokensGen) Update(ctx context.Context, m *model.PasswordResetTokens) error {
	query := "UPDATE password_reset_tokens SET token=:token,created_at=:created_at WHERE email=:email"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *PasswordResetTokensGen) UpdateSelective(ctx context.Context, m *model.PasswordResetTokens) error {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, m.Email)
	query := "UPDATE password_reset_tokens SET " + strings.Join(sets, ",") + " WHERE email=?"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *PasswordResetTokensGen) Delete(ctx context.Context, email string) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM password_reset_tokens WHERE email=?", email)
	return err
}

func (d *PasswordResetTokensGen) FindByID(ctx context.Context, email string) (*model.PasswordResetTokens, error) {
	var m model.PasswordResetTokens
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM password_reset_tokens WHERE email=?", email)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *PasswordResetTokensGen) FindByIds(ctx context.Context, ids []string) ([]model.PasswordResetTokens, error) {
	var list []model.PasswordResetTokens
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM password_reset_tokens WHERE email IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *PasswordResetTokensGen) DeleteByIds(ctx context.Context, ids []string) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM password_reset_tokens WHERE email IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *PasswordResetTokensGen) UpdateByIds(ctx context.Context, ids []string, fields map[string]any) (int64, error) {
	if len(ids) == 0 || len(fields) == 0 { return 0, nil }
	var sets []string
	var args []any
	for k, v := range fields { sets = append(sets, k+"=?"); args = append(args, v) }
	query := "UPDATE password_reset_tokens SET " + strings.Join(sets, ",") + " WHERE email IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
}

func (d *PasswordResetTokensGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *PasswordResetTokensGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *PasswordResetTokensGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *PasswordResetTokensGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *PasswordResetTokensGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}

func (d *PasswordResetTokensGen) nonNilFields(m *model.PasswordResetTokens) ([]string, []string, []any) {
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
