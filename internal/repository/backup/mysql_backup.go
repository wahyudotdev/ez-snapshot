package backup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

type MySqlBackup struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Output   string
}

func (m MySqlBackup) Dump(ctx context.Context, opts ...DumpDbOpts) (string, error) {
	// build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		m.User, m.Password, m.Host, m.Port, m.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()

	// pick filename if none
	output := m.Output
	if output == "" {
		output = fmt.Sprintf("%s_%s.sql", m.Database, time.Now().Format("20060102_150405"))
	}

	outfile, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer outfile.Close()

	// write header
	fmt.Fprintf(outfile, "-- MySQL Dump\n-- Database: %s\n-- Generated at: %s\n\n",
		m.Database, time.Now().Format(time.RFC3339))

	// get tables
	tables, err := getTables(ctx, db)
	if err != nil {
		return "", err
	}

	for _, tbl := range tables {
		// schema
		var createStmt string
		if err := db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE `%s`", tbl)).
			Scan(&tbl, &createStmt); err != nil {
			return "", err
		}
		fmt.Fprintf(outfile, "--\n-- Table structure for `%s`\n--\n", tbl)
		fmt.Fprintf(outfile, "DROP TABLE IF EXISTS `%s`;\n", tbl)
		fmt.Fprintln(outfile, createStmt+";")
		fmt.Fprintln(outfile)

		// dump rows
		if err := dumpTableData(ctx, db, outfile, tbl); err != nil {
			return "", err
		}
	}

	return output, nil
}

func getTables(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tbl string
		if err := rows.Scan(&tbl); err != nil {
			return nil, err
		}
		tables = append(tables, tbl)
	}
	return tables, nil
}

func dumpTableData(ctx context.Context, db *sql.DB, outfile *os.File, tbl string) error {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM `%s`", tbl))
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	fmt.Fprintf(outfile, "--\n-- Dumping data for table `%s`\n--\n", tbl)

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}

		values := make([]string, len(cols))
		for i, v := range vals {
			if v == nil {
				values[i] = "NULL"
			} else {
				values[i] = fmt.Sprintf("'%s'", escape(fmt.Sprint(v)))
			}
		}

		insertStmt := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s);",
			tbl, strings.Join(cols, "`, `"), strings.Join(values, ", "))
		fmt.Fprintln(outfile, insertStmt)
	}

	fmt.Fprintln(outfile)
	return nil
}

func escape(s string) string {
	// naive escaping, enough for basic text
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "''")
	return s
}
