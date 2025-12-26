package model

type UserScores struct {
	Id    int32 `db:"id"`
	Uid   int32 `db:"uid"`
	Score int32 `db:"score"`
}

func (UserScores) TableName() string { return "user_scores" }
