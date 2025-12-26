package model

type UserDetail struct {
	Uid     int32  `db:"uid"`
	Address string `db:"address"`
}

func (UserDetail) TableName() string { return "user_detail" }
