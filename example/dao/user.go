package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type UserDAO struct {
	*gen.UserGen
}

func NewUserDAO(db *dbx.DB) *UserDAO {
	return &UserDAO{UserGen: gen.NewUserGen(db)}
}
