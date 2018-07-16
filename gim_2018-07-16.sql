# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 118.89.108.25 (MySQL 5.7.22)
# Database: gim
# Generation Time: 2018-07-16 09:10:41 +0000
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
  `user_id` varchar(40) NOT NULL,
  `device_token` varchar(80) NOT NULL DEFAULT '' COMMENT '设备token',
  `device_id` varchar(80) NOT NULL DEFAULT '' COMMENT '设备uuid',
  `user_agent` varchar(64) DEFAULT '' COMMENT '设备描述',
  `platform` tinyint(1) NOT NULL DEFAULT '0' COMMENT '设备类型(iphone android,web),1/2/3',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1. 推送，0 不推送',
  `unique_mac_uuid` varchar(40) DEFAULT '',
  `environment` tinyint(1) NOT NULL DEFAULT '0',
  `sound` int(1) NOT NULL DEFAULT '1',
  `show_detail` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `device_id` (`device_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table im_location
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_location`;

CREATE TABLE `im_location` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` varchar(40) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



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
  `user_Agent` varchar(64) NOT NULL,
  `logout_at` datetime DEFAULT NULL,
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '登录态：-1 被踢出 ，0登出，1使用中,',
  `forbidden` char(1) NOT NULL DEFAULT '0',
  `device_id` varchar(80) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table im_message
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_message`;

CREATE TABLE `im_message` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增',
  `sender` int(10) NOT NULL COMMENT '发送人(用户ID)',
  `receiver` int(10) NOT NULL COMMENT '接收人(可以是用户，群，讨论组)',
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '保存为二进制流',
  `time_stamp` int(10) NOT NULL COMMENT '发送日期',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '1发送成功',
  `update_at` int(10) DEFAULT NULL COMMENT '状态修改时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



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
  `user_id` varchar(40) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `account` varchar(64) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '账号',
  `password` varchar(64) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '密码-md5',
  `nick` varchar(64) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '用户昵称',
  `sign` varchar(255) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '个人签名',
  `avatar` varchar(255) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '状态 0离线，1在线，2离开，3请勿打扰，4忙碌，5Q我吧，6隐身',
  `create_at` datetime NOT NULL COMMENT '注册日期',
  `update_at` datetime NOT NULL COMMENT '更新日期',
  `app_id` int(6) NOT NULL COMMENT '所属应用编号',
  `forbidden` int(1) DEFAULT '0',
  `origin_password` varchar(32) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `longitude` varchar(14) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '经度',
  `latitude` varchar(14) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '维度',
  `l_time_stamp` varchar(14) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '定位时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `account` (`account`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table im_user_authorization
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_user_authorization`;

CREATE TABLE `im_user_authorization` (
  `host_id` varchar(40) NOT NULL DEFAULT '',
  `guest_id` varchar(40) NOT NULL DEFAULT '',
  `a_type` tinyint(4) NOT NULL COMMENT '类型 1 定位授权',
  `status` tinyint(4) NOT NULL COMMENT '0 未授权，1 授权',
  PRIMARY KEY (`host_id`,`guest_id`,`a_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table im_user_location_path
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_user_location_path`;

CREATE TABLE `im_user_location_path` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `p_identifier` varchar(64) NOT NULL DEFAULT '',
  `user_id` varchar(40) NOT NULL DEFAULT '',
  `l_time_stamp` varchar(14) NOT NULL DEFAULT '',
  `latitude` varchar(14) NOT NULL DEFAULT '',
  `longitude` varchar(14) NOT NULL DEFAULT '',
  `l_type` tinyint(2) DEFAULT '0' COMMENT '0 gps 1 图片gps',
  `device_id` varchar(80) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `p_identifier` (`p_identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table im_user_relation
# ------------------------------------------------------------

DROP TABLE IF EXISTS `im_user_relation`;

CREATE TABLE `im_user_relation` (
  `user_id` varchar(40) NOT NULL DEFAULT '',
  `f_user_id` varchar(40) NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL COMMENT '-3拉黑，-2 删除，-1 拒绝，0 申请，1 接受',
  `receiver_id` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`user_id`,`f_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
