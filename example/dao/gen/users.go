package gen

import (
	"context"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type UsersGen struct {
	db *dbx.DB
}

func NewUsersGen(db *dbx.DB) *UsersGen {
	return &UsersGen{db: db}
}

func (d *UsersGen) Insert(ctx context.Context, m *model.Users) error {
	query := "INSERT INTO users (id,name,email,email_verified_at,password,remember_token,created_at,updated_at) VALUES (:id,:name,:email,:email_verified_at,:password,:remember_token,:created_at,:updated_at)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UsersGen) Update(ctx context.Context, m *model.Users) error {
	query := "UPDATE users SET name=:name,email=:email,email_verified_at=:email_verified_at,password=:password,remember_token=:remember_token,created_at=:created_at,updated_at=:updated_at WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
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

func (d *UsersGen) FindByIds(ctx context.Context, ids []int64) ([]model.Users, error) {
	var list []model.Users
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM users WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
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
