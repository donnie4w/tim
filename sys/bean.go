// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package sys

type RemoteNode struct {
	Addr      string
	UUID      int64
	CSNUM     int32
	Host      string
	AdminAddr string
	StatDesc  string
}

type Server interface {
	Serve() (err error)
	Close() (err error)
}

type tldb struct {
	Addr      string `json:"addr"`
	Auth      string `json:"auth"`
	Tls       bool   `json:"tls"`
	ExtentMax int    `json:"extent"`
}

type ConfBean struct {
	Seed                 int64     `json:"seed"`
	Salt                 string    `json:"salt"`
	Tldb                 *tldb     `json:"tldb"`
	TldbExtent           []*tldb   `json:"tldb.extent"`
	Property             *property `json:"sql.property"`
	Security             *security `json:"security"`
	Notice               *notice   `json:"notice"`
	NoDBAuth             *noDBAuth `json:"nodbauth"`
	Listen               int       `json:"im.listen"`
	Ssl_crt              string    `json:"ssl_certificate"`
	Ssl_crt_key          string    `json:"ssl_certificate_key"`
	ConnectLimit         int64     `json:"connectLimit"`
	Memlimit             int       `json:"memlimit"`
	MaxMessageSize       int       `json:"maxmessagesize"`
	Public               string    `json:"public.node"`
	CacheExpireTime      int       `json:"cache.expiretime"`
	Pwd                  string    `json:"cluser.pwd"`
	NodeMaxlength        *int      `json:"node.maxlength"`
	EncryptKey           string    `json:"cluser.encryptkey"`
	ClusListen           string    `json:"cluser.listen"`
	AdminListen          string    `json:"admin.listen"`
	AdminTls             bool      `json:"admin.tls"`
	AdmListen            *string   `json:"adm.listen"`
	Init                 bool      `json:"init"`
	Bind                 *string   `json:"bind"`
	PingTo               int64     `json:"ping.timeout"`
	Keystore             *string   `json:"keystore"`
	MaxBackup            *int      `json:"maxbackup"`
	TaskLimit            *int      `json:"tasklimit"`
	DeviceLimit          int       `json:"device.limit"`
	DevicetypeLimit      int       `json:"devicetype.limit"`
	MessageNoauth        bool      `json:"message.noauth"`
	PresenceOfflineBlock bool      `json:"presence.offline.block"`
}

type property struct {
	Tim_mysql_connection      string `json:"tim.mysql.connection"`
	Tim_postgreSQL_connection string `json:"tim.postgreSQL.connection"`
	Tim_oracle_connection     string `json:"tim.oracle.connection"`
	Tim_sqlserver_connection  string `json:"tim.sqlserver.connection"`

	Tim_mysql_connection_mod      string `json:"tim.mysql.connection.mod"`
	Tim_postgreSQL_connection_mod string `json:"tim.postgreSQL.connection.mod"`
	Tim_oracle_connection_mod     string `json:"tim.oracle.connection.mod"`
	Tim_sqlserver_connection_mod  string `json:"tim.sqlserver.connection.mod"`

	Tim_sql_login string `json:"tim.sql.login"`
	Tim_sql_token string `json:"tim.sql.token"`

	Tim_sql_savemessage     string `json:"tim.sql.message.save"`
	Tim_sql_getmessage      string `json:"tim.sql.message.get"`
	Tim_sql_getmessage_byid string `json:"tim.sql.message.get.byid"`
	Tim_sql_delmessage_byid string `json:"tim.sql.message.del.byid"`

	Tim_sql_offlinemsg_save       string `json:"tim.sql.offlinemsg.save"`
	Tim_sql_offlinemsg_save_mid   string `json:"tim.sql.offlinemsg.save.mid"`
	Tim_sql_offlinemsg_save_nomid string `json:"tim.sql.offlinemsg.save.nomid"`
	Tim_sql_offlinemsg_get        string `json:"tim.sql.offlinemsg.get"`
	Tim_sql_offlinemsg_del        string `json:"tim.sql.offlinemsg.del"`
	Tim_sql_offlinemsg_delin      string `json:"tim.sql.offlinemsg.del.in"`

	Tim_sql_authuser  string `json:"tim.sql.user.auth"`
	Tim_sql_existuser string `json:"tim.sql.user.exist"`
	Tim_sql_authroom  string `json:"tim.sql.room.auth"`
	Tim_sql_existroom string `json:"tim.sql.room.exist"`

	Tim_sql_roster       string `json:"tim.sql.roster"`
	Tim_sql_roster_add   string `json:"tim.sql.roster.add"`
	Tim_sql_roster_rm    string `json:"tim.sql.roster.rm"`
	Tim_sql_roster_block string `json:"tim.sql.roster.block"`

	Tim_sql_roomroster string `json:"tim.sql.room.roster"`
	Tim_sql_userroom   string `json:"tim.sql.user.room"`
}

type security struct {
	MaxDatalimit   int64   `json:"maxdata"`
	ReqHzSecond    int     `json:"reqhzsecond"`
	ForBidRegister bool    `json:"forbid.register"`
	ForBidToken    bool    `json:"forbid.token"`
	ConnectAuthUrl *string `json:"connectauth.url"`
	ForBidIface    []int8  `json:"forbid.ifaces"`
}

type notice struct {
	Loginstat *string `json:"loginstat.url"`
}

type openssl struct {
	PrivateBytes []byte
	PublicBytes  []byte
	PublicPath   string
	PrivatePath  string
}

type istat interface {
	CReq() int64
	CReqDo()
	CReqDone()

	CPros() int64
	CProsDo()
	CProsDone()

	Tx() int64
	TxDo()
	TxDone()

	Ibs() int64
	Ib(int64)

	Obs() int64
	Ob(int64)
}

type noDBAuth struct {
	Url      *string `json:"url"`
	TimName  string  `json:"timname"`
	TimPwd   string  `json:"timpwd"`
	UserName string  `json:"username"`
	Password string  `json:"password"`
}
