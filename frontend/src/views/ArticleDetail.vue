<template>
  <div class="article-detail-container" v-if="article">
    <Navbar />
    
    <div class="article-content">
      <div class="article-card card">
        <div class="article-header">
          <div class="header-top">
            <button class="back-btn" @click="goBack">
              <span class="back-icon">←</span>
              <span class="back-text">返回</span>
            </button>
            <div class="action-bar" v-if="canEditOrDelete">
              <button class="action-btn edit-btn" @click="goEdit">
                <span class="btn-icon">✏️</span>
                编辑
              </button>
              <button class="action-btn delete-btn" @click="onDelete">
                <span class="btn-icon">🗑️</span>
                删除
              </button>
            </div>
          </div>
          
          <h1 class="article-title">{{ article.title }}</h1>
          
          <div class="article-meta">
            <div class="meta-item">
              <span class="meta-icon">👤</span>
              <span class="meta-label">作者：</span>
              <span class="meta-value">{{ article.author.nickname || article.author.username }}</span>
              <button
                v-if="showFollowBtn"
                class="follow-btn"
                :class="{ following: isFollowing }"
                @click="toggleFollow"
              >
                <span class="follow-icon">{{ isFollowing ? '✓' : '+' }}</span>
                {{ isFollowing ? '已关注' : '关注' }}
              </button>
            </div>
            
            <div class="meta-item">
              <span class="meta-icon">📂</span>
              <span class="meta-label">分类：</span>
              <span class="meta-value">{{ categoryName }}</span>
            </div>
            
            <div class="meta-item">
              <span class="meta-icon">📅</span>
              <span class="meta-label">发布时间：</span>
              <span class="meta-value">{{ formatDate(article.publish_time || article.created_at) }}</span>
            </div>
            
            <div class="meta-item">
              <span class="meta-icon">👁️</span>
              <span class="meta-label">浏览量：</span>
              <span class="meta-value">{{ article.view_count || 0 }}</span>
            </div>
          </div>
        </div>
        
        <div class="article-body">
          <div class="cover-section" v-if="article.cover">
            <img v-if="isImage(article.cover)" :src="fullUrl(article.cover)" class="cover-img" />
            <a v-else :href="fullUrl(article.cover)" target="_blank" class="file-link">
              <span class="file-icon">📎</span>
              下载附件：{{ getFileName(article.cover) }}
            </a>
          </div>
          
          <div class="content" v-html="article.content"></div>
        </div>
        
        <div class="article-footer">
          <div class="tags-section" v-if="article.tags && article.tags.length">
            <span class="tags-label">标签：</span>
            <div class="tags-list">
              <span
                v-for="tag in article.tags"
                :key="tag.id"
                class="article-tag"
              >
                #{{ tag.name }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <CommentList :article-id="article.id" v-if="article" />
    </div>
    
    <!-- 悬浮点赞按钮 -->
    <div class="floating-like-btn">
      <button
        class="like-btn"
        :class="{ liked: liked }"
        @click="handleLike"
        :disabled="likeLoading"
      >
        <span class="like-icon">{{ liked ? '❤️' : '🤍' }}</span>
        <span class="like-text">{{ liked ? '已点赞' : '点赞' }}</span>
        <span class="like-count">{{ likeCount }}</span>
      </button>
    </div>
    
    <!-- 删除确认对话框 -->
    <div class="modal-overlay" v-if="deleteDialogVisible" @click="deleteDialogVisible = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">确认删除</h3>
        </div>
        <div class="modal-body">
          <p>确定要删除这篇文章吗？此操作不可撤销。</p>
        </div>
        <div class="modal-footer">
          <button class="modal-btn cancel-btn" @click="deleteDialogVisible = false">取消</button>
          <button class="modal-btn confirm-btn" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import CommentList from '../components/CommentList.vue';
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import request from '../api/request';
import { getCategories } from '../api/category';
import { getLikeStatus, like, unlike } from '../api/like';
import { ArrowLeft, Promotion } from '@element-plus/icons-vue';
import { formatDate } from '../utils/date';
import Navbar from '../components/Navbar.vue';

