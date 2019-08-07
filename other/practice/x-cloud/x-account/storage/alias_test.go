package storage

import (
	"testing"
)

func OpenTestAliasTable(t *testing.T) *AliasTable {
	db := OpenTestDB(t)
	if err := DropAliasTable(db); err != nil {
		t.Fatalf("DropAliasTable: %v", err)
	}
	tb, err := CreateAliasTable(db)
	if err != nil {
		t.Fatalf("CreateAliasTable: %v", err)
	}
	return tb
}

func TestAliasTableCreateGet(t *testing.T) {
	tb := OpenTestAliasTable(t)
	tests := []struct {
		create bool
		alias  Alias
	}{
		{
			create: true,
			alias:  Alias{Type: "phone", Name: "135123456", UserID: "1"},
		},
		{
			create: true,
			alias:  Alias{Type: "email", Name: "zhanghui@ablecloud.cn", UserID: "2"},
		},
		{
			create: true,
			alias:  Alias{Type: "username", Name: "zhanghui", UserID: "2"},
		},
		{
			create: false,
			alias:  Alias{Type: "weixin", Name: "openid", UserID: "2"},
		},
	}
	for i, tt := range tests {
		if tt.create {
			if err := tb.Create(tt.alias); err != nil {
				t.Fatalf("%d: Create: %v", i, err)
			}
			t.Logf("%d: create %v alias", i, tt.alias)
		}
	}
	for i, tt := range tests {
		a, ok, err := tb.Get(tt.alias.Type, tt.alias.Name)
		if err != nil {
			t.Fatalf("%d: Get: %v", i, err)
		}
		if got, want := ok, tt.create; got != want {
			t.Fatalf("%d: Get: ok: got %v, want %v", i, got, want)
		}
		if ok {
			if got, want := a, tt.alias; got != want {
				t.Fatalf("%d: Get: alias: got %v, want %v", i, got, want)
			}
			t.Logf("%d: get %v alias", i, a)
		}
	}
}
