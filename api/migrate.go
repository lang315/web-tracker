package main

import (
	"github.com/go-pg/pg/orm"
	"log"
)

func migrate()  {
	createTableIfNotExists(&User{})
}

func createTableIfNotExists(tbl interface{}) {
	if err := db.CreateTable(tbl, &orm.CreateTableOptions{IfNotExists: true}); err != nil {
		log.Fatalln(err)
	}
}
