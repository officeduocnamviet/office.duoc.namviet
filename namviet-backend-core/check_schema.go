package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgresql://postgres.nndohdttxohgoxuudwta:namviet-admin-super-key@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var columns []string
	db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'inventory_transactions'").Scan(&columns)
	fmt.Println("Columns in inventory_transactions:", columns)
}
