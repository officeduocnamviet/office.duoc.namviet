package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/jackc/pgx/v5"
)

func main() {
	dsn := "postgresql://postgres.elrvxcpbrudalnecnqap:namviet123a@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	// Query tables and columns
	query := `
		SELECT 
			c.table_name,
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default
		FROM 
			information_schema.columns c
		JOIN 
			information_schema.tables t ON c.table_name = t.table_name AND c.table_schema = t.table_schema
		WHERE 
			c.table_schema = 'public' AND t.table_type = 'BASE TABLE'
		ORDER BY 
			c.table_name, c.ordinal_position;
	`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	tables := make(map[string][]string)
	for rows.Next() {
		var tableName, columnName, dataType, isNullable string
		var columnDefault *string
		err := rows.Scan(&tableName, &columnName, &dataType, &isNullable, &columnDefault)
		if err != nil {
			log.Fatalf("Scan failed: %v\n", err)
		}

		def := ""
		if columnDefault != nil {
			def = " DEFAULT " + *columnDefault
		}

		line := fmt.Sprintf("- **%s**: %s (Nullable: %s)%s", columnName, dataType, isNullable, def)
		tables[tableName] = append(tables[tableName], line)
	}

	var tableNames []string
	for name := range tables {
		tableNames = append(tableNames, name)
	}
	sort.Strings(tableNames)

	file, err := os.Create("../database_schema.md")
	if err != nil {
		log.Fatalf("Create file failed: %v\n", err)
	}
	defer file.Close()

	file.WriteString("# Database Schema\n\n")
	for _, name := range tableNames {
		file.WriteString(fmt.Sprintf("### Table: %s\n", name))
		for _, col := range tables[name] {
			file.WriteString(col + "\n")
		}
		file.WriteString("\n")
	}

	fmt.Println("Schema successfully exported to ../database_schema.md")
}
