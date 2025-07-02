<template>
  <div class="admin-categories">
    <!-- 操作栏 -->
    <div class="action-bar">
      <h2>分类管理</h2>
      <button class="add-btn" @click="showAddModal = true">
        <i class="fas fa-plus"></i>
        添加分类
      </button>
    </div>

    <!-- 分类列表 -->
    <div class="categories-grid">
      <div 
        v-for="category in categories" 
        :key="category.id" 
        class="category-card"
      >
        <div class="category-header">
          <h3>{{ category.name }}</h3>
          <div class="category-actions">
            <button class="btn-edit" @click="editCategory(category)">
              <i class="fas fa-edit"></i>
            </button>
            <button class="btn-delete" @click="deleteCategory(category.id)">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
        
        <div class="category-content">
          <p class="category-description">{{ category.description || '暂无描述' }}</p>
          <div class="category-stats">
            <span class="stat-item">
              <i class="fas fa-newspaper"></i>
              {{ category.article_count || 0 }} 篇文章
            </span>
            <span class="stat-item">
              <i class="fas fa-eye"></i>
              {{ category.view_count || 0 }} 次浏览
            </span>
          </div>
        </div>
        
        <div class="category-footer">
          <span class="category-status" :class="category.status">
            {{ getStatusText(category.status) }}
          </span>
          <span class="category-date">
            创建于 {{ formatDate(category.created_at) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 添加/编辑分类弹窗 -->
    <div v-if="showAddModal || showEditModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditModal ? '编辑分类' : '添加分类' }}</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>分类名称 *</label>
            <input 
              v-model="formData.name" 
              type="text" 
              class="form-input" 
              placeholder="请输入分类名称"
            />
          </div>
          
          <div class="form-group">
            <label>分类描述</label>
            <textarea 
              v-model="formData.description" 
              class="form-textarea" 
              rows="3"
              placeholder="请输入分类描述"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>状态</label>
            <select v-model="formData.status" class="form-select">
              <option value="active">启用</option>
              <option value="inactive">禁用</option>
            </select>
          </div>
          
          <div class="form-group">
            <label>排序</label>
            <input 
              v-model.number="formData.sort_order" 
              type="number" 
              class="form-input" 
              placeholder="数字越小排序越靠前"
            />
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button class="btn-save" @click="saveCategory" :disabled="!formData.name">
            {{ showEditModal ? '更新' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import request from '../api/request';
import { formatDate } from '../utils/date';

// 响应式数据
const categories = ref<any[]>([]);
const showAddModal = ref(false);
const showEditModal = ref(false);
const formData = ref({
  name: '',
  description: '',
  status: 'active',
  sort_order: 0
});

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await request.get('/api/v1/admin/categories', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    categories.value = response.data || [];
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '获取分类列表失败');
  }
};

// 编辑分类
const editCategory = (category: any) => {
  formData.value = { ...category };
  showEditModal.value = true;
};

// 保存分类
const saveCategory = async () => {
  if (!formData.value.name.trim()) {
    ElMessage.warning('请输入分类名称');
    return;
  }

  try {
    if (showEditModal.value) {
      // 更新分类
      await request.put(`/api/v1/admin/categories/${formData.value.id}`, formData.value, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
      });
      ElMessage.success('分类更新成功');
    } else {
      // 创建分类
      await request.post('/api/v1/admin/categories', formData.value, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
      });
      ElMessage.success('分类创建成功');
    }
    
    closeModal();
    fetchCategories();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '操作失败');
  }
};

// 删除分类
const deleteCategory = async (categoryId: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该分类吗？删除后该分类下的文章将变为未分类状态。', '确认删除', {
      type: 'warning'
    });
    
    await request.delete(`/api/v1/admin/categories/${categoryId}`, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('分类删除成功');
    fetchCategories();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.error || '删除失败');
    }
  }
};

// 关闭弹窗
const closeModal = () => {
  showAddModal.value = false;
  showEditModal.value = false;
  formData.value = {
    name: '',
    description: '',
    status: 'active',
    sort_order: 0
  };
};

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: { [key: string]: string } = {
    active: '启用',
    inactive: '禁用'
  };
  return statusMap[status] || status;
};

onMounted(() => {
  fetchCategories();
});
</script>

<style scoped>
.admin-categories {
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

.action-bar h2 {
  margin: 0;
  color: #333;
  font-size: 1.5rem;
}

.add-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
  padding: 12px 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.category-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid #f0f0f0;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  border-color: #667eea;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.category-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
  font-weight: 600;
}

.category-actions {
  display: flex;
  gap: 8px;
}

.btn-edit, .btn-delete {
  background: none;
  border: none;
  padding: 6px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-edit {
  color: #667eea;
}

.btn-edit:hover {
  background: rgba(102, 126, 234, 0.1);
}

.btn-delete {
  color: #ff6b6b;
}

.btn-delete:hover {
  background: rgba(255, 107, 107, 0.1);
}

.category-content {
  margin-bottom: 16px;
}

.category-description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin: 0 0 12px 0;
}

.category-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 12px;
}

.stat-item i {
  color: #667eea;
}

.category-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.category-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.category-status.active {
  background: #51cf66;
  color: white;
}

.category-status.inactive {
  background: #adb5bd;
  color: white;
}

.category-date {
  color: #999;
  font-size: 12px;
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

.form-input, .form-select, .form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
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

.btn-save:hover:not(:disabled) {
  background: #5a6fd8;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .add-btn {
    justify-content: center;
  }
  
  .categories-grid {
    grid-template-columns: 1fr;
  }
  
  .category-stats {
    flex-direction: column;
    gap: 8px;
  }
  
  .category-footer {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
}
</style> 