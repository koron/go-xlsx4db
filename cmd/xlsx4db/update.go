package main

import (
	"database/sql"

	xlsx4db "github.com/koron/go-xlsx4db"
	"github.com/tealeg/xlsx"
)

func init() {
	validModes["update"] = update
}

func update(db *sql.DB, filename string, tables []string) error {
	xf, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	err = xlsx4db.Update(db, xf, tables...)
	if err != nil {
		return err
	}
	return nil
}
