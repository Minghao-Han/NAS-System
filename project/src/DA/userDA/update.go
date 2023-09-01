package userDA

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"nas/project/src/Entities"
)

func Update(user Entities.User) (rowsAffected int64, err error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//推迟数据库连接的关闭
	defer db.Close()
	stmt, err := db.Prepare("UPDATE  user SET user_name=?, password=? ,total_capacity=?,margin=? WHERE user_id=?")
	if err != nil {
		return
	}
	//执行修改操作
	rs, err := stmt.Exec(user.UserId, user.UserName, user.Password, user.Capacity, user.Margin)
	if err != nil {
		return
	}
	//返回插入的id
	rowsAffected, err = rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	return

}
