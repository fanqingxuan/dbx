package gen

import (
	"context"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type UserScoresGen struct {
	db *dbx.DB
}

func NewUserScoresGen(db *dbx.DB) *UserScoresGen {
	return &UserScoresGen{db: db}
}

func (d *UserScoresGen) Insert(ctx context.Context, m *model.UserScores) error {
	query := "INSERT INTO user_scores (id,uid,score) VALUES (:id,:uid,:score)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserScoresGen) Update(ctx context.Context, m *model.UserScores) error {
	query := "UPDATE user_scores SET uid=:uid,score=:score WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserScoresGen) Delete(ctx context.Context, id int32) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM user_scores WHERE id=?", id)
	return err
}

func (d *UserScoresGen) FindByID(ctx context.Context, id int32) (*model.UserScores, error) {
	var m model.UserScores
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM user_scores WHERE id=?", id)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *UserScoresGen) FindByIds(ctx context.Context, ids []int32) ([]model.UserScores, error) {
	var list []model.UserScores
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM user_scores WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *UserScoresGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *UserScoresGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *UserScoresGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *UserScoresGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *UserScoresGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}
