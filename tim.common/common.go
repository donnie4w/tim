/**
 * donnie4w@gmail.com  tim server
 */
package common

import (
	"tim.conf"
)

/*版本*/
var VersionName = "tim 1.0"
var VersionCode = 2
var Author = "wuxiaodong"
var Email = "donnie4w@gmail.com"

var CF = &conf.ConfBean{KV: make(map[string]string, 0), Db_Exsit: 1, MustAuth: 1}

var ClusterConf = &conf.ClusterBean{IsCluster: 1}
