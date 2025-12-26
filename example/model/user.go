package model

import "time"

type User struct {
	Id        int64      `db:"id"`
	Name      string     `db:"name"`
	Status    *int8      `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Age       *int32     `db:"age"`
}

func (User) TableName() string { return "user" }
