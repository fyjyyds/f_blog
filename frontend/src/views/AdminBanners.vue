<template>
  <div class="admin-banners">
    <!-- 操作栏 -->
    <div class="action-bar">
      <h2>轮播图管理</h2>
      <button class="add-btn" @click="showAddModal = true">
        <i class="fas fa-plus"></i>
        添加轮播图
      </button>
    </div>

    <!-- 轮播图列表 -->
    <div class="banners-grid">
      <div 
        v-for="banner in banners" 
        :key="banner.id" 
        class="banner-card"
      >
        <div class="banner-image">
          <img :src="getBannerUrl(banner.image)" :alt="banner.title" />
          <div class="banner-overlay">
            <div class="banner-actions">
              <button class="btn-edit" @click="editBanner(banner)">
                <i class="fas fa-edit"></i>
              </button>
              <button class="btn-delete" @click="deleteBanner(banner.id)">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
        </div>
        
        <div class="banner-content">
          <h3>{{ banner.title }}</h3>
          <p class="banner-description">{{ banner.description || '暂无描述' }}</p>
          
          <div class="banner-info">
            <div class="info-item">
              <span class="label">链接：</span>
              <a :href="banner.link_url" target="_blank" class="link-url">
                {{ banner.link_url || '无链接' }}
              </a>
            </div>
            
            <div class="info-item">
              <span class="label">排序：</span>
              <span class="sort-order">{{ banner.sort_order || 0 }}</span>
            </div>
            
            <div class="info-item">
              <span class="label">状态：</span>
              <span class="status-badge" :class="banner.status">
                {{ getStatusText(banner.status) }}
              </span>
            </div>
          </div>
          
          <div class="banner-footer">
            <span class="banner-date">
              创建于 {{ formatDate(banner.created_at) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑轮播图弹窗 -->
    <div v-if="showAddModal || showEditModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditModal ? '编辑轮播图' : '添加轮播图' }}</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>标题 *</label>
            <input 
              v-model="formData.title" 
              type="text" 
              class="form-input" 
              placeholder="请输入轮播图标题"
            />
          </div>
          
          <div class="form-group">
            <label>描述</label>
            <textarea 
              v-model="formData.description" 
              class="form-textarea" 
              rows="3"
              placeholder="请输入轮播图描述"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label>图片 *</label>
            <div class="image-upload">
              <input 
                ref="fileInput"
                type="file" 
                accept="image/*" 
                @change="handleImageUpload"
                class="file-input"
              />
              <div class="upload-area" @click="$refs.fileInput.click()">
                <div v-if="!imagePreview" class="upload-placeholder">
                  <i class="fas fa-cloud-upload-alt"></i>
                  <p>点击上传图片</p>
                </div>
                <img v-else :src="imagePreview" alt="预览" class="image-preview" />
              </div>
            </div>
          </div>
          
          <div class="form-group">
            <label>链接地址</label>
            <input 
              v-model="formData.link_url" 
              type="url" 
              class="form-input" 
              placeholder="请输入链接地址（可选）"
            />
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
          <button class="btn-save" @click="saveBanner" :disabled="!formData.title || !imagePreview">
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

const getBannerUrl = (path: string) => {
  if (!path) return '';
  if (path.startsWith('http')) return path;
  return 'http://localhost:8080' + path;
};

// 响应式数据
const banners = ref<any[]>([]);
const showAddModal = ref(false);
const showEditModal = ref(false);
const imagePreview = ref('');
const formData = ref({
  title: '',
  description: '',
  image: '',
  link_url: '',
  sort_order: 0,
  status: 'active'
});

// 获取轮播图列表
const fetchBanners = async () => {
  try {
    const response = await request.get('/api/v1/admin/banners', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    banners.value = response.data || [];
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '获取轮播图列表失败');
  }
};

// 编辑轮播图
const editBanner = (banner: any) => {
  formData.value = { ...banner };
  imagePreview.value = getBannerUrl(banner.image);
  showEditModal.value = true;
};

// 处理图片上传
const handleImageUpload = (event: any) => {
  const file = event.target.files[0];
  if (file) {
    // 创建预览
    const reader = new FileReader();
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string;
    };
    reader.readAsDataURL(file);
    
    // 上传到服务器
    uploadImage(file);
  }
};

// 上传图片到服务器
const uploadImage = async (file: File) => {
  try {
    const form = new FormData();
    form.append('file', file);
    const response = await request.post('/api/v1/upload', form, {
      headers: { 
        'Content-Type': 'multipart/form-data',
        Authorization: 'Bearer ' + localStorage.getItem('token')
      }
    });
    formData.value.image = response.data.url;
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '图片上传失败');
  }
};

// 保存轮播图
const saveBanner = async () => {
  if (!formData.value.title.trim()) {
    ElMessage.warning('请输入轮播图标题');
    return;
  }

  if (!imagePreview.value) {
    ElMessage.warning('请上传轮播图图片');
    return;
  }

  try {
    if (showEditModal.value) {
      // 更新轮播图
      await request.put(`/api/v1/admin/banners/${formData.value.id}`, formData.value, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
      });
      ElMessage.success('轮播图更新成功');
    } else {
      // 创建轮播图
      await request.post('/api/v1/admin/banners', formData.value, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
      });
      ElMessage.success('轮播图创建成功');
    }
    
    closeModal();
    fetchBanners();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '操作失败');
  }
};

// 删除轮播图
const deleteBanner = async (bannerId: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该轮播图吗？', '确认删除', {
      type: 'warning'
    });
    
    await request.delete(`/api/v1/admin/banners/${bannerId}`, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('轮播图删除成功');
    fetchBanners();
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
  imagePreview.value = '';
  formData.value = {
    title: '',
    description: '',
    image: '',
    link_url: '',
    sort_order: 0,
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
  fetchBanners();
});
</script>

<style scoped>
.admin-banners {
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

.banners-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.banner-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.banner-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.banner-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.banner-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.banner-card:hover .banner-overlay {
  opacity: 1;
}

.banner-actions {
  display: flex;
  gap: 8px;
}

.btn-edit, .btn-delete {
  background: rgba(255, 255, 255, 0.9);
  border: none;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #333;
}

.btn-edit:hover {
  background: white;
  color: #667eea;
}

.btn-delete:hover {
  background: white;
  color: #ff6b6b;
}

.banner-content {
  padding: 20px;
}

.banner-content h3 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 1.2rem;
  font-weight: 600;
}

.banner-description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin: 0 0 16px 0;
}

.banner-info {
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.label {
  font-weight: 500;
  color: #333;
  min-width: 60px;
}

.link-url {
  color: #667eea;
  text-decoration: none;
  word-break: break-all;
}

.link-url:hover {
  text-decoration: underline;
}

.sort-order {
  color: #666;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: #51cf66;
  color: white;
}

.status-badge.inactive {
  background: #adb5bd;
  color: white;
}

.banner-footer {
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.banner-date {
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
  width: 600px;
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

.image-upload {
  position: relative;
}

.file-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.upload-area {
  border: 2px dashed #ddd;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #f8f9fa;
}

.upload-area:hover {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.upload-placeholder {
  color: #666;
}

.upload-placeholder i {
  font-size: 48px;
  margin-bottom: 12px;
  color: #ddd;
}

.upload-placeholder p {
  margin: 0;
  font-size: 14px;
}

.image-preview {
  max-width: 100%;
  max-height: 200px;
  border-radius: 6px;
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
  
  .banners-grid {
    grid-template-columns: 1fr;
  }
  
  .banner-info {
    flex-direction: column;
  }
  
  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .label {
    min-width: auto;
  }
}
</style> 