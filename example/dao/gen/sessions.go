package gen

import (
	"context"
	"reflect"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type SessionsGen struct {
	db *dbx.DB
}

func NewSessionsGen(db *dbx.DB) *SessionsGen {
	return &SessionsGen{db: db}
}

func (d *SessionsGen) Insert(ctx context.Context, m *model.Sessions) error {
	query := "INSERT INTO sessions (id,user_id,ip_address,user_agent,payload,last_activity) VALUES (:id,:user_id,:ip_address,:user_agent,:payload,:last_activity)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *SessionsGen) InsertSelective(ctx context.Context, m *model.Sessions) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO sessions (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *SessionsGen) Update(ctx context.Context, m *model.Sessions) error {
	query := "UPDATE sessions SET user_id=:user_id,ip_address=:ip_address,user_agent=:user_agent,payload=:payload,last_activity=:last_activity WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *SessionsGen) UpdateSelective(ctx context.Context, m *model.Sessions) error {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, m.Id)
	query := "UPDATE sessions SET " + strings.Join(sets, ",") + " WHERE id=?"
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

func (d *SessionsGen) FindByIds(ctx context.Context, ids []string) ([]model.Sessions, error) {
	var list []model.Sessions
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

func (d *SessionsGen) UpdateByIds(ctx context.Context, ids []string, fields map[string]any) (int64, error) {
	if len(ids) == 0 || len(fields) == 0 { return 0, nil }
	var sets []string
	var args []any
	for k, v := range fields { sets = append(sets, k+"=?"); args = append(args, v) }
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

func (d *SessionsGen) nonNilFields(m *model.Sessions) ([]string, []string, []any) {
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
