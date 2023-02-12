/**
 * donnie4w@gmail.com  tim server
 */
package model

import (
	"tim/DB"

	"github.com/donnie4w/gdao"
)

func init() {
	gdao.SetDB(DB.Master)
	gdao.SetAdapterType(gdao.MYSQL)
}
