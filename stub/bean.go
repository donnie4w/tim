// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package stub

type RemoteNode struct {
	Addr      string
	UUID      int64
	CSNUM     int32
	Host      string
	AdminAddr string
	StatDesc  string
}

type ConfBean struct {
	MaskSeed             string      `json:"maskseed"`
	Salt                 string      `json:"salt"`
	Tldb                 *tldb       `json:"tldb"`
	TldbExtent           []*tldb     `json:"tldb_extent"`
	InlineDB             *InlineDB   `json:"inlinedb"`
	InlineExtent         []*InlineDB `json:"inlinedb_extent"`
	ExternalDB           *externalDB `json:"externaldb"`
	NoDB                 *bool       `json:"nodb"`
	Security             *security   `json:"security"`
	PprofAddr            string      `json:"pprof_addr"`
	ClientListen         string      `json:"client_listen"`
	Ssl_crt              string      `json:"ssl_certificate"`
	Ssl_crt_key          string      `json:"ssl_certificate_key"`
	ConnectLimit         int64       `json:"connect_limit"`
	Memlimit             int         `json:"memlimit"`
	MaxMessageSize       int         `json:"maxmessagesize"`
	CacheAuthExpire      int64       `json:"cache_authexpire"`
	NodeMaxlength        *int        `json:"node_maxlength"`
	CsListen             string      `json:"cluster_listen"`
	CsAccess             string      `json:"cluster_access"`
	WebAdminListen       string      `json:"webadmin_listen"`
	WebAdminNoTls        bool        `json:"webadmin_nouse_tls"`
	AdmListen            *string     `json:"server_api_listen"`
	NoInit               bool        `json:"noinit"`
	PingTo               int64       `json:"ping_timeout"`
	Keystore             *string     `json:"keystore"`
	MaxBackup            *int        `json:"maxbackup"`
	RequestRate          int64       `json:"request_rate"`
	DeviceLimit          int         `json:"device_limit"`
	DevicetypeLimit      int         `json:"devicetype_limit"`
	MessageNoAuth        bool        `json:"message_noauth"`
	PresenceOfflineBlock bool        `json:"presence_offline_block"`
	TimAdminAuth         bool        `json:"timAdminAuth"`
	TTL                  uint64      `json:"ttl"`
	TokenTimeout         int64       `json:"tokenTimeout"`
	BlockAPI             []int8      `json:"blockapi"`
	Raftx                *Raftx      `json:"raftx"`
	Rax                  *Rax        `json:"rax"`
	Etcd                 *Etcd       `json:"etcd"`
	Redis                *Redis      `json:"redis"`
	ZooKeeper            *ZooKeeper  `json:"zookeeper"`
	//Notice               *notice     `json:"notice"`
	//NoDBAuth             *noDBAuth `json:"nodbauth"`
	Public     string  `json:"public.node"`
	Pwd        string  `json:"cluster.pwd"`
	EncryptKey string  `json:"cluster.encryptkey"`
	Bind       *string `json:"bind"`
}

type externalDB struct {
	Tim_mysql_connection      string `json:"tim.mysql.connection"`
	Tim_postgresql_connection string `json:"tim.postgresql.connection"`
	Tim_oracle_connection     string `json:"tim.oracle.connection"`
	Tim_sqlserver_connection  string `json:"tim.sqlserver.connection"`

	Tim_mysql_connection_mod      string `json:"tim.mysql.connection.mod"`
	Tim_postgresql_connection_mod string `json:"tim.postgresql.connection.mod"`
	Tim_oracle_connection_mod     string `json:"tim.oracle.connection.mod"`
	Tim_sqlserver_connection_mod  string `json:"tim.sqlserver.connection.mod"`

	Tim_sql_login string `json:"tim.sql.login"`
	Tim_sql_token string `json:"tim.sql.token"`

	Tim_sql_savemessage     string `json:"tim.sql.message.save"`
	Tim_sql_getmessage      string `json:"tim.sql.message.get"`
	Tim_sql_getchatid_byid  string `json:"tim.sql.message.chatid.byid"`
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

type InlineDB struct {
	Tim_mysql_connection      string `json:"tim.mysql.connection"`
	Tim_postgresql_connection string `json:"tim.postgresql.connection"`
	Tim_sqlserver_connection  string `json:"tim.sqlserver.connection"`
	Tim_sqlite_connection     string `json:"tim.sqlite.connection"`

	Tim_mysql_connection_mod      string `json:"tim.mysql.connection.mod"`
	Tim_postgresql_connection_mod string `json:"tim.postgresql.connection.mod"`
	Tim_sqlserver_connection_mod  string `json:"tim.sqlserver.connection.mod"`

	ExtentMax int `json:"extent"`
}

type security struct {
	MaxDatalimit   int64   `json:"maxdata"`
	ReqHzSecond    int     `json:"reqhzsecond"`
	ForBidRegister bool    `json:"forbid_register"`
	ForBidToken    bool    `json:"forbid_token"`
	ConnectAuthUrl *string `json:"connectauth_url"`
}

type Openssl struct {
	PrivateBytes []byte
	PublicBytes  []byte
	PublicPath   string
	PrivatePath  string
}

//type notice struct {
//	Loginstat *string `json:"loginstat.url"`
//}

//type noDBAuth struct {
//	Url      *string `json:"url"`
//	TimName  string  `json:"timname"`
//	TimPwd   string  `json:"timpwd"`
//	UserName string  `json:"username"`
//	Password string  `json:"password"`
//}

type tldb struct {
	Addr      string `json:"addr"`
	Auth      string `json:"auth"`
	Tls       bool   `json:"tls"`
	ExtentMax int    `json:"extent"`
}

type Raftx struct {
	ListenAddr string   `json:"listen"`
	Peers      []string `json:"peers"`
}

type Redis struct {
	Addr     []string `json:"addr"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	DB       int      `json:"db"`
	Protocol int      `json:"protocol"`
}

type Etcd struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	CAFile    string   `json:"cafile"`   // CA certificate path
	CertFile  string   `json:"certfile"` // Path of the client certificate
	KeyFile   string   `json:"keyfile"`  // Path of the client private key
}

type ZooKeeper struct {
	Servers        []string `json:"servers"`
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	SessionTimeout int      `json:"session_timeout"`
}

type Rax struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}
