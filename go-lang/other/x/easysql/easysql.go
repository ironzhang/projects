package easysql

import (
	"fmt"
	"strings"
)

type functionItem struct {
	name  string
	value interface{}
}

type columnItem struct {
	name  string
	value interface{}
}

type whereItem struct {
	op    string
	where string
	args  []interface{}
}

type limitItem struct {
	offset int
	limit  int
}

type orderItem struct {
	sort string
	cols []string
}

type EasySQL struct {
	op        string
	table     string
	err       error
	functions []functionItem
	columns   []columnItem
	wheres    []whereItem
	order     *orderItem
	limit     *limitItem
}

func InsertInto(table string) *EasySQL {
	return &EasySQL{op: "INSERT", table: table}
}

func DeleteFrom(table string) *EasySQL {
	return &EasySQL{op: "DELETE", table: table}
}

func Update(table string) *EasySQL {
	return &EasySQL{op: "UPDATE", table: table}
}

func Replace(table string) *EasySQL {
	return &EasySQL{op: "REPLACE", table: table}
}

func SelectFrom(table string) *EasySQL {
	return &EasySQL{op: "SELECT", table: table}
}

func (s *EasySQL) Err() error {
	return s.err
}

func (s *EasySQL) Function(fun string, arg interface{}) *EasySQL {
	s.functions = append(s.functions, functionItem{name: fun, value: arg})
	return s
}

func (s *EasySQL) Functions(funs string, args ...interface{}) *EasySQL {
	items := strings.Split(funs, ",")
	if len(items) != len(args) {
		s.err = opErrorf("Functions", "arguments is unexpect, %q %v", funs, args)
		return s
	}
	for i := 0; i < len(args); i++ {
		s.Function(strings.TrimSpace(items[i]), args[i])
	}
	return s
}

func (s *EasySQL) Column(col string, arg interface{}) *EasySQL {
	s.columns = append(s.columns, columnItem{name: col, value: arg})
	return s
}

func (s *EasySQL) Columns(cols string, args ...interface{}) *EasySQL {
	items := strings.Split(cols, ",")
	if len(items) != len(args) {
		s.err = opErrorf("Columns", "arguments is unexpect, %q %v", cols, args)
		return s
	}
	for i := 0; i < len(args); i++ {
		s.Column(strings.TrimSpace(items[i]), args[i])
	}
	return s
}

func (s *EasySQL) Where(where string, args ...interface{}) *EasySQL {
	return s.AndWhere(where, args...)
}

func (s *EasySQL) AndWhere(where string, args ...interface{}) *EasySQL {
	s.wheres = append(s.wheres, whereItem{op: "AND", where: where, args: args})
	return s
}

func (s *EasySQL) OrWhere(where string, args ...interface{}) *EasySQL {
	s.wheres = append(s.wheres, whereItem{op: "OR", where: where, args: args})
	return s
}

func (s *EasySQL) OrderBy(asc bool, cols ...string) *EasySQL {
	if len(cols) <= 0 {
		s.err = opErrorf("OrderBy", "no order by column")
	}
	var sort string
	if asc {
		sort = "ASC"
	} else {
		sort = "DESC"
	}
	s.order = &orderItem{sort: sort, cols: cols}
	return s
}

func (s *EasySQL) Limit(offset, limit int) *EasySQL {
	s.limit = &limitItem{offset: offset, limit: limit}
	return s
}

func (s *EasySQL) Vars() (vars []interface{}) {
	for _, item := range s.functions {
		vars = append(vars, item.value)
	}
	for _, item := range s.columns {
		vars = append(vars, item.value)
	}
	return vars
}

func (s *EasySQL) Query() (query string, args []interface{}, err error) {
	if s.err != nil {
		return "", nil, err
	}

	switch s.op {
	case "INSERT":
		return s.insertQuery()
	case "DELETE":
		return s.deleteQuery()
	case "UPDATE":
		return s.updateQuery()
	case "REPLACE":
		return s.replaceQuery()
	case "SELECT":
		return s.selectQuery()
	}

	return "", nil, opErrorf("Query", "%q unsupported", s.op)
}

func (s *EasySQL) insertQuery() (query string, args []interface{}, err error) {
	var columns []string
	for _, item := range s.columns {
		columns = append(columns, item.name)
	}
	values := ""
	n := len(s.columns)
	if n > 0 {
		values = "?" + strings.Repeat(",?", n-1)
	}
	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", s.table, strings.Join(columns, ","), values)
	args = s.columnArgs()
	return
}

func (s *EasySQL) deleteQuery() (query string, args []interface{}, err error) {
	query = fmt.Sprintf("DELETE FROM %s", s.table)
	if len(s.wheres) > 0 {
		query = query + " WHERE " + s.whereQuery()
		args = s.whereArgs()
	}
	return
}

func (s *EasySQL) updateQuery() (query string, args []interface{}, err error) {
	var items []string
	for _, item := range s.columns {
		items = append(items, fmt.Sprintf("%s=?", item.name))
	}
	query = fmt.Sprintf("UPDATE %s SET %s", s.table, strings.Join(items, ", "))
	args = s.columnArgs()
	if len(s.wheres) > 0 {
		query = query + " WHERE " + s.whereQuery()
		args = append(args, s.whereArgs()...)
	}
	return
}

func (s *EasySQL) replaceQuery() (query string, args []interface{}, err error) {
	var columns []string
	for _, item := range s.columns {
		columns = append(columns, item.name)
	}
	values := ""
	n := len(s.columns)
	if n > 0 {
		values = "?" + strings.Repeat(",?", n-1)
	}
	query = fmt.Sprintf("REPLACE INTO %s (%s) VALUES (%s)", s.table, strings.Join(columns, ","), values)
	args = s.columnArgs()
	return
}

func (s *EasySQL) selectQuery() (query string, args []interface{}, err error) {
	var items []string
	for _, item := range s.functions {
		items = append(items, item.name)
	}
	for _, item := range s.columns {
		items = append(items, item.name)
	}
	query = fmt.Sprintf("SELECT %s FROM %s", strings.Join(items, ","), s.table)
	if len(s.wheres) > 0 {
		query = query + " WHERE " + s.whereQuery()
		args = s.whereArgs()
	}
	if s.order != nil {
		query = query + " ORDER BY " + s.orderQuery() + " " + s.order.sort
	}
	if s.limit != nil {
		query = query + " LIMIT ?,?"
		args = append(args, s.limit.offset, s.limit.limit)
	}
	return
}

func (s *EasySQL) whereQuery() string {
	var wheres []string
	for _, item := range s.wheres {
		if len(wheres) > 0 {
			wheres = append(wheres, item.op)
		}
		where := fmt.Sprintf("(%s)", item.where)
		wheres = append(wheres, where)
	}
	return strings.Join(wheres, " ")
}

func (s *EasySQL) orderQuery() string {
	var columns []string
	for _, column := range s.order.cols {
		columns = append(columns, column)
	}
	return strings.Join(columns, ",")
}

func (s *EasySQL) whereArgs() []interface{} {
	var args []interface{}
	for _, item := range s.wheres {
		args = append(args, item.args...)
	}
	return args
}

func (s *EasySQL) columnArgs() []interface{} {
	var args []interface{}
	for _, item := range s.columns {
		args = append(args, item.value)
	}
	return args
}

func (s *EasySQL) String() string {
	query, _, _ := s.Query()
	return query
}
