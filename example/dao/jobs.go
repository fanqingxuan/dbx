package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type JobsDAO struct {
	*gen.JobsGen
}

func NewJobsDAO(db *dbx.DB) *JobsDAO {
	return &JobsDAO{JobsGen: gen.NewJobsGen(db)}
}
