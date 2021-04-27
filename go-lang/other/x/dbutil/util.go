package dbutil

import (
	"database/sql"
	"fmt"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

// Options 数据库选项
type Options struct {
	Addr        string
	User        string
	Password    string
	Database    string
	MaxOpenConn int
	MaxIdleConn int
}

func DefaultOptions(user, password, database string) Options {
	return Options{
		Addr:        "localhost:3306",
		User:        user,
		Password:    password,
		Database:    database,
		MaxOpenConn: 50,
		MaxIdleConn: 20,
	}
}

// OpenDatabase 打开数据库
func OpenDatabase(driver string, opts Options) (*sql.DB, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", opts.User, opts.Password, opts.Addr, opts.Database)
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(opts.MaxOpenConn)
	db.SetMaxIdleConns(opts.MaxIdleConn)
	return db, nil
}

// CreateTable 创建数据表
func CreateTable(db *sql.DB, createSQL string) error {
	_, err := db.Exec(createSQL)
	return err
}

// DropTable 删除数据表
func DropTable(db *sql.DB, table string) error {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	return err
}

// TruncateTable 清空数据表
func TruncateTable(db *sql.DB, table string) error {
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
	return err
}

// ShowTables 显示数据表
func ShowTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var table string
	var tables []string
	for rows.Next() {
		if err = rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// Count 统计数据表记录数
func Count(db *sql.DB, table string) (count int, err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if err = db.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

/*
// Exec 执行SQL语句
func Exec(db *sql.DB, esql *easysql.EasySQL) (sql.Result, error) {
	query, args, err := esql.Query()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

// Query 执行SQL查询
func Query(db *sql.DB, esql *easysql.EasySQL) (*sql.Rows, error) {
	query, args, err := esql.Query()
	if err != nil {
		return nil, err
	}
	return db.Query(query, args...)
}

// QueryRow 查询一行数据
func QueryRow(db *sql.DB, esql *easysql.EasySQL) error {
	query, args, err := esql.Query()
	if err != nil {
		return err
	}
	return db.QueryRow(query, args...).Scan(esql.Vars()...)
}
*/
