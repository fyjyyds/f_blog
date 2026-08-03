<template>
  <div class="popular-articles-section">
    <div class="section-header">
      <h3 class="section-title">🔥 热门文章</h3>
      <div class="update-info" v-if="lastUpdate">
        <span class="update-text">最近更新</span>
        <span class="update-time">{{ formatDate(lastUpdate) }}</span>
      </div>
    </div>

    <div v-if="articles.length > 0" class="articles-container">
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
            <span class="author">{{ article.author?.nickname || article.author?.username || '匿名' }}</span>
            <span class="divider">·</span>
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
              :style="{ backgroundColor: tag.color || '#5b6abf' }"
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
      <p>暂无热门文章</p>
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

const fetchPopularArticles = async () => {
  loading.value = true;
  try {
    const response = await getPopularArticles(8);
    articles.value = response.data.data;
    lastUpdate.value = response.data.updated_at;
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取热门文章失败');
  } finally {
    loading.value = false;
  }
};

const calculatePopularityScore = (article: any) => {
  return (article.view_count || 0) * 1 + (article.like_count || 0) * 3 + (article.comment_count || 0) * 2;
};

const getRankClass = (index: number) => {
  if (index === 0) return 'rank-gold';
  if (index === 1) return 'rank-silver';
  if (index === 2) return 'rank-bronze';
  return '';
};

const goToArticle = (id: number) => {
  router.push(`/article/${id}`);
};

onMounted(() => {
  fetchPopularArticles();
});
</script>

<style scoped>
.popular-articles-section {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0;
}

.update-text {
  font-size: 11px;
  color: #9ca3af;
}

.articles-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.article-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border: 1px solid #f3f4f6;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.article-item:hover {
  background: #f9fafb;
}

.rank-badge {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  flex-shrink: 0;
  background: #f3f4f6;
  color: #6b7280;
}

.rank-gold {
  background: #fef3c7;
  color: #92400e;
}

.rank-silver {
  background: #e5e7eb;
  color: #374151;
}

.rank-bronze {
  background: #fed7aa;
  color: #9a3412;
}

.article-content {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 14px;
  font-weight: 500;
  color: #1a1a2e;
  margin: 0 0 6px 0;
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
  margin-bottom: 8px;
  font-size: 12px;
  color: #9ca3af;
}

.author {
  color: #5b6abf;
}

.article-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 12px;
  color: #9ca3af;
}

.popularity-score {
  color: #f59e0b;
}

.article-tags {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
  margin-top: 8px;
}

.tag {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 10px;
  color: #fff;
}

.more-tags {
  font-size: 10px;
  color: #9ca3af;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 20px;
  gap: 12px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #e5e7eb;
  border-top: 2px solid #5b6abf;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  color: #9ca3af;
  font-size: 13px;
  margin: 0;
}

.empty-state {
  text-align: center;
  padding: 32px 20px;
  color: #9ca3af;
  font-size: 13px;
}
</style>
