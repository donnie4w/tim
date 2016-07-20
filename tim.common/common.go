/**
 * donnie4w@gmail.com  tim server
 */
package common

import (
	"tim.conf"
)

/*版本*/
var _Version = "tim 1.0"
var _Author = "wuxiaodong"
var _Email = "donnie4w@gmail.com;donnie4wu@qq.com"

var ConfBean = &conf.ConfBean{KV: make(map[string]string, 0), Db_Exsit: 1}
