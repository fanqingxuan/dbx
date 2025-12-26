package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type UserDetailDAO struct {
	*gen.UserDetailGen
}

func NewUserDetailDAO(db *dbx.DB) *UserDetailDAO {
	return &UserDetailDAO{UserDetailGen: gen.NewUserDetailGen(db)}
}
