package model

import "time"

type PasswordResetTokens struct {
	Email     string     `db:"email"`
	Token     string     `db:"token"`
	CreatedAt *time.Time `db:"created_at"`
}

func (PasswordResetTokens) TableName() string { return "password_reset_tokens" }
