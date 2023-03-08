package storage

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func OpenTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return db
}
