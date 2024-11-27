package main

import (
	"database/sql"
	"log"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

var conn *sql.DB

func OracleDBConf(Body body) *sql.DB {
	port := 1521
	connStr := go_ora.BuildUrl(Body.Host, port, Body.Service, Body.Username, Body.Password, nil)

	conn, err := sql.Open("oracle", connStr)

	conn.SetMaxOpenConns(10) // Max simultaneous connections
	conn.SetMaxIdleConns(5)  // Max idle connections
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	return conn

}

func CloseDB() {
	log.Println("Connection closed")
	if conn != nil {
		conn.Close()
	}
}
