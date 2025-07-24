<template>
  <div class="home-container">
    <Navbar />
    <BannerCarousel />
    
    <div class="main-content">
      <div class="content-wrapper">
        <!-- 左侧主体 -->
        <div class="main-left">
          <!-- 筛选栏 -->
          <div class="filter-section card">
            <div class="filter-header">
              <h3 class="section-title">🔍 文章筛选</h3>
              <div class="filter-controls">
                <div class="filter-group">
                  <label class="filter-label">分类</label>
                  <select v-model="categoryId" class="filter-select" @change="() => fetchArticles()">
                    <option value="">全部分类</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                      {{ cat.name }}
                    </option>
                  </select>
                </div>
                <div class="filter-group">
                  <label class="filter-label">标签</label>
                  <div class="tag-selector">
                    <div 
                      v-for="tag in tags" 
                      :key="tag.id"
                      class="tag-option"
                      :class="{ active: tagIds.includes(tag.id) }"
                      @click="toggleTag(tag.id)"
                    >
                      {{ tag.name }}
                    </div>
                  </div>
                </div>
                <div class="filter-group">
                  <label class="filter-label">排序</label>
                  <div class="sort-buttons">
                    <button 
                      class="sort-btn"
                      :class="{ active: sortType === 'new' }"
                      @click="setSort('new')"
                    >
                      🕒 最新
                    </button>
                    <button 
                      class="sort-btn"
                      :class="{ active: sortType === 'hot' }"
                      @click="setSort('hot')"
                    >
                      🔥 最热
                    </button>
                    <button 
                      class="sort-btn"
                      :class="{ active: sortType === 'comment' }"
                      @click="setSort('comment')"
                    >
                      💬 评论
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- 文章列表 -->
          <div class="articles-section">
            <div class="articles-header">
              <h2 class="articles-title">
                📚 文章列表 
                <span v-if="hasFilters" class="filter-indicator">
                  (筛选结果: {{ total }})
                </span>
                <span v-else>
                  (共 {{ total }} 篇)
                </span>
              </h2>
              <div v-if="hasFilters" class="filter-summary">
                <span class="filter-item" v-if="categoryId">
                  📂 {{ getCategoryName(categoryId) }}
                </span>
                <span class="filter-item" v-if="tagIds.length > 0">
                  🏷️ {{ getTagNames(tagIds).join(', ') }}
                </span>
                <span class="filter-item" v-if="sortType !== 'new'">
                  🔄 {{ getSortName(sortType) }}
                </span>
              </div>
            </div>
            <div class="articles-grid">
              <div 
                v-for="article in articles" 
                :key="article.id"
                class="article-card card"
                @click="goDetail(article.id)"
              >
                <div class="article-header">
                  <h3 class="article-title">{{ article.title }}</h3>
                  <div class="article-meta">
                    <span class="meta-item">👁️ {{ article.view_count }}</span>
                    <span class="meta-item">💬 {{ article.comment_count }}</span>
                    <span class="meta-item">❤️ {{ article.like_count || 0 }}</span>
                  </div>
                </div>
                <div class="article-summary">{{ article.summary }}</div>
                <div class="article-footer">
                  <div class="article-tags" v-if="article.tags && article.tags.length">
                    <span 
                      v-for="tag in article.tags" 
                      :key="tag.id"
                      class="article-tag"
                    >
                      #{{ tag.name }}
                    </span>
                  </div>
                  <div class="article-actions">
                    <button 
                      class="like-btn"
                      :class="{ 'liked': articleLiked[article.id] }"
                      @click.stop="toggleArticleLike(article.id)"
                      :disabled="!isLoggedIn"
                    >
                      <span class="like-icon">{{ articleLiked[article.id] ? '❤️' : '🤍' }}</span>
                      <span class="like-count">{{ article.like_count || 0 }}</span>
                    </button>
                    <div class="article-date">
                      {{ formatDate(article.created_at) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="articles.length === 0" class="empty-state">
              <div class="empty-icon">📝</div>
              <h3 class="empty-title">暂无文章</h3>
              <p class="empty-desc">尝试调整筛选条件或稍后再来查看</p>
            </div>
          </div>
        </div>
        <!-- 右侧侧栏 -->
        <div class="main-right">
          <PopularArticles />
        </div>
      </div>
    </div>
    <!-- 固定右下角分页组件 -->
    <div class="fixed-pagination" v-if="totalPages > 1">
      <button :disabled="page === 1" @click="changePage(page-1)">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页</span>
      <button :disabled="page === totalPages" @click="changePage(page+1)">下一页</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import Navbar from '../components/Navbar.vue';
import BannerCarousel from './BannerCarousel.vue';
import PopularArticles from '../components/PopularArticles.vue';
import { ref, onMounted, computed } from 'vue';
import { listArticles } from '../api/article';
import { useRouter } from 'vue-router';
import request from '../api/request';
import { getLikeStatus, like, unlike } from '../api/like';
import { ElMessage } from 'element-plus';

const articles = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(8); // 默认8
const totalPages = computed(() => Math.ceil(total.value / pageSize.value));
const sortType = ref('new');
const router = useRouter();
const categories = ref<any[]>([]);
const categoryId = ref('');
const tags = ref<any[]>([]);
const tagIds = ref<any[]>([]);

// 检测是否有筛选条件
const hasFilters = computed(() => {
  return categoryId.value || tagIds.value.length > 0 || sortType.value !== 'new';
});

// 点赞相关
const articleLiked = ref<Record<number, boolean>>({});
const isLoggedIn = computed(() => !!localStorage.getItem('token'));

onMounted(async () => {
  await fetchContentSettings();
  fetchCategories();
  fetchTags();
  fetchArticles();
});

async function fetchContentSettings() {
  try {
    const res = await request.get('/api/v1/admin/settings');
    if (res.data && res.data.content && res.data.content.articlesPerPage) {
      pageSize.value = parseInt(res.data.content.articlesPerPage, 10) || 8;
    }
  } catch (e) {
    // 忽略错误，使用默认
  }
}

async function fetchCategories() {
  try {
    const res = await request.get('/api/v1/categories');
    categories.value = res.data;
  } catch (error) {
    console.error('获取分类失败:', error);
  }
}

async function fetchTags() {
  try {
    const res = await request.get('/api/v1/tags');
    tags.value = res.data;
  } catch (error) {
    console.error('获取标签失败:', error);
  }
}

async function fetchArticles(p = page.value) {
  try {
    const params: any = { sort: sortType.value, page: p, page_size: pageSize.value };
    if (categoryId.value) params.category_id = categoryId.value;
    if (tagIds.value && tagIds.value.length > 0) params.tag_ids = tagIds.value.join(',');
    const query = Object.entries(params)
      .map(([k, v]) => `${k}=${encodeURIComponent(String(v))}`)
      .join('&');
    const res = await listArticles(query ? `?${query}` : '');
    articles.value = res.data.data;
    total.value = res.data.total;
    page.value = res.data.page;
    // 获取所有文章的点赞状态
    if (isLoggedIn.value) {
      for (const article of articles.value) {
        await fetchArticleLikeStatus(article.id);
      }
    }
  } catch (error) {
    console.error('获取文章失败:', error);
  }
}

function changePage(p: number) {
  if (p < 1 || p > totalPages.value) return;
  page.value = p;
  fetchArticles(p);
}

function toggleTag(tagId: number) {
  const index = tagIds.value.indexOf(tagId);
  if (index > -1) {
    tagIds.value.splice(index, 1);
  } else {
    tagIds.value.push(tagId);
  }
  page.value = 1;
  fetchArticles(1);
}

function setSort(type: string) {
  sortType.value = type;
  page.value = 1;
  fetchArticles(1);
}

function goDetail(id: number) {
  router.push(`/article/${id}`);
}

// 获取文章点赞状态
async function fetchArticleLikeStatus(articleId: number) {
  if (!isLoggedIn.value) return;
  
  try {
    const response = await getLikeStatus('article', articleId);
    articleLiked.value[articleId] = response.data.liked;
  } catch (error) {
    console.error('获取文章点赞状态失败:', error);
  }
}

// 切换文章点赞状态
async function toggleArticleLike(articleId: number) {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录');
    return;
  }
  
  try {
    if (articleLiked.value[articleId]) {
      await unlike('article', articleId);
    } else {
      await like('article', articleId);
    }
    articleLiked.value[articleId] = !articleLiked.value[articleId];
    
    // 更新文章列表中的点赞数
    const article = articles.value.find(a => a.id === articleId);
    if (article) {
      if (articleLiked.value[articleId]) {
        article.like_count = (article.like_count || 0) + 1;
      } else {
        article.like_count = Math.max(0, (article.like_count || 0) - 1);
      }
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败');
  }
}

function formatDate(dateString: string) {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}

// 获取分类名称
function getCategoryName(categoryId: string) {
  const category = categories.value.find(cat => cat.id.toString() === categoryId);
  return category ? category.name : '未知分类';
}

// 获取标签名称列表
function getTagNames(tagIds: number[]) {
  return tagIds.map(id => {
    const tag = tags.value.find(tag => tag.id === id);
    return tag ? tag.name : '未知标签';
  });
}

// 获取排序方式名称
function getSortName(sortType: string) {
  const sortNames: Record<string, string> = {
    'new': '最新',
    'hot': '最热',
    'comment': '评论最多'
  };
  return sortNames[sortType] || '最新';
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-wrapper {
  display: flex;
  gap: 32px;
  align-items: flex-start;
}

.main-left {
  flex: 1 1 0;
  min-width: 0;
}

.main-right {
  width: 340px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.filter-section {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 24px;
}

.filter-header {
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.filter-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  align-items: flex-start;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 200px;
}

.filter-label {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.filter-select {
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 10px 12px;
  color: white;
  font-size: 14px;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.filter-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.filter-select option {
  background: #1a1a2e;
  color: white;
}

.tag-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-option {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  padding: 6px 12px;
  font-size: 12px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.tag-option:hover {
  background: rgba(102, 126, 234, 0.2);
  border-color: #667eea;
  transform: translateY(-1px);
}

.tag-option.active {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: #667eea;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.sort-buttons {
  display: flex;
  gap: 8px;
}

.sort-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 8px 12px;
  color: white;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.sort-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  border-color: #667eea;
  transform: translateY(-1px);
}

.sort-btn.active {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: #667eea;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.articles-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.articles-header {
  padding: 0 8px;
}

.articles-title {
  font-size: 24px;
  font-weight: 700;
  color: white;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-indicator {
  color: #667eea;
  font-size: 20px;
  font-weight: 600;
}

.filter-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 8px;
  padding: 0 8px;
}

.filter-item {
  background: rgba(102, 126, 234, 0.2);
  color: #667eea;
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid rgba(102, 126, 234, 0.3);
  backdrop-filter: blur(10px);
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
}

.article-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.article-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, #667eea, transparent);
}

.article-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  border-color: rgba(102, 126, 234, 0.3);
}

.article-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 16px;
}

.article-title {
  font-size: 18px;
  font-weight: 700;
  color: white;
  line-height: 1.4;
  margin: 0;
  flex: 1;
}

.article-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex-shrink: 0;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  backdrop-filter: blur(10px);
}

.article-summary {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin-bottom: 16px;
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.article-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.like-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 6px 10px;
  color: white;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.like-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}

.like-btn.liked {
  background: rgba(255, 107, 107, 0.2);
  border-color: #ff6b6b;
  color: #ff6b6b;
}

.like-btn.liked:hover:not(:disabled) {
  background: rgba(255, 107, 107, 0.3);
}

.like-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.like-icon {
  font-size: 14px;
}

.like-count {
  font-weight: 600;
  font-size: 11px;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}

.article-tag {
  background: rgba(102, 126, 234, 0.2);
  color: #667eea;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.article-date {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  flex-shrink: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: white;
  margin: 0 0 8px 0;
}

.empty-desc {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  margin: 0;
}

.fixed-pagination {
  position: fixed;
  right: 40px;
  bottom: 40px;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(40, 40, 60, 0.95);
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0,0,0,0.15);
  padding: 12px 24px;
}
.fixed-pagination button {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  color: white;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}
.fixed-pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
@media (max-width: 900px) {
  .content-wrapper {
    flex-direction: column;
    gap: 0;
  }
  .main-right {
    width: 100%;
    margin-top: 24px;
  }
  .fixed-pagination {
    right: 10px;
    bottom: 10px;
    padding: 8px 12px;
    font-size: 12px;
  }
}
</style>