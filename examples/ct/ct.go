package main

import (
	"database/sql"
	"fmt"

	"github.com/koron/go-xlsx4db"
	_ "github.com/lib/pq"
)

const dbname = `postgres://vagrant:db1234@127.0.0.1/vagrant?sslmode=disable`

func main() {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	tables, err := xlsx4db.FetchTables(db)
	if err != nil {
		panic(err)
	}
	for _, t := range tables {
		rows, err := db.Query("SELECT * FROM " + t)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		ctypes, err := rows.ColumnTypes()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%q:\n", t)
		for _, ct := range ctypes {
			fmt.Printf("  %q:\n", ct.Name())
			if dt := ct.DatabaseTypeName(); dt != "" {
				fmt.Printf("    DatabaseTypeName: %q\n", dt)
			}
			if n, ok := ct.Nullable(); ok {
				fmt.Printf("    Nullable: %t\n", n)
			}
			if l, ok := ct.Length(); ok {
				fmt.Printf("	Length: %d\n", l)
			}
			if precision, scale, ok := ct.DecimalSize(); ok {
				fmt.Printf("    DecimalSize:\n")
				fmt.Printf("      precision: %d\n", precision)
				fmt.Printf("      scale:     %d\n", scale)
			}
			if st := ct.ScanType().Name(); st != "" {
				fmt.Printf("    ScanType: %q\n", st)
			}
		}
	}
}
