package types

import "time"

type Token struct {
	UserID        string
	AccessToken   string
	AccessExpire  time.Time
	RefreshToken  string
	RefreshExpire time.Time
}
