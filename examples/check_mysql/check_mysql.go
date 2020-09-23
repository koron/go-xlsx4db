package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

const dbname = `vagrant:db1234@/vagrant`

func isMySQL(db *sql.DB) bool {
	return reflect.ValueOf(db.Driver()).Type().String() == "*mysql.MySQLDriver"
}

func isPostgreSQL(db *sql.DB) bool {
	return reflect.ValueOf(db.Driver()).Type().String() == "*pq.drv"
}

func run() error {
	db, err := sql.Open("mysql", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf("isPostgreSQL()=%t\n", isPostgreSQL(db))
	fmt.Printf("isMySQL()=%t\n", isMySQL(db))
	err = db.Ping()
	if err != nil {
		log.Print("HERE:0")
		return err
	}
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Print("HERE:1")
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Print("HERE:2")
			return err
		}
		fmt.Printf("  %s\n", name)
	}
	err = rows.Err()
	if err != nil {
		log.Print("HERE:3")
		return err
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
}
