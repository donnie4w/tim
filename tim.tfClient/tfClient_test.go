package tfClient

import (
	"fmt"
	"testing"

	//	"git.apache.org/thrift.git/lib/go/thrift"
	//	"github.com/donnie4w/go-logger/logger"
	//	. "tim.common"
	. "tim.protocol"
)

func TestRemote(t *testing.T) {
	tid := NewTid()

	tid.Name = "734604"
	pwd := "e10adc3949ba59abbe56e057f20f883e"
	HttpClient(func(client *ITimClient) {
		r, er := client.TimRemoteUserAuth(tid, pwd)
		if er == nil && r != nil {
			fmt.Println(r)
			if r.ExtraMap != nil {
				if password, ok := r.ExtraMap["password"]; ok {
					if pwd == password {
						fmt.Print("ok")
					}
				}
				if extraAuth, ok := r.ExtraMap["extraAuth"]; ok {
					if pwd == extraAuth {
						fmt.Print("ok2")
					}
				}
			}
		}
	})
}
