/**
 * donnie4w@gmail.com  tim server
 */
package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"

	"github.com/donnie4w/dom4g"
	"github.com/donnie4w/go-logger/logger"
)

/**配置结构对象*/

type ClusterBean struct {
	RedisAddr   string //  redis ip:port
	RedisPwd    string //
	RedisDB     int    //
	RequestAddr string // 访问地址
	RequestType string // 访问类型
	Domain      string // 域名
	Username    string // 登陆名
	Password    string // 登陆密码
	IsCluster   int    // 1集群 0不集群
	Interflow   int    // 合流信息发送 0不合流  1合流
	Keytimeout  int    // key 过期时间
}

func (cb *ClusterBean) Init(filexml string) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Init error", err)
			logger.Error(string(debug.Stack()))
			b = false
		}
	}()
	if !isExist(filexml) {
		b = false
		return
	}
	xmlconfig, err := os.Open(filexml)
	if err != nil {
		panic(fmt.Sprint("xmlconfig is error:", err.Error()))
		os.Exit(0)
	}
	config, err := ioutil.ReadAll(xmlconfig)
	if err != nil {
		panic(fmt.Sprint("config is error:", err.Error()))
		os.Exit(1)
	}
	dom, err := dom4g.LoadByXml(string(config))
	if err == nil {
		nodes := dom.AllNodes()
		if nodes != nil {
			fmt.Println(`======================cluster start======================`)
			i := 0
			for _, node := range nodes {
				name := node.Name()
				value := node.Value
				v := reflect.ValueOf(cb).Elem().FieldByName(name)
				if v.CanSet() {
					fmt.Println("set====>", name, value)
					switch v.Type().Name() {
					case "string":
						v.Set(reflect.ValueOf(value))
					case "int":
						i, _ := strconv.Atoi(value)
						v.Set(reflect.ValueOf(i))
					default:
						fmt.Println("other type:", v.Type().Name(), ">>>", name)
					}
				} else {
					fmt.Println("no set====>", name, value)
					i++
				}
			}
			fmt.Println(`=======================cluster end=======================`)
			if i > 0 {
				fmt.Println("not set number:", i)
			}
		}
	}
	b = true
	return
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
