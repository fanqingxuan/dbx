package model

import "time"

type Migrations struct {
	Id        int64      `db:"id"`
	Migration string     `db:"migration"`
	Batch     int64      `db:"batch"`
	CreatedAt *time.Time `db:"created_at"`
}

func (Migrations) TableName() string { return "migrations" }