const route = useRoute();
const router = useRouter();
const article = ref<any>(null);
const categories = ref<{ id: number; name: string }[]>([]);
const deleteDialogVisible = ref(false);

const liked = ref(false);
const likeCount = ref(0);
const likeLoading = ref(false);
const isFollowing = ref(false);
const showFollowBtn = computed(() => {
  const token = localStorage.getItem('token');
  if (!token) return false;
  const userId = Number(localStorage.getItem('user_id'));
  return article.value && article.value.author && article.value.author.id !== userId;
});

onMounted(async () => {
  try {
    // 获取文章详情
    const res = await request.get(`/api/v1/articles/${route.params.id}`);
    article.value = res.data.data;
    // 获取分类列表
    const catRes = await getCategories();
    categories.value = catRes.data;
    // 获取点赞状态
    await fetchLikeStatus();
    await fetchIsFollowing();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '加载失败');
  }
});

const fetchLikeStatus = async () => {
  if (!article.value) return;
  const res = await getLikeStatus('article', article.value.id);
  liked.value = res.data.liked;
  likeCount.value = res.data.count;
};

const fetchIsFollowing = async () => {
  if (!showFollowBtn.value) return;
  const res = await request.get('/api/v1/user/followings');
  isFollowing.value = res.data.some((f: any) => f.following_id === article.value.author.id);
};

const handleLike = async () => {
  if (!localStorage.getItem('token')) {
    ElMessage.warning('请先登录');
    return;
  }
  if (likeLoading.value) return;
  likeLoading.value = true;
  try {
    let res;
    if (liked.value) {
      res = await unlike('article', article.value.id);
    } else {
      res = await like('article', article.value.id);
    }
    liked.value = res.data.liked;
    likeCount.value = res.data.count;
    window.dispatchEvent(new Event('refresh-user-stats'));
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '操作失败');
  } finally {
    likeLoading.value = false;
  }
};

const toggleFollow = async () => {
  try {
    if (!isFollowing.value) {
      await request.post('/api/v1/user/follow', { following_id: article.value.author.id });
      isFollowing.value = true;
      ElMessage.success('关注成功');
    } else {
      await request.delete('/api/v1/user/follow', { data: { following_id: article.value.author.id } });
      isFollowing.value = false;
      ElMessage.success('已取消关注');
    }
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '操作失败');
  }
};

const categoryName = computed(() => {
  if (!article.value || !categories.value.length) return '';
  const cat = categories.value.find(c => c.id === article.value.category_id);
  return cat ? cat.name : '';
});

function isImage(url: string) {
  return /\.(jpg|jpeg|png|gif)$/i.test(url);
}
function getFileName(url: string) {
  return url.split('/').pop();
}
function fullUrl(url: string) {
  if (!url) return '';
  if (url.startsWith('http')) return url;
  return 'http://localhost:8080' + url;
}

function goBack() {
  router.back();
}

// 权限判断：作者本人或管理员可编辑/删除
const canEditOrDelete = computed(() => {
  const userId = Number(localStorage.getItem('user_id'));
  const role = localStorage.getItem('role');
  const isAuthor = article.value && article.value.author && article.value.author.id === userId;
  const isAdmin = role === 'admin';
  
  console.log('权限判断:', {
    userId,
    role,
    articleAuthorId: article.value?.author?.id,
    isAuthor,
    isAdmin,
    result: isAuthor || isAdmin
  });
  
  return article.value && (isAuthor || isAdmin);
});

function goEdit() {
  router.push(`/edit-article/${article.value.id}`);
}

function onDelete() {
  deleteDialogVisible.value = true;
}
async function confirmDelete() {
  try {
    await request.delete(`/api/v1/articles/${article.value.id}`);
    ElMessage.success('删除成功');
    router.push('/');
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '删除失败');
  } finally {
    deleteDialogVisible.value = false;
  }
}
</script>

<style scoped>
.article-detail-container {
  min-height: 100vh;
}

