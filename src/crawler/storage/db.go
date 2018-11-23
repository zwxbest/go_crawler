package storage

import (
	"database/sql"
	"log"
	"fmt"
)

var Db sql.DB
var InsetIntoCategory = "insert into category(name,url,parent_id) values (?,?,?);"

// 获取表数据
func Get() {
	rows, err := Db.Query("select * from user;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	cloumns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	// for rows.Next() {
	//  err := rows.Scan(&cloumns[0], &cloumns[1], &cloumns[2])
	//  if err != nil {
	//      log.Fatal(err)
	//  }
	//  fmt.Println(cloumns[0], cloumns[1], cloumns[2])
	// }
	values := make([]sql.RawBytes, len(cloumns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(cloumns[i], ": ", value)
		}
		fmt.Println("------------------")
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// 插入数据
func Insert(category NewsCategory) {

	stmt, err := Db.Prepare(InsetIntoCategory)
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(category.Url, 0)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

// 删除数据
func Delete() {
	stmt, err := Db.Prepare("DELETE FROM user WHERE name='python'")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec()
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

// 更新数据
func Update(db *sql.DB) {
	stmt, err := db.Prepare("UPDATE user SET age=27 WHERE name='python'")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec()
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}
