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
  `category` varchar(32) NOT NULL DEFAULT '' COMMENT 'category',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT 'remark',
  `payment_method` varchar(32) NOT NULL DEFAULT '' COMMENT 'payment method',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  `deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'record deleted',
  PRIMARY KEY (`id`),
  KEY `user_cate_trans_time_idx` (`user_no`,`category`,`deleted`,`trans_time`),
  KEY `user_trans_time_idx` (`user_no`,`deleted`,`trans_time`),
  KEY `user_cate_trans_id_idx` (`user_no`,`category`,`trans_id`,`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Cashflow';

CREATE TABLE `cashflow_statistics` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no',
  `agg_type` varchar(10) NOT NULL DEFAULT '' COMMENT 'aggregation type: MONTHLY, YEARLY, WEEKLY',
  `agg_range` varchar(10) NOT NULL DEFAULT '' COMMENT 'aggregation range, year or month value',
  `agg_value` decimal(22,8) DEFAULT '0.00000000' COMMENT 'amount',
  `currency` varchar(6) DEFAULT '' COMMENT 'currency',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted',
  PRIMARY KEY (`id`),
  KEY `user_agg_type_range_idx` (`user_no`, `agg_type`, `agg_range`, `currency`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Cashflow Statistics';