package gen

import (
	"context"

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

func (d *UserGen) Update(ctx context.Context, m *model.User) error {
	query := "UPDATE user SET name=:name,status=:status,created_at=:created_at,updated_at=:updated_at,deleted_at=:deleted_at,age=:age WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
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
