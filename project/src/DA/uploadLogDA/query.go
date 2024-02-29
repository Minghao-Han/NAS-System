package uploadDA

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"nas/project/src/Entities"
	"nas/project/src/Utils"
)

//func Query() (uploadlogs []Entities.UploadLog, err error) {
//	db, err := sql.Open("sqlite3", dbPath)
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//	defer func(db *sql.DB) {
//		err := db.Close()
//		if err != nil {
//			return
//		}
//	}(db)
//	rows, err := db.Query("SELECT * FROM upload_log")
//	if err != nil {
//		fmt.Printf("query data failed: %v\n", err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var uploadLog Entities.UploadLog
//		//只取查询的第一个
//		err := rows.Scan(&uploadLog.Id, &uploadLog.Uploader, &uploadLog.Path, &uploadLog.Finished, &uploadLog.ReceivedBytes, &uploadLog.Size)
//		if err != nil {
//			return nil, err
//		}
//		uploadlogs = append(uploadlogs, uploadLog)
//	}
//	return uploadlogs, nil
//}

func FindById(id int) (*Entities.UploadLog, error) {
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
	rows, err := db.Query("SELECT * FROM upload_log where id = ?", id)
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(rows)
	if rows.Next() {
		var uploadLog Entities.UploadLog
		//只取查询的第一个
		err := rows.Scan(&uploadLog.Id, &uploadLog.Uploader, &uploadLog.Path, &uploadLog.Finished, &uploadLog.ReceivedBytes, &uploadLog.Size, &uploadLog.ClientFilePath)
		if err != nil {
			return nil, err
		}
		return &uploadLog, nil
	}
	return nil, err
}

func FindByUploader(uploader string) (uploadlogs []Entities.UploadLog, err error) {
	if safe := Utils.SQLInjectionDetector(uploader); !safe {
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
	rows, err := db.Query("SELECT * FROM upload_log where uploader = ?", uploader)
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(rows)

	for rows.Next() {
		var uploadLog Entities.UploadLog
		//只取查询的第一个
		err := rows.Scan(&uploadLog.Id, &uploadLog.Uploader, &uploadLog.Path, &uploadLog.Finished, &uploadLog.ReceivedBytes, &uploadLog.Size)
		if err != nil {
			return nil, err
		}
		uploadlogs = append(uploadlogs, uploadLog)
	}
	return uploadlogs, nil
}

func FindUnfinishedByUploader(uploader string) (uploadlogs []Entities.UploadLog, err error) {
	if safe := Utils.SQLInjectionDetector(uploader); !safe {
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
	rows, err := db.Query("SELECT * FROM upload_log where uploader = ? and finished=0", uploader)
	if err != nil {
		fmt.Printf("query data failed: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(rows)

	for rows.Next() {
		var uploadLog Entities.UploadLog
		//只取查询的第一个
		err := rows.Scan(&uploadLog.Id, &uploadLog.Uploader, &uploadLog.Path, &uploadLog.Finished, &uploadLog.ReceivedBytes, &uploadLog.Size)
		if err != nil {
			return nil, err
		}
		uploadlogs = append(uploadlogs, uploadLog)
	}
	return uploadlogs, nil
}
