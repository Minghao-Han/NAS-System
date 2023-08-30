package utils

import (
	"database/sql"
	"fmt"
)

func Query() {
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

	defer rows.Close()

	for rows.Next() {
		var user_id int
		var user_name string
		var password string
		var capacity int
		var margin int

		err := rows.Scan(&user_id, &user_name, &password, &capacity, &margin)
		if err != nil {
			fmt.Printf("get data failed: %v\n", err)
			return
		}

		fmt.Printf("%d\t%s\t%s\t%d\t%d\n", user_id, user_name, password, capacity, margin)
	}

	if err != nil {
		fmt.Printf("create table failed: %v\n", err)
	}
	fmt.Println("database connected")
}
