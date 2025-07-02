-- 初始化文章点赞数（如果还没有初始化过）
UPDATE articles SET like_count = (
  SELECT COUNT(*) FROM likes 
  WHERE target_type = 'article' AND target_id = articles.id
) WHERE like_count = 0;

-- 初始化评论点赞数（如果还没有初始化过）
UPDATE comments SET like_count = (
  SELECT COUNT(*) FROM likes 
  WHERE target_type = 'comment' AND target_id = comments.id
) WHERE like_count = 0;

-- 确保所有已发布文章都有正确的状态
UPDATE articles SET status = 'published' WHERE status IS NULL AND publish_time IS NOT NULL;

-- 为热门文章添加索引以提高查询性能（MySQL 8.0+标准写法，无DESC/WHERE）
CREATE INDEX IF NOT EXISTS idx_articles_popularity ON articles (view_count, like_count, comment_count);
CREATE INDEX IF NOT EXISTS idx_articles_status ON articles (status); 