package model

import "time"

type Users struct {
	Id              int64      `db:"id"`
	Name            string     `db:"name"`
	Email           string     `db:"email"`
	EmailVerifiedAt *time.Time `db:"email_verified_at"`
	Password        string     `db:"password"`
	RememberToken   *string    `db:"remember_token"`
	CreatedAt       *time.Time `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

func (Users) TableName() string { return "users" }
