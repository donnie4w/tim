// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package stub

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/tim/log"
)

type RemoteNode struct {
	Addr      string
	UUID      int64
	CSNUM     int32
	Host      string
	AdminAddr string
	StatDesc  string
}

type ConfBean struct {
	MaskSeed             string         `json:"maskseed"`
	Salt                 string         `json:"salt"`
	Tldb                 *tldb          `json:"tldb"`
	TldbExtent           []*tldb        `json:"tldb_extent"`
	Mongodb              *mongodb       `json:"mongodb"`
	MongodbExtent        []*mongodb     `json:"mongodb_extent"`
	Cassandra            *cassandra     `json:"cassandra"`
	CassandraExtent      []*cassandra   `json:"cassandra_extent"`
	InlineDB             *InlineDB      `json:"inlinedb"`
	InlineExtent         []*InlineDB    `json:"inlinedb_extent"`
	ExternalDB           *externalDB    `json:"externaldb"`
	AdminAccount         *AdminAccount  `json:"admin_account_init"`
	NoDB                 *bool          `json:"nodb"`
	Security             *security      `json:"security"`
	PprofAddr            string         `json:"pprof_addr"`
	UuidBits             int            `json:"uuid_bits"`
	ClientListen         string         `json:"client_listen"`
	Ssl_crt              string         `json:"ssl_certificate"`
	Ssl_crt_key          string         `json:"ssl_certificate_key"`
	ConnectLimit         int64          `json:"connect_limit"`
	Memlimit             int            `json:"memlimit"`
	MaxMessageSize       int            `json:"maxmessagesize"`
	CacheAuthExpire      int64          `json:"cache_authexpire"`
	NodeMaxlength        *int           `json:"node_maxlength"`
	CsListen             string         `json:"cluster_listen"`
	CsAccess             string         `json:"cluster_access"`
	WebAdminListen       string         `json:"webadmin_listen"`
	WebAdminNoTls        bool           `json:"webadmin_nouse_tls"`
	AdmListen            *string        `json:"server_api_listen"`
	NoInit               bool           `json:"noinit"`
	PingTo               int64          `json:"ping_timeout"`
	KsDir                string         `json:"ks_dir"`
	MaxBackup            *int           `json:"maxbackup"`
	RequestRate          int64          `json:"request_rate"`
	DeviceLimit          int            `json:"device_limit"`
	DevicetypeLimit      int            `json:"devicetype_limit"`
	MessageNoAuth        bool           `json:"message_noauth"`
	PresenceOfflineBlock bool           `json:"presence_offline_block"`
	UseTimDomain         bool           `json:"useTimDomain"`
	TTL                  uint64         `json:"ttl"`
	TokenTimeout         int64          `json:"tokenTimeout"`
	DeactivateApi        *deactivateApi `json:"deactivate_api"`
	Raftx                *Raftx         `json:"raftx"`
	Rax                  *Rax           `json:"rax"`
	Etcd                 *Etcd          `json:"etcd"`
	Redis                *Redis         `json:"redis"`
	ZooKeeper            *ZooKeeper     `json:"zookeeper"`
	//BlockAPI            []int8         `json:"blockapi"`
	//Notice              *notice     	 `json:"notice"`
	//NoDBAuth            *noDBAuth 	 `json:"nodbauth"`
	Public     string  `json:"public_node"`
	Pwd        string  `json:"cluster_pwd"`
	EncryptKey string  `json:"cluster_encryptkey"`
	Bind       *string `json:"bind"`
}

type externalDB struct {
	MYSQL        *Connect `json:"mysql"`
	POSTGRESQL   *Connect `json:"postgresql"`
	SQLSERVER    *Connect `json:"sqlserver"`
	ORACLE       *Connect `json:"oracle"`
	MARIADB      *Connect `json:"mariadb"`
	OCEANBASE    *Connect `json:"oceanbase"`
	TIDB         *Connect `json:"tidb"`
	GREENPLUM    *Connect `json:"greenplum"`
	OPENGAUSS    *Connect `json:"opengauss"`
	COCKROACHDB  *Connect `json:"cockroachdb"`
	ENTERPRISEDB *Connect `json:"enterprisedb"`

	Sql_login string `json:"sql.login"`
	Sql_token string `json:"sql.token"`

	Sql_message_save string `json:"sql.message.save"`
	Sql_message_get  string `json:"sql.message.get"`
	Sql_message_fid  string `json:"sql.message.fid"`
	Sql_message_del  string `json:"sql.message.del"`

	Sql_offlinemsg_save       string `json:"sql.offlinemsg.save"`
	Sql_offlinemsg_save_mid   string `json:"sql.offlinemsg.save.mid"`
	Sql_offlinemsg_save_nomid string `json:"sql.offlinemsg.save.nomid"`
	Sql_offlinemsg_get        string `json:"sql.offlinemsg.get"`
	Sql_offlinemsg_del        string `json:"sql.offlinemsg.del"`
	Sql_offlinemsg_delin      string `json:"sql.offlinemsg.delin"`

	Sql_authuser  string `json:"sql.user.auth"`
	Sql_existuser string `json:"sql.user.exist"`
	Sql_authroom  string `json:"sql.room.auth"`
	Sql_existroom string `json:"sql.room.exist"`

	Sql_roster       string `json:"sql.roster"`
	Sql_roster_add   string `json:"sql.roster.add"`
	Sql_roster_rm    string `json:"sql.roster.rm"`
	Sql_roster_block string `json:"sql.roster.block"`

	Sql_roomroster string `json:"sql.room.roster"`
	Sql_userroom   string `json:"sql.user.room"`
}

