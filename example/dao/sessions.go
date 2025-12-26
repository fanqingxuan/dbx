package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type SessionsDAO struct {
	*gen.SessionsGen
}

func NewSessionsDAO(db *dbx.DB) *SessionsDAO {
	return &SessionsDAO{SessionsGen: gen.NewSessionsGen(db)}
}
