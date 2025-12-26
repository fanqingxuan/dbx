package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type UserScoresDAO struct {
	*gen.UserScoresGen
}

func NewUserScoresDAO(db *dbx.DB) *UserScoresDAO {
	return &UserScoresDAO{UserScoresGen: gen.NewUserScoresGen(db)}
}
