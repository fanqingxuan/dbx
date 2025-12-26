package gen

import (
	"context"

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

func (d *PasswordResetTokensGen) Update(ctx context.Context, m *model.PasswordResetTokens) error {
	query := "UPDATE password_reset_tokens SET token=:token,created_at=:created_at WHERE email=:email"
	_, err := d.db.NamedExecContext(ctx, query, m)
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
