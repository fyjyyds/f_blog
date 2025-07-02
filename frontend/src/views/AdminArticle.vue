<template>
  <div class="admin-articles">
    <!-- 搜索和过滤栏 -->
    <div class="action-bar">
      <div class="search-section">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索文章标题或内容..."
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
          <option value="published">已发布</option>
          <option value="rejected">已驳回</option>
          <option value="draft">草稿</option>
          <option value="private">私密</option>
        </select>
        
        <select v-model="categoryFilter" @change="handleFilter" class="filter-select">
          <option value="">全部分类</option>
          <option v-for="category in categories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
        </select>
      </div>
    </div>

    <!-- 批量操作栏 -->
    <div class="batch-actions" v-if="selectedArticles.length > 0">
      <span class="selected-count">已选择 {{ selectedArticles.length }} 篇文章</span>
      <div class="batch-buttons">
        <button class="btn-approve" @click="batchApprove">
          <i class="fas fa-check"></i>
          批量通过
        </button>
        <button class="btn-reject" @click="batchReject">
          <i class="fas fa-times"></i>
          批量驳回
        </button>
        <button class="btn-delete" @click="batchDelete">
          <i class="fas fa-trash"></i>
          批量删除
        </button>
      </div>
    </div>

    <!-- 文章列表 -->
    <div class="articles-table">
      <table>
        <thead>
          <tr>
            <th width="40">
              <input 
                type="checkbox" 
                :checked="isAllSelected" 
                @change="toggleSelectAll"
              />
            </th>
            <th>ID</th>
            <th>标题</th>
            <th>作者</th>
            <th>分类</th>
            <th>状态</th>
            <th>浏览量</th>
            <th>点赞数</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="article in filteredArticles" :key="article.id">
            <td>
              <input 
                type="checkbox" 
                :checked="selectedArticles.includes(article.id)"
                @change="toggleSelect(article.id)"
              />
            </td>
            <td>{{ article.id }}</td>
            <td class="article-title">
              <div class="title-content">
                <span class="title-text">{{ article.title }}</span>
                <div class="title-tags" v-if="article.tags && article.tags.length">
                  <span 
                    v-for="tag in article.tags.slice(0, 3)" 
                    :key="tag.id" 
                    class="tag-badge"
                    :style="{ background: tag.color || '#667eea' }"
                  >
                    {{ tag.name }}
                  </span>
                </div>
              </div>
            </td>
            <td>
              <div class="author-info">
                <img 
                  :src="article.author?.avatar || '/static/avatar/default.jpg'" 
                  :alt="article.author?.username"
                  class="author-avatar"
                />
                <span>{{ article.author?.username || '未知用户' }}</span>
              </div>
            </td>
            <td>
              <span class="category-badge" v-if="article.category">
                {{ article.category.name }}
              </span>
              <span v-else class="no-category">未分类</span>
            </td>
            <td>
              <span class="status-badge" :class="article.status">
                {{ getStatusText(article.status) }}
              </span>
            </td>
            <td>{{ article.view_count || 0 }}</td>
            <td>{{ article.like_count || 0 }}</td>
            <td>{{ formatDate(article.created_at) }}</td>
            <td>
              <div class="action-buttons">
                <button 
                  v-if="article.status === 'pending'" 
                  class="btn-approve" 
                  @click="reviewArticle(article.id, 'published')"
                  title="通过审核"
                >
                  <i class="fas fa-check"></i>
                </button>
                
                <button 
                  v-if="article.status === 'pending'" 
                  class="btn-reject" 
                  @click="rejectArticle(article.id)"
                  title="驳回文章"
                >
                  <i class="fas fa-times"></i>
                </button>
                
                <button 
                  v-if="article.status === 'rejected'" 
                  class="btn-approve" 
                  @click="reviewArticle(article.id, 'published')"
                  title="重新通过"
                >
                  <i class="fas fa-redo"></i>
                </button>
                
                <button class="btn-view" @click="viewArticle(article)" title="查看文章">
                  <i class="fas fa-eye"></i>
                </button>
                
                <button class="btn-edit" @click="editArticle(article)" title="编辑文章">
                  <i class="fas fa-edit"></i>
                </button>
                
                <button class="btn-delete" @click="deleteArticle(article.id)" title="删除文章">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
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

    <!-- 查看文章弹窗 -->
    <div v-if="showViewModal" class="modal-overlay" @click="closeViewModal">
      <div class="modal-content article-view" @click.stop>
        <div class="modal-header">
          <h3>{{ viewArticleData.title }}</h3>
          <button class="close-btn" @click="closeViewModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="article-meta">
            <div class="meta-item">
              <span class="label">作者：</span>
              <span>{{ viewArticleData.author?.username }}</span>
            </div>
            <div class="meta-item">
              <span class="label">分类：</span>
              <span>{{ viewArticleData.category?.name || '未分类' }}</span>
            </div>
            <div class="meta-item">
              <span class="label">状态：</span>
              <span class="status-badge" :class="viewArticleData.status">
                {{ getStatusText(viewArticleData.status) }}
              </span>
            </div>
            <div class="meta-item">
              <span class="label">创建时间：</span>
              <span>{{ formatDate(viewArticleData.created_at) }}</span>
            </div>
          </div>
          
          <div class="article-content">
            <h4>摘要</h4>
            <p>{{ viewArticleData.summary || '暂无摘要' }}</p>
            
            <h4>内容</h4>
            <div class="content-text" v-html="viewArticleData.content"></div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeViewModal">关闭</button>
          <button 
            v-if="viewArticleData.status === 'pending'" 
            class="btn-approve" 
            @click="reviewArticle(viewArticleData.id, 'published')"
          >
            通过审核
          </button>
          <button 
            v-if="viewArticleData.status === 'pending'" 
            class="btn-reject" 
            @click="rejectArticle(viewArticleData.id)"
          >
            驳回文章
          </button>
        </div>
      </div>
    </div>

    <!-- 编辑文章弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑文章</h3>
          <button class="close-btn" @click="closeEditModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>标题 *</label>
            <input v-model="editForm.title" type="text" class="form-input" />
          </div>
          
          <div class="form-group">
            <label>摘要</label>
            <textarea v-model="editForm.summary" class="form-textarea" rows="3"></textarea>
          </div>
          
          <div class="form-group">
            <label>内容</label>
            <textarea v-model="editForm.content" class="form-textarea" rows="8"></textarea>
          </div>
          
          <div class="form-group">
            <label>分类</label>
            <select v-model="editForm.category_id" class="form-select">
              <option value="">选择分类</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label>状态</label>
            <select v-model="editForm.status" class="form-select">
              <option value="pending">待审核</option>
              <option value="published">已发布</option>
              <option value="rejected">已驳回</option>
              <option value="draft">草稿</option>
              <option value="private">私密</option>
            </select>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeEditModal">取消</button>
          <button class="btn-save" @click="saveArticle">保存</button>
        </div>
      </div>
    </div>

    <!-- 驳回原因弹窗 -->
    <div v-if="showRejectModal" class="modal-overlay" @click="closeRejectModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>驳回文章</h3>
          <button class="close-btn" @click="closeRejectModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>驳回原因 *</label>
            <textarea 
              v-model="rejectReason" 
              class="form-textarea" 
              rows="4"
              placeholder="请输入驳回原因，将通知作者"
            ></textarea>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeRejectModal">取消</button>
          <button class="btn-reject" @click="confirmReject" :disabled="!rejectReason.trim()">
            确认驳回
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import request from '../api/request';
import { formatDate } from '../utils/date';

