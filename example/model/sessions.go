package model

type Sessions struct {
	Id           string  `db:"id"`
	UserId       *int64  `db:"user_id"`
	IpAddress    *string `db:"ip_address"`
	UserAgent    *string `db:"user_agent"`
	Payload      string  `db:"payload"`
	LastActivity int32   `db:"last_activity"`
}

func (Sessions) TableName() string { return "sessions" }
