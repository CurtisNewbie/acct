-- Initialize schema

create database acct;

use acct;

CREATE TABLE `cashflow` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no',
  `direction` varchar(6) NOT NULL DEFAULT '' COMMENT 'flow direction: IN / OUT',
  `trans_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'transaction time',
  `trans_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'transaction id',
  `counterparty` varchar(255) DEFAULT '' COMMENT 'counterparty',
  `amount` decimal(22,8) DEFAULT '0.00000000' COMMENT 'amount',
  `currency` varchar(6) DEFAULT '' COMMENT 'currency',
  `extra` json DEFAULT NULL COMMENT 'extra info about the transaction',
  `category` varchar(128) NOT NULL DEFAULT '' COMMENT 'category name',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT 'remark',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  `deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'record deleted',
  PRIMARY KEY (`id`),
  KEY `user_cate_trans_time_idx` (`user_no`,`category`,`trans_time`,`deleted`),
  KEY `user_trans_time_idx` (`user_no`,`trans_time`,`deleted`),
  KEY `user_trans_id_idx` (`user_no`,`trans_id`,`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Cash flow';

