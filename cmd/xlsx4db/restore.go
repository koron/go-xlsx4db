package main

import (
	"database/sql"

	xlsx4db "github.com/koron/go-xlsx4db"
	"github.com/tealeg/xlsx"
)

func init() {
	validModes["restore"] = restore
}

func restore(db *sql.DB, filename string, tables []string) error {
	xf, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	err = xlsx4db.Restore(db, xf, true, tables...)
	if err != nil {
		return err
	}
	return nil
}
