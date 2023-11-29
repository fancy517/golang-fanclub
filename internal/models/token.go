package models

import "time"

type Token struct {
	UserID int       `json:"user_id"`
	Hash   string    `json:"hash"`
	Expiry time.Time `json:"expiry"`
}
