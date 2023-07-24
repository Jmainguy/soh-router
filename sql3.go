package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// TestItem : A test Item
type TestItem struct {
	URL  string
	Name string
}

// InitDB : Create database if it doesnt exist
func InitDB(filepath string) *sql.DB {
	file := fmt.Sprintf("file:%v?cache=shared&mode=rwc", filepath)
	// Initialize the database
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}

	// Check if the database connection is successful
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Println("Database initialized successfully")
	return db
}

// CreateTable : Create sqlite table
func CreateTable(db *sql.DB) {
	// create table if not exists
	sqlTable := `
    CREATE TABLE IF NOT EXISTS docker_pool(
        url string primary key,
        name string
    );
    `
	_, err := db.Exec(sqlTable)
	check(err)
}

// ReadItem : Read item from database
func ReadItem(db *sql.DB) (url string) {
	sqlReadall := `
    SELECT Url FROM docker_pool limit 1;
    `

	stmt, err := db.Prepare(sqlReadall)
	check(err)
	err = stmt.QueryRow().Scan(&url)
	check(err)
	stmt.Close()
	return url
}

// StoreItem : Store item in database
func StoreItem(db *sql.DB, items []TestItem) {
	sqlAdditem := `
    INSERT OR REPLACE INTO docker_pool(
        Url,
        Name
    ) values(?, ?)
    `

	routeSQL, err := db.Prepare(sqlAdditem)
	check(err)

	for _, item := range items {
		tx, err := db.Begin()
		check(err)
		_, err = tx.Stmt(routeSQL).Exec(item.URL, item.Name)
		if err != nil {
			log.Println(err)
			log.Println("doing rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
			check(err)
		}
	}
}

// DelItem : Delete item in database
func DelItem(db *sql.DB, url string) {
	sqlDelitem := `
    DELETE FROM docker_pool where Url = ?
    `
	stmt, err := db.Prepare(sqlDelitem)
	check(err)
	check(err)
	_, err = stmt.Exec(url)
	check(err)
	stmt.Close()
}

// DelName : delete item with a name from database
func DelName(db *sql.DB, name string) {
	routeSQL, err := db.Prepare("delete from docker_pool where name = ?;")
	check(err)
	tx, err := db.Begin()
	check(err)
	_, err = tx.Stmt(routeSQL).Exec(name)
	if err != nil {
		log.Println(err)
		log.Println("doing rollback")
		tx.Rollback()
	} else {
		err = tx.Commit()
		check(err)
	}
}
