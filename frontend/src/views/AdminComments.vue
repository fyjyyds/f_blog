<template>
  <div class="admin-comments">
    <!-- 搜索和过滤栏 -->
    <div class="action-bar">
      <div class="search-section">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索评论内容或用户名..."
          class="search-input"
          @input="handleSearch"
        />
        <button class="search-btn" @click="handleSearch">
          <i class="fas fa-search"></i>
        </button>
      </div>
      
      <div class="filter-section">
        <select v-model="statusFilter" @change="handleFilter" class="filter-select">
          <option value="">全部状态</option>
          <option value="pending">待审核</option>
          <option value="approved">已通过</option>
          <option value="rejected">已拒绝</option>
        </select>
        
        <select v-model="articleFilter" @change="handleFilter" class="filter-select">
          <option value="">全部文章</option>
          <option v-for="article in articles" :key="article.id" :value="article.id">
            {{ article.title }}
          </option>
        </select>
      </div>
    </div>

    <!-- 评论列表 -->
    <div class="comments-list">
      <div 
        v-for="comment in filteredComments" 
        :key="comment.id" 
        class="comment-card"
        :class="comment.status"
      >
        <div class="comment-header">
          <div class="comment-author">
            <img 
              :src="comment.user?.avatar || '/static/avatar/default.jpg'" 
              :alt="comment.user?.username"
              class="author-avatar"
            />
            <div class="author-info">
              <span class="author-name">{{ comment.user?.username || '未知用户' }}</span>
              <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
            </div>
          </div>
          
          <div class="comment-status">
            <span class="status-badge" :class="comment.status">
              {{ getStatusText(comment.status) }}
            </span>
          </div>
        </div>
        
        <div class="comment-content">
          <div class="article-info">
            <span class="article-title">文章：{{ comment.article?.title || '未知文章' }}</span>
          </div>
          
          <div class="comment-text">
            <p v-if="comment.parent_id" class="reply-to">
              回复 @{{ getParentCommentAuthor(comment.parent_id) }}：
            </p>
            <p>{{ comment.content }}</p>
          </div>
        </div>
        
        <div class="comment-actions">
          <button 
            v-if="comment.status === 'pending'" 
            class="btn-approve" 
            @click="approveComment(comment.id)"
          >
            <i class="fas fa-check"></i>
            通过
          </button>
          
          <button 
            v-if="comment.status === 'pending'" 
            class="btn-reject" 
            @click="rejectComment(comment.id)"
          >
            <i class="fas fa-times"></i>
            拒绝
          </button>
          
          <button class="btn-edit" @click="editComment(comment)">
            <i class="fas fa-edit"></i>
            编辑
          </button>
          
          <button class="btn-delete" @click="deleteComment(comment.id)">
            <i class="fas fa-trash"></i>
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <button 
        :disabled="currentPage === 1" 
        @click="changePage(currentPage - 1)"
        class="page-btn"
      >
        <i class="fas fa-chevron-left"></i>
      </button>
      
      <span class="page-info">
        第 {{ currentPage }} 页，共 {{ totalPages }} 页
      </span>
      
      <button 
        :disabled="currentPage === totalPages" 
        @click="changePage(currentPage + 1)"
        class="page-btn"
      >
        <i class="fas fa-chevron-right"></i>
      </button>
    </div>

    <!-- 编辑评论弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑评论</h3>
          <button class="close-btn" @click="closeEditModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>评论内容</label>
            <textarea 
              v-model="editForm.content" 
              class="form-textarea" 
              rows="4"
              placeholder="请输入评论内容"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>状态</label>
            <select v-model="editForm.status" class="form-select">
              <option value="pending">待审核</option>
              <option value="approved">已通过</option>
              <option value="rejected">已拒绝</option>
            </select>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeEditModal">取消</button>
          <button class="btn-save" @click="saveComment">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import request from '../api/request';
import { formatDate } from '../utils/date';

// 响应式数据
const comments = ref<any[]>([]);
const articles = ref<any[]>([]);
const searchQuery = ref('');
const statusFilter = ref('');
const articleFilter = ref('');
const currentPage = ref(1);
const totalPages = ref(1);
const showEditModal = ref(false);
const editForm = ref<any>({});

// 获取评论列表
const fetchComments = async () => {
  try {
    const params = {
      page: currentPage.value,
      limit: 10,
      search: searchQuery.value,
      status: statusFilter.value,
      article_id: articleFilter.value
    };
    
    const response = await request.get('/api/v1/admin/comments', {
      params,
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    comments.value = response.data.comments || [];
    totalPages.value = response.data.totalPages || 1;
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '获取评论列表失败');
  }
};

// 获取文章列表
const fetchArticles = async () => {
  try {
    const response = await request.get('/api/v1/admin/articles', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    articles.value = response.data || [];
  } catch (error: any) {
    console.error('获取文章列表失败:', error);
  }
};

// 过滤评论
const filteredComments = computed(() => {
  let filtered = comments.value;
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(comment => 
      comment.content.toLowerCase().includes(query) ||
      comment.user?.username.toLowerCase().includes(query)
    );
  }
  
  if (statusFilter.value) {
    filtered = filtered.filter(comment => comment.status === statusFilter.value);
  }
  
  if (articleFilter.value) {
    filtered = filtered.filter(comment => comment.article_id === parseInt(articleFilter.value));
  }
  
  return filtered;
});

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1;
  fetchComments();
};

