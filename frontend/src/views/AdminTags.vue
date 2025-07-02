<template>
  <div class="admin-tags">
    <!-- 操作栏 -->
    <div class="action-bar">
      <h2>标签管理</h2>
      <button class="add-btn" @click="showAddModal = true">
        <i class="fas fa-plus"></i>
        添加标签
      </button>
    </div>

    <!-- 标签列表 -->
    <div class="tags-grid">
      <div 
        v-for="tag in tags" 
        :key="tag.id" 
        class="tag-card"
      >
        <div class="tag-header">
          <div class="tag-name">
            <span class="tag-color" :style="{ background: tag.color || '#667eea' }"></span>
            <h3>{{ tag.name }}</h3>
          </div>
          <div class="tag-actions">
            <button class="btn-edit" @click="editTag(tag)">
              <i class="fas fa-edit"></i>
            </button>
            <button class="btn-delete" @click="deleteTag(tag.id)">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
        
        <div class="tag-content">
          <p class="tag-description">{{ tag.description || '暂无描述' }}</p>
          <div class="tag-stats">
            <span class="stat-item">
              <i class="fas fa-newspaper"></i>
              {{ tag.article_count || 0 }} 篇文章
            </span>
            <span class="stat-item">
              <i class="fas fa-eye"></i>
              {{ tag.view_count || 0 }} 次浏览
            </span>
          </div>
        </div>
        
        <div class="tag-footer">
          <span class="tag-status" :class="tag.status">
            {{ getStatusText(tag.status) }}
          </span>
          <span class="tag-date">
            创建于 {{ formatDate(tag.created_at) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 添加/编辑标签弹窗 -->
    <div v-if="showAddModal || showEditModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditModal ? '编辑标签' : '添加标签' }}</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>标签名称 *</label>
            <input 
              v-model="formData.name" 
              type="text" 
              class="form-input" 
              placeholder="请输入标签名称"
            />
          </div>
          
          <div class="form-group">
            <label>标签描述</label>
            <textarea 
              v-model="formData.description" 
              class="form-textarea" 
              rows="3"
              placeholder="请输入标签描述"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>标签颜色</label>
            <div class="color-picker">
              <input 
                v-model="formData.color" 
                type="color" 
                class="color-input"
              />
              <span class="color-preview" :style="{ background: formData.color || '#667eea' }"></span>
            </div>
          </div>
          
          <div class="form-group">
            <label>状态</label>
            <select v-model="formData.status" class="form-select">
              <option value="active">启用</option>
              <option value="inactive">禁用</option>
            </select>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button class="btn-save" @click="saveTag" :disabled="!formData.name">
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
const tags = ref<any[]>([]);
const showAddModal = ref(false);
const showEditModal = ref(false);
const formData = ref({
  name: '',
  description: '',
  color: '#667eea',
  status: 'active'
});

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await request.get('/api/v1/admin/tags', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    tags.value = response.data || [];
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '获取标签列表失败');
  }
};

// 编辑标签
const editTag = (tag: any) => {
  formData.value = { ...tag };
  showEditModal.value = true;
};

// 保存标签
const saveTag = async () => {
  if (!formData.value.name.trim()) {
    ElMessage.warning('请输入标签名称');
    return;
  }

  try {
    if (showEditModal.value) {
      // 更新标签
      await request.put(`/api/v1/admin/tags/${formData.value.id}`, formData.value, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
      });
      ElMessage.success('标签更新成功');
    } else {
      // 创建标签
      await request.post('/api/v1/admin/tags', formData.value, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
      });
      ElMessage.success('标签创建成功');
    }
    
    closeModal();
    fetchTags();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '操作失败');
  }
};

// 删除标签
const deleteTag = async (tagId: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该标签吗？删除后该标签将从所有文章中移除。', '确认删除', {
      type: 'warning'
    });
    
    await request.delete(`/api/v1/admin/tags/${tagId}`, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('标签删除成功');
    fetchTags();
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
    color: '#667eea',
    status: 'active'
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
  fetchTags();
});
</script>

<style scoped>
.admin-tags {
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

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.tag-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid #f0f0f0;
}

.tag-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  border-color: #667eea;
}

.tag-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.tag-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.tag-color {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.tag-name h3 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
  font-weight: 600;
}

.tag-actions {
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

.tag-content {
  margin-bottom: 16px;
}

.tag-description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin: 0 0 12px 0;
}

.tag-stats {
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

.tag-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.tag-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.tag-status.active {
  background: #51cf66;
  color: white;
}

.tag-status.inactive {
  background: #adb5bd;
  color: white;
}

.tag-date {
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

.color-picker {
  display: flex;
  align-items: center;
  gap: 12px;
}

.color-input {
  width: 50px;
  height: 40px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  background: none;
}

.color-input::-webkit-color-swatch-wrapper {
  padding: 0;
}

.color-input::-webkit-color-swatch {
  border: none;
  border-radius: 6px;
}

.color-preview {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  border: 2px solid #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
  
  .tags-grid {
    grid-template-columns: 1fr;
  }
  
  .tag-stats {
    flex-direction: column;
    gap: 8px;
  }
  
  .tag-footer {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
}
</style> 