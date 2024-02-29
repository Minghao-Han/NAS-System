package uploadDA

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"nas/project/src/Entities"
)

// 定义用户信息结构体

func Insert(uploadLog Entities.UploadLog) (Id int, err error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	//推迟数据库连接的关闭
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(db)

	stmt, err := db.Prepare("INSERT INTO upload_log(uploader, path, finished, received_bytes, size,client_file_path) VALUES (?, ?, ?, ?,?, ?)")
	if err != nil {
		fmt.Printf("insert data failed: %v\n", err)
	}

	//执行插入操作
	rs, err := stmt.Exec(uploadLog.Uploader, uploadLog.Path, uploadLog.Finished, uploadLog.ReceivedBytes, uploadLog.Size, uploadLog.ClientFilePath)
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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(stmt)
	return

}
