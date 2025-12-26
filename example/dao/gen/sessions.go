package gen

import (
	"context"

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

func (d *SessionsGen) Update(ctx context.Context, m *model.Sessions) error {
	query := "UPDATE sessions SET user_id=:user_id,ip_address=:ip_address,user_agent=:user_agent,payload=:payload,last_activity=:last_activity WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
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
