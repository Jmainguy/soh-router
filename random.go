package main

import (
	"database/sql"
	"io"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func forward(db *sql.DB, conn net.Conn) {
	var client net.Conn
	var err error
	var attempts int
	for {
		attempts++
		// Pull a container from the pool, and connect to it
		target := pullDockerFromPool(db)
		client, err = net.Dial("tcp", target)
		if err != nil {
			// We were unable to connect to the container, try a different one
			log.Printf("Attempt #%d for target has failed\n", attempts)
			log.Println(err)
		} else {
			// End for loop, we got a good target
			break
		}
	}
	log.Printf("Connected to localhost %v\n", conn)
	// Add another host to pool
	go dockerStuff(db)
	go func() {
		defer func() {
			if err := client.Close(); err != nil {
				log.Println("Error closing client:", err)
			}
		}()

		defer func() {
			if err := conn.Close(); err != nil {
				log.Println("Error closing conn:", err)
			}
		}()

		_, err := io.Copy(client, conn)
		check(err)
	}()
	go func() {
		defer func() {
			if err := client.Close(); err != nil {
				log.Println("Error closing client:", err)
			}
		}()

		defer func() {
			if err := conn.Close(); err != nil {
				log.Println("Error closing conn:", err)
			}
		}()

		_, err := io.Copy(conn, client)
		check(err)
	}()
}
