package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ironzhang/x/dbutil"
)

var DefaultOptions = dbutil.DefaultOptions("root", "123456", "wordbook")
