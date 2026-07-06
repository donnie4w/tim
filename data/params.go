// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

type Sqlite string

func (s Sqlite) CreateSql() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS timmessage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chatid BLOB NOT NULL,
    fid INTEGER NOT NULL,
    stanza BLOB,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timmessage_timeseries ON timmessage (timeseries);
CREATE INDEX IF NOT EXISTS idx_timmessage_chatid ON timmessage (chatid);`,

		`CREATE TABLE IF NOT EXISTS timuser (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid INTEGER NOT NULL UNIQUE,
    pwd VARCHAR(60) NOT NULL,
    createtime INTEGER NOT NULL,
    ubean BLOB,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timuser_timeseries ON timuser (timeseries);`,

		`CREATE TABLE IF NOT EXISTS timgroup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gtype INTEGER NOT NULL,
    uuid INTEGER NOT NULL UNIQUE,
    createtime INTEGER NOT NULL,
    status INTEGER NOT NULL,
    rbean BLOB,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timgroup_timeseries ON timgroup (timeseries);`,

		`CREATE TABLE IF NOT EXISTS timoffline (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid INTEGER NOT NULL,
    chatid BLOB DEFAULT NULL,
    stanza BLOB DEFAULT NULL,
    mid INTEGER DEFAULT NULL,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timoffline_timeseries ON timoffline (timeseries);
CREATE INDEX IF NOT EXISTS idx_timoffline_uuid ON timoffline (uuid);`,

		`CREATE TABLE IF NOT EXISTS timrelate (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid BLOB NOT NULL,
    status INTEGER NOT NULL,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timrelate_timeseries ON timrelate (timeseries);
CREATE INDEX IF NOT EXISTS idx_timrelate_Uuid ON timrelate (uuid);`,

		`CREATE TABLE IF NOT EXISTS timroster (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    unikid BLOB NOT NULL UNIQUE,
    uuid INTEGER NOT NULL,
    tuuid INTEGER NOT NULL,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timroster_timeseries ON timroster (timeseries);
CREATE INDEX IF NOT EXISTS idx_timroster_uuid ON timroster (uuid);`,

		`CREATE TABLE IF NOT EXISTS timmucroster (
    id INTEGER PRIMARY KEY 	AUTOINCREMENT,
    unikid BLOB NOT NULL UNIQUE,
    uuid INTEGER NOT NULL,
    tuuid INTEGER NOT NULL,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timmucroster_timeseries ON timmucroster (timeseries);
CREATE INDEX IF NOT EXISTS idx_timmucroster_uuid ON timmucroster (uuid);`,

		`CREATE TABLE IF NOT EXISTS timblock (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    unikid BLOB NOT NULL UNIQUE,
    uuid INTEGER NOT NULL,
    tuuid INTEGER NOT NULL,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timblock_timeseries ON timblock (timeseries);
CREATE INDEX IF NOT EXISTS idx_timblock_uuid ON timblock (uuid);`,

		`CREATE TABLE IF NOT EXISTS timblockroom (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    unikid BLOB NOT NULL UNIQUE,
    uuid INTEGER NOT NULL,
    tuuid INTEGER NOT NULL,
    timeseries INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timblockroom_timeseries ON timblockroom (timeseries);
CREATE INDEX IF NOT EXISTS idx_timblockroom_uuid ON timblockroom (uuid);`,

		`CREATE TABLE IF NOT EXISTS timdomain (
   	id INTEGER PRIMARY KEY AUTOINCREMENT,
  	adminaccount TEXT NOT NULL,
  	adminpassword TEXT NOT NULL,
  	timdomain TEXT NOT NULL,
  	createtime INTEGER NOT NULL,
  	timeseries INTEGER NOT NULL
);
CREATE UNIQUE INDEX uk_timdomain_adminaccount ON timdomain (adminaccount);
CREATE UNIQUE INDEX uk_timdomain_timdomain ON timdomain (timdomain);
CREATE INDEX idx_timdomain_timeseries ON timdomain (timeseries);`,
	}
}

type PostgreSql string

func (s PostgreSql) CreateSql() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS timmessage (
    id BIGSERIAL PRIMARY KEY,
    chatid BYTEA NOT NULL CHECK (length(chatid) = 16),
    fid INT NOT NULL,
    stanza BYTEA,
    timeseries BIGINT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_timmessage_timeseries ON timmessage (timeseries);`,
		`CREATE INDEX IF NOT EXISTS idx_timmessage_chatid ON timmessage USING hash(chatid) ;`,

		`CREATE TABLE IF NOT EXISTS timuser (
    id BIGSERIAL PRIMARY KEY,
    uuid BIGINT NOT NULL UNIQUE,
    pwd VARCHAR(60) NOT NULL,
    createtime BIGINT NOT NULL,
    ubean BYTEA,
    timeseries BIGINT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_timuser_timeseries ON timuser (timeseries);`,

		`CREATE TABLE IF NOT EXISTS timgroup (
    id BIGSERIAL PRIMARY KEY,
    gtype SMALLINT NOT NULL,
    uuid BIGINT NOT NULL UNIQUE,
    createtime BIGINT NOT NULL,
    status SMALLINT NOT NULL,
    rbean BYTEA,
    timeseries BIGINT NOT NULL
);`,

		`CREATE INDEX IF NOT EXISTS idx_timgroup_timeseries ON timgroup (timeseries);`,

		`CREATE TABLE IF NOT EXISTS timoffline (
    id BIGSERIAL PRIMARY KEY,
    uuid BIGINT NOT NULL,
    chatid BIGINT DEFAULT NULL,
    stanza BYTEA DEFAULT NULL,
    mid BIGINT DEFAULT NULL,
    timeseries BIGINT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_timoffline_timeseries ON timoffline (timeseries);`,
		`CREATE INDEX IF NOT EXISTS idx_timoffline_uuid ON timoffline USING hash(uuid);`,

		`CREATE TABLE IF NOT EXISTS timrelate (
    id BIGSERIAL PRIMARY KEY,
    uuid BYTEA NOT NULL CHECK (length(uuid) = 16) UNIQUE,
    status SMALLINT NOT NULL,
    timeseries BIGINT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_timrelate_timeseries ON timrelate (timeseries);`,
		`CREATE INDEX IF NOT EXISTS idx_timrelate_uuid ON timrelate USING hash(uuid);`,

		`CREATE TABLE IF NOT EXISTS timroster (
    id BIGSERIAL PRIMARY KEY,
    unikid BYTEA NOT NULL CHECK (length(unikid) = 16) UNIQUE,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL
);`,

		`CREATE INDEX IF NOT EXISTS idx_timroster_timeseries ON timroster (timeseries)`,
		`CREATE INDEX IF NOT EXISTS idx_timroster_uuid ON timroster USING hash(uuid);`,

		`CREATE TABLE IF NOT EXISTS timmucroster (
    id BIGSERIAL PRIMARY KEY,
    unikid BYTEA NOT NULL CHECK (length(unikid) = 16) UNIQUE,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL
);`,

		`CREATE INDEX IF NOT EXISTS idx_timmucroster_timeseries ON timmucroster (timeseries);`,
		`CREATE INDEX IF NOT EXISTS idx_timmucroster_uuid ON timmucroster USING hash(uuid);`,

		`CREATE TABLE IF NOT EXISTS timblock (
    id BIGSERIAL PRIMARY KEY,
    unikid BYTEA NOT NULL CHECK (length(unikid) = 16) UNIQUE,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_timblock_timeseries ON timblock (timeseries);`,
		`CREATE INDEX IF NOT EXISTS idx_timblock_uuid ON timblock USING hash(uuid);`,

		`CREATE TABLE IF NOT EXISTS timblockroom (
    id BIGSERIAL PRIMARY KEY,
    unikid BYTEA NOT NULL CHECK (length(unikid) = 16) UNIQUE,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL
);`,

		`CREATE INDEX IF NOT EXISTS idx_timblockroom_timeseries ON timblockroom (timeseries);`,
		`CREATE INDEX IF NOT EXISTS idx_timblockroom_uuid ON timblockroom USING hash(uuid);`,

		`CREATE TABLE IF NOT EXISTS timdomain (
  id BIGSERIAL PRIMARY KEY,
  adminaccount VARCHAR(128) NOT NULL,
  adminpassword VARCHAR(64) NOT NULL,
  timdomain VARCHAR(64) NOT NULL,
  createtime BIGINT NOT NULL,
  timeseries BIGINT NOT NULL,
  CONSTRAINT uk_timdomain_adminaccount UNIQUE (adminaccount),
  CONSTRAINT uk_timdomain_timdomain UNIQUE (timdomain)
);`,
		`CREATE INDEX idx_timuser_timeseries ON timdomain (timeseries);`,
	}
}

type Mysql string

func (s Mysql) CreateSql() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS timmessage (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chatid BINARY(16) NOT NULL,
    fid int NOT NULL,
    stanza BLOB,
    timeseries BIGINT NOT NULL,
    KEY idx_timmessage_timeseries (timeseries),
    KEY idx_timmessage_chatid (chatid)
);`,

		`CREATE TABLE IF NOT EXISTS timuser (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid BIGINT NOT NULL,
    pwd VARCHAR(60) NOT NULL,
    createtime BIGINT NOT NULL,
    ubean BLOB,
    timeseries BIGINT NOT NULL,
    KEY idx_timuser_timeseries (timeseries),
    UNIQUE INDEX uk_timuser_uuid (uuid) 
);`,

		`CREATE TABLE IF NOT EXISTS timgroup (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    gtype TINYINT NOT NULL,
    uuid BIGINT NOT NULL,
    createtime BIGINT NOT NULL,
    status TINYINT NOT NULL,
    rbean BLOB,
    timeseries BIGINT NOT NULL,
    KEY idx_timgroup_timeseries (timeseries),
    UNIQUE INDEX uk_timgroup_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timoffline (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid BIGINT NOT NULL,
    chatid BIGINT DEFAULT NULL,
    stanza BLOB,
    mid BIGINT DEFAULT NULL,
    timeseries BIGINT NOT NULL,
    KEY idx_timoffline_timeseries (timeseries),
    KEY idx_timoffline_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timrelate (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid BINARY(16) NOT NULL,
    status TINYINT NOT NULL,
    timeseries BIGINT NOT NULL,
    KEY idx_timrelate_timeseries (timeseries),
    UNIQUE INDEX uk_timrelate_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timroster (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    unikid BINARY(16) NOT NULL,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL,
    UNIQUE INDEX uk_timroster_unikid (unikid),
    KEY idx_timroster_timeseries (timeseries),
    KEY idx_timroster_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timmucroster (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    unikid BINARY(16) NOT NULL,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL,
    UNIQUE INDEX uk_timmucroster_unikid (unikid),
    KEY idx_timmucroster_timeseries (timeseries),
    KEY idx_timmucroster_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timblock (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    unikid BINARY(16) NOT NULL,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL,
    UNIQUE INDEX uk_timblock_unikid (unikid),
    KEY idx_timblock_timeseries (timeseries),
    KEY idx_timblock_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timblockroom (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    unikid BINARY(16) NOT NULL,
    uuid BIGINT NOT NULL,
    tuuid BIGINT NOT NULL,
    timeseries BIGINT NOT NULL,
    UNIQUE INDEX uk_timblockroom_unikid (unikid),
    KEY idx_timblockroom_timeseries (timeseries),
    KEY idx_timblockroom_uuid (uuid)
);`,

		`CREATE TABLE IF NOT EXISTS timdomain (
	  id BIGINT NOT NULL AUTO_INCREMENT,
	  adminaccount VARCHAR(128) NOT NULL,
	  adminpassword VARCHAR(64) NOT NULL,
	  timdomain VARCHAR(64) NOT NULL,
	  createtime BIGINT NOT NULL,
	  timeseries BIGINT NOT NULL,
	  PRIMARY KEY (id),
	  UNIQUE KEY uk_timdomain_adminaccount (adminaccount),
	  UNIQUE KEY uk_timdomain_timdomain (timdomain),
	  KEY idx_timdomain_timeseries (timeseries)
);`,
	}
}

type SqlServer string

func (s SqlServer) CreateSql() []string {
	return []string{
		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timmessage]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timmessage] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [chatid] BINARY(16) NOT NULL,
		[fid] INT NOT NULL,
        [stanza] VARBINARY(MAX),
        [timeseries] BIGINT NOT NULL
    );
    CREATE INDEX [idx_timmessage_timeseries] ON [timmessage] ([timeseries]);
    CREATE INDEX [idx_timmessage_chatid] ON [timmessage] ([chatid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timuser]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timuser] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [uuid] BIGINT NOT NULL,
        [pwd] VARCHAR(60) NOT NULL,
        [createtime] BIGINT NOT NULL,
        [ubean] VARBINARY(MAX),
        [timeseries] BIGINT NOT NULL
    );
    CREATE INDEX [idx_timuser_timeseries] ON [timuser] ([timeseries]);
    CREATE UNIQUE INDEX [uq_timuser_uuid] ON [timuser] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timgroup]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timgroup] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [gtype] TINYINT NOT NULL,
        [uuid] BIGINT NOT NULL,
        [createtime] BIGINT NOT NULL,
        [status] TINYINT NOT NULL,
        [rbean] VARBINARY(MAX),
        [timeseries] BIGINT NOT NULL
    );
    CREATE INDEX [idx_timgroup_timeseries] ON [timgroup] ([timeseries]);
    CREATE UNIQUE INDEX [uq_timgroup_uuid] ON [timgroup] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timoffline]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timoffline] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [uuid] BIGINT NOT NULL,
        [chatid] BIGINT DEFAULT NULL,
        [stanza] VARBINARY(MAX) DEFAULT NULL,
        [mid] BIGINT DEFAULT NULL,
        [timeseries] BIGINT NOT NULL
    );
    CREATE INDEX [idx_timoffline_timeseries] ON [timoffline] ([timeseries]);
    CREATE INDEX [idx_timoffline_uuid] ON [timoffline] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timrelate]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timrelate] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [uuid] BINARY(16) NOT NULL,
        [status] TINYINT NOT NULL,
        [timeseries] BIGINT NOT NULL
    );
    CREATE INDEX [idx_timrelate_timeseries] ON [timrelate] ([timeseries]);
    CREATE UNIQUE [INDEX uq_timrelate_uuid] ON [timrelate] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timroster]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timroster] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [unikid] BINARY(16) NOT NULL,
        [uuid] BIGINT NOT NULL,
        [tuuid] BIGINT NOT NULL,
        [timeseries] BIGINT NOT NULL
    );
    CREATE UNIQUE INDEX [uq_timroster_unikid] ON [timroster] ([unikid]);
    CREATE INDEX [idx_timroster_timeseries] ON [timroster] ([timeseries]);
    CREATE INDEX [idx_timroster_uuid] ON [timroster] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timmucroster]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timmucroster] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [unikid] BINARY(16) NOT NULL,
        [uuid] BIGINT NOT NULL,
        [tuuid] BIGINT NOT NULL,
        [timeseries] BIGINT NOT NULL
    );
    CREATE UNIQUE INDEX [uq_timmucroster_unikid] ON [timmucroster] ([unikid]);
    CREATE INDEX [idx_timmucroster_timeseries] ON [timmucroster] ([timeseries]);
    CREATE INDEX [idx_timmucroster_uuid] ON [timmucroster] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timblock]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timblock] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [unikid] BINARY(16) NOT NULL,
        [uuid] BIGINT NOT NULL,
        [tuuid] BIGINT NOT NULL,
        [timeseries] BIGINT NOT NULL
    );
    CREATE UNIQUE INDEX [uq_timblock_unikid] ON [timblock] ([unikid]);
    CREATE INDEX [idx_timblock_timeseries] ON [timblock] ([timeseries]);
    CREATE INDEX [idx_timblock_uuid] ON [timblock] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timblockroom]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timblockroom] (
        [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
        [unikid] BINARY(16) NOT NULL,
        [uuid] BIGINT NOT NULL,
        [tuuid] BIGINT NOT NULL,
        [timeseries] BIGINT NOT NULL
    );
    CREATE UNIQUE INDEX [uq_timblockroom_unikid] ON [timblockroom] ([unikid]);
    CREATE INDEX [idx_timblockroom_timeseries] ON [timblockroom] ([timeseries]);
    CREATE INDEX [idx_timblockroom_uuid] ON [timblockroom] ([uuid]);
END;`,

		`IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[timdomain]') AND type in (N'U'))
BEGIN
    CREATE TABLE [timdomain] (
        [id] BIGINT NOT NULL IDENTITY(1,1),
        [adminaccount] NVARCHAR(128) NOT NULL,
        [adminpassword] NVARCHAR(64) NOT NULL,
        [timdomain] NVARCHAR(64) NOT NULL,
        [createtime] BIGINT NOT NULL,
        [timeseries] BIGINT NOT NULL,
        PRIMARY KEY ([id]),
        UNIQUE ([adminaccount]),
        UNIQUE ([timdomain])
    );
    CREATE INDEX [idx_timdomain_timeseries] ON [timdomain] ([timeseries]);
END;`,
	}
}

type Oracle string

func (s Oracle) CreateSql() []string {
	return []string{
		`CREATE TABLE timmessage (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,  
    chatid RAW(16) NOT NULL, 
    fid NUMBER NOT NULL,
    stanza BLOB,
    timeseries NUMBER NOT NULL, 
    CONSTRAINT idx_timmessage_timeseries INDEX (timeseries),
    CONSTRAINT idx_timmessage_chatid INDEX (chatid)
);`,

		`CREATE TABLE timuser (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    uuid NUMBER NOT NULL,
    pwd VARCHAR2(60) NOT NULL,  
    createtime NUMBER NOT NULL,
    ubean BLOB,
    timeseries NUMBER NOT NULL,
    CONSTRAINT idx_timuser_timeseries INDEX (timeseries),
    CONSTRAINT uk_timuser_uuid UNIQUE (uuid)
);`,

		`CREATE TABLE timgroup (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    gtype NUMBER(3) NOT NULL,  
    uuid NUMBER NOT NULL,
    createtime NUMBER NOT NULL,
    status NUMBER(3) NOT NULL,  
    rbean BLOB,
    timeseries NUMBER NOT NULL,
    CONSTRAINT idx_timgroup_timeseries INDEX (timeseries),
    CONSTRAINT uk_timgroup_uuid UNIQUE (uuid)
);`,

		`CREATE TABLE timoffline (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    uuid NUMBER NOT NULL,
    chatid NUMBER DEFAULT NULL, 
    stanza BLOB,
    mid NUMBER DEFAULT NULL,
    timeseries NUMBER NOT NULL,
    CONSTRAINT idx_timoffline_timeseries INDEX (timeseries),
    CONSTRAINT idx_timoffline_uuid INDEX (uuid)
);`,

		`CREATE TABLE timrelate (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    uuid RAW(16) NOT NULL,  
    status NUMBER(3) NOT NULL,  
    timeseries NUMBER NOT NULL,
    CONSTRAINT idx_timrelate_timeseries INDEX (timeseries),
    CONSTRAINT uk_timrelate_uuid UNIQUE (uuid)
);`,

		`CREATE TABLE timroster (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    unikid RAW(16) NOT NULL, 
    uuid NUMBER NOT NULL,
    tuuid NUMBER NOT NULL,
    timeseries NUMBER NOT NULL,
    CONSTRAINT uk_timroster_unikid UNIQUE (unikid),
    CONSTRAINT idx_timroster_timeseries INDEX (timeseries),
    CONSTRAINT idx_timroster_uuid INDEX (uuid)
);`,

		`CREATE TABLE timmucroster (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    unikid RAW(16) NOT NULL, 
    uuid NUMBER NOT NULL,
    tuuid NUMBER NOT NULL,
    timeseries NUMBER NOT NULL,
    CONSTRAINT uk_timmucroster_unikid UNIQUE (unikid),
    CONSTRAINT idx_timmucroster_timeseries INDEX (timeseries),
    CONSTRAINT idx_timmucroster_uuid INDEX (uuid)
);`,

		`CREATE TABLE timblock (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    unikid RAW(16) NOT NULL, 
    uuid NUMBER NOT NULL,
    tuuid NUMBER NOT NULL,
    timeseries NUMBER NOT NULL,
    CONSTRAINT uk_timblock_unikid UNIQUE (unikid),
    CONSTRAINT idx_timblock_timeseries INDEX (timeseries),
    CONSTRAINT idx_timblock_uuid INDEX (uuid)
);`,

		`CREATE TABLE timblockroom (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    unikid RAW(16) NOT NULL, 
    uuid NUMBER NOT NULL,
    tuuid NUMBER NOT NULL,
    timeseries NUMBER NOT NULL,
    CONSTRAINT uk_timblockroom_unikid UNIQUE (unikid),
    CONSTRAINT idx_timblockroom_timeseries INDEX (timeseries),
    CONSTRAINT idx_timblockroom_uuid INDEX (uuid)
);`,

		`CREATE TABLE timdomain (
    id NUMBER PRIMARY KEY, 
    adminaccount VARCHAR2(128) NOT NULL,
    adminpassword VARCHAR2(64) NOT NULL,
    timdomain VARCHAR2(64) NOT NULL,
    createtime NUMBER NOT NULL,
    timeseries NUMBER NOT NULL,
    CONSTRAINT uk_timdomain_adminaccount UNIQUE (adminaccount),
    CONSTRAINT uk_timdomain_timdomain UNIQUE (timdomain),
    CONSTRAINT idx_timdomain_timeseries INDEX (timeseries)
);`,
	}
}

type MongoDB string

func (s MongoDB) CreateSql() []string {
	return []string{`{
    "collection": "timmessage",
    "fields": {
      "_id": "ObjectId",
      "mid": "NumberLong",
      "chatid": "Binary",
      "fid": "NumberInt",
      "stanza": "Binary",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "timeseries": 1 } },
      { "key": { "chatid": 1 } },
      { "key": { "mid": 1 } }
    ]
  }`,
		`{
    "collection": "timuser",
    "fields": {
      "_id": "ObjectId",
      "uuid": "NumberLong",
      "pwd": "String",
      "createtime": "NumberLong",
      "ubean": "Binary",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 }, "unique": true }
    ]
  }`,
		`{
    "collection": "timgroup",
    "fields": {
      "_id": "ObjectId",
      "gtype": "NumberInt",
      "uuid": "NumberLong",
      "createtime": "NumberLong",
      "status": "NumberInt",
      "rbean": "Binary",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 }, "unique": true }
    ]
  }`,
		`{
    "collection": "timoffline",
    "fields": {
      "_id": "ObjectId",
      "uuid": "NumberLong",
      "chatid": "Binary",
      "stanza": "Binary",
      "mid": "NumberLong",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 } }
    ]
  }`,
		`{
    "collection": "timrelate",
    "fields": {
      "_id": "ObjectId",
      "uuid": "Binary",
      "status": "NumberInt",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 }, "unique": true }
    ]
  }`,
		`{
    "collection": "timroster",
    "fields": {
      "_id": "ObjectId",
      "unikid": "Binary",
      "uuid": "NumberLong",
      "tuuid": "NumberLong",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "unikid": 1 }, "unique": true },
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 } }
    ]
  }`,
		`{
    "collection": "timmucroster",
    "fields": {
      "_id": "ObjectId",
      "unikid": "Binary",
      "uuid": "NumberLong",
      "tuuid": "NumberLong",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "unikid": 1 }, "unique": true },
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 } }
    ]
  }`,
		`{
    "collection": "timblock",
    "fields": {
      "_id": "ObjectId",
      "unikid": "Binary",
      "uuid": "NumberLong",
      "tuuid": "NumberLong",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "unikid": 1 }, "unique": true },
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 } }
    ]
  }`,
		`{
    "collection": "timblockroom",
    "fields": {
      "_id": "ObjectId",
      "unikid": "Binary",
      "uuid": "NumberLong",
      "tuuid": "NumberLong",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "unikid": 1 }, "unique": true },
      { "key": { "timeseries": 1 } },
      { "key": { "uuid": 1 } }
    ]
  }`,
		`{
    "collection": "timdomain",
    "fields": {
      "_id": "ObjectId",
      "adminaccount": "String",
      "adminpassword": "String",
      "timdomain": "String",
      "createtime": "NumberLong",
      "timeseries": "NumberLong"
    },
    "indexes": [
      { "key": { "adminaccount": 1 }, "unique": true },
      { "key": { "timdomain": 1 }, "unique": true },
      { "key": { "timeseries": 1 } }
    ]
  }`,
	}
}

type Cassandra string

func (c Cassandra) CreateSql() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS timmessage (
    id BIGINT,
    chatid BINARY(16),
    fid INT,
    stanza BLOB,
    timeseries BIGINT,
    PRIMARY KEY (chatid,id) 
);`,
		`CREATE TABLE IF NOT EXISTS timuser (
    id BIGINT,
    uuid BIGINT,
    pwd TEXT,
    createtime BIGINT,
    ubean BLOB,
    timeseries BIGINT,
    PRIMARY KEY (uuid,id)
);`,
		`CREATE TABLE IF NOT EXISTS timgroup (
    id BIGINT,
    gtype TINYINT,
    uuid BIGINT,
    createtime BIGINT,
    status TINYINT,
    rbean BLOB,
    timeseries BIGINT,
    PRIMARY KEY (uuid,id)
);`,
		`CREATE TABLE IF NOT EXISTS timoffline (
    id BIGINT,
    uuid BIGINT,
    chatid BIGINT,
    stanza BLOB,
    mid BIGINT,
    timeseries BIGINT,
    PRIMARY KEY (uuid,id)
);`,
		`CREATE TABLE IF NOT EXISTS timrelate (
    id BIGINT,
    uuid BINARY(16),
    status TINYINT,
    timeseries BIGINT,
    PRIMARY KEY (uuid,id)
);`,
		`CREATE TABLE IF NOT EXISTS timroster (
    id BIGINT,
    unikid BINARY(16),
    uuid BIGINT,
    tuuid BIGINT,
    timeseries BIGINT,
    PRIMARY KEY (uuid,unikid)
);`,
		`CREATE TABLE IF NOT EXISTS timmucroster (
    id BIGINT,
    unikid BINARY(16),
    uuid BIGINT,
    tuuid BIGINT,
    timeseries BIGINT,
    PRIMARY KEY (uuid,unikid)
);`,
		`CREATE TABLE IF NOT EXISTS timblock (
    id BIGINT,
    unikid BINARY(16),
    uuid BIGINT,
    tuuid BIGINT,
    timeseries BIGINT,
    PRIMARY KEY (uuid,unikid)
);`,
		`CREATE TABLE IF NOT EXISTS timblockroom (
    id BIGINT,
    unikid BINARY(16),
    uuid BIGINT,
    tuuid BIGINT,
    timeseries BIGINT,
    PRIMARY KEY (uuid,unikid)
);`,
		`CREATE TABLE IF NOT EXISTS timdomain (
    id BIGINT,
    adminaccount TEXT,
    adminpassword TEXT,
    timdomain TEXT,
    createtime BIGINT,
    timeseries BIGINT,
    PRIMARY KEY (adminaccount)
);`,
	}
}