type InlineDB struct {
	MYSQL        *Connect `json:"mysql"`
	POSTGRESQL   *Connect `json:"postgresql"`
	SQLSERVER    *Connect `json:"sqlserver"`
	ORACLE       *Connect `json:"oracle"`
	SQLITE       *Connect `json:"sqlite"`
	MARIADB      *Connect `json:"mariadb"`
	OCEANBASE    *Connect `json:"oceanbase"`
	TIDB         *Connect `json:"tidb"`
	GREENPLUM    *Connect `json:"greenplum"`
	OPENGAUSS    *Connect `json:"opengauss"`
	COCKROACHDB  *Connect `json:"cockroachdb"`
	ENTERPRISEDB *Connect `json:"enterprisedb"`
	ExtentMax    int      `json:"extent"`
}

type security struct {
	MaxDatalimit   int64   `json:"maxdata"`
	ReqHzSecond    int     `json:"reqhzsecond"`
	ConnectAuthUrl *string `json:"connectauth_url"`
}

type deactivateApi struct {
	TIMREGISTER        bool `json:"TIMREGISTER"`
	TIMTOKEN           bool `json:"TIMTOKEN"`
	TIMAUTH            bool `json:"TIMAUTH"`
	TIMOFFLINEMSG      bool `json:"TIMOFFLINEMSG"`
	TIMOFFLINEMSGEND   bool `json:"TIMOFFLINEMSGEND"`
	TIMBROADPRESENCE   bool `json:"TIMBROADPRESENCE"`
	TIMLOGOUT          bool `json:"TIMLOGOUT"`
	TIMPULLMESSAGE     bool `json:"TIMPULLMESSAGE"`
	TIMVROOM           bool `json:"TIMVROOM"`
	TIMBUSINESS        bool `json:"TIMBUSINESS"`
	TIMNODES           bool `json:"TIMNODES"`
	TIMMESSAGE         bool `json:"TIMMESSAGE"`
	TIMPRESENCE        bool `json:"TIMPRESENCE"`
	TIMREVOKEMESSAGE   bool `json:"TIMREVOKEMESSAGE"`
	TIMBURNMESSAGE     bool `json:"TIMBURNMESSAGE"`
	TIMSTREAM          bool `json:"TIMSTREAM"`
	TIMBIGSTRING       bool `json:"TIMBIGSTRING"`
	TIMBIGBINARY       bool `json:"TIMBIGBINARY"`
	TIMBIGBINARYSTREAM bool `json:"TIMBIGBINARYSTREAM"`
}

type Openssl struct {
	PrivateBytes []byte
	PublicBytes  []byte
	PublicPath   string
	PrivatePath  string
}

type AdminAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tldb struct {
	Addr      string `json:"addr"`
	Auth      string `json:"auth"`
	Tls       bool   `json:"tls"`
	ExtentMax int    `json:"extent"`
}

type Connect struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pwd    string `json:"password"`
	DBname string `json:"dbname"`
	Dsn    string `json:"dsn"`
}

func (c *Connect) DSN(dbType base.DBType) (driver string, dsn string) {
	switch dbType {
	case gdao.ORACLE:
		if c.Dsn == "" {
			c.Dsn = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, c.User, c.Pwd, c.Host, c.Port, c.DBname)
		}
		driver, dsn = "godror", c.Dsn
	case gdao.SQLSERVER:
		if c.Dsn == "" {
			c.Dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", c.User, c.Pwd, c.Host, c.Port, c.DBname)
		}
		driver, dsn = "sqlserver", c.Dsn
	case gdao.SQLITE:
		if c.Dsn == "" {
			c.Dsn = c.DBname
		}
		driver, dsn = "sqlite3", c.Dsn
	case gdao.MYSQL, gdao.MARIADB, gdao.OCEANBASE, gdao.TIDB:
		if c.Dsn == "" {
			c.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pwd, c.Host, c.Port, c.DBname)
		}
		driver, dsn = "mysql", c.Dsn
	case gdao.POSTGRESQL, gdao.GREENPLUM, gdao.OPENGAUSS, gdao.COCKROACHDB, gdao.ENTERPRISEDB:
		if c.Dsn == "" {
			c.Dsn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", c.User, c.Pwd, c.DBname, c.Host, c.Port)
		}
		driver, dsn = "postgres", c.Dsn
	}

	log.DelayPrint(c.dbname(dbType), "  DSN[", dsn, "]")
	return
}

func (c *Connect) dbname(dbType base.DBType) string {
	switch dbType {
	case gdao.ORACLE:
		return "ORACLE"
	case gdao.SQLSERVER:
		return "SQLSERVER"
	case gdao.SQLITE:
		return "SQLITE"
	case gdao.MYSQL:
		return "MYSQL"
	case gdao.MARIADB:
		return "MARIADB"
	case gdao.OCEANBASE:
		return "OCEANBASE"
	case gdao.TIDB:
		return "TIDB"
	case gdao.POSTGRESQL:
		return "POSTGRESQL"
	case gdao.GREENPLUM:
		return "GREENPLUM"
	case gdao.OPENGAUSS:
		return "OPENGAUSS"
	case gdao.COCKROACHDB:
		return "OPENGAUSS"
	case gdao.ENTERPRISEDB:
		return "ENTERPRISEDB"
	}
	return ""
}

type mongodb struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	AuthDB    string `json:"authDB"`
	DbName    string `json:"dbname"`
	ExtentMax int    `json:"extent"`
}

type cassandra struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Hosts    []string `json:"hosts"`
	Port     string   `json:"port"`
	Keyspace string   `json:"keyspace"`
}

type Raftx struct {
	ListenAddr string   `json:"listen"`
	Peers      []string `json:"endpoints"`
}

type Redis struct {
	Addr     []string `json:"endpoints"`
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
	Servers        []string `json:"endpoints"`
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	SessionTimeout int      `json:"session_timeout"`
}

type Rax struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}
