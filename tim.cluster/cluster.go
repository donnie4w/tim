package cluster

import (
	"github.com/donnie4w/go-logger/logger"
	. "tim.conf"
)

var Cluster *ClusterBean = new(ClusterBean)
var flag bool = false

func InitCluster(filexml string) {
	if Cluster.Init(filexml) {
		Redis.initPool()
		flag = Redis.Ping()
		logger.Info("iscluster:", flag)
	}
}

func IsCluster() bool {
	return false
}

func Add2CluserServer(loginname string) (err error) {
	return
}

func MustSendToRemote(loginname string) (addrmap map[string]string) {
	return
}
