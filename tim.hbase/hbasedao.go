/**
 * donnie4w@gmail.com  tim server
 */
package hbase

import (
	"github.com/donnie4w/go-logger/logger"
)

//'tim_serialno','tablename','id'	//信息内容表
type Tim_serialno struct {
	tablename string
	id        int64
}

//'tim_message','id','stamp','chatid','fromuser','touser','msgtype','msgmode','gname','small','large','stanza','createtime','Index'	//信息内容表
type Tim_message struct {
	Id          int64  `#id`
	Stamp       string `stamp`
	Chatid      string `chatid`
	Fromuser    string `fromuser`
	Touser      string `touser`
	Msgtype     string `msgtype`
	Msgmode     string `msgmode`
	Gname       string `gname`
	Small       string `small`
	Large       string `large`
	Stanza      string `stanza`
	Createtime  string `createtime`
	IndexChatid string `idx_`
}

func (t *Tim_message) Tablename() string {
	return "tim_message"
}
func (t *Tim_message) Insert() (row int64, err error) {
	tablename := t.Tablename()
	row, err = getSerialNo(tablename, "id")
	if err != nil {
		logger.Error("insert error:", err.Error())
		return
	} else {
		saveObject(t, tablename, row)
	}
	return
}
func (t *Tim_message) Update(row int64) (err error) {
	tablename := t.Tablename()
	saveObject(t, tablename, row)
	return
}

func (t *Tim_message) Updates(rows []int64) (err error) {
	tablename := t.Tablename()
	saveObjects(t, tablename, rows)
	return
}

func (t *Tim_message) UpdateByBeans(beans []*Bean) (err error) {
	tablename := t.Tablename()
	saveObjectByBeans(t, tablename, beans)
	return
}

//'tim_offline','id','mid','domain','username','stamp','fromuser','msgtype','msgmode','gname','message_size','stanza','createtime','index'	//离线消息存储表
type Tim_offline struct {
	Id                  int64  `#id`
	Mid                 string `mid`
	Domain              string `domain`
	Username            string `username`
	Stamp               string `stamp`
	Fromuser            string `fromuser`
	Msgtype             string `msgtype`
	Msgmode             string `msgmode`
	Gname               string `gname`
	Message_size        string `message_size`
	Stanza              string `stanza`
	Createtime          string `createtime`
	IndexMid            string `idx_`
	IndexDomainUsername string `idx_`
}

func (t *Tim_offline) Tablename() string {
	return "tim_offline"
}

func (t *Tim_offline) Insert() (row int64, err error) {
	tablename := t.Tablename()
	row, err = getSerialNo(tablename, "id")
	if err != nil {
		return
	} else {
		saveObject(t, tablename, row)
	}
	return
}
func (t *Tim_offline) Update(row int64) (err error) {
	tablename := t.Tablename()
	saveObject(t, tablename, row)
	return
}

func (t *Tim_offline) Updates(rows []int64) (err error) {
	tablename := t.Tablename()
	saveObjects(t, tablename, rows)
	return
}

func (t *Tim_offline) Delete(row int64) (err error) {
	DeleteRow(t.Tablename(), row)
	return
}

func (t *Tim_offline) Deletes(rows []int64) (err error) {
	DeleteRows(t.Tablename(), rows)
	return
}

func (t *Tim_offline) DeleteByBean(beans []*Bean) (err error) {
	//	DeleteRows(t.Tablename(), rows)
	err = DeleteFromQualifier(t.Tablename(), beans)
	return
}

//// 'tim_domain','id','domain','createtime','remark' //域名表
//type Tim_domain struct {
//	Id         int64  `#id`
//	Domain     string `domain`
//	Createtime string `createtime`
//	Remark     string `remark`
//}

//func (t *Tim_domain) tablename() string {
//	return "tim_domain"
//}

//func (t *Tim_domain) Insert() (row int64, err error) {
//	tablename := t.tablename()
//	row, err = getSerialNo(tablename, "id")
//	if err != nil {
//		return
//	} else {
//		saveObject(t, tablename, row)
//	}
//	return
//}
//func (t *Tim_domain) Update(row int64) (err error) {
//	tablename := t.tablename()
//	saveObject(t, tablename, row)
//	return
//}

