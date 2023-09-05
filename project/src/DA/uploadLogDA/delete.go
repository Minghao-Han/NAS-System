package uploadDA

import (
	"database/sql"
	"log"
	"nas/project/src/Utils"
)

var dbPath = Utils.DefaultConfigReader().Get("Database:dbPath").(string)

// 通过id删除
func DeleteById(id int) (rows int, err error) {
	//1.操作数据库
	db, err := sql.Open("sqlite3", dbPath)
	//错误检查
	if err != nil {
		log.Fatal(err.Error())
	}
	//推迟数据库连接的关闭
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM upload_log WHERE id=?")
	if err != nil {
		log.Fatalln(err)
	}

	rs, err := stmt.Exec(id)
	if err != nil {
		log.Fatalln(err)
	}
	//删除的行数
	row, err := rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	rows = int(row)
	return
}

func DeleteByPath(path string) (rows int, err error) {
	//1.操作数据库
	db, err := sql.Open("sqlite3", dbPath)
	//错误检查
	if err != nil {
		log.Fatal(err.Error())
	}
	//推迟数据库连接的关闭
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM upload_log WHERE path=?")
	if err != nil {
		log.Fatalln(err)
	}

	rs, err := stmt.Exec(path)
	if err != nil {
		log.Fatalln(err)
	}
	//删除的行数
	row, err := rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	rows = int(row)
	return
}
