package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net"
)

func check(e error) {
	if e != nil {
		log.Println(e)
		//panic(e)
	}
}

func forward(db *sql.DB, conn net.Conn) {
	target := pullDockerFromPool(db)
	client, err := net.Dial("tcp", target)
	if err != nil {
		check(err)
	}
	log.Printf("Connected to localhost %v\n", conn)
	// Add another host to pool
	go dockerStuff(db)
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}
