package gen

import (
	"context"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/model"
)

type JobsGen struct {
	db *dbx.DB
}

func NewJobsGen(db *dbx.DB) *JobsGen {
	return &JobsGen{db: db}
}

func (d *JobsGen) Insert(ctx context.Context, m *model.Jobs) error {
	query := "INSERT INTO jobs (id,queue,payload,attempts,reserve_time,available_time,create_time) VALUES (:id,:queue,:payload,:attempts,:reserve_time,:available_time,:create_time)"
	_, err := d.db.NamedExecContext(ctx, query, m)
	return err
}

func (d *JobsGen) Update(ctx context.Context, m *model.Jobs) error {
	query := "UPDATE jobs SET queue=:queue,payload=:payload,attempts=:attempts,reserve_time=:reserve_time,available_time=:available_time,create_time=:create_time WHERE id=:id"
	_, err := d.db.NamedExecContext(ctx, query, m)
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

func (d *JobsGen) FindByIds(ctx context.Context, ids []int32) ([]model.Jobs, error) {
	var list []model.Jobs
	if len(ids) == 0 { return list, nil }
	query, args, _ := d.db.In("SELECT * FROM jobs WHERE id IN (?)", ids)
	err := d.db.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
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
