package main

import (
	"database/sql"
	"log"

	"github.com/koron/go-xlsx4db"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

const dbname = `postgres://vagrant:db1234@127.0.0.1/vagrant?sslmode=disable`

func main() {
	xf, err := xlsx.OpenFile("tmp/in.xlsx")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = xlsx4db.Restore(db, xf, true, "users")
	if err != nil {
		log.Fatal(err)
	}
}
