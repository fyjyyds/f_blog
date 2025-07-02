/*
 Navicat Premium Dump SQL

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80042 (8.0.42)
 Source Host           : localhost:3306
 Source Schema         : f_blog

 Target Server Type    : MySQL
 Target Server Version : 80042 (8.0.42)
 File Encoding         : 65001

 Date: 02/07/2025 16:41:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article_tags
-- ----------------------------
DROP TABLE IF EXISTS `article_tags`;
CREATE TABLE `article_tags`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `article_id` bigint UNSIGNED NULL DEFAULT NULL,
  `tag_id` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_article_tag`(`article_id` ASC, `tag_id` ASC) USING BTREE,
  INDEX `idx_article_id`(`article_id` ASC) USING BTREE,
  INDEX `idx_tag_id`(`tag_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章标签关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article_tags
-- ----------------------------
INSERT INTO `article_tags` VALUES (1, 12, 1, '2025-06-30 23:15:42.000');
INSERT INTO `article_tags` VALUES (2, 12, 2, '2025-06-30 23:15:42.000');
INSERT INTO `article_tags` VALUES (3, 12, 4, '2025-06-30 23:15:42.000');
INSERT INTO `article_tags` VALUES (4, 13, 1, '2025-06-30 23:26:20.000');
INSERT INTO `article_tags` VALUES (5, 13, 2, '2025-06-30 23:26:20.000');
INSERT INTO `article_tags` VALUES (6, 13, 3, '2025-06-30 23:26:20.000');
INSERT INTO `article_tags` VALUES (7, 13, 4, '2025-06-30 23:26:20.000');
INSERT INTO `article_tags` VALUES (8, 13, 5, '2025-06-30 23:26:20.000');
INSERT INTO `article_tags` VALUES (9, 14, 2, '2025-07-01 18:25:38.000');
INSERT INTO `article_tags` VALUES (10, 14, 3, '2025-07-01 18:25:38.000');
INSERT INTO `article_tags` VALUES (11, 14, 4, '2025-07-01 18:25:38.000');
INSERT INTO `article_tags` VALUES (12, 15, 1, '2025-07-01 22:31:11.000');
INSERT INTO `article_tags` VALUES (13, 15, 2, '2025-07-01 22:31:11.000');
INSERT INTO `article_tags` VALUES (14, 15, 3, '2025-07-01 22:31:11.000');
INSERT INTO `article_tags` VALUES (18, 17, 6, '2025-07-02 14:56:55.638');
INSERT INTO `article_tags` VALUES (19, 17, 7, '2025-07-02 14:56:55.642');
INSERT INTO `article_tags` VALUES (20, 18, 4, '2025-07-02 15:05:02.669');
INSERT INTO `article_tags` VALUES (21, 18, 5, '2025-07-02 15:05:02.673');
INSERT INTO `article_tags` VALUES (22, 18, 7, '2025-07-02 15:05:02.674');
INSERT INTO `article_tags` VALUES (23, 19, 2, '2025-07-02 15:28:10.240');
INSERT INTO `article_tags` VALUES (24, 19, 5, '2025-07-02 15:28:10.241');

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `author_id` bigint UNSIGNED NULL DEFAULT NULL,
  `category_id` bigint UNSIGNED NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `view_count` bigint UNSIGNED NULL DEFAULT NULL,
  `like_count` bigint UNSIGNED NULL DEFAULT NULL,
  `comment_count` bigint UNSIGNED NULL DEFAULT NULL,
  `publish_time` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_author_id`(`author_id` ASC) USING BTREE,
  INDEX `idx_category_id`(`category_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_publish_time`(`publish_time` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_view_count`(`view_count` ASC) USING BTREE,
  INDEX `idx_like_count`(`like_count` ASC) USING BTREE,
  INDEX `idx_articles_popularity`(`view_count` ASC, `like_count` ASC, `comment_count` ASC) USING BTREE,
  INDEX `idx_articles_status`(`status` ASC) USING BTREE,
  FULLTEXT INDEX `ft_title_content`(`title`, `content`),
  CONSTRAINT `fk_articles_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_articles_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of articles
-- ----------------------------
INSERT INTO `articles` VALUES (2, '我的第一篇文章', '这是文章内容，支持Markdown。', '文章摘要', '/static/cover/cover1.jpg', 1, 1, 'published', 17, 2, 5, '2025-06-28 15:33:15.000', '2025-06-28 15:33:15.000', '2025-07-02 12:01:36.448');
INSERT INTO `articles` VALUES (5, 'sdas ', 'sadasdasdas', 'asdasas', '/static/upload/file_1751186263616133900.png', 2, 3, 'published', 10, 1, 18, '2025-06-29 16:37:44.000', '2025-06-29 16:37:44.000', '2025-07-02 13:25:25.980');
INSERT INTO `articles` VALUES (6, 'fds', 'sadasdsa', 'sdasdsad', '', 4, 1, 'published', 8, 3, 2, '2025-06-29 18:24:21.000', '2025-06-29 18:24:21.000', '2025-07-02 12:49:55.323');
INSERT INTO `articles` VALUES (7, '地方说', '啊实打实的', '实打实', '', 5, 2, 'published', 3, 4, 2, '2025-06-30 11:14:28.000', '2025-06-30 11:14:28.000', '2025-07-02 12:00:27.118');
INSERT INTO `articles` VALUES (8, 'sad', 'asdas', 'asdas', '/static/upload/file_1751271502349756000.txt', 6, 1, 'published', 0, 2, 5, '2025-06-30 16:18:24.000', '2025-06-30 16:18:24.000', '2025-06-30 17:55:49.000');
INSERT INTO `articles` VALUES (9, 'sdaas', 'asdas', 'asdas', '', 2, 2, 'published', 13, 1, 0, '2025-06-30 16:42:27.000', '2025-06-30 16:42:27.000', '2025-07-02 12:49:25.293');
INSERT INTO `articles` VALUES (10, '顶顶顶顶顶', '撒旦发射点发射点', '阿松大', '/static/upload/file_1751286804912684500.py', 2, 2, 'published', 2, 0, 0, '2025-06-30 20:33:27.000', '2025-06-30 20:33:27.000', '2025-06-30 22:42:39.383');
INSERT INTO `articles` VALUES (11, 'fsdf', 'sadasd', 'asdas', '/static/upload/file_1751294884680230500.png', 2, 1, 'published', 7, 0, 0, '2025-06-30 22:48:05.768', '2025-06-30 22:48:05.768', '2025-07-02 14:53:33.669');
INSERT INTO `articles` VALUES (13, 'sdfs', 'saadasd', 'asdas', '/static/upload/file_1751297179534870700.jpg', 2, 1, 'published', 12, 1, 5, '2025-06-30 23:26:20.387', '2025-06-30 23:26:20.387', '2025-07-02 15:30:25.037');
INSERT INTO `articles` VALUES (15, 'dfgd', 'dsdasd', 'fdg', '/static/upload/file_1751380259154640100.docx', 15, 2, 'published', 18, 3, 0, '2025-07-01 22:31:10.745', '2025-07-01 22:31:10.745', '2025-07-02 12:01:08.066');
INSERT INTO `articles` VALUES (16, '盛大收购德国', '撒发射点发生', '胜多负少', '/static/upload/file_1751435132168856100.txt', 17, 1, 'published', 6, 2, 0, '2025-07-02 13:45:34.179', '2025-07-02 13:45:34.179', '2025-07-02 15:27:25.286');
INSERT INTO `articles` VALUES (17, 'asd', 'asd', 'ad', '', 19, 2, 'published', 1, 0, 0, '2025-07-02 14:56:55.635', '2025-07-02 14:56:55.635', '2025-07-02 14:56:55.675');
INSERT INTO `articles` VALUES (18, 'df', 'hjg', 'sdf', '/static/upload/file_1751439900402187800.doc', 19, 2, 'published', 1, 0, 0, '2025-07-02 15:05:56.832', '2025-07-02 15:05:02.664', '2025-07-02 15:05:56.832');
INSERT INTO `articles` VALUES (19, 'aa', 'aa', 'aaa', '', 19, 1, 'published', 0, 0, 0, '2025-07-02 15:29:13.947', '2025-07-02 15:28:10.235', '2025-07-02 15:29:13.947');

-- ----------------------------
-- Table structure for audit_logs
-- ----------------------------
DROP TABLE IF EXISTS `audit_logs`;
CREATE TABLE `audit_logs`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '操作者ID',
  `action` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `target_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标类型',
  `target_id` bigint UNSIGNED NOT NULL COMMENT '目标ID',
  `old_data` json NULL COMMENT '操作前数据',
  `new_data` json NULL COMMENT '操作后数据',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '用户代理',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_action`(`action` ASC) USING BTREE,
  INDEX `idx_target`(`target_type` ASC, `target_id` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '审核日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of audit_logs
-- ----------------------------

-- ----------------------------
-- Table structure for banners
-- ----------------------------
DROP TABLE IF EXISTS `banners`;
CREATE TABLE `banners`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '轮播图ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `image` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `link` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `sort_order` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_sort`(`sort` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '轮播图表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of banners
-- ----------------------------
INSERT INTO `banners` VALUES (1, '欢迎使用F_Blog', '/static/banners/banner1.jpg', '/', 1, 'active', '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `banners` VALUES (2, '技术分享平台', '/static/banners/banner2.jpg', '/articles', 2, 'active', '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `banners` VALUES (3, '汽车人变身', '/static/upload/file_1751425389624351500.jpg', '', 0, 'active', '2025-07-02 11:03:11.420', '2025-07-02 11:03:11.420', 0);

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `parent_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '父分类ID',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `article_count` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章数量',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `sort_order` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_sort`(`sort` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO `categories` VALUES (1, '技术', '技术相关文章', 0, 1, 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `categories` VALUES (2, '生活', '生活相关文章', 0, 2, 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `categories` VALUES (3, 'Go语言', 'Go语言相关', 1, 1, 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `categories` VALUES (4, 'Web开发', 'Web开发相关', 1, 2, 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `categories` VALUES (5, '数据库', '数据库相关', 1, 3, 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `categories` VALUES (18, '前端', 'vue技术', 0, 0, 'active', 0, '2025-07-02 10:21:42.622', '2025-07-02 10:21:42.622', 0);

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `article_id` bigint UNSIGNED NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `parent_id` bigint UNSIGNED NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `like_count` bigint UNSIGNED NULL DEFAULT NULL,
  `reply_count` bigint UNSIGNED NULL DEFAULT NULL,
  `reply_to_user` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_article_id`(`article_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE,
  CONSTRAINT `fk_comments_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 40 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comments
-- ----------------------------
INSERT INTO `comments` VALUES (2, 2, 3, 0, '这是一条评论', 'approved', 1, 0, NULL, '2025-06-29 16:55:19.000', '2025-07-02 12:01:21.593');
INSERT INTO `comments` VALUES (3, 5, 2, 0, 'sdkashdasjkd', 'approved', 0, 0, NULL, '2025-06-29 17:11:29.000', '2025-06-29 17:11:29.000');
INSERT INTO `comments` VALUES (4, 5, 2, 0, 'fsdskfhds', 'approved', 0, 0, NULL, '2025-06-29 17:11:43.000', '2025-06-29 17:11:43.000');
INSERT INTO `comments` VALUES (5, 5, 2, 3, 'sdas', 'approved', 0, 0, NULL, '2025-06-29 17:11:49.000', '2025-06-29 17:11:49.000');
INSERT INTO `comments` VALUES (6, 5, 2, 5, 'cdf', 'approved', 0, 0, NULL, '2025-06-29 17:11:53.000', '2025-06-29 17:11:53.000');
INSERT INTO `comments` VALUES (7, 5, 3, 6, 'asdas', 'approved', 0, 0, NULL, '2025-06-29 17:12:10.000', '2025-06-29 17:12:10.000');
INSERT INTO `comments` VALUES (8, 5, 3, 6, 'sdasd', 'approved', 0, 0, NULL, '2025-06-29 17:14:27.000', '2025-06-29 17:14:27.000');
INSERT INTO `comments` VALUES (9, 5, 3, 7, 'sdasda', 'approved', 0, 0, NULL, '2025-06-29 17:14:31.000', '2025-06-29 17:14:31.000');
INSERT INTO `comments` VALUES (10, 5, 2, 7, 'saadasd', 'approved', 0, 0, NULL, '2025-06-29 17:21:07.000', '2025-06-29 17:21:07.000');
INSERT INTO `comments` VALUES (11, 5, 2, 3, 'fs', 'approved', 0, 0, NULL, '2025-06-29 17:31:08.000', '2025-06-29 17:31:08.000');
INSERT INTO `comments` VALUES (13, 5, 2, 3, 'asdasd', 'approved', 0, 0, NULL, '2025-06-29 17:41:27.000', '2025-06-29 17:41:27.000');
INSERT INTO `comments` VALUES (14, 5, 2, 3, 'sadas', 'approved', 0, 0, NULL, '2025-06-29 17:41:29.000', '2025-06-29 17:41:29.000');
INSERT INTO `comments` VALUES (15, 5, 2, 3, 'sadsa', 'approved', 0, 0, NULL, '2025-06-29 17:41:31.000', '2025-06-29 17:41:31.000');
INSERT INTO `comments` VALUES (16, 5, 2, 3, 'sedasd', 'approved', 0, 0, NULL, '2025-06-29 17:44:36.000', '2025-06-29 17:44:36.000');
INSERT INTO `comments` VALUES (17, 5, 3, 3, 'sdasd', 'approved', 0, 0, NULL, '2025-06-29 17:45:03.000', '2025-06-29 17:45:03.000');
INSERT INTO `comments` VALUES (18, 5, 3, 3, 'dasd', 'approved', 0, 0, NULL, '2025-06-29 17:48:56.000', '2025-06-29 17:48:56.000');
INSERT INTO `comments` VALUES (19, 5, 3, 3, '@2 asdas', 'approved', 0, 0, NULL, '2025-06-29 17:53:20.000', '2025-06-29 17:53:20.000');
INSERT INTO `comments` VALUES (20, 5, 2, 3, 'dasd', 'approved', 0, 0, '2', '2025-06-29 18:03:35.000', '2025-06-29 18:03:35.000');
INSERT INTO `comments` VALUES (21, 2, 3, 0, '这是一条评论', 'approved', 0, 0, '', '2025-06-30 10:21:08.000', '2025-06-30 10:21:08.000');
INSERT INTO `comments` VALUES (22, 2, 3, 1, '回复内容', 'approved', 0, 0, '', '2025-06-30 10:21:42.000', '2025-06-30 10:21:42.000');
INSERT INTO `comments` VALUES (23, 6, 5, 0, '实打实打算', 'approved', 0, 0, '', '2025-06-30 11:09:29.000', '2025-06-30 11:09:29.000');
INSERT INTO `comments` VALUES (24, 6, 5, 23, '士大夫大师傅但是', 'approved', 0, 0, '', '2025-06-30 11:09:40.000', '2025-06-30 11:09:40.000');
INSERT INTO `comments` VALUES (25, 8, 6, 0, 'asdas', 'approved', 0, 0, '', '2025-06-30 16:18:43.000', '2025-06-30 16:18:43.000');
INSERT INTO `comments` VALUES (26, 8, 6, 25, 'dasd', 'approved', 0, 0, '', '2025-06-30 16:18:46.000', '2025-06-30 16:18:46.000');
INSERT INTO `comments` VALUES (27, 8, 6, 25, 'sdasd', 'approved', 0, 0, '6', '2025-06-30 16:21:25.000', '2025-06-30 16:21:25.000');
INSERT INTO `comments` VALUES (28, 8, 6, 25, 'asdsa', 'approved', 0, 0, '', '2025-06-30 16:21:30.000', '2025-06-30 16:21:30.000');
INSERT INTO `comments` VALUES (29, 8, 2, 25, 'asda', 'approved', 0, 0, '6', '2025-06-30 16:21:51.000', '2025-06-30 16:21:51.000');
INSERT INTO `comments` VALUES (30, 7, 2, 0, 'dfsdf ', 'approved', 0, 0, '', '2025-06-30 16:29:13.000', '2025-06-30 16:29:13.000');
INSERT INTO `comments` VALUES (31, 7, 2, 30, 'fdsf', 'approved', 0, 0, '', '2025-06-30 16:29:20.000', '2025-06-30 16:29:20.000');
INSERT INTO `comments` VALUES (32, 5, 2, 3, 'dasd', 'approved', 0, 0, 'ffff', '2025-06-30 16:30:01.000', '2025-06-30 16:30:01.000');
INSERT INTO `comments` VALUES (33, 13, 15, 0, 'asdsad', 'approved', 1, 0, '', '2025-07-01 22:31:36.318', '2025-07-02 11:59:51.875');
INSERT INTO `comments` VALUES (34, 13, 15, 33, 'asdasd', 'approved', 2, 0, '', '2025-07-01 22:31:40.509', '2025-07-02 15:30:26.754');
INSERT INTO `comments` VALUES (35, 13, 15, 33, 'asdas', 'approved', 1, 0, 'qazz', '2025-07-01 22:31:45.226', '2025-07-02 12:49:44.038');
INSERT INTO `comments` VALUES (36, 13, 15, 33, 'fsdf', 'approved', 0, 0, '', '2025-07-01 22:31:51.994', '2025-07-01 22:31:51.994');
INSERT INTO `comments` VALUES (37, 13, 15, 33, 'sdfsd', 'approved', 0, 0, 'qazz', '2025-07-01 22:32:00.564', '2025-07-01 22:32:00.564');
INSERT INTO `comments` VALUES (38, 2, 17, 2, 'asdas', 'approved', 0, 0, '', '2025-07-02 12:01:31.184', '2025-07-02 12:01:31.184');
INSERT INTO `comments` VALUES (39, 2, 17, 2, 'asdas', 'approved', 0, 0, 'admin', '2025-07-02 12:01:36.443', '2025-07-02 12:01:36.443');

-- ----------------------------
-- Table structure for comments_copy1
-- ----------------------------
DROP TABLE IF EXISTS `comments_copy1`;
CREATE TABLE `comments_copy1`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `article_id` bigint UNSIGNED NOT NULL COMMENT '文章ID',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '评论者ID',
  `parent_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '父评论ID',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `status` enum('pending','approved','rejected') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '评论状态',
  `like_count` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `reply_count` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '回复数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_article_id`(`article_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comments_copy1
-- ----------------------------

-- ----------------------------
-- Table structure for follows
-- ----------------------------
DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关注ID',
  `follower_id` bigint UNSIGNED NULL DEFAULT NULL,
  `following_id` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_follower_following`(`follower_id` ASC, `following_id` ASC) USING BTREE,
  INDEX `idx_follower_id`(`follower_id` ASC) USING BTREE,
  INDEX `idx_following_id`(`following_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '关注表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of follows
-- ----------------------------
INSERT INTO `follows` VALUES (1, 3, 2, '2025-06-29 10:09:09.000');
INSERT INTO `follows` VALUES (2, 3, 3, '2025-06-30 10:09:29.000');
INSERT INTO `follows` VALUES (3, 3, 1, '2025-06-30 10:09:36.000');
INSERT INTO `follows` VALUES (6, 5, 2, '2025-06-30 10:35:00.000');
INSERT INTO `follows` VALUES (7, 5, 1, '2025-06-29 10:35:08.000');
INSERT INTO `follows` VALUES (9, 2, 5, '2025-06-30 16:30:56.000');
INSERT INTO `follows` VALUES (12, 2, 6, '2025-06-30 16:42:10.000');

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `target_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `target_id` bigint UNSIGNED NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_target`(`user_id` ASC, `target_type` ASC, `target_id` ASC) USING BTREE,
  INDEX `idx_target`(`target_type` ASC, `target_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '点赞表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of likes
-- ----------------------------
INSERT INTO `likes` VALUES (20, 2, 'article', 6, '2025-06-30 16:17:24.000');
INSERT INTO `likes` VALUES (21, 3, 'article', 7, '2025-06-30 16:17:35.000');
INSERT INTO `likes` VALUES (22, 6, 'article', 8, '2025-06-30 16:18:41.000');
INSERT INTO `likes` VALUES (23, 2, 'article', 7, '2025-06-30 16:29:22.000');
INSERT INTO `likes` VALUES (26, 2, 'article', 9, '2025-06-30 17:32:10.000');
INSERT INTO `likes` VALUES (27, 2, 'article', 8, '2025-06-30 17:32:18.000');
INSERT INTO `likes` VALUES (28, 2, 'article', 5, '2025-06-30 17:32:21.000');
INSERT INTO `likes` VALUES (30, 0, 'article', 13, '2025-07-01 22:32:29.598');
INSERT INTO `likes` VALUES (32, 0, 'article', 2, '2025-07-01 22:33:18.239');
INSERT INTO `likes` VALUES (33, 0, 'article', 15, '2025-07-01 22:37:38.262');
INSERT INTO `likes` VALUES (35, 16, 'article', 15, '2025-07-01 22:42:54.240');
INSERT INTO `likes` VALUES (36, 16, 'article', 6, '2025-07-01 22:43:06.386');
INSERT INTO `likes` VALUES (37, 16, 'article', 7, '2025-07-01 22:43:17.043');
INSERT INTO `likes` VALUES (38, 15, 'article', 7, '2025-07-01 22:44:12.209');
INSERT INTO `likes` VALUES (39, 15, 'article', 2, '2025-07-01 22:44:17.861');
INSERT INTO `likes` VALUES (41, 17, 'comment', 33, '2025-07-02 11:59:51.870');
INSERT INTO `likes` VALUES (45, 17, 'article', 15, '2025-07-02 12:00:39.123');
INSERT INTO `likes` VALUES (49, 17, 'comment', 2, '2025-07-02 12:01:21.586');
INSERT INTO `likes` VALUES (50, 17, 'comment', 34, '2025-07-02 12:49:35.798');
INSERT INTO `likes` VALUES (51, 17, 'comment', 35, '2025-07-02 12:49:44.033');
INSERT INTO `likes` VALUES (52, 17, 'article', 6, '2025-07-02 12:49:55.318');
INSERT INTO `likes` VALUES (53, 19, 'article', 16, '2025-07-02 14:53:29.327');
INSERT INTO `likes` VALUES (54, 17, 'article', 16, '2025-07-02 15:27:25.282');
INSERT INTO `likes` VALUES (55, 19, 'comment', 34, '2025-07-02 15:30:26.749');

-- ----------------------------
-- Table structure for notifications
-- ----------------------------
DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知ID',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '接收用户ID',
  `type` enum('comment','like','follow','system','review') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知类型',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知内容',
  `data` json NULL COMMENT '通知数据',
  `is_read` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已读',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_type`(`type` ASC) USING BTREE,
  INDEX `idx_is_read`(`is_read` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通知表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of notifications
-- ----------------------------
INSERT INTO `notifications` VALUES (1, 19, 'review', '文章审核通过', '你的文章《aa》审核通过', '{\"status\": \"published\", \"article_id\": 19}', 1, '2025-07-02 15:29:14', '2025-07-02 15:29:53');
INSERT INTO `notifications` VALUES (2, 15, 'like', '你的评论被点赞', '有人点赞了你的评论', '{\"comment_id\": 34}', 0, '2025-07-02 15:30:27', '2025-07-02 15:30:27');

-- ----------------------------
-- Table structure for settings
-- ----------------------------
DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'string',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_settings_key`(`key` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of settings
-- ----------------------------
INSERT INTO `settings` VALUES (1, 'siteKeywords', '技术分享；博客；交流', 'string', '2025-07-02 11:21:38.295', '2025-07-02 13:34:07.414');
INSERT INTO `settings` VALUES (2, 'siteLogo', '', 'string', '2025-07-02 11:21:38.300', '2025-07-02 13:34:07.417');
INSERT INTO `settings` VALUES (3, 'icp', '', 'string', '2025-07-02 11:21:38.303', '2025-07-02 13:34:07.402');
INSERT INTO `settings` VALUES (4, 'siteName', 'F_Blog', 'string', '2025-07-02 11:21:38.308', '2025-07-02 13:34:07.406');
INSERT INTO `settings` VALUES (5, 'siteDescription', '一个现代化的博客系统', 'string', '2025-07-02 11:21:38.312', '2025-07-02 13:34:07.411');
INSERT INTO `settings` VALUES (6, 'email_smtpHost', 'smtp.example.com', 'string', '2025-07-02 11:21:38.316', '2025-07-02 11:21:38.316');
INSERT INTO `settings` VALUES (7, 'email_smtpPort', '587', 'string', '2025-07-02 11:21:38.320', '2025-07-02 11:21:38.320');
INSERT INTO `settings` VALUES (8, 'email_emailUser', 'admin@example.com', 'string', '2025-07-02 11:21:38.323', '2025-07-02 11:21:38.323');
INSERT INTO `settings` VALUES (9, 'email_emailPassword', '', 'string', '2025-07-02 11:21:38.326', '2025-07-02 11:21:38.326');
INSERT INTO `settings` VALUES (10, 'email_senderName', 'F_Blog', 'string', '2025-07-02 11:21:38.330', '2025-07-02 11:21:38.330');
INSERT INTO `settings` VALUES (11, 'security_minPasswordLength', '8', 'string', '2025-07-02 11:21:38.334', '2025-07-02 13:34:07.429');
INSERT INTO `settings` VALUES (12, 'security_maxLoginAttempts', '5', 'string', '2025-07-02 11:21:38.337', '2025-07-02 13:34:07.432');
INSERT INTO `settings` VALUES (13, 'security_lockoutDuration', '30', 'string', '2025-07-02 11:21:38.341', '2025-07-02 13:34:07.436');
INSERT INTO `settings` VALUES (14, 'security_jwtExpireHours', '24', 'string', '2025-07-02 11:21:38.345', '2025-07-02 13:34:07.439');
INSERT INTO `settings` VALUES (15, 'security_enableCaptcha', 'true', 'string', '2025-07-02 11:21:38.349', '2025-07-02 13:34:07.420');
INSERT INTO `settings` VALUES (16, 'security_enableTwoFactor', 'true', 'string', '2025-07-02 11:21:38.352', '2025-07-02 13:34:07.425');
INSERT INTO `settings` VALUES (17, 'content_enableCommentLikes', 'true', 'string', '2025-07-02 11:21:38.356', '2025-07-02 13:34:07.454');
INSERT INTO `settings` VALUES (18, 'content_enableArticleLikes', 'true', 'string', '2025-07-02 11:21:38.360', '2025-07-02 13:34:07.458');
INSERT INTO `settings` VALUES (19, 'content_articlesPerPage', '8', 'string', '2025-07-02 11:21:38.364', '2025-07-02 13:34:07.462');
INSERT INTO `settings` VALUES (20, 'content_summaryLength', '200', 'string', '2025-07-02 11:21:38.368', '2025-07-02 13:34:07.443');
INSERT INTO `settings` VALUES (21, 'content_commentModeration', 'manual', 'string', '2025-07-02 11:21:38.371', '2025-07-02 13:34:07.447');
INSERT INTO `settings` VALUES (22, 'content_allowAnonymousComments', 'false', 'string', '2025-07-02 11:21:38.374', '2025-07-02 13:34:07.450');

-- ----------------------------
-- Table structure for tags
-- ----------------------------
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `article_count` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章数量',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `color` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '标签表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tags
-- ----------------------------
INSERT INTO `tags` VALUES (1, 'Go', 'Go语言', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `tags` VALUES (2, 'Golang', 'Go语言', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `tags` VALUES (3, 'Web', 'Web开发', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `tags` VALUES (4, 'MySQL', 'MySQL数据库', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `tags` VALUES (5, 'Redis', 'Redis缓存', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `tags` VALUES (6, 'Vue', 'Vue.js', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);
INSERT INTO `tags` VALUES (7, 'React', 'React.js', 'active', 0, '2025-06-28 12:58:20.000', '2025-06-28 12:58:20.000', NULL);

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色描述',
  `permissions` json NULL COMMENT '权限列表',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_roles
-- ----------------------------

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `nickname` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `bio` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `avatar` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `gender` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `birthday` datetime(3) NULL DEFAULT NULL,
  `website` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `github` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `twitter` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `linkedin` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `role` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'user',
  `status` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'active',
  `email_verified` tinyint(1) NULL DEFAULT 0,
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `email_verify_token` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `email_verify_expire` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `uk_email`(`email` ASC) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `email`(`email` ASC) USING BTREE,
  INDEX `idx_role`(`role` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'testuser', 'testuser@example.com', '$2a$10$48y5zddEogyizUSCot5RmeDr0AczRjVa.cWrQ0DACA2XssekHOBFW', '新昵称', '新的个人简介', '/static/avatar/avatar_1.jpg', 'female', '1999-12-31 00:00:00.000', 'https://myblog.com', 'mygithub', 'mytwitter', 'mylinkedin', 'admin', 'active', 0, NULL, '2025-06-28 13:15:52.000', '2025-06-28 15:18:09.000', '', NULL);
INSERT INTO `users` VALUES (2, 'fan', '@qq.com', '$2a$10$X36fTS9CBUmGNqoeduTTaOqfKZ34tMES1pCS4b00wFK2a5IyPK4PC', 'qwe', 'qwewq', '/static/avatar/avatar_2.png', 'male', '2025-06-03 00:00:00.000', 'sada', NULL, NULL, NULL, 'user', 'active', 0, NULL, '2025-06-28 13:32:39.000', '2025-06-30 22:49:25.778', '99d850b3935162d7421859f7fe9a48db', NULL);
INSERT INTO `users` VALUES (3, 'fanx', '@qq.com', '$2a$10$EEQRBldAbKiUQhV60LaESO0y/.bTRgpJbFYkpBvwyEY3F6eZdvk1u', 'ffff', '新简介', '/static/avatar/avatar_3.jpg', 'male', '2000-01-01 00:00:00.000', 'https://example.com', NULL, NULL, NULL, 'user', 'active', 0, NULL, '2025-06-28 13:37:37.000', '2025-06-30 09:51:02.000', NULL, NULL);
INSERT INTO `users` VALUES (4, 'addd', 'f14@qq.com', '$2a$10$V3mjJxEUDIxvdS.DhpGzOe2fBbGxrDE4B.adjPuEss3p1ZYTiC58e', '', '', '', 'other', NULL, '', '', '', '', 'user', 'active', 0, NULL, '2025-06-29 18:23:48.000', '2025-06-29 18:23:48.000', NULL, NULL);
INSERT INTO `users` VALUES (5, 'lkk', '12456@qq.com', '$2a$10$009t.a0.WYPNBJR0VKIec.RRcFgYogxpgdzom4QdP1.WQ4MNcOa3u', 'fky', 'sadsad', '/static/avatar/avatar_5.png', 'male', '2002-06-07 00:00:00.000', 'asdasd', '', '', '', 'user', 'active', 0, NULL, '2025-06-30 10:30:38.000', '2025-06-30 11:00:37.000', NULL, NULL);
INSERT INTO `users` VALUES (6, 'qwe', '1235678@qq.com', '$2a$10$pIIbIoVsge.icSEDrsrel.IfIhuMeD3JvjKndQKJFea1y7NCRjoPq', '', '', '', 'other', NULL, '', '', '', '', 'user', 'active', 0, NULL, '2025-06-30 16:18:01.000', '2025-06-30 16:18:01.000', NULL, NULL);
INSERT INTO `users` VALUES (7, 'ddasd', '12389@qq.com', '$2a$10$D3hlDPukhlHhepIDRtg.4e6.y7haIKrOpj.ennqe9/eQsv.nx.HKm', '', '', '', 'other', NULL, '', '', '', '', 'user', 'active', 0, NULL, '2025-06-30 22:41:56.877', '2025-06-30 22:41:56.877', NULL, NULL);
INSERT INTO `users` VALUES (10, 'rxy', 'qwer@163.com', '$2a$10$602cCwL4siXdmTadK23WzehTj0a86cFQWaik85gEaPEcg48jGYZPi', '', '', '', '', NULL, '', '', '', '', 'user', 'pending', 0, NULL, '2025-07-01 18:27:08.635', '2025-07-01 18:27:08.635', 'bd886168cb7506d24525f802d94f86e8', '2025-07-02 18:27:08.635');
INSERT INTO `users` VALUES (11, 'zyt', '16511@qq.com', '$2a$10$njamR.DNQztQGv5/thfw0uD9aDDn2ckNcTWEoUEBsC0u18BJaxa.m', '', '', '', '', NULL, '', '', '', '', 'user', 'pending', 0, NULL, '2025-07-01 18:42:51.758', '2025-07-01 18:42:51.758', 'a7886279547f5247829e92467fd5aabb', '2025-07-02 18:42:51.758');
INSERT INTO `users` VALUES (12, 'tfc', '11631@qq.com', '$2a$10$v2HFVh48WylPZrDB16Sv/.PZPYcDzHX.NfrWUxnIuRyBaRuOmlnlK', '', '', '', '', NULL, '', '', '', '', 'user', 'pending', 0, NULL, '2025-07-01 18:57:31.069', '2025-07-01 18:57:31.069', 'efb5ceb0bf19f3ab2158bb43eebd873a', '2025-07-02 18:57:31.069');
INSERT INTO `users` VALUES (13, 'wslk', '@qq.com', '$2a$10$amvx7nR82X/o7kiyFt83fepmJg4.PpAkOnAJ4pOi5raql9phADcYS', '', '', '', '', NULL, '', '', '', '', 'user', 'pending', 0, NULL, '2025-07-01 19:05:13.420', '2025-07-01 19:05:13.420', '99d850b3935162d7421859f7fe9a48db', '2025-07-02 19:05:13.420');
INSERT INTO `users` VALUES (14, 'qaz', '161@qq.com', '$2a$10$RbT3MfCqX1Ld/FVY2TGzX.dV7O8/7CNFG6Ai4h5kFK0o8OcNSN6oe', '', '', '', '', NULL, '', '', '', '', 'user', 'pending', 0, NULL, '2025-07-01 19:25:59.791', '2025-07-01 19:25:59.791', '00da01da786104c78ab275a6efb8600f', '2025-07-02 19:25:59.791');
INSERT INTO `users` VALUES (15, 'qazz', '11@qq.com', '$2a$10$7WOYID7hb98KlN9wb/ymdu7XN3xcS4VOV6mr.A6Z8xE5TqceiB4wa', '', '', '', '', NULL, '', '', '', '', 'user', 'active', 1, NULL, '2025-07-01 19:27:36.773', '2025-07-01 22:56:37.532', '', NULL);
INSERT INTO `users` VALUES (16, 'godf', 'f11@163.com', '$2a$10$feI2FcWRRjCqpPPmwn/G6O/Qtl8ALA4ahf8WIYFaVIwTg2wS0880e', 'asdas', 'asdasd', '/static/avatar/avatar_16.jpg', 'female', '2025-07-31 08:00:00.000', 'asd', '', '', '', 'user', 'active', 1, NULL, '2025-07-01 22:36:30.114', '2025-07-01 22:37:26.859', '', NULL);
INSERT INTO `users` VALUES (17, 'admin', 'admin@example.com', '$2a$10$PI0pRk.BsKX8pHcw5qcQre3erMBLQepDxV5d.hm7SuOSoI/ifILe6', '', '', '/static/avatar/avatar_17.jpg', '', NULL, '', '', '', '', 'admin', 'active', 1, NULL, '2025-07-01 23:20:03.537', '2025-07-02 12:01:53.064', '', NULL);
INSERT INTO `users` VALUES (18, 'system', 'system@example.com', 'dummy', '系统', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'admin', 'active', 1, NULL, '2025-07-02 13:23:42.000', '2025-07-02 13:23:42.000', NULL, NULL);
INSERT INTO `users` VALUES (19, 'faf', 'f161@163.com', '$2a$10$KS.e04zShpFStQQiG00N3OGZS4i/F4.Hu15FVqBMScXytD9SfxvOS', 'dg', 'asda', '/static/avatar/avatar_19.jpg', 'male', '2025-07-30 08:00:00.000', 'sadsa', '', '', '', 'user', 'active', 1, NULL, '2025-07-02 14:40:03.536', '2025-07-02 14:41:20.720', '', NULL);

SET FOREIGN_KEY_CHECKS = 1;
