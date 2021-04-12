package storage

import (
	"database/sql"

	"github.com/ironzhang/x/easysql"
)

const PriorityTableName = "priorities"

type PriorityTable struct {
	db *sql.DB
}

func NewPriorityTable(db *sql.DB) *PriorityTable {
	return new(PriorityTable).Init(db)
}

func (t *PriorityTable) Init(db *sql.DB) *PriorityTable {
	t.db = db
	return t
}

func (t *PriorityTable) Set(word string, priority int) error {
	esql := easysql.Replace(PriorityTableName).Columns("word,priority", word, priority)
	query, args, err := esql.Query()
	if err != nil {
		return err
	}
	_, err = t.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (t *PriorityTable) Get(word string) (priority int, ok bool, err error) {
	esql := easysql.SelectFrom(PriorityTableName).Columns("priority", &priority).Where("word=?", word)
	query, args, err := esql.Query()
	if err != nil {
		return 0, false, err
	}
	err = t.db.QueryRow(query, args...).Scan(esql.Vars()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}
	return priority, true, nil
}
