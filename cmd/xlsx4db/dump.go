package main

import (
	"database/sql"

	xlsx4db "github.com/koron/go-xlsx4db"
	"github.com/tealeg/xlsx"
)

func init() {
	validModes["dump"] = dump
}

func dump(db *sql.DB, filename string, tables []string) error {
	xf := xlsx.NewFile()
	err := xlsx4db.Dump(xf, db, tables...)
	if err != nil {
		return err
	}
	err = xf.Save(filename)
	if err != nil {
		return err
	}
	return nil
}
