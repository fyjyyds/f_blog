-- 初始化文章点赞数
UPDATE articles SET like_count = (
  SELECT COUNT(*) FROM likes 
  WHERE target_type = 'article' AND target_id = articles.id
);

-- 初始化评论点赞数
UPDATE comments SET like_count = (
  SELECT COUNT(*) FROM likes 
  WHERE target_type = 'comment' AND target_id = comments.id
); 