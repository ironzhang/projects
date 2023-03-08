package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/account")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(50)
	if err = db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}

	for {
		var uid int
		var phone string
		if err = db.QueryRow("select uid, phone from 1_account_info where uid=1").Scan(&uid, &phone); err != nil {
			log.Printf("query row fail: %v", err)
		} else {
			log.Printf("query row success: uid[%d], phone[%s]", uid, phone)
		}
		time.Sleep(time.Second)
	}
}