// 响应式数据
const articles = ref<any[]>([]);
const categories = ref<any[]>([]);
const loading = ref(false);
const searchQuery = ref('');
const statusFilter = ref('');
const categoryFilter = ref('');
const selectedArticles = ref<number[]>([]);
const currentPage = ref(1);
const pageSize = ref(20);
const totalPages = ref(1);

// 弹窗状态
const showViewModal = ref(false);
const showEditModal = ref(false);
const showRejectModal = ref(false);
const viewArticleData = ref<any>({});
const editForm = ref<any>({});
const rejectReason = ref('');
const rejectArticleId = ref<number | null>(null);

// 计算属性
const filteredArticles = computed(() => {
  let filtered = articles.value;
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(article => 
      article.title?.toLowerCase().includes(query) ||
      article.content?.toLowerCase().includes(query) ||
      article.summary?.toLowerCase().includes(query)
    );
  }
  
  if (statusFilter.value) {
    filtered = filtered.filter(article => article.status === statusFilter.value);
  }
  
  if (categoryFilter.value) {
    filtered = filtered.filter(article => article.category_id === parseInt(categoryFilter.value));
  }
  
  return filtered;
});

const isAllSelected = computed(() => {
  return filteredArticles.value.length > 0 && 
         filteredArticles.value.every(article => selectedArticles.value.includes(article.id));
});

// 方法
async function fetchArticles() {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value || undefined,
      category_id: categoryFilter.value || undefined,
      search: searchQuery.value || undefined
    };
    
    const res = await request.get('/api/v1/admin/articles', { params });
    articles.value = res.data.articles || res.data;
    totalPages.value = Math.ceil((res.data.total || articles.value.length) / pageSize.value);
  } catch (err: any) {
    console.error('获取文章失败:', err);
  }
  loading.value = false;
}

