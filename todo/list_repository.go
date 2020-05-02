package todo

import (
	"database/sql"
	"errors"
	"fmt"
)

// ListRepository is an interface
type ListRepository interface {
	SaveItem(ID string, name string) (err error)
	FindItem(ID string) (res string, err error)
	DeleteItem(ID string) (err error)
	Init() (err error)
}

// MySQLListRepository impl of Listrepo
type MySQLListRepository struct {
	// conn sql.Conn
	db *sql.DB
}

// NewListRepository is a constructor
func NewListRepository(db *sql.DB) (repo ListRepository) {
	repo = MySQLListRepository{db: db}
	return repo
}

func (repo MySQLListRepository) SaveItem(ID string, name string) (err error) {
	_, err = repo.db.Exec("Insert into `items` (`ID`, `name`) values (?, ?)", ID, name)
	return
}

func (repo MySQLListRepository) FindItem(ID string) (res string, err error) {
	var row *sql.Row
	row = repo.db.QueryRow("SELECT `name` FROM `items` WHERE `ID` = ?; ", ID)
	err = row.Scan(&res)
	return
}

func (repo MySQLListRepository) DeleteItem(ID string) (err error) {
	var res sql.Result
	var affected int64
	res, err = repo.db.Exec("DELETE FROM `items` WHERE `ID` = ?; ", ID)
	affected, err = res.RowsAffected()
	if err != nil {
		return
	}

	if affected == 0 {
		return errors.New("invalid ID")
	}

	return
}

func (repo MySQLListRepository) Init() (err error) {
	fmt.Println("hey")
	// var result sql.Result
	_, err = repo.db.Exec(`CREATE TABLE IF NOT EXISTS items (ID VARCHAR(255), name VARCHAR(255)) ENGINE=INNODB;`)
	return err
}
