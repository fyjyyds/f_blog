# F_Blog 数据库设计文档

## 数据库概述

- **数据库类型**: MySQL 8.0+
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_unicode_ci
- **存储引擎**: InnoDB

## 表结构设计

### 1. 用户表 (users)

```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password` varchar(255) NOT NULL COMMENT '密码(加密)',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `bio` text COMMENT '个人简介',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `gender` enum('male','female','other') DEFAULT NULL COMMENT '性别',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `website` varchar(255) DEFAULT NULL COMMENT '个人网站',
  `github` varchar(50) DEFAULT NULL COMMENT 'GitHub用户名',
  `twitter` varchar(50) DEFAULT NULL COMMENT 'Twitter用户名',
  `linkedin` varchar(50) DEFAULT NULL COMMENT 'LinkedIn用户名',
  `role` enum('user','admin') NOT NULL DEFAULT 'user' COMMENT '用户角色',
  `status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '用户状态',
  `email_verified` tinyint(1) NOT NULL DEFAULT '0' COMMENT '邮箱是否验证',
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

### 2. 用户角色表 (user_roles)

```sql
CREATE TABLE `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `permissions` json DEFAULT NULL COMMENT '权限列表',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色表';
```

### 3. 文章表 (articles)

```sql
CREATE TABLE `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `title` varchar(200) NOT NULL COMMENT '文章标题',
  `content` longtext NOT NULL COMMENT '文章内容',
  `summary` text COMMENT '文章摘要',
  `cover` varchar(255) DEFAULT NULL COMMENT '封面图片URL',
  `author_id` bigint unsigned NOT NULL COMMENT '作者ID',
  `category_id` bigint unsigned DEFAULT NULL COMMENT '分类ID',
  `status` enum('draft','published','private','pending','rejected') NOT NULL DEFAULT 'draft' COMMENT '文章状态',
  `view_count` int unsigned NOT NULL DEFAULT '0' COMMENT '阅读量',
  `like_count` int unsigned NOT NULL DEFAULT '0' COMMENT '点赞数',
  `comment_count` int unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
  `publish_time` timestamp NULL DEFAULT NULL COMMENT '发布时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_publish_time` (`publish_time`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_view_count` (`view_count`),
  KEY `idx_like_count` (`like_count`),
  FULLTEXT KEY `ft_title_content` (`title`,`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';
```

### 4. 分类表 (categories)

```sql
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `description` varchar(255) DEFAULT NULL COMMENT '分类描述',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT '父分类ID',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态',
  `article_count` int unsigned NOT NULL DEFAULT '0' COMMENT '文章数量',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';
```

### 5. 标签表 (tags)

```sql
CREATE TABLE `tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `name` varchar(50) NOT NULL COMMENT '标签名称',
  `description` varchar(255) DEFAULT NULL COMMENT '标签描述',
  `status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态',
  `article_count` int unsigned NOT NULL DEFAULT '0' COMMENT '文章数量',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';
```

### 6. 文章标签关联表 (article_tags)

```sql
CREATE TABLE `article_tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `article_id` bigint unsigned NOT NULL COMMENT '文章ID',
  `tag_id` bigint unsigned NOT NULL COMMENT '标签ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_article_tag` (`article_id`,`tag_id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';
```

### 7. 评论表 (comments)

```sql
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `article_id` bigint unsigned NOT NULL COMMENT '文章ID',
  `user_id` bigint unsigned NOT NULL COMMENT '评论者ID',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT '父评论ID',
  `content` text NOT NULL COMMENT '评论内容',
  `status` enum('pending','approved','rejected') NOT NULL DEFAULT 'pending' COMMENT '评论状态',
  `like_count` int unsigned NOT NULL DEFAULT '0' COMMENT '点赞数',
  `reply_count` int unsigned NOT NULL DEFAULT '0' COMMENT '回复数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';
```

### 8. 点赞表 (likes)

```sql
CREATE TABLE `likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `target_type` enum('article','comment') NOT NULL COMMENT '点赞目标类型',
  `target_id` bigint unsigned NOT NULL COMMENT '点赞目标ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_target` (`user_id`,`target_type`,`target_id`),
  KEY `idx_target` (`target_type`,`target_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';
```

### 9. 关注表 (follows)

```sql
CREATE TABLE `follows` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '关注ID',
  `follower_id` bigint unsigned NOT NULL COMMENT '关注者ID',
  `following_id` bigint unsigned NOT NULL COMMENT '被关注者ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_follower_following` (`follower_id`,`following_id`),
  KEY `idx_follower_id` (`follower_id`),
  KEY `idx_following_id` (`following_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关注表';
```

### 10. 通知表 (notifications)

```sql
CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '通知ID',
  `user_id` bigint unsigned NOT NULL COMMENT '接收用户ID',
  `type` enum('comment','like','follow','system','review') NOT NULL COMMENT '通知类型',
  `title` varchar(100) NOT NULL COMMENT '通知标题',
  `content` text NOT NULL COMMENT '通知内容',
  `data` json DEFAULT NULL COMMENT '通知数据',
  `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';
```

### 11. 轮播图表 (banners)

```sql
CREATE TABLE `banners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '轮播图ID',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `image` varchar(255) NOT NULL COMMENT '图片URL',
  `link` varchar(255) DEFAULT NULL COMMENT '链接地址',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='轮播图表';
```

### 12. 审核日志表 (audit_logs)

```sql
CREATE TABLE `audit_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint unsigned NOT NULL COMMENT '操作者ID',
  `action` varchar(50) NOT NULL COMMENT '操作类型',
  `target_type` varchar(50) NOT NULL COMMENT '目标类型',
  `target_id` bigint unsigned NOT NULL COMMENT '目标ID',
  `old_data` json DEFAULT NULL COMMENT '操作前数据',
  `new_data` json DEFAULT NULL COMMENT '操作后数据',
  `ip` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_action` (`action`),
  KEY `idx_target` (`target_type`,`target_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审核日志表';
```

## 索引设计

### 主要索引
1. **用户表**: 用户名、邮箱唯一索引，角色、状态、创建时间普通索引
2. **文章表**: 作者、分类、状态、发布时间、阅读量、点赞数索引，标题内容全文索引
3. **评论表**: 文章、用户、父评论、状态、创建时间索引
4. **点赞表**: 用户目标唯一索引，目标类型ID索引
5. **关注表**: 关注者被关注者唯一索引
6. **通知表**: 用户、类型、已读状态、创建时间索引

### 性能优化
1. 使用合适的字段类型和长度
2. 为常用查询字段创建索引
3. 使用复合索引优化多字段查询
4. 全文索引支持文章搜索
5. 合理使用外键约束（可选）

## 数据初始化

### 默认分类
```sql
INSERT INTO `categories` (`name`, `description`, `parent_id`, `sort`) VALUES
('技术', '技术相关文章', 0, 1),
('生活', '生活相关文章', 0, 2),
('Go语言', 'Go语言相关', 1, 1),
('Web开发', 'Web开发相关', 1, 2),
('数据库', '数据库相关', 1, 3);
```

### 默认标签
```sql
INSERT INTO `tags` (`name`, `description`) VALUES
('Go', 'Go语言'),
('Golang', 'Go语言'),
('Web', 'Web开发'),
('MySQL', 'MySQL数据库'),
('Redis', 'Redis缓存'),
('Vue', 'Vue.js'),
('React', 'React.js');
```

### 默认轮播图
```sql
   INSERT INTO `banners` (`title`, `image`, `link`, `sort`) VALUES
   ('欢迎使用F_Blog', '/static/banners/banner1.jpg', '/', 1),
   ('技术分享平台', '/static/banners/banner2.jpg', '/articles', 2);
```

## 数据库维护

### 定期维护任务
1. **数据备份**: 每日自动备份
2. **索引优化**: 定期分析索引使用情况
3. **数据清理**: 清理过期日志和临时数据
4. **性能监控**: 监控慢查询和连接数

### 数据迁移
1. 使用数据库迁移工具管理表结构变更
2. 保留数据迁移历史记录
3. 测试环境验证迁移脚本
4. 生产环境执行前备份数据 