package main

import (
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqldb, err := config()
	check(err)
	db := InitDB(sqldb)
	CreateTable(db)
	listener, err := net.Listen("tcp", "0.0.0.0:8085")
	go keep10InPool(db)
	check(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection %v\n", conn)
		go forward(db, conn)
	}
}
