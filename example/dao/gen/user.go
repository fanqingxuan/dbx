package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type UserGen struct {
	db *dbx.DB
}

func NewUserGen(db *dbx.DB) *UserGen {
	return &UserGen{db: db}
}

func (d *UserGen) Insert(ctx context.Context, m do.User) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO user (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
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

func (d *UserGen) FindByIds(ctx context.Context, ids []int64) ([]*model.User, error) {
	var list []*model.User
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

func (d *UserGen) UpdateById(ctx context.Context, m do.User, id int64) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, id)
	query := "UPDATE user SET " + strings.Join(sets, ",") + " WHERE id=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserGen) UpdateByIds(ctx context.Context, m do.User, ids []int64) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
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

func (d *UserGen) nonNilFields(m do.User) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Id != nil { cols = append(cols, "id"); vals = append(vals, "?"); args = append(args, m.Id) }
	if m.Name != nil { cols = append(cols, "name"); vals = append(vals, "?"); args = append(args, m.Name) }
	if m.Status != nil { cols = append(cols, "status"); vals = append(vals, "?"); args = append(args, m.Status) }
	if m.CreatedAt != nil { cols = append(cols, "created_at"); vals = append(vals, "?"); args = append(args, m.CreatedAt) }
	if m.UpdatedAt != nil { cols = append(cols, "updated_at"); vals = append(vals, "?"); args = append(args, m.UpdatedAt) }
	if m.DeletedAt != nil { cols = append(cols, "deleted_at"); vals = append(vals, "?"); args = append(args, m.DeletedAt) }
	if m.Age != nil { cols = append(cols, "age"); vals = append(vals, "?"); args = append(args, m.Age) }
	return cols, vals, args
}
