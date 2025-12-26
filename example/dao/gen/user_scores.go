package gen

import (
	"context"
	"reflect"
	"strings"

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

func (d *UserScoresGen) InsertSelective(ctx context.Context, m *model.UserScores) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO user_scores (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *UserScoresGen) Update(ctx context.Context, m *model.UserScores) error {
	query := "UPDATE user_scores SET uid=:uid,score=:score WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *UserScoresGen) UpdateSelective(ctx context.Context, m *model.UserScores) error {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, m.Id)
	query := "UPDATE user_scores SET " + strings.Join(sets, ",") + " WHERE id=?"
	_, err := d.db.ExecCtx(ctx, query, args...)
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

func (d *UserScoresGen) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM user_scores WHERE id IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *UserScoresGen) UpdateByIds(ctx context.Context, ids []int32, fields map[string]any) (int64, error) {
	if len(ids) == 0 || len(fields) == 0 { return 0, nil }
	var sets []string
	var args []any
	for k, v := range fields { sets = append(sets, k+"=?"); args = append(args, v) }
	query := "UPDATE user_scores SET " + strings.Join(sets, ",") + " WHERE id IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
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

func (d *UserScoresGen) nonNilFields(m *model.UserScores) ([]string, []string, []any) {
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
