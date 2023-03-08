package storage

import "github.com/jinzhu/gorm"

type TokenTable struct {
	db *gorm.DB
}

func CreateTokenTable(db *gorm.DB) (*TokenTable, error) {
	if err := db.AutoMigrate(Token{}).Error; err != nil {
		return nil, err
	}
	return &TokenTable{db: db}, nil
}

func DropTokenTable(db *gorm.DB) error {
	if err := db.DropTableIfExists(Token{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *TokenTable) Create(tk Token) error {
	if err := t.db.Create(&tk).Error; err != nil {
		return err
	}
	return nil
}

func (t *TokenTable) Update(tk Token) error {
	if err := t.db.Model(&tk).Update(&tk).Error; err != nil {
		return err
	}
	return nil
}

func (t *TokenTable) Get(userID string) (tk Token, ok bool, err error) {
	if err = t.db.Where(&Token{UserID: userID}).First(&tk).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return tk, false, nil
		}
		return tk, false, err
	}
	return tk, true, nil
}
