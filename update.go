package xlsx4db

import (
	"database/sql"
	"fmt"

	"github.com/tealeg/xlsx"
)

// Update updates (UPDATE or INSERT) tables from XLSX file.
func Update(db *sql.DB, xf *xlsx.File, tables ...string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	sheets := xf.Sheets
	if len(tables) > 0 {
		sheets = make([]*xlsx.Sheet, 0, len(tables))
		for _, t := range tables {
			if xs, ok := xf.Sheet[t]; ok {
				sheets = append(sheets, xs)
			}
		}
		//tables, err = FetchTables(db)
		//if err != nil {
		//	return err
		//}
	}
	for _, xs := range sheets {
		err := updateTable(db, tx, xs, xs.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateTable(db *sql.DB, tx *sql.Tx, xs *xlsx.Sheet, table string) error {
	cols := xs.Rows[0].Cells
	columns := make([]string, len(cols))
	for i, xc := range cols {
		columns[i] = xc.Value
	}
	q, err := BuildUpsertQuery(db, table, columns)
	if err != nil {
		return err
	}
	st, err := tx.Prepare(q)
	if err != nil {
		return fmt.Errorf("prepare(%q) failed: %s", q, err.Error())
	}
	args := make([]interface{}, len(cols))
	for _, xr := range xs.Rows[1:] {
		for i := range args {
			args[i], err = cellToValue(xr.Cells[i])
			if err != nil {
				return err
			}
		}
		args2 := append(args, args...)
		_, err := st.Exec(args2...)
		if err != nil {
			return err
		}
	}
	return nil
}
