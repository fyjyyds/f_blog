-- 为comments表添加reply_to_user字段
ALTER TABLE `comments` ADD COLUMN `reply_to_user` VARCHAR(100) DEFAULT NULL COMMENT '回复给哪个用户' AFTER `reply_count`; 
ALTER TABLE `comments` ADD COLUMN `reply_to_user` VARCHAR(100) DEFAULT NULL COMMENT '回复给哪个用户' AFTER `reply_count`; 