package storage

import (
	"database/sql"
	"encoding/json"

	"github.com/ironzhang/wordbook/cores/types"
	"github.com/ironzhang/x/easysql"
)

const WordTableName = "words"

type WordTable struct {
	db *sql.DB
}

func NewWordTable(db *sql.DB) *WordTable {
	return new(WordTable).Init(db)
}

func (t *WordTable) Init(db *sql.DB) *WordTable {
	t.db = db
	return t
}

func (t *WordTable) Set(w types.Word) error {
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	esql := easysql.Replace(WordTableName).Columns("word,symbol", w.Word, data)
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

func (t *WordTable) Get(word string) (w types.Word, ok bool, err error) {
	var data []byte
	esql := easysql.SelectFrom(WordTableName).Column("symbol", &data).Where("word=?", word)
	query, args, err := esql.Query()
	if err != nil {
		return w, false, err
	}
	err = t.db.QueryRow(query, args...).Scan(esql.Vars()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return w, false, nil
		}
		return w, false, err
	}
	err = json.Unmarshal(data, &w)
	if err != nil {
		return w, false, err
	}
	return w, true, nil
}

func (t *WordTable) Remove(word string) error {
	esql := easysql.DeleteFrom(WordTableName).Where("word=?", word)
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

func sliceCap(limit int) int {
	if limit > 0 {
		return limit
	}
	return 10
}

func (t *WordTable) List(offset, limit int) ([]types.Word, error) {
	var data []byte
	esql := easysql.SelectFrom(WordTableName).Column("symbol", &data).Limit(offset, limit)
	query, args, err := esql.Query()
	if err != nil {
		return nil, err
	}
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	vars := esql.Vars()
	words := make([]types.Word, 0, sliceCap(limit))
	for rows.Next() {
		if err = rows.Scan(vars...); err != nil {
			return nil, err
		}
		var w types.Word
		if err = json.Unmarshal(data, &w); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, nil
}

func (t *WordTable) ListByPriority(offset, limit int) ([]types.Word, error) {
	var data []byte
	esql := easysql.SelectFrom(WordTableName + "," + PriorityTableName)
	esql = esql.Column("symbol", &data)
	esql = esql.Where("words.word = priorities.word")
	esql = esql.OrderBy(false, "priority")
	esql = esql.Limit(offset, limit)
	query, args, err := esql.Query()
	if err != nil {
		return nil, err
	}
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	vars := esql.Vars()
	words := make([]types.Word, 0, sliceCap(limit))
	for rows.Next() {
		if err = rows.Scan(vars...); err != nil {
			return nil, err
		}
		var w types.Word
		if err = json.Unmarshal(data, &w); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, nil
}
