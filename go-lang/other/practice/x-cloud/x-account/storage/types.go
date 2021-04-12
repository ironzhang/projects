package storage

import (
	"time"
)

type Alias struct {
	Type   string `gorm:"type:varchar(32);not null;primary_key"`
	Name   string `gorm:"type:varchar(64);not null;primary_key"`
	UserID string `gorm:"type:varchar(128);not null"`
}

func (p Alias) TableName() string {
	return "alias"
}

type Account struct {
	UserID    string `gorm:"type:varchar(128);not null;primary_key"`
	Password  string `gorm:"type:varchar(128);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Account) TableName() string {
	return "account"
}

type Token struct {
	UserID        string `gorm:"type:varchar(128);not null;primary_key"`
	AccessToken   string `gorm:"type:varchar(64);not null"`
	RefreshToken  string `gorm:"type:varchar(64);not null"`
	AccessExpire  time.Time
	RefreshExpire time.Time
}

func (p Token) TableName() string {
	return "token"
}
