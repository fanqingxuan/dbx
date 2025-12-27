package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type SessionsGen struct {
	db *dbx.DB
}

func NewSessionsGen(db *dbx.DB) *SessionsGen {
	return &SessionsGen{db: db}
}

func (d *SessionsGen) Insert(ctx context.Context, m do.Sessions) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO sessions (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *SessionsGen) Delete(ctx context.Context, id string) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM sessions WHERE id=?", id)
	return err
}

func (d *SessionsGen) FindByID(ctx context.Context, id string) (*model.Sessions, error) {
	var m model.Sessions
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM sessions WHERE id=?", id)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *SessionsGen) FindByIds(ctx context.Context, ids []string) ([]*model.Sessions, error) {
	var list []*model.Sessions
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM sessions WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *SessionsGen) DeleteByIds(ctx context.Context, ids []string) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM sessions WHERE id IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *SessionsGen) UpdateById(ctx context.Context, m do.Sessions, id string) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, id)
	query := "UPDATE sessions SET " + strings.Join(sets, ",") + " WHERE id=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *SessionsGen) UpdateByIds(ctx context.Context, m do.Sessions, ids []string) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	query := "UPDATE sessions SET " + strings.Join(sets, ",") + " WHERE id IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
}

func (d *SessionsGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *SessionsGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *SessionsGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *SessionsGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *SessionsGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}

func (d *SessionsGen) nonNilFields(m do.Sessions) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Id != nil { cols = append(cols, "id"); vals = append(vals, "?"); args = append(args, m.Id) }
	if m.UserId != nil { cols = append(cols, "user_id"); vals = append(vals, "?"); args = append(args, m.UserId) }
	if m.IpAddress != nil { cols = append(cols, "ip_address"); vals = append(vals, "?"); args = append(args, m.IpAddress) }
	if m.UserAgent != nil { cols = append(cols, "user_agent"); vals = append(vals, "?"); args = append(args, m.UserAgent) }
	if m.Payload != nil { cols = append(cols, "payload"); vals = append(vals, "?"); args = append(args, m.Payload) }
	if m.LastActivity != nil { cols = append(cols, "last_activity"); vals = append(vals, "?"); args = append(args, m.LastActivity) }
	return cols, vals, args
}
