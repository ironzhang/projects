package account

import (
	"errors"
	"time"

	"github.com/ironzhang/pearls/uuid"
	"github.com/ironzhang/practice/x-cloud/x-account/storage"
	"github.com/ironzhang/practice/x-cloud/x-account/types"
	"github.com/jinzhu/gorm"
)

type Manager struct {
	account *storage.AccountTable
	alias   *storage.AliasTable
	token   *storage.TokenTable
}

func NewManager(db *gorm.DB) (*Manager, error) {
	account, err := storage.CreateAccountTable(db)
	if err != nil {
		return nil, err
	}
	alias, err := storage.CreateAliasTable(db)
	if err != nil {
		return nil, err
	}
	token, err := storage.CreateTokenTable(db)
	if err != nil {
		return nil, err
	}
	return &Manager{
		account: account,
		alias:   alias,
		token:   token,
	}, nil
}

func (m *Manager) Register(typ, name, password string) error {
	_, ok, err := m.alias.Get(typ, name)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("this account has resistered")
	}

	userID := uuid.New().String()
	if err = m.account.Create(storage.Account{UserID: userID, Password: password}); err != nil {
		return err
	}
	if err = m.alias.Create(storage.Alias{Type: typ, Name: name, UserID: userID}); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Unregister() {
}

func (m *Manager) Login(typ, name, password string) (types.Token, error) {
	alias, ok, err := m.alias.Get(typ, name)
	if err != nil {
		return types.Token{}, err
	}
	if !ok {
		return types.Token{}, errors.New("this account has not resistered")
	}

	account, ok, err := m.account.Get(alias.UserID)
	if err != nil {
		return types.Token{}, err
	}
	if !ok {
		return types.Token{}, errors.New("this account has not resistered")
	}
	if password != account.Password {
		return types.Token{}, errors.New("the password is wrong")
	}

	t, err := m.login(account.UserID)
	if err != nil {
		return types.Token{}, err
	}
	return types.Token{
		UserID:        t.UserID,
		AccessToken:   t.AccessToken,
		AccessExpire:  t.AccessExpire,
		RefreshToken:  t.RefreshToken,
		RefreshExpire: t.RefreshExpire,
	}, nil
}

func (m *Manager) Logout() {
}

func (m *Manager) login(userID string) (storage.Token, error) {
	token, ok, err := m.token.Get(userID)
	if err != nil {
		return token, err
	}
	if !ok {
		token, err = m.createToken(userID)
		if err != nil {
			return token, err
		}
		return token, nil
	}
	token, err = m.refreshToken(token)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (m *Manager) createToken(userID string) (storage.Token, error) {
	t := storage.Token{
		UserID:        userID,
		AccessToken:   uuid.New().String(),
		RefreshToken:  uuid.New().String(),
		AccessExpire:  time.Now().Add(time.Hour),
		RefreshExpire: time.Now().Add(24 * time.Hour),
	}
	if err := m.token.Create(t); err != nil {
		return t, err
	}
	return t, nil
}

func (m *Manager) refreshToken(t storage.Token) (storage.Token, error) {
	now := time.Now()
	if t.AccessExpire.Before(now) {
		t.AccessToken = uuid.New().String()
	}
	if t.RefreshExpire.Before(now) {
		t.RefreshToken = uuid.New().String()
	}
	t.AccessExpire = time.Now().Add(time.Hour)
	t.RefreshExpire = time.Now().Add(24 * time.Hour)
	if err := m.token.Update(t); err != nil {
		return t, err
	}
	return t, nil
}
