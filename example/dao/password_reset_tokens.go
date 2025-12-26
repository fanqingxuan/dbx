package dao

import (
	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/fanqingxuan/dbx/example/dao/gen"
)

type PasswordResetTokensDAO struct {
	*gen.PasswordResetTokensGen
}

func NewPasswordResetTokensDAO(db *dbx.DB) *PasswordResetTokensDAO {
	return &PasswordResetTokensDAO{PasswordResetTokensGen: gen.NewPasswordResetTokensGen(db)}
}
