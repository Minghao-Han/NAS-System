package userDA

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"nas/project/src/DA/Entities"
)

func Query() (users []Entities.User, err error) {
	db, err := sql.Open("sqlite3", "./nas.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	//_, err = db.Exec("INSERT INTO user(user_id, user_name, password,total_capacity,margin) VALUES (?, ?, ?, ?, ?)", 1, "张三", 111, 2, 2)
	//if err != nil {
	//	fmt.Printf("insert data failed: %v\n", err)
	//}
	rows, err := db.Query("SELECT * FROM user ")
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return
	}

	for rows.Next() {
		var myUser Entities.User
		//遍历表中所有行的信息
		rows.Scan(&myUser.UserId, &myUser.UserName, &myUser.Password, &myUser.Capacity, &myUser.Margin)
		//将user添加到users中
		users = append(users, myUser)
	}

	defer rows.Close()
	return users, err
}
