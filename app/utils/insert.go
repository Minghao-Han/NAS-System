package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// 定义用户信息结构体
type User struct {
	UserId   int    `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Capacity int    `json:"capacity" form:"capacity"`
	Margin   int    `json:"margin" form:"margin"`
}

func Insert(user User) (Id int, err error) {
	db, err := sql.Open("sqlite3", "./nas.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	//推迟数据库连接的关闭
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user(user_id, user_name, password,total_capacity,margin) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("insert data failed: %v\n", err)
	}

	//执行插入操作
	rs, err := stmt.Exec(user.UserId, user.UserName, user.Password, user.Capacity, user.Margin)
	if err != nil {
		return
	}

	//返回插入的id
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}
	//将id类型转换
	Id = int(id)
	defer stmt.Close()
	return

}
