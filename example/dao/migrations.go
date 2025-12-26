package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type MigrationsDAO struct {
	*gen.MigrationsGen
}

func NewMigrationsDAO(db *dbx.DB) *MigrationsDAO {
	return &MigrationsDAO{MigrationsGen: gen.NewMigrationsGen(db)}
}
