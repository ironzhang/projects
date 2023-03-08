package storage

import "github.com/jinzhu/gorm"

type AliasTable struct {
	db *gorm.DB
}

func CreateAliasTable(db *gorm.DB) (*AliasTable, error) {
	if err := db.AutoMigrate(Alias{}).Error; err != nil {
		return nil, err
	}
	return &AliasTable{db: db}, nil
}

func DropAliasTable(db *gorm.DB) error {
	if err := db.DropTableIfExists(Alias{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *AliasTable) Create(a Alias) error {
	if err := t.db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

func (t *AliasTable) Get(typ, name string) (a Alias, ok bool, err error) {
	if err = t.db.Where(&Alias{Type: typ, Name: name}).First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return a, false, nil
		}
		return a, false, err
	}
	return a, true, nil
}

func (t *AliasTable) Remove(typ, name string) error {
	return nil
}

func (t *AliasTable) RemoveByUserID(userID string) error {
	return nil
}
