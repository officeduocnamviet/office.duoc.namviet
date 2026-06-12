package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := "postgresql://postgres.iudkexocalqdhxuyjacu:Longlong123%40a@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?pgbouncer=true"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT table_name, column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_schema = 'public' 
ORDER BY table_name, ordinal_position;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	f, _ := os.Create("../database_schema.md")
	defer f.Close()

	f.WriteString("# Database Schema\n")
	var currentTable string
	for rows.Next() {
		var tableName, colName, dataType, isNullable string
		var colDef sql.NullString
		rows.Scan(&tableName, &colName, &dataType, &isNullable, &colDef)

		if tableName != currentTable {
			f.WriteString(fmt.Sprintf("\n### Table: %s\n", tableName))
			currentTable = tableName
		}
		defStr := ""
		if colDef.Valid {
			defStr = fmt.Sprintf(" DEFAULT %s", colDef.String)
		}
		f.WriteString(fmt.Sprintf("- **%s**: %s (Nullable: %s)%s\n", colName, dataType, isNullable, defStr))
	}
}
