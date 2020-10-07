package xlsx4db

import (
	"database/sql"

	"github.com/tealeg/xlsx"
)

var (
	headerStyle = xlsx.NewStyle()
	nullStyle   = xlsx.NewStyle()
)

func init() {
	headerStyle.Fill = *xlsx.NewFill("solid", "000000", "")
	headerStyle.Font.Color = "FFFFFFFF"
	headerStyle.ApplyFill = true
	headerStyle.ApplyFont = true

	nullStyle.Fill = *xlsx.NewFill("solid", "cccccc", "")
	nullStyle.ApplyFill = true
}

// Dump dumps tables to XLSX file.
func Dump(xf *xlsx.File, db *sql.DB, tables ...string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if len(tables) == 0 {
		tables, err = FetchTables(db)
		if err != nil {
			return err
		}
	}
	var quote string
	if isPostgreSQL(db) {
		quote = `"`
	}
	for _, t := range tables {
		xs, err := xf.AddSheet(t)
		if err != nil {
			return err
		}
		err = dumpTable(xs, tx, t, quote)
		if err != nil {
			return err
		}
	}
	return nil
}

func dumpTable(xs *xlsx.Sheet, tx *sql.Tx, table string, quote string) error {
	if quote != "" {
		table = quote + table + quote
	}
	rows, err := tx.Query("SELECT * FROM " + table)
	if err != nil {
		return err
	}
	defer rows.Close()
	ctypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	// convert column types and add as the header to sheet.
	var (
		h1   = xs.AddRow()
		vals = make([]interface{}, len(ctypes))
	)
	for i, ct := range ctypes {
		c := h1.AddCell()
		c.SetString(ct.Name())
		c.SetStyle(headerStyle)
		vals[i] = new(interface{})
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
			// NULL value's bgcolor should not be default (white).
			w := v.(*interface{})
			if *w == nil {
				c.SetString(nullLabel)
				c.SetStyle(nullStyle)
				continue
			}
			c.SetValue(*w)
		}
	}
	return nil
}
