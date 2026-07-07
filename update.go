package xlsx4db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tealeg/xlsx"
)

// Update updates (UPDATE or INSERT) tables from XLSX file.
func Update(db *sql.DB, xf *xlsx.File, tables ...string) error {
	return UpdateContext(context.Background(), db, xf, tables...)
}

// UpdateContext updates (UPDATE or INSERT) tables from XLSX file with context.Context.
func UpdateContext(ctx context.Context, db *sql.DB, xf *xlsx.File, tables ...string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	sheets := xf.Sheets
	if len(tables) > 0 {
		sheets = make([]*xlsx.Sheet, 0, len(tables))
		for _, t := range tables {
			if xs, ok := xf.Sheet[t]; ok {
				sheets = append(sheets, xs)
			}
		}
	}
	for _, xs := range sheets {
		err := updateTable(ctx, db, tx, xs, xs.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func updateTable(ctx context.Context, db *sql.DB, tx *sql.Tx, xs *xlsx.Sheet, table string) error {
	cols := xs.Rows[0].Cells
	columns := make([]string, len(cols))
	for i, xc := range cols {
		columns[i] = xc.Value
	}
	q, err := BuildUpsertQuery(db, table, columns)
	if err != nil {
		return err
	}
	st, err := tx.PrepareContext(ctx, q)
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
		_, err := st.ExecContext(ctx, args2...)
		if err != nil {
			return err
		}
	}
	return nil
}
