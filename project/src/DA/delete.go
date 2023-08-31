package DA

import (
	"database/sql"
	"log"
)

// 通过id删除
func Del(id int) (rows int, err error) {
	//1.操作数据库
	db, err := sql.Open("sqlite3", "./nas.db")
	//错误检查
	if err != nil {
		log.Fatal(err.Error())
	}
	//推迟数据库连接的关闭
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM user WHERE user_id=?")
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
