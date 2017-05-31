package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"strings"
)

type modeMain func(db *sql.DB, filename string, tables []string) error

var validModes = map[string]modeMain{}

func main() {
	err := main2()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func main2() error {
	var (
		driver   string
		dbname   string
		mode     string
		xlsxfile string
		tables   string
	)

	// initiate flags.
	flag.StringVar(&driver, "driver", "", `DB driver: "mysql" or "postgres"`)
	flag.StringVar(&dbname, "dbname", "", `DB source string, example:
	  * MySQL: "{user}:{pass}@{addr}/{name}"
	  * PostgreSQL: "postgres://{user}:{pass}@{addr}/{name}?sslmode=disable"`)
	flag.StringVar(&mode, "mode", "", `Mode: "dump", "restore" or "update"`)
	flag.StringVar(&xlsxfile, "xlsx", "", `Excel file name to operate`)
	flag.StringVar(&tables, "tables", "", `OPTION: table names to operate`)
	flag.Parse()

	// check value of flags.
	if _, ok := validDrivers[driver]; !ok {
		return fmt.Errorf("invalid driver (-driver): %q", driver)
	}
	if dbname == "" {
		return errors.New("empty `-dbname`")
	}
	proc, ok := validModes[mode]
	if !ok {
		return fmt.Errorf("invalid mode (-mode): %q", mode)
	}
	if xlsxfile == "" {
		return errors.New("empty `-xlsx`")
	}
	tableList := parseTables(tables)

	db, err := sql.Open(driver, dbname)
	if err != nil {
		return err
	}
	defer db.Close()

	// dispatch to main of each modes.
	err = proc(db, xlsxfile, tableList)
	if err != nil {
		return err
	}
	return nil
}

func parseTables(s string) []string {
	tables := strings.Split(s, ",")
	if len(tables) > 0 && tables[0] == "" {
		return tables[1:]
	}
	return tables
}
