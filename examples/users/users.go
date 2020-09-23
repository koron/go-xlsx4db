package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

const dbname = `postgres://vagrant:db1234@127.0.0.1/vagrant?sslmode=disable`

func main() {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	ctypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}
	vals := make([]interface{}, len(ctypes))
	for i, ct := range ctypes {
		vals[i] = reflect.New(ct.ScanType()).Interface()
	}

	for rows.Next() {
		err := rows.Scan(vals...)
		if err != nil {
			panic(err)
		}
		fmt.Println("-")
		for i, ct := range ctypes {
			v := *(vals[i].(*interface{}))
			if v == nil {
				fmt.Printf("  %q: (NULL)\n", ct.Name())
				continue
			}
			w := reflect.ValueOf(v)
			t := w.Type()
			fmt.Printf("  %q:\n", ct.Name())
			fmt.Printf("    value: %#v\n", v)
			fmt.Printf("    type: %q\n", t)
		}
	}
}
