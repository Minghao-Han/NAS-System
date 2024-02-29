package uploadDA

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"nas/project/src/Entities"
)

func Update(uploadLog Entities.UploadLog) (rowsAffected int64, err error) {
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
	stmt, err := db.Prepare("UPDATE  upload_log SET uploader=?, path=? ,finished=?,finished=?,received_bytes=?,size=?,client_file_path=? WHERE id=?")
	if err != nil {
		return
	}
	//执行修改操作
	rs, err := stmt.Exec(uploadLog.Id, uploadLog.Uploader, uploadLog.Path, uploadLog.Finished, uploadLog.ReceivedBytes, uploadLog.Size, uploadLog.ClientFilePath)
	if err != nil {
		return
	}
	//返回插入的id
	rowsAffected, err = rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(stmt)
	return

}
