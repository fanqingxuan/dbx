package gen

import (
	"context"
	"strings"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/do"
	"github.com/fanqingxuan/dbx/example/model"
)

type JobsGen struct {
	db *dbx.DB
}

func NewJobsGen(db *dbx.DB) *JobsGen {
	return &JobsGen{db: db}
}

func (d *JobsGen) Insert(ctx context.Context, m do.Jobs) error {
	cols, vals, args := d.nonNilFields(m)
	if len(cols) == 0 { return nil }
	query := "INSERT INTO jobs (" + strings.Join(cols, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	_, err := d.db.ExecCtx(ctx, query, args...)
	return err
}

func (d *JobsGen) Delete(ctx context.Context, id int32) error {
	_, err := d.db.ExecCtx(ctx, "DELETE FROM jobs WHERE id=?", id)
	return err
}

func (d *JobsGen) FindByID(ctx context.Context, id int32) (*model.Jobs, error) {
	var m model.Jobs
	err := d.db.QueryRowCtx(ctx, &m, "SELECT * FROM jobs WHERE id=?", id)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *JobsGen) FindByIds(ctx context.Context, ids []int32) ([]*model.Jobs, error) {
	var list []*model.Jobs
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM jobs WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (d *JobsGen) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	query, args, _ := d.db.In("DELETE FROM jobs WHERE id IN (?)", ids)
	return d.ExecCtx(ctx, query, args...)
}

func (d *JobsGen) UpdateById(ctx context.Context, m do.Jobs, id int32) (int64, error) {
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	args = append(args, id)
	query := "UPDATE jobs SET " + strings.Join(sets, ",") + " WHERE id=?"
	return d.ExecCtx(ctx, query, args...)
}

func (d *JobsGen) UpdateByIds(ctx context.Context, m do.Jobs, ids []int32) (int64, error) {
	if len(ids) == 0 { return 0, nil }
	cols, _, args := d.nonNilFields(m)
	if len(cols) == 0 { return 0, nil }
	var sets []string
	for _, c := range cols { sets = append(sets, c+"=?") }
	query := "UPDATE jobs SET " + strings.Join(sets, ",") + " WHERE id IN (?)"
	query, inArgs, _ := d.db.In(query, ids)
	args = append(args, inArgs...)
	return d.ExecCtx(ctx, query, args...)
}

func (d *JobsGen) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowsCtx(ctx, dest, query, args...)
}

func (d *JobsGen) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryRowCtx(ctx, dest, query, args...)
}

func (d *JobsGen) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryValueCtx(ctx, dest, query, args...)
}

func (d *JobsGen) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.QueryColumnCtx(ctx, dest, query, args...)
}

func (d *JobsGen) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := d.db.ExecCtx(ctx, query, args...)
	if err != nil { return 0, err }
	return result.RowsAffected()
}

func (d *JobsGen) nonNilFields(m do.Jobs) ([]string, []string, []any) {
	var cols, vals []string
	var args []any
	if m.Id != nil { cols = append(cols, "id"); vals = append(vals, "?"); args = append(args, m.Id) }
	if m.Queue != nil { cols = append(cols, "queue"); vals = append(vals, "?"); args = append(args, m.Queue) }
	if m.Payload != nil { cols = append(cols, "payload"); vals = append(vals, "?"); args = append(args, m.Payload) }
	if m.Attempts != nil { cols = append(cols, "attempts"); vals = append(vals, "?"); args = append(args, m.Attempts) }
	if m.ReserveTime != nil { cols = append(cols, "reserve_time"); vals = append(vals, "?"); args = append(args, m.ReserveTime) }
	if m.AvailableTime != nil { cols = append(cols, "available_time"); vals = append(vals, "?"); args = append(args, m.AvailableTime) }
	if m.CreateTime != nil { cols = append(cols, "create_time"); vals = append(vals, "?"); args = append(args, m.CreateTime) }
	return cols, vals, args
}
