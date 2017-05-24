package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var validDrivers = map[string]struct{}{}

func init() {
	for _, n := range sql.Drivers() {
		validDrivers[n] = struct{}{}
	}
}
