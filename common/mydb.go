package common

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type MyDB struct {
	Name string
	DB   *sql.DB
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (mydb *MyDB) Open() error {
	if mydb.DB != nil {
		return fmt.Errorf("Database already open")
	}
	if mydb.Name == "" {
		return fmt.Errorf("Database path not set")
	}

	if !fileExists(mydb.Name) {
		return fmt.Errorf("File does not exist")
	}

	db, err := sql.Open("sqlite3", mydb.Name)
	if err != nil {
		return err
	}

	mydb.DB = db
	return nil
}

func (mydb *MyDB) Close() error {
	err := mydb.DB.Close()
	if err != nil {
		return err
	}
	mydb.DB = nil
	return nil
}

func (mydb *MyDB) HasTable(tableName string) (bool, error) {
	if mydb.DB == nil {
		return false, fmt.Errorf("Database not opened")
	}

	db := mydb.DB

	rows, err := db.Query("SELECT COUNT(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", tableName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	ok := rows.Next()
	if !ok {
		return false, fmt.Errorf("Internal error")
	}
	var cnt int
	err = rows.Scan(&cnt)
	if err != nil {
		return false, err
	}
	ok = rows.Next()
	if ok {
		return false, fmt.Errorf("Internal error")
	}

	err = rows.Err()
	if err != nil {
		return false, err
	}

	if cnt < 0 || cnt > 1 {
		return false, fmt.Errorf("Internal error")
	}

	return cnt == 1, nil
}

func (mydb *MyDB) DropTable(tableName string) error {
	if mydb.DB == nil {
		return fmt.Errorf("Database not opened")
	}

	_, err := mydb.DB.Exec(`DROP TABLE "` + tableName + `"`) // TODO: why not use "?" here?
	if err != nil {
		return err
	}

	return nil
}