.article-content {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
}

.article-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 28px;
  margin-bottom: 24px;
}

.article-header {
  margin-bottom: 24px;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 7px 12px;
  color: #374151;
  font-size: 13px;
  cursor: pointer;
}
.back-btn:hover {
  background: #e5e7eb;
}

.back-icon {
  font-size: 18px;
  font-weight: bold;
}

.action-bar {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.edit-btn {
  background: #5b6abf;
  color: #fff;
}
.edit-btn:hover {
  background: #4a59ae;
}

.delete-btn {
  background: #ef4444;
  color: #fff;
}
.delete-btn:hover {
  background: #dc2626;
}

.btn-icon {
  font-size: 14px;
}

.article-title {
  font-size: 36px;
  font-weight: 800;
  color: white;
  line-height: 1.3;
  margin: 0 0 24px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
}

.meta-label {
  color: #9ca3af;
}

.meta-value {
  color: #374151;
  font-weight: 500;
}

.follow-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 4px 12px;
  color: #374151;
  font-size: 12px;
  cursor: pointer;
  margin-left: 8px;
}
.follow-btn:hover {
  border-color: #5b6abf;
  color: #5b6abf;
}
.follow-btn.following {
  background: #5b6abf;
  border-color: #5b6abf;
  color: #fff;
}

.follow-icon {
  font-size: 12px;
  font-weight: bold;
}

.article-body {
  margin-bottom: 32px;
}

.cover-section {
  margin-bottom: 24px;
}

.cover-img {
  width: 100%;
  max-width: 600px;
  border-radius: 8px;
}

.file-link {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  color: #5b6abf;
  text-decoration: none;
}
.file-link:hover {
  background: #e5e7eb;
}

.content {
  font-size: 15px;
  line-height: 1.8;
  color: #1f2937;
  word-break: break-word;
}

.content :deep(h1),
.content :deep(h2),
.content :deep(h3),
.content :deep(h4),
.content :deep(h5),
.content :deep(h6) {
  color: white;
  margin-top: 24px;
  margin-bottom: 16px;
}

.content :deep(p) {
  margin-bottom: 16px;
}

.content :deep(code) {
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.content :deep(pre) {
  background: rgba(0, 0, 0, 0.3);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
}

.article-footer {
  border-top: 1px solid #e5e7eb;
  padding-top: 20px;
}

.tags-section {
  display: flex;
  align-items: center;
  gap: 10px;
}

.tags-label {
  color: #9ca3af;
  font-size: 13px;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.article-tag {
  background: #eef0ff;
  color: #5b6abf;
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
}

.floating-like-btn {
  position: fixed;
  right: 32px;
  bottom: 32px;
  z-index: 1000;
}

.like-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  padding: 10px 16px;
  color: #374151;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}
.like-btn:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}
.like-btn.liked {
  background: #ef4444;
  border-color: #ef4444;
  color: #fff;
}

.like-count {
  background: #f3f4f6;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-content {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 24px;
  max-width: 400px;
  width: 90%;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 12px 0;
}

.modal-body p {
  color: #6b7280;
  margin: 0;
  line-height: 1.6;
}

.modal-footer {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.modal-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.cancel-btn {
  background: #f3f4f6;
  color: #374151;
  border: 1px solid #e5e7eb;
}
.cancel-btn:hover {
  background: #e5e7eb;
}

.confirm-btn {
  background: #ef4444;
  color: #fff;
}
.confirm-btn:hover {
  background: #dc2626;
}

@media (max-width: 768px) {
  .article-content {
    padding: 16px;
  }
  
  .article-card {
    padding: 20px;
  }
  
  .article-title {
    font-size: 24px;
  }
  
  .article-meta {
    flex-direction: column;
    gap: 12px;
  }
  
  .header-top {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .action-bar {
    width: 100%;
    justify-content: flex-end;
  }
  
  .floating-like-btn {
    right: 20px;
    bottom: 20px;
  }
  
  .like-btn {
    padding: 10px 16px;
    font-size: 12px;
  }
}
</style>
 