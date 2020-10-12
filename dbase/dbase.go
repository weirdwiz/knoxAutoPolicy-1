package dbase

import (
	"database/sql"
	"os"

	//imported
	_ "github.com/go-sql-driver/mysql"
)

//Conn connects to
func Conn() (db *sql.DB) {
	dbDriver := os.Getenv("NETWORKFLOW_DB_DRIVER")
	dbUser := os.Getenv("NETWORKFLOW_DB_USER")
	dbPass := os.Getenv("NETWORKFLOW_DB_PASS")
	dbName := os.Getenv("NETWORKFLOW_DB_NAME")

	// for test
	dbDriver = "mysql"
	dbUser = "root"
	dbPass = "password"
	dbName = "flow_management"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(localhost:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}