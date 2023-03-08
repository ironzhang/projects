package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ironzhang/practice/x-cloud/x-account/account"
	"github.com/ironzhang/practice/x-cloud/x-account/handler"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func main() {
	db, err := gorm.Open("mysql", "root:123456@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("gorm.Open: %v", err)
	}
	a, err := account.NewManager(db)
	if err != nil {
		log.Fatalf("account.NewManager: %v", err)
	}
	h := handler.New(a)
	e := echo.New()
	handler.Register(e, h)
	e.Logger.Fatal(e.Start(":8000"))
}
