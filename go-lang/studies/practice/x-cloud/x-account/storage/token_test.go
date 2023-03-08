package storage

import (
	"testing"
	"time"
)

func OpenTestTokenTable(t *testing.T) *TokenTable {
	db := OpenTestDB(t)
	if err := DropTokenTable(db); err != nil {
		t.Fatalf("DropTokenTable: %v", err)
	}
	tb, err := CreateTokenTable(db)
	if err != nil {
		t.Fatalf("CreateTokenTable: %v", err)
	}
	return tb
}

func TestTokenTable(t *testing.T) {
	var err error
	tb := OpenTestTokenTable(t)

	err = tb.Create(Token{
		UserID:        "1",
		AccessToken:   "AccessToken",
		RefreshToken:  "RefreshToken",
		AccessExpire:  time.Now().Add(time.Hour),
		RefreshExpire: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	err = tb.Create(Token{
		UserID:        "2",
		AccessToken:   "AccessToken",
		RefreshToken:  "RefreshToken",
		AccessExpire:  time.Now().Add(time.Hour),
		RefreshExpire: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	err = tb.Update(Token{
		UserID:        "1",
		AccessToken:   "NewAccessToken",
		RefreshToken:  "NewRefreshToken",
		AccessExpire:  time.Now().Add(2 * time.Hour),
		RefreshExpire: time.Now().Add(48 * time.Hour),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	tk, ok, err := tb.Get("2")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !ok {
		t.Fatalf("Get: not found")
	}
	t.Logf("Token: %v", tk)
}
