<template>
  <div class="popular-articles-section">
    <div class="section-header">
      <h3 class="section-title">
        <span class="title-icon">🔥</span>
        热门文章
      </h3>
      <div class="update-info">
        <span class="update-text">每日更新</span>
        <span class="update-time" v-if="lastUpdate">{{ lastUpdate }}</span>
      </div>
    </div>
    
    <div class="articles-container" v-if="articles.length > 0">
      <div 
        v-for="(article, index) in articles" 
        :key="article.id"
        class="article-item"
        @click="goToArticle(article.id)"
      >
        <div class="rank-badge" :class="getRankClass(index)">
          {{ index + 1 }}
        </div>
        
        <div class="article-content">
          <h4 class="article-title">{{ article.title }}</h4>
          <div class="article-meta">
            <span class="author">{{ article.author?.nickname || article.author?.username }}</span>
            <span class="divider">•</span>
            <span class="date">{{ formatDate(article.created_at) }}</span>
          </div>
          
          <div class="article-stats">
            <div class="stat-item">
              <span class="stat-icon">👁️</span>
              <span class="stat-value">{{ article.view_count || 0 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">❤️</span>
              <span class="stat-value">{{ article.like_count || 0 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">💬</span>
              <span class="stat-value">{{ article.comment_count || 0 }}</span>
            </div>
            <div class="stat-item popularity-score">
              <span class="stat-icon">⭐</span>
              <span class="stat-value">{{ calculatePopularityScore(article) }}</span>
            </div>
          </div>
          
          <div class="article-tags" v-if="article.tags && article.tags.length">
            <span 
              v-for="tag in article.tags.slice(0, 2)" 
              :key="tag.id"
              class="tag"
              :style="{ backgroundColor: tag.color || '#667eea' }"
            >
              {{ tag.name }}
            </span>
            <span v-if="article.tags.length > 2" class="more-tags">
              +{{ article.tags.length - 2 }}
            </span>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p class="loading-text">加载热门文章...</p>
    </div>
    
    <div v-else class="empty-state">
      <div class="empty-icon">📝</div>
      <p class="empty-text">暂无热门文章</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { getPopularArticles } from '../api/article';
import { formatDate } from '../utils/date';

const router = useRouter();
const articles = ref<any[]>([]);
const loading = ref(false);
const lastUpdate = ref('');

// 获取热门文章
const fetchPopularArticles = async () => {
  loading.value = true;
  try {
    const response = await getPopularArticles(8); // 获取8篇热门文章
    articles.value = response.data.data;
    lastUpdate.value = response.data.updated_at;
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取热门文章失败');
  } finally {
    loading.value = false;
  }
};

// 计算热度分数
const calculatePopularityScore = (article: any) => {
  const viewScore = (article.view_count || 0) * 1;
  const likeScore = (article.like_count || 0) * 3;
  const commentScore = (article.comment_count || 0) * 2;
  return viewScore + likeScore + commentScore;
};

// 获取排名样式
const getRankClass = (index: number) => {
  if (index === 0) return 'rank-gold';
  if (index === 1) return 'rank-silver';
  if (index === 2) return 'rank-bronze';
  return 'rank-normal';
};

// 跳转到文章详情
const goToArticle = (id: number) => {
  router.push(`/article/${id}`);
};

onMounted(() => {
  fetchPopularArticles();
});
</script>

<style scoped>
.popular-articles-section {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 24px;
  backdrop-filter: blur(20px);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  color: white;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-icon {
  font-size: 24px;
}

.update-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
}

.update-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.update-time {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.articles-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.article-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.3), transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.article-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

.article-item:hover::before {
  opacity: 1;
}

.rank-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  color: white;
  flex-shrink: 0;
}

.rank-gold {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: #333;
}

.rank-silver {
  background: linear-gradient(135deg, #c0c0c0, #e5e5e5);
  color: #333;
}

.rank-bronze {
  background: linear-gradient(135deg, #cd7f32, #daa520);
  color: white;
}

.rank-normal {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.article-content {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 16px;
  font-weight: 600;
  color: white;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 12px;
}

.author {
  color: #667eea;
  font-weight: 500;
}

.divider {
  color: rgba(255, 255, 255, 0.3);
}

.date {
  color: rgba(255, 255, 255, 0.5);
}

.article-stats {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.stat-icon {
  font-size: 14px;
}

.stat-value {
  font-weight: 600;
}

.popularity-score {
  color: #ffd700;
  font-weight: 700;
}

.article-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.tag {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 10px;
  font-weight: 500;
  color: white;
  background: rgba(102, 126, 234, 0.2);
}

.more-tags {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  gap: 16px;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  gap: 12px;
}

.empty-icon {
  font-size: 32px;
  opacity: 0.5;
}

.empty-text {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  margin: 0;
}

@media (max-width: 768px) {
  .popular-articles-section {
    padding: 16px;
  }
  
  .section-header {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
  
  .article-item {
    padding: 12px;
    gap: 12px;
  }
  
  .article-stats {
    gap: 12px;
  }
  
  .rank-badge {
    width: 28px;
    height: 28px;
    font-size: 12px;
  }
  
  .article-title {
    font-size: 14px;
  }
}
</style> 