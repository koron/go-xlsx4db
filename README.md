# XLSX for Database

Dump and restore RDBMS by using Excel (XLSX)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/go-xlsx4db)](https://pkg.go.dev/github.com/koron/go-xlsx4db)
[![Actions/Go](https://github.com/koron/go-xlsx4db/workflows/Go/badge.svg)](https://github.com/koron/go-xlsx4db/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/go-xlsx4db)](https://goreportcard.com/report/github.com/koron/go-xlsx4db)

## Usage as a command

```console
$ go install github.com/koron/go-xlsx4db/cmd/xlsx4db@latest
```

### Options

```console
$ ./xlsx4db -h
Usage of xlsx4db:
  -dbname string
        DB source string, example:
          * MySQL: "{user}:{pass}@{addr}/{name}"
          * PostgreSQL: "postgres://{user}:{pass}@{addr}/{name}?sslmode=disable"
  -driver string
        DB driver: "mysql" or "postgres"
  -mode string
        Mode: "dump" or "restore"
  -tables string
        OPTION: table names to dump/restore
  -xlsx string
        Excel file name to operate
```

## Usage as a package

```console
$ go get github.com/koron/go-xlsx4db@latest
```

## Tips

### Tips in Japanese

*   値を `NULL` にするには

    セルの内容を `(NULL)` とし、背景を白以外の色で塗りつぶす。

### Tips in English

*   How to make a column `NULL`

    Make a cell content as `(NULL)` and fill background with the color except
    white.

## Misc

*   How does retrieve table names?
    *   PostgreSQL: `SELECT relname FROM pg_stat_user_tables`
    *   MySQL: `SHOW TABLES` or `SHOW TABLES FROM {db_name}`

*   What dbname does MySQL accept?
    *   `vagrant:db1234@tcp(127.0.0.1:3306)/vagrant`