// 过滤处理
const handleFilter = () => {
  currentPage.value = 1;
  fetchComments();
};

// 分页处理
const changePage = (page: number) => {
  currentPage.value = page;
  fetchComments();
};

// 通过评论
const approveComment = async (commentId: number) => {
  try {
    await request.post(`/api/v1/admin/comments/${commentId}/approve`, {}, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('评论已通过');
    fetchComments();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '操作失败');
  }
};

// 拒绝评论
const rejectComment = async (commentId: number) => {
  try {
    await request.post(`/api/v1/admin/comments/${commentId}/reject`, {}, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('评论已拒绝');
    fetchComments();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '操作失败');
  }
};

// 编辑评论
const editComment = (comment: any) => {
  editForm.value = { ...comment };
  showEditModal.value = true;
};

// 保存评论
const saveComment = async () => {
  try {
    await request.put(`/api/v1/admin/comments/${editForm.value.id}`, editForm.value, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('评论更新成功');
    closeEditModal();
    fetchComments();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '更新失败');
  }
};

// 删除评论
const deleteComment = async (commentId: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该评论吗？此操作不可恢复！', '确认删除', {
      type: 'error'
    });
    
    await request.delete(`/api/v1/admin/comments/${commentId}`, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('评论删除成功');
    fetchComments();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.error || '删除失败');
    }
  }
};

// 关闭编辑弹窗
const closeEditModal = () => {
  showEditModal.value = false;
  editForm.value = {};
};

// 获取父评论作者
const getParentCommentAuthor = (parentId: number) => {
  const parentComment = comments.value.find(c => c.id === parentId);
  return parentComment?.user?.username || '未知用户';
};

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: { [key: string]: string } = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  };
  return statusMap[status] || status;
};

onMounted(() => {
  fetchComments();
  fetchArticles();
});
</script>

<style scoped>
.admin-comments {
  max-width: 1200px;
  margin: 0 auto;
}

.action-bar {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.search-section {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-input {
  padding: 10px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  width: 300px;
  font-size: 14px;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.search-btn {
  background: #667eea;
  border: none;
  color: white;
  padding: 10px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-btn:hover {
  background: #5a6fd8;
}

.filter-section {
  display: flex;
  gap: 10px;
}

.filter-select {
  padding: 10px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  font-size: 14px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.comment-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border-left: 4px solid #ddd;
}

.comment-card.pending {
  border-left-color: #ffc107;
}

.comment-card.approved {
  border-left-color: #28a745;
}

.comment-card.rejected {
  border-left-color: #dc3545;
}

.comment-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.comment-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.comment-date {
  color: #666;
  font-size: 12px;
}

.comment-status {
  display: flex;
  align-items: center;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.pending {
  background: #fff3cd;
  color: #856404;
}

.status-badge.approved {
  background: #d4edda;
  color: #155724;
}

.status-badge.rejected {
  background: #f8d7da;
  color: #721c24;
}

.comment-content {
  margin-bottom: 16px;
}

.article-info {
  margin-bottom: 8px;
}

.article-title {
  color: #667eea;
  font-size: 14px;
  font-weight: 500;
}

.comment-text {
  color: #333;
  line-height: 1.6;
}

.reply-to {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}

.comment-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.btn-approve, .btn-reject, .btn-edit, .btn-delete {
  background: none;
  border: none;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.btn-approve {
  color: #28a745;
  background: rgba(40, 167, 69, 0.1);
}

.btn-approve:hover {
  background: rgba(40, 167, 69, 0.2);
}

.btn-reject {
  color: #dc3545;
  background: rgba(220, 53, 69, 0.1);
}

.btn-reject:hover {
  background: rgba(220, 53, 69, 0.2);
}

.btn-edit {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.btn-edit:hover {
  background: rgba(102, 126, 234, 0.2);
}

.btn-delete {
  color: #dc3545;
  background: rgba(220, 53, 69, 0.1);
}

.btn-delete:hover {
  background: rgba(220, 53, 69, 0.2);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
}

.page-btn {
  background: white;
  border: 1px solid #ddd;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.page-btn:hover:not(:disabled) {
  background: #f8f9fa;
  border-color: #667eea;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  color: #666;
  font-size: 14px;
}

/* 弹窗样式 */
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
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 500px;
  max-width: 90vw;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 18px;
  color: #666;
  cursor: pointer;
  padding: 4px;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  color: #333;
  font-weight: 500;
}

.form-select, .form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
}

.form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-cancel, .btn-save {
  padding: 10px 20px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.btn-cancel {
  background: #f8f9fa;
  color: #666;
}

.btn-cancel:hover {
  background: #e9ecef;
}

.btn-save {
  background: #667eea;
  color: white;
}

.btn-save:hover {
  background: #5a6fd8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .search-section {
    justify-content: center;
  }
  
  .search-input {
    width: 100%;
  }
  
  .filter-section {
    justify-content: center;
  }
  
  .comment-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .comment-actions {
    justify-content: center;
  }
}
</style> 