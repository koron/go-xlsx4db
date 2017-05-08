package xlsx4db

import (
	"database/sql"

	"github.com/tealeg/xlsx"
)

// Restore restores tables from XLSX file.
func Restore(db *sql.DB, xf *xlsx.File, refresh bool, tables ...string) error {
	return nil
}
