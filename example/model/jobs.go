package model

type Jobs struct {
	Id            int32   `db:"id"`
	Queue         *string `db:"queue"`
	Payload       *string `db:"payload"`
	Attempts      *int8   `db:"attempts"`
	ReserveTime   *int32  `db:"reserve_time"`
	AvailableTime *int32  `db:"available_time"`
	CreateTime    *int32  `db:"create_time"`
}

func (Jobs) TableName() string { return "jobs" }
