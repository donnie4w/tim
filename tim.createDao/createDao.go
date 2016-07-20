/**
 * donnie4w@gmail.com  tim server
 */
package main

import (
	"database/sql"
	"fmt"

	"github.com/donnie4w/gdao"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/tim")
	if err != nil {
		fmt.Println("any error on open database ", err.Error())
		return
	}
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	gdao.SetDB(db)
	gdao.SetAdapterType(gdao.MYSQL)
}

func main() {
	gbs, _ := gdao.ExecuteQuery("show tables")
	for _, g := range gbs {
		for _, v := range g.FieldBeens {
			fmt.Println(v.Name(), v.Index(), v.Value())
			g := v.Value()
			err := gdao.CreateDaoFile(g.(string), "dao", "D:/liteIDEspace/tim/src/tim.dao")
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
