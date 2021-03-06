# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.21)
# Database: gim
# Generation Time: 2018-05-21 03:50:13 +0000
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
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增',
  `sender` int(10) NOT NULL COMMENT '发送人(用户ID)',
  `receiver` int(10) NOT NULL COMMENT '接收人(可以是用户，群，讨论组)',
  `content` varbinary(4096) NOT NULL DEFAULT '' COMMENT '保存为二进制流',
  `time_stamp` int(10) NOT NULL COMMENT '发送日期',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1发送成功',
  `update_at` int(10) DEFAULT NULL COMMENT '状态修改时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



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




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
