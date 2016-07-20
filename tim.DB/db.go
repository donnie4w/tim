/**
 * donnie4w@gmail.com  tim server
 */
package DB

import (
	"database/sql"
	"os"

	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	. "tim.common"
)

var Master *sql.DB

func Init() {
	initmaster()
}

func initmaster() {
	if Master == nil {
		logger.Info("master init")
		dataSourceName, maxOpenConns, maxIdleConns := ConfBean.GetDB()
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			logger.Info("any error on open database ", err.Error())
			os.Exit(1)
			return
		}
		db.SetMaxOpenConns(maxOpenConns)
		db.SetMaxIdleConns(maxIdleConns)
		Master = db
	}
}
