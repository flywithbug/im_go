# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.7.21)
# Database: gim
# Generation Time: 2018-05-12 23:59:07 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table im_device
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_device`;

CREATE TABLE `im_device` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `u_id` int(10) NOT NULL,
  `device_token` varchar(64) NOT NULL DEFAULT '' COMMENT '设备token',
  `device_id` varchar(64) NOT NULL DEFAULT '' COMMENT '设备uuid',
  `platform` tinyint(1) NOT NULL DEFAULT '1' COMMENT '设备类型(iphone android,web)',
  `description` varchar(64) NOT NULL DEFAULT '' COMMENT '设备描述',
  `status` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `device_id` (`device_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table im_login
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_login`;

CREATE TABLE `im_login` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '数据库自增',
  `u_id` int(10) NOT NULL,
  `user_id` varchar(40) NOT NULL DEFAULT '' COMMENT '用户ID',
  `app_id` mediumint(6) NOT NULL DEFAULT '10',
  `token` varchar(40) NOT NULL DEFAULT '' COMMENT '用户token',
  `login_at` datetime NOT NULL COMMENT '登录日期',
  `login_ip` varchar(32) NOT NULL COMMENT '用户登录IP',
  `logout_at` datetime DEFAULT NULL,
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '登录态：-1 被踢出 ，0登出，1使用中,',
  `forbidden` char(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table im_message
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_message`;

CREATE TABLE `im_message` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '数据库自增',
  `sender` int(10) NOT NULL COMMENT '发送人(用户ID)',
  `receiver` int(10) NOT NULL COMMENT '接收人(可以是用户，群，讨论组)',
  `content` varbinary(4096) NOT NULL DEFAULT '' COMMENT '保存为二进制流',
  `time_stamp` int(10) NOT NULL COMMENT '发送日期',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '消息状态 0对方离线，未发送，1送达，2已读，3 撤回 4 删除',
  `update_at` int(10) DEFAULT NULL COMMENT '状态修改时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `im_message` WRITE;
/*!40000 ALTER TABLE `im_message` DISABLE KEYS */;

INSERT INTO `im_message` (`id`, `sender`, `receiver`, `content`, `time_stamp`, `status`, `update_at`)
VALUES
	(10186,10002,10001,X'6D736737',1526009065,0,NULL),
	(10199,10002,10001,X'6D73673131353236313031333036',1526101306,0,NULL),
	(10185,10002,10001,X'6D736736',1526009064,0,NULL),
	(10184,10002,10001,X'6D736735',1526009062,0,NULL),
	(10198,10002,10001,X'6D73673131353236313031323938',1526101298,0,NULL),
	(10182,10002,10001,X'6D736733',1526009025,0,NULL),
	(10183,10002,10001,X'6D736734',1526009044,0,NULL),
	(10180,10002,10001,X'6D736731',1526009014,0,NULL),
	(10181,10002,10001,X'6D736732',1526009023,0,NULL),
	(10187,10002,10001,X'6D736738',1526009067,0,NULL),
	(10197,10002,10001,X'6D73673131353236313030393433',1526100943,0,NULL);

/*!40000 ALTER TABLE `im_message` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table im_relationship
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_relationship`;

CREATE TABLE `im_relationship` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '数据库自增',
  `u_id` int(10) NOT NULL COMMENT '用户Id',
  `relation_id` varchar(40) NOT NULL DEFAULT '' COMMENT 'u_id+friend_id',
  `friend_id` int(11) NOT NULL COMMENT '好友申请人',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 好友申请,1 好友关系，-1 黑名单:  ',
  `remark` varchar(20) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `relation_id` (`relation_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `im_relationship` WRITE;
/*!40000 ALTER TABLE `im_relationship` DISABLE KEYS */;

INSERT INTO `im_relationship` (`id`, `u_id`, `relation_id`, `friend_id`, `status`, `remark`)
VALUES
	(10017,10001,'10001-10002',10002,0,''),
	(10020,10001,'0b83e3965f12affa4371beaa267c3071',10002,0,'');

/*!40000 ALTER TABLE `im_relationship` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table im_room
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_room`;

CREATE TABLE `im_room` (
  `id` int(10) NOT NULL COMMENT '群的唯一标识',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '群名称',
  `creator` int(10) NOT NULL COMMENT '创建者 user_id',
  `create_at` datetime NOT NULL COMMENT '创建日期',
  `user_num` int(5) NOT NULL DEFAULT '100' COMMENT '群允许的用户数量',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table im_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_user`;

CREATE TABLE `im_user` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '唯一标识,自增',
  `user_id` varchar(40) NOT NULL,
  `account` varchar(64) NOT NULL DEFAULT '' COMMENT '账号',
  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '密码-md5',
  `nick` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `sign` varchar(255) NOT NULL DEFAULT '' COMMENT '个人签名',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '状态 0离线，1在线，2离开，3请勿打扰，4忙碌，5Q我吧，6隐身',
  `create_at` datetime NOT NULL COMMENT '注册日期',
  `update_at` datetime NOT NULL COMMENT '更新日期',
  `app_id` int(6) NOT NULL COMMENT '所属应用编号',
  `forbidden` char(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `account` (`account`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `im_user` WRITE;
/*!40000 ALTER TABLE `im_user` DISABLE KEYS */;

INSERT INTO `im_user` (`id`, `user_id`, `account`, `password`, `nick`, `sign`, `avatar`, `status`, `create_at`, `update_at`, `app_id`, `forbidden`)
VALUES
	(10001,'d5f75fbc-4f64-4f78-b320-2ca770847800','ori','16b1c83de8f9518e673838b2d6ea75dc','ori','','http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg',0,'2018-05-02 20:20:59','2018-05-02 20:20:59',10,'0'),
	(10002,'80d07eb7-09d5-4332-aa4d-01990a291dfd','ori1','16b1c83de8f9518e673838b2d6ea75dc','ori','','http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg',0,'2018-05-03 16:16:26','2018-05-03 16:16:26',10,'0'),
	(10003,'fe9609b4-d22e-4672-b27c-1c13a3849f37','ori2','16b1c83de8f9518e673838b2d6ea75dc','ori','','http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg',0,'2018-05-03 16:18:02','2018-05-03 16:18:02',10,'0'),
	(10004,'df417cac-15e1-430a-8858-27cec879f267','ori3','16b1c83de8f9518e673838b2d6ea75dc','ori','','http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg',0,'2018-05-04 18:22:04','2018-05-04 18:22:04',10,'0');

/*!40000 ALTER TABLE `im_user` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
