package xlsx4db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tealeg/xlsx"
)

const (
	nullLabel  = "(NULL)"
	emptyColor = "FFFFFFFF"
)

// Restore restores tables from XLSX file.
func Restore(db *sql.DB, xf *xlsx.File, refresh bool, tables ...string) error {
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
		tables, err = FetchTables(db)
		if err != nil {
			return err
		}
	}
	for _, xs := range sheets {
		err := restoreTable(db, tx, xs, xs.Name, refresh)
		if err != nil {
			return err
		}
	}
	return nil
}

func restoreTable(db *sql.DB, tx *sql.Tx, xs *xlsx.Sheet, table string, refresh bool) error {
	if refresh {
		_, err := tx.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}
	cols := xs.Rows[0].Cells
	columns := make([]string, len(cols))
	for i, xc := range cols {
		columns[i] = xc.Value
	}
	q, err := BuildInsertQuery(db, table, columns)
	if err != nil {
		return err
	}
	st, err := tx.Prepare(q)
	if err != nil {
		return fmt.Errorf("Prepare(%q) failed: %s", q, err.Error())
	}
	args := make([]interface{}, len(cols))
	for _, xr := range xs.Rows[1:] {
		for i, _ := range args {
			args[i], err = cellToValue(xr.Cells[i])
			if err != nil {
				return err
			}
		}
		_, err := st.Exec(args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func isCellNull(xc *xlsx.Cell) bool {
	if xc.Value != nullLabel {
		return false
	}
	xs := xc.GetStyle()
	fg := ""
	if xs.ApplyFill {
		fg = xs.Fill.FgColor
	}
	return fg != "" && fg != emptyColor
}

func cellToValue(xc *xlsx.Cell) (interface{}, error) {
	switch xc.Type() {
	case xlsx.CellTypeString:
		if isCellNull(xc) {
			return nil, nil
		}
		return xc.String(), nil
	case xlsx.CellTypeStringFormula:
		return xc.Value, nil
	case xlsx.CellTypeNumeric:
		if xc.NumFmt != "general" {
			return xc.GetTime(false)
		}
		return xc.Value, nil
	case xlsx.CellTypeBool:
		return xc.Bool(), nil
	case xlsx.CellTypeInline:
		return nil, errors.New("not support: inline")
	case xlsx.CellTypeError:
		return nil, fmt.Errorf("cell error: %s", xc.Value)
	case xlsx.CellTypeDate:
		return xc.GetTime(false)
	}
	return nil, errors.New("unknown cell type")
}
