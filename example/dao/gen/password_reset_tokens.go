package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type PasswordResetTokensGen struct {
	db *dbx.DB
}

func NewPasswordResetTokensGen(db *dbx.DB) *PasswordResetTokensGen {
	return &PasswordResetTokensGen{db: db}
}

func (d *PasswordResetTokensGen) Insert(ctx context.Context, m do.PasswordResetTokens) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO password_reset_tokens (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
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

func (d *PasswordResetTokensGen) FindByIds(ctx context.Context, ids []string) ([]*model.PasswordResetTokens, error) {
	var list []*model.PasswordResetTokens
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

func (d *PasswordResetTokensGen) UpdateById(ctx context.Context, m do.PasswordResetTokens, email string) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, email)
	query := "UPDATE password_reset_tokens SET " + strings.Join(sets, ",") + " WHERE email=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *PasswordResetTokensGen) UpdateByIds(ctx context.Context, m do.PasswordResetTokens, ids []string) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
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

func (d *PasswordResetTokensGen) nonNilFields(m do.PasswordResetTokens) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Email != nil { cols = append(cols, "email"); vals = append(vals, "?"); args = append(args, m.Email) }
	if m.Token != nil { cols = append(cols, "token"); vals = append(vals, "?"); args = append(args, m.Token) }
	if m.CreatedAt != nil { cols = append(cols, "created_at"); vals = append(vals, "?"); args = append(args, m.CreatedAt) }
	return cols, vals, args
}
