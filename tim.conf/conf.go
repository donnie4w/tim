/**
 * donnie4w@gmail.com  tim server
 */
package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/donnie4w/dom4g"
)

/**配置结构对象*/

type ConfBean struct {
	Port int    //端口    7272
	Addr string //绑定ip

	Logdir  string //日志路径
	LogName string //日志名字

	Db_Exsit          int    //数据库操作
	Db_dataSourceName string // 数据库连接
	Db_MaxOpenConns   int    // 最大连接
	Db_MaxIdleConns   int    // 最大闲置连接

	HeartBeat int //ping 客户端心跳时间 秒

	HttpPort int //http端口

	Presence int //出席

	KV map[string]string //系统key value

	ConfirmAck int //发送消息是否需要回执 1需要 0不需要
}

/**设置Ip信息*/
func (cf *ConfBean) SetIp(port int, addr string) {
	cf.Port, cf.Addr = port, addr
}
func (cf *ConfBean) GetIp() (port int, addr string) {
	return cf.Port, cf.Addr
}

func (cf *ConfBean) GetHttpPort() (port int) {
	if cf.HttpPort > 0 {
		port = cf.HttpPort
	} else {
		port = 3939
	}
	return
}

/**设置日志信息*/
func (cf *ConfBean) SetLog(logdir, logname string) {
	cf.Logdir, cf.LogName = logdir, logname
}
func (cf *ConfBean) GetLog() (logdir string, logname string) {
	return cf.Logdir, cf.LogName
}

/**数据库设置*/
func (cf *ConfBean) SetDB(dataSourceName string, maxOpenConns, maxIdleConns int) {
	cf.Db_dataSourceName, cf.Db_MaxOpenConns, cf.Db_MaxIdleConns = dataSourceName, maxOpenConns, maxIdleConns
}

func (cf *ConfBean) GetDB() (dataSourceName string, maxOpenConns, maxIdleConns int) {
	return cf.Db_dataSourceName, cf.Db_MaxOpenConns, cf.Db_MaxIdleConns
}

func (cf *ConfBean) GetKV(keyword string, defaultValue string) (value string) {
	if v, ok := cf.KV[keyword]; ok {
		value = v
	} else {
		value = defaultValue
	}
	return
}

func (cf *ConfBean) Init(filexml string) {
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
			fmt.Println(`======================conf start======================`)
			i := 0
			for _, node := range nodes {
				name := node.Name()
				value := node.Value
				v := reflect.ValueOf(cf).Elem().FieldByName(name)
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
			fmt.Println(`=======================conf end=======================`)
			if i > 0 {
				fmt.Println("not set number:", i)
			}
		}
	}
}
