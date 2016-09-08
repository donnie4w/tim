
DROP TABLE IF EXISTS `tim_config`;

CREATE TABLE `tim_config` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `keyword` varchar(64) NOT NULL COMMENT '键',
  `valuestr` varchar(64) NOT NULL COMMENT '值',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  `remark` varchar(100) NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='系统配置表';

/*Table structure for table `tim_domain` */

DROP TABLE IF EXISTS `tim_domain`;

CREATE TABLE `tim_domain` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `domain` varchar(64) NOT NULL COMMENT '域名',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  `remark` varchar(100) NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='域名表';

/*Table structure for table `tim_message` */

DROP TABLE IF EXISTS `tim_message`;

CREATE TABLE `tim_message` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `stamp` varchar(20) NOT NULL COMMENT '时间戳毫秒',
  `chatid` varchar(64) NOT NULL COMMENT '聊天ID',
  `fromuser` varchar(64) NOT NULL COMMENT '发信者Id',
  `touser` varchar(64) NOT NULL COMMENT '接收者Id',
  `msgtype` int(2) NOT NULL DEFAULT '1' COMMENT '1文字2图片3语音4视频',
  `msgmode` int(2) NOT NULL DEFAULT '1' COMMENT '类型 1chat 2group',
  `gname` varchar(64) NOT NULL DEFAULT '' COMMENT '群用户发信者Id',
  `small` int(1) NOT NULL DEFAULT '1' COMMENT '有效信息-小号',
  `large` int(1) NOT NULL DEFAULT '1' COMMENT '有效信息-大号',
  `stanza` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '信息体',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `tm_chatid` (`chatid`,`small`,`large`),
  KEY `tm_chatid_stamp` (`stamp`,`chatid`)
) ENGINE=InnoDB AUTO_INCREMENT=283068 DEFAULT CHARSET=utf8 COMMENT='信息内容表';

/*Table structure for table `tim_mucmember` */

DROP TABLE IF EXISTS `tim_mucmember`;

CREATE TABLE `tim_mucmember` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `roomtid` varchar(64) NOT NULL COMMENT '聊天TID',
  `domain` varchar(64) NOT NULL COMMENT '域名',
  `tidname` varchar(20) NOT NULL COMMENT 'tidname',
  `type` int(1) NOT NULL COMMENT '用户类型 0:普通用户 1管理者 2创建者',
  `nickname` varchar(32) NOT NULL COMMENT '昵称',
  `affiliation` int(4) NOT NULL COMMENT '等级',
  `updatetime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '最后修改时间',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `roommemberid_tid` (`roomtid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='tim房间用户信息表';

/*Table structure for table `tim_mucmessage` */

DROP TABLE IF EXISTS `tim_mucmessage`;

CREATE TABLE `tim_mucmessage` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `stamp` varchar(20) NOT NULL COMMENT '时间戳毫秒',
  `fromuser` varchar(64) NOT NULL COMMENT '发信者Id',
  `roomtidname` varchar(64) NOT NULL COMMENT '房间tidname',
  `domain` varchar(64) NOT NULL COMMENT '域名',
  `msgtype` int(2) NOT NULL DEFAULT '1' COMMENT '1文字2图片3语音4视频',
  `stanza` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '信息体',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `tm_mucfromuser` (`fromuser`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='房间信息内容表';

/*Table structure for table `tim_mucoffline` */

DROP TABLE IF EXISTS `tim_mucoffline`;

CREATE TABLE `tim_mucoffline` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `mid` int(10) NOT NULL COMMENT '消息mid',
  `domain` varchar(64) NOT NULL COMMENT '域名',
  `username` varchar(64) NOT NULL COMMENT '用户名称',
  `stamp` varchar(20) NOT NULL COMMENT '时间戳毫秒',
  `roomid` varchar(64) NOT NULL COMMENT '房间Id',
  `msgtype` int(2) NOT NULL DEFAULT '1' COMMENT '1文字2图片3语音4视频5其它',
  `message_size` int(10) NOT NULL COMMENT '消息的大小，字节',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `tm_mucusername` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=1679 DEFAULT CHARSET=utf8 COMMENT='房间离线消息存储表';

/*Table structure for table `tim_mucroom` */

DROP TABLE IF EXISTS `tim_mucroom`;

CREATE TABLE `tim_mucroom` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `roomtid` varchar(64) NOT NULL COMMENT '聊天TID',
  `theme` varchar(64) NOT NULL COMMENT '当前房间主题',
  `name` varchar(64) NOT NULL COMMENT '房间名',
  `domain` varchar(64) NOT NULL COMMENT '域名',
  `password` varchar(32) NOT NULL COMMENT '房间密码',
  `maxusers` int(10) NOT NULL COMMENT '用户个数最大值',
  `description` varchar(255) NOT NULL COMMENT '房间描述',
  `updatetime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '最后修改时间',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT 'tim房间信息表',
  PRIMARY KEY (`id`),
  KEY `room_tid` (`roomtid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='tim房间信息表';

/*Table structure for table `tim_offline` */

DROP TABLE IF EXISTS `tim_offline`;

CREATE TABLE `tim_offline` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `mid` int(10) NOT NULL COMMENT '消息mid',
  `domain` varchar(64) NOT NULL COMMENT '域名',
  `username` varchar(64) NOT NULL COMMENT '用户名称',
  `stamp` varchar(20) NOT NULL COMMENT '时间戳毫秒',
  `fromuser` varchar(64) NOT NULL COMMENT '发信者Id',
  `msgtype` int(2) NOT NULL DEFAULT '1' COMMENT '1文字2图片3语音4视频',
  `msgmode` int(2) NOT NULL DEFAULT '1' COMMENT '类型 1chat 2group',
  `gname` varchar(64) NOT NULL DEFAULT '' COMMENT '群用户发信者Id',
  `message_size` int(10) NOT NULL COMMENT '消息的大小，字节',
  `stanza` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '信息体',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `tm_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='离线消息存储表';

/*Table structure for table `tim_property` */

DROP TABLE IF EXISTS `tim_property`;

CREATE TABLE `tim_property` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `keyword` varchar(64) NOT NULL COMMENT '键',
  `valueint` int(64) NOT NULL DEFAULT '0' COMMENT '属性值int',
  `valuestr` varchar(255) NOT NULL DEFAULT '' COMMENT '属性值string',
  `remark` varchar(100) NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='系统属性表';

/*Table structure for table `tim_roster` */

DROP TABLE IF EXISTS `tim_roster`;

CREATE TABLE `tim_roster` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `loginname` varchar(64) NOT NULL COMMENT '用户标识',
  `username` varchar(64) NOT NULL COMMENT '用户登录名',
  `rostername` varchar(64) NOT NULL COMMENT '关系用户登陆名',
  `rostertype` varchar(64) NOT NULL DEFAULT '' COMMENT '关系类型',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  `remarknick` varchar(64) NOT NULL COMMENT '备注名',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_loginname` (`loginname`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='花名册表';

/*Table structure for table `tim_user` */

DROP TABLE IF EXISTS `tim_user`;

CREATE TABLE `tim_user` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `loginname` varchar(64) NOT NULL COMMENT '用户标识',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `usernick` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `plainpassword` varchar(64) NOT NULL DEFAULT '' COMMENT '密码明文',
  `encryptedpassword` varchar(64) NOT NULL DEFAULT '' COMMENT '加密密码',
  `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间',
  `updatetime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_loginname` (`loginname`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='用户表';
