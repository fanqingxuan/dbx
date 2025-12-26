package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type UsersDAO struct {
	*gen.UsersGen
}

func NewUsersDAO(db *dbx.DB) *UsersDAO {
	return &UsersDAO{UsersGen: gen.NewUsersGen(db)}
}
