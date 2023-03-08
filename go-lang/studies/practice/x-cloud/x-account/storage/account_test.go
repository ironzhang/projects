package storage

import (
	"testing"
)

func OpenTestAccountTable(t *testing.T) *AccountTable {
	db := OpenTestDB(t)
	if err := DropAccountTable(db); err != nil {
		t.Fatalf("DropAccountTable: %v", err)
	}
	tb, err := CreateAccountTable(db)
	if err != nil {
		t.Fatalf("CreateAccountTable: %v", err)
	}
	return tb
}

func TestAccountTableCreateGet(t *testing.T) {
	tb := OpenTestAccountTable(t)
	tests := []struct {
		create  bool
		account Account
	}{
		{
			create:  true,
			account: Account{UserID: "1", Password: "123456"},
		},
		{
			create:  true,
			account: Account{UserID: "2", Password: "123456"},
		},
		{
			create:  false,
			account: Account{UserID: "3", Password: "123456"},
		},
	}
	for i, tt := range tests {
		if tt.create {
			if err := tb.Create(tt.account); err != nil {
				t.Fatalf("%d: Create: %v", i, err)
			}
			t.Logf("%d: create %v account", i, tt.account)
		}
	}
	for i, tt := range tests {
		a, ok, err := tb.Get(tt.account.UserID)
		if err != nil {
			t.Fatalf("%d: Get: %v", i, err)
		}
		if got, want := ok, tt.create; got != want {
			t.Fatalf("%d: Get: ok: got %v, want %v", i, got, want)
		}
		if ok {
			tt.account.CreatedAt = a.CreatedAt
			tt.account.UpdatedAt = a.UpdatedAt
			if got, want := a, tt.account; got != want {
				t.Fatalf("%d: Get: account: got %v, want %v", i, got, want)
			}
			t.Logf("%d: get %v account", i, a)
		}
	}
}

func TestAccountTableSetPassword(t *testing.T) {
	tb := OpenTestAccountTable(t)
	accounts := []Account{
		{UserID: "1", Password: "123456"},
		{UserID: "2", Password: "123456"},
		{UserID: "3", Password: "123456"},
	}
	for i, a := range accounts {
		if err := tb.Create(a); err != nil {
			t.Fatalf("%d: Create: %v", i, err)
		}
		t.Logf("%d: create %v account", i, a)
	}

	tests := []struct {
		userID   string
		password string
	}{
		{userID: "1", password: ""},
		{userID: "2", password: "1"},
		{userID: "3", password: "2"},
	}
	for i, tt := range tests {
		if err := tb.SetPassword(tt.userID, tt.password); err != nil {
			t.Fatalf("%d: SetPassword: %v", i, err)
		}
	}
}
