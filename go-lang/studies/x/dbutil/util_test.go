package dbutil

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestOpenDatabase(t *testing.T) {
	opts := Options{
		Addr:     "localhost:3306",
		User:     "root",
		Password: "123456",
		Database: "test",
	}
	db, err := OpenDatabase("mysql", opts)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	db.Close()
}

func OpenTestDatabase() *sql.DB {
	opts := Options{
		Addr:     "localhost:3306",
		User:     "root",
		Password: "123456",
		Database: "test",
	}
	db, err := OpenDatabase("mysql", opts)
	if err != nil {
		panic(err)
	}
	return db
}

func LookupTable(db *sql.DB, table string) (bool, error) {
	tables, err := ShowTables(db)
	if err != nil {
		return false, err
	}
	for _, tb := range tables {
		if tb == table {
			return true, nil
		}
	}
	return false, nil
}

func TestCreateDropTable(t *testing.T) {
	var tb = "test_tb1"
	var createSQL = `
	CREATE TABLE IF NOT EXISTS test_tb1 (
		id BIGINT AUTO_INCREMENT NOT NULL,
		c1 VARCHAR(64) NOT NULL,
		c2 VARCHAR(64) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	db := OpenTestDatabase()
	if err := CreateTable(db, createSQL); err != nil {
		t.Fatalf("create table: %v", err)
	}
	found, err := LookupTable(db, tb)
	if err != nil {
		t.Fatalf("lookup table: %v", err)
	}
	if !found {
		t.Fatalf("%s table not create", tb)
	}

	if err := DropTable(db, tb); err != nil {
		t.Fatalf("drop table: %v", err)
	}
	found, err = LookupTable(db, tb)
	if err != nil {
		t.Fatalf("lookup table: %v", err)
	}
	if found {
		t.Fatalf("%s table not drop", tb)
	}
}

func TestTruncateTable(t *testing.T) {
	var tb = "test_tb2"
	var createSQL = `
	CREATE TABLE IF NOT EXISTS test_tb2 (
		id BIGINT AUTO_INCREMENT NOT NULL,
		c1 VARCHAR(64) NOT NULL,
		c2 VARCHAR(64) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	db := OpenTestDatabase()
	if err := CreateTable(db, createSQL); err != nil {
		t.Fatalf("create table: %v", err)
	}

	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s SET c1='c1', c2='c2'", tb))
	if err != nil {
		t.Fatalf("insert into: %v", err)
	}
	count, err := Count(db, tb)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count <= 0 {
		t.Fatalf("count(%d) <= 0", count)
	} else {
		t.Logf("count: %d", count)
	}

	if err = TruncateTable(db, tb); err != nil {
		t.Fatalf("truncate table: %v", err)
	}
	count, err = Count(db, tb)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count > 0 {
		t.Fatalf("count(%d) > 0", count)
	} else {
		t.Logf("count: %d", count)
	}
}