async function fetchCategories() {
  try {
    const res = await request.get('/api/v1/categories');
    categories.value = res.data;
  } catch (err: any) {
    console.error('获取分类失败:', err);
  }
}

function handleSearch() {
  currentPage.value = 1;
  fetchArticles();
}

function handleFilter() {
  currentPage.value = 1;
  fetchArticles();
}

function changePage(page: number) {
  currentPage.value = page;
  fetchArticles();
}

function toggleSelect(articleId: number) {
  const index = selectedArticles.value.indexOf(articleId);
  if (index > -1) {
    selectedArticles.value.splice(index, 1);
  } else {
    selectedArticles.value.push(articleId);
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedArticles.value = [];
  } else {
    selectedArticles.value = filteredArticles.value.map(article => article.id);
  }
}

async function reviewArticle(id: number, status: string) {
  try {
    await request.post(`/api/v1/admin/articles/${id}/review`, { status });
    await fetchArticles();
    showViewModal.value = false;
    showEditModal.value = false;
  } catch (err: any) {
    console.error('审核失败:', err);
  }
}

function rejectArticle(id: number) {
  rejectArticleId.value = id;
  rejectReason.value = '';
  showRejectModal.value = true;
}

async function confirmReject() {
  if (!rejectReason.value.trim()) return;
  
  try {
    await request.post(`/api/v1/admin/articles/${rejectArticleId.value}/review`, {
      status: 'rejected',
      reason: rejectReason.value
    });
    await fetchArticles();
    showRejectModal.value = false;
    showViewModal.value = false;
  } catch (err: any) {
    console.error('驳回失败:', err);
  }
}

function viewArticle(article: any) {
  viewArticleData.value = { ...article };
  showViewModal.value = true;
}

function editArticle(article: any) {
  editForm.value = { ...article };
  showEditModal.value = true;
}

async function saveArticle() {
  try {
    await request.put(`/api/v1/admin/articles/${editForm.value.id}`, editForm.value);
    await fetchArticles();
    showEditModal.value = false;
  } catch (err: any) {
    console.error('保存失败:', err);
  }
}

async function deleteArticle(id: number) {
  if (!confirm('确定要删除该文章吗？')) return;
  
  try {
    await request.delete(`/api/v1/admin/articles/${id}`);
    await fetchArticles();
  } catch (err: any) {
    console.error('删除失败:', err);
  }
}

async function batchApprove() {
  if (!confirm(`确定要通过这 ${selectedArticles.value.length} 篇文章吗？`)) return;
  
  try {
    await Promise.all(
      selectedArticles.value.map(id => 
        request.post(`/api/v1/admin/articles/${id}/review`, { status: 'published' })
      )
    );
    selectedArticles.value = [];
    await fetchArticles();
  } catch (err: any) {
    console.error('批量通过失败:', err);
  }
}

async function batchReject() {
  if (!confirm(`确定要驳回这 ${selectedArticles.value.length} 篇文章吗？`)) return;
  
  try {
    await Promise.all(
      selectedArticles.value.map(id => 
        request.post(`/api/v1/admin/articles/${id}/review`, { status: 'rejected' })
      )
    );
    selectedArticles.value = [];
    await fetchArticles();
  } catch (err: any) {
    console.error('批量驳回失败:', err);
  }
}

async function batchDelete() {
  if (!confirm(`确定要删除这 ${selectedArticles.value.length} 篇文章吗？`)) return;
  
  try {
    await Promise.all(
      selectedArticles.value.map(id => 
        request.delete(`/api/v1/admin/articles/${id}`)
      )
    );
    selectedArticles.value = [];
    await fetchArticles();
      } catch (err: any) {
    console.error('批量删除失败:', err);
  }
}

function closeViewModal() {
  showViewModal.value = false;
  viewArticleData.value = {};
}

function closeEditModal() {
  showEditModal.value = false;
  editForm.value = {};
}

function closeRejectModal() {
  showRejectModal.value = false;
  rejectReason.value = '';
  rejectArticleId.value = null;
}

function getStatusText(status: string) {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    published: '已发布',
    rejected: '已驳回',
    draft: '草稿',
    private: '私密'
  };
  return statusMap[status] || status;
}

onMounted(() => {
  fetchArticles();
  fetchCategories();
});
</script>