////'tim_config','id','keyword','valuestr','createtime','remark' //配置表
//type Tim_config struct {
//	Id         int64  `#id`
//	Keyword    string `keyword`
//	Valuestr   string `valuestr`
//	Createtime string `createtime`
//	Remark     string `remark`
//}

//func (t *Tim_config) tablename() string {
//	return "tim_config"
//}

//func (t *Tim_config) Insert() (row int64, err error) {
//	tablename := t.tablename()
//	row, err = getSerialNo(tablename, "id")
//	if err != nil {
//		return
//	} else {
//		saveObject(t, tablename, row)
//	}
//	return
//}

//func (t *Tim_config) Update(row int64) (err error) {
//	tablename := t.tablename()
//	saveObject(t, tablename, row)
//	return
//}

////'tim_property','id','keyword','valueint','valuestr','remark'  //系统属性表
//type Tim_property struct {
//	Id       int64  `#id`
//	Keyword  string `keyword`
//	Valueint string `valueint`
//	Valuestr string `valuestr`
//	Remark   string `remark`
//}

//func (t *Tim_property) tablename() string {
//	return "tim_property"
//}

//func (t *Tim_property) Insert() (row int64, err error) {
//	tablename := t.tablename()
//	row, err = getSerialNo(tablename, "id")
//	if err != nil {
//		return
//	} else {
//		saveObject(t, tablename, row)
//	}
//	return
//}

//func (t *Tim_property) Update(row int64) (err error) {
//	tablename := t.tablename()
//	saveObject(t, tablename, row)
//	return
//}

//create 'tim_mucmessage','id','stamp','fromuser','roomtidname','domain','msgtype','stanza','createtime','index'   //房间信息内容表
type Tim_mucmessage struct {
	Id                  int64  `#id`
	Stamp               string `stamp`
	Fromuser            string `fromuser`
	Roomtidname         string `roomtidname`
	Domain              string `domain`
	Msgtype             string `msgtype`
	Stanza              string `stanza`
	Createtime          string `createtime`
	IndexFromuserDomain string `idx_`
}

func (t *Tim_mucmessage) Tablename() string {
	return "tim_mucmessage"
}

func (t *Tim_mucmessage) Insert() (row int64, err error) {
	tablename := t.Tablename()
	row, err = getSerialNo(tablename, "id")
	if err != nil {
		return
	} else {
		saveObject(t, tablename, row)
	}
	return
}

func (t *Tim_mucmessage) Update(row int64) (err error) {
	tablename := t.Tablename()
	saveObject(t, tablename, row)
	return
}

func (t *Tim_mucmessage) Delete(row int64) (err error) {
	DeleteRow(t.Tablename(), row)
	return
}

func (t *Tim_mucmessage) Deletes(rows []int64) (err error) {
	DeleteRows(t.Tablename(), rows)
	return
}

//create 'tim_mucoffline','id','mid','domain','username','stamp','roomid','msgtype','message_size','createtime','index'  //房间离线消息存储表
type Tim_mucoffline struct {
	Id                  int64  `#id`
	Mid                 string `mid`
	Domain              string `domain`
	Username            string `username`
	Stamp               string `stamp`
	Roomid              string `roomid`
	Msgtype             string `msgtype`
	Message_size        string `message_size`
	Createtime          string `createtime`
	IndexMid            string `idx_`
	IndexDomainUsername string `idx_`
}

func (t *Tim_mucoffline) Tablename() string {
	return "tim_mucoffline"
}

func (t *Tim_mucoffline) Insert() (row int64, err error) {
	tablename := t.Tablename()
	row, err = getSerialNo(tablename, "id")
	if err != nil {
		return
	} else {
		saveObject(t, tablename, row)
	}
	return
}

func (t *Tim_mucoffline) Update(row int64) (err error) {
	tablename := t.Tablename()
	saveObject(t, tablename, row)
	return
}

func (t *Tim_mucoffline) Delete(row int64) (err error) {
	DeleteRow(t.Tablename(), row)
	return
}

func (t *Tim_mucoffline) Deletes(rows []int64) (err error) {
	DeleteRows(t.Tablename(), rows)
	return
}

func (t *Tim_mucoffline) DeleteByBean(beans []*Bean) (err error) {
	err = DeleteFromQualifier(t.Tablename(), beans)
	return
}
