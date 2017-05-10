package xlsx4db

import (
	"database/sql"
	"reflect"

	"github.com/tealeg/xlsx"
)

// Dump dumps tables to XLSX file.
func Dump(xf *xlsx.File, db *sql.DB, tables ...string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		tables, err = FetchTables(db)
		if err != nil {
			return err
		}
	}
	for _, t := range tables {
		xs, err := xf.AddSheet(t)
		if err != nil {
			return err
		}
		err = dumpTable(xs, tx, t)
	}
	err = tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func dumpTable(xs *xlsx.Sheet, tx *sql.Tx, table string) error {
	rows, err := tx.Query("SELECT * FROM " + table)
	if err != nil {
		return err
	}
	defer rows.Close()
	ctypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	// TODO: column type specific operations.
	// convert column types and add as the header to sheet.
	var (
		h1   = xs.AddRow()
		vals = make([]interface{}, len(ctypes))
	)
	for i, ct := range ctypes {
		h1.AddCell().SetString(ct.Name())
		vals[i] = reflect.New(ct.ScanType()).Interface()
	}
	// convert values to xlsx'x cells
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			return err
		}
		xr := xs.AddRow()
		for _, v := range vals {
			c := xr.AddCell()
			// TODO: NULL value as special.
			c.SetValue(*(v.(*interface{})))
		}
	}
	return nil
}
