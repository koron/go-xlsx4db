package main

import (
	"database/sql"

	"github.com/koron/go-xlsx4db"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

const dbname = `postgres://vagrant:db1234@127.0.0.1/vagrant?sslmode=disable`

func main() {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	xf := xlsx.NewFile()
	err = xlsx4db.Dump(xf, db, "users")
	if err != nil {
		panic(err)
	}
	err = xf.Save("tmp/out.xlsx")
	if err != nil {
		panic(err)
	}
}
