package storage

import (
	"github.com/jinzhu/gorm"
)

type AccountTable struct {
	db *gorm.DB
}

func CreateAccountTable(db *gorm.DB) (*AccountTable, error) {
	if err := db.AutoMigrate(Account{}).Error; err != nil {
		return nil, err
	}
	return &AccountTable{db: db}, nil
}

func DropAccountTable(db *gorm.DB) error {
	if err := db.DropTableIfExists(Account{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *AccountTable) Create(a Account) error {
	if err := t.db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

func (t *AccountTable) Get(userID string) (a Account, ok bool, err error) {
	if err = t.db.Where(&Account{UserID: userID}).First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return a, false, nil
		}
		return a, false, err
	}
	return a, true, nil
}

func (t *AccountTable) SetPassword(userID, password string) error {
	if err := t.db.Model(&Account{UserID: userID}).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}
