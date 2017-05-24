# XLSX for Database

## Tips

### Tips in Japanese

*   値を `NULL` にするには

    セルの内容を `(NULL)` とし、背景を白以外の色で塗りつぶす。

## Tips in English

*   How to make a column `NULL`

    Make a cell content as `(NULL)` and fill background with the color except
    white.

## Misc

*   How does retrieve table names?
    *   PostgreSQL: `SELECT relname FROM pg_stat_user_tables`
    *   MySQL: `SHOW TABLES` or `SHOW TABLES FROM {db_name}`
