package easysql_test

import (
	"fmt"
	"os"

	"github.com/ironzhang/x/easysql"
)

type Account struct {
	Id       int64
	Name     string
	Password string
	Status   int8
}

//INSERT INTO tb_account (name, password, status) VALUES("iron", "123456", 1);
func ExampleInsertQuery1() {
	esql := easysql.InsertInto("tb_account").Columns("name, password, status", "iron", "123456", 1)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// INSERT INTO tb_account (name,password,status) VALUES (?,?,?)
	// [iron 123456 1]
}

//DELETE FROM tb_account where id=1;
func ExampleDeleteQuery1() {
	esql := easysql.DeleteFrom("tb_account").Where("id=?", 1)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// DELETE FROM tb_account WHERE (id=?)
	// [1]
}

//UPDATE tb_account SET name="iron", password="123456" where id=1;
func ExampleUpdateQuery1() {
	esql := easysql.Update("tb_account").Columns("name, password", "iron", "123456").Where("id=?", 1)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// UPDATE tb_account SET name=?, password=? WHERE (id=?)
	// [iron 123456 1]
}

//UPDATE `tb_user` SET `image_url`="image" WHERE `uid`=1
func ExampleUpdateQuery2() {
	esql := easysql.Update("tb_user").Column("image_url", "image").Where("uid=?", 1)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// UPDATE tb_user SET image_url=? WHERE (uid=?)
	// [image 1]
}

// REPLACE INTO `tb_account` (id, name, password) VALUES (1, "iron", "123456");
func ExampleReplaceQuery1() {
	esql := easysql.Replace("tb_account").Columns("id, name, password", 1, "iron", "123456")
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// REPLACE INTO tb_account (id,name,password) VALUES (?,?,?)
	// [1 iron 123456]
}

//SELECT id, name, password FROM tb_account LIMIT 0, 100;
func ExampleSelectQuery1() {
	a := Account{}
	esql := easysql.SelectFrom("tb_account").Columns("id , name , password", &a.Id, &a.Name, &a.Password).Limit(0, 100)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// SELECT id,name,password FROM tb_account LIMIT ?,?
	// [0 100]
}

//SELECT COUNT(*) FROM tb_account where status>0 and status<10;
func ExampleSelectQuery2() {
	var count int
	esql := easysql.SelectFrom("tb_account").Functions("COUNT(*)", &count).Where("status>?", 0).AndWhere("status<?", 10)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// SELECT COUNT(*) FROM tb_account WHERE (status>?) AND (status<?)
	// [0 10]
}

// SELECT id, name, password FROM tb_account ORDER BY id ASC LIMIT 0, 1;
func ExampleSelectQuery3() {
	a := Account{}
	esql := easysql.SelectFrom("tb_account").Columns("id, name, password", &a.Id, &a.Name, &a.Password).OrderBy(true, "id").Limit(0, 1)
	query, args, err := esql.Query()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(query)
	fmt.Println(args)

	// output:
	// SELECT id,name,password FROM tb_account ORDER BY id ASC LIMIT ?,?
	// [0 1]
}
