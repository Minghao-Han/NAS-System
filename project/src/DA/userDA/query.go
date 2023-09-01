package userDA

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"nas/project/src/Entities"
	"nas/project/src/Utils"
)

func Query() (users []Entities.User, err error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

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
		err := rows.Scan(&myUser.UserId, &myUser.UserName, &myUser.Password, &myUser.Capacity, &myUser.Margin)
		if err != nil {
			return nil, err
		}
		//将user添加到users中
		users = append(users, myUser)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	return users, err
}

func FindById(userId int) (*Entities.User, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)
	rows, err := db.Query("SELECT user_id,user_name,total_capacity,margin FROM user where user_id = ?", userId)
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user Entities.User
		//遍历表中所有行的信息
		err := rows.Scan(&user.UserId, &user.UserName, &user.Capacity, &user.Margin)
		if err != nil {
			return nil, err
		}
		user.Password = ""
		return &user, nil
	}
	return nil, err
}

func FindByUsername(username string) (*Entities.User, error) {
	if safe := Utils.SQLInjectionDetector(username); !safe {
		return nil, fmt.Errorf("detected sql injection attack")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)
	rows, err := db.Query("SELECT user_id,user_name,total_capacity,margin FROM user where user_name = ?", username)
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user Entities.User
		//遍历表中所有行的信息
		err := rows.Scan(&user.UserId, &user.UserName, &user.Capacity, &user.Margin)
		if err != nil {
			return nil, err
		}
		user.Password = ""
		return &user, nil
	}
	return nil, err
}

func GetPassword(username string) ([]byte, error) {
	if safe := Utils.SQLInjectionDetector(username); !safe {
		return nil, fmt.Errorf("detected sql injection attack")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)
	rows, err := db.Query("SELECT password FROM user where user_name = ?", username)
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		pwStr := new(string)
		//遍历表中所有行的信息
		err := rows.Scan(pwStr)
		if err != nil {
			return nil, err
		}
		pwBytes := []byte(*pwStr)
		return pwBytes, nil
	}
	return nil, err
}