<style scoped>
.admin-articles {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

/* 搜索和过滤栏 */
.action-bar {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  max-width: 400px;
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #e1e5e9;
  border-radius: 25px;
  font-size: 14px;
  transition: all 0.3s ease;
  background: white;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.search-btn {
  padding: 12px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.filter-section {
  display: flex;
  gap: 10px;
}

.filter-select {
  padding: 10px 15px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  background: rgb(37, 35, 35);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-select:focus {
  outline: none;
  border-color: #667eea;
}

/* 批量操作栏 */
.batch-actions {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 15px 20px;
  margin-bottom: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.selected-count {
  font-weight: 600;
  color: #667eea;
}

.batch-buttons {
  display: flex;
  gap: 10px;
}

.batch-buttons button {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 5px;
}

.btn-approve {
  background: linear-gradient(135deg, #4ade80 0%, #22c55e 100%);
  color: white;
}

.btn-approve:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(34, 197, 94, 0.4);
}

.btn-reject {
  background: linear-gradient(135deg, #f87171 0%, #ef4444 100%);
  color: white;
}

.btn-reject:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(239, 68, 68, 0.4);
}

.btn-delete {
  background: linear-gradient(135deg, #f97316 0%, #ea580c 100%);
  color: white;
}

.btn-delete:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(234, 88, 12, 0.4);
}

/* 文章表格 */
.articles-table {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.articles-table table {
  width: 100%;
  border-collapse: collapse;
}

.articles-table th {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 15px 12px;
  text-align: left;
  font-weight: 600;
  font-size: 14px;
}

.articles-table td {
  padding: 12px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 14px;
}

.articles-table tr:hover {
  background: rgba(102, 126, 234, 0.05);
}

.articles-table tr:last-child td {
  border-bottom: none;
}

/* 文章标题 */
.article-title {
  max-width: 300px;
}

.title-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.title-text {
  font-weight: 600;
  color: #1e293b;
  line-height: 1.4;
}

.title-tags {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.tag-badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  color: white;
  font-weight: 500;
}

/* 作者信息 */
.author-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

/* 分类标签 */
.category-badge {
  padding: 4px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 15px;
  font-size: 12px;
  font-weight: 500;
}

.no-category {
  color: #94a3b8;
  font-style: italic;
}

/* 状态标签 */
.status-badge {
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 12px;
  font-weight: 500;
  text-align: center;
  display: inline-block;
  min-width: 80px;
}

.status-badge.pending {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
}

.status-badge.published {
  background: linear-gradient(135deg, #4ade80 0%, #22c55e 100%);
  color: white;
}

.status-badge.rejected {
  background: linear-gradient(135deg, #f87171 0%, #ef4444 100%);
  color: white;
}

.status-badge.draft {
  background: linear-gradient(135deg, #94a3b8 0%, #64748b 100%);
  color: white;
}

.status-badge.private {
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
  color: white;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.action-buttons button {
  padding: 6px 10px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  height: 32px;
}

.btn-view {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
}

.btn-view:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(37, 99, 235, 0.4);
}

.btn-edit {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.btn-edit:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(217, 119, 6, 0.4);
}

/* 分页 */
.pagination {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 15px 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
}

.page-btn {
  padding: 8px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.page-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-weight: 600;
  color: #667eea;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 15px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  animation: modalSlideIn 0.3s ease;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-50px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.modal-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  padding: 5px;
  border-radius: 5px;
  transition: background 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.modal-body {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid #e1e5e9;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 表单样式 */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #374151;
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s ease;
  background: rgb(37, 36, 36);
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

/* 按钮样式 */
.btn-cancel {
  padding: 10px 20px;
  background: #f1f5f9;
  color: #64748b;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-cancel:hover {
  background: #e2e8f0;
}

.btn-save {
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-save:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

/* 文章查看弹窗特殊样式 */
.article-view {
  width: 800px;
}

.article-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
  margin-bottom: 20px;
  padding: 15px;
  background: #f8fafc;
  border-radius: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.meta-item .label {
  font-weight: 600;
  color: #64748b;
  min-width: 80px;
}

.article-content h4 {
  color: #1e293b;
  margin: 20px 0 10px 0;
  font-size: 16px;
  font-weight: 600;
}

.article-content p {
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 15px;
}

.content-text {
  background: #f8fafc;
  padding: 15px;
  border-radius: 8px;
  line-height: 1.6;
  color: #374151;
  max-height: 300px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-section {
    max-width: none;
  }
  
  .filter-section {
    justify-content: center;
  }
  
  .batch-actions {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .articles-table {
    overflow-x: auto;
  }
  
  .articles-table table {
    min-width: 800px;
  }
  
  .action-buttons {
    flex-direction: column;
    gap: 3px;
  }
  
  .modal-content {
    margin: 20px;
    max-width: calc(100vw - 40px);
  }
  
  .article-view {
    width: auto;
  }
  
  .article-meta {
    grid-template-columns: 1fr;
  }
}
</style> 