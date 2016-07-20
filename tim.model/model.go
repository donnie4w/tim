/**
 * donnie4w@gmail.com  tim server
 */
package model

import (
	"github.com/donnie4w/gdao"
	"tim.DB"
)

func init() {
	gdao.SetDB(DB.Master)
	gdao.SetAdapterType(gdao.MYSQL)
}
