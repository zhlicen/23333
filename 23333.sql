/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50505
Source Host           : localhost:3306
Source Database       : app_23333

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2017-03-28 17:23:47
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for pms_albums
-- ----------------------------
DROP TABLE IF EXISTS `23333_user`;
CREATE TABLE `23333_user` (
  `uid` varchar(20) NOT NULL,
  `username` varchar(20) NOT NULL DEFAULT '',
  `mobile` varchar(30) DEFAULT '',
  `email` varchar(30) DEFAULT '',
  `regtime` varchar(30) DEFAULT '',
  `password` varchar(255) DEFAULT '',
  PRIMARY KEY (`uid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='User Table';

