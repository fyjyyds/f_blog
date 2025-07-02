<template>
  <div class="admin-settings">
    <!-- 设置导航 -->
    <div class="settings-nav">
      <button 
        v-for="tab in tabs" 
        :key="tab.key"
        class="nav-tab"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        <i :class="tab.icon"></i>
        {{ tab.title }}
      </button>
    </div>

    <!-- 设置内容 -->
    <div class="settings-content">
      <!-- 基本设置 -->
      <div v-if="activeTab === 'basic'" class="settings-section">
        <h3>基本设置</h3>
        <div class="form-group">
          <label>网站名称</label>
          <input v-model="basicSettings.siteName" type="text" class="form-input" />
        </div>
        
        <div class="form-group">
          <label>网站描述</label>
          <textarea v-model="basicSettings.siteDescription" class="form-textarea" rows="3"></textarea>
        </div>
        
        <div class="form-group">
          <label>网站关键词</label>
          <input v-model="basicSettings.siteKeywords" type="text" class="form-input" placeholder="用逗号分隔" />
        </div>
        
        <div class="form-group">
          <label>网站Logo</label>
          <div class="image-upload">
            <input ref="logoInput" type="file" accept="image/*" @change="handleLogoUpload" class="file-input" />
            <div class="upload-area" @click="$refs.logoInput.click()">
              <div v-if="!basicSettings.siteLogo" class="upload-placeholder">
                <i class="fas fa-image"></i>
                <p>点击上传Logo</p>
              </div>
              <img v-else :src="basicSettings.siteLogo" alt="Logo" class="image-preview" />
            </div>
          </div>
        </div>
        
        <div class="form-group">
          <label>备案信息</label>
          <input v-model="basicSettings.icp" type="text" class="form-input" />
        </div>
      </div>

      <!-- 安全设置 -->
      <div v-if="activeTab === 'security'" class="settings-section">
        <h3>安全设置</h3>
        
        <div class="form-group">
          <label>密码最小长度</label>
          <input v-model.number="securitySettings.minPasswordLength" type="number" class="form-input" min="6" max="20" />
        </div>
        
        <div class="form-group">
          <label>登录失败锁定次数</label>
          <input v-model.number="securitySettings.maxLoginAttempts" type="number" class="form-input" min="3" max="10" />
        </div>
        
        <div class="form-group">
          <label>锁定时间（分钟）</label>
          <input v-model.number="securitySettings.lockoutDuration" type="number" class="form-input" min="5" max="60" />
        </div>
        
        <div class="form-group">
          <label>JWT过期时间（小时）</label>
          <input v-model.number="securitySettings.jwtExpireHours" type="number" class="form-input" min="1" max="168" />
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="securitySettings.enableCaptcha" type="checkbox" />
            启用验证码
          </label>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="securitySettings.enableTwoFactor" type="checkbox" />
            启用双因素认证
          </label>
        </div>
      </div>

      <!-- 内容设置 -->
      <div v-if="activeTab === 'content'" class="settings-section">
        <h3>内容设置</h3>
        
        <div class="form-group">
          <label>每页显示文章数</label>
          <input v-model.number="contentSettings.articlesPerPage" type="number" class="form-input" min="5" max="50" />
        </div>
        
        <div class="form-group">
          <label>文章摘要长度</label>
          <input v-model.number="contentSettings.summaryLength" type="number" class="form-input" min="50" max="500" />
        </div>
        
        <div class="form-group">
          <label>评论审核</label>
          <select v-model="contentSettings.commentModeration" class="form-select">
            <option value="auto">自动通过</option>
            <option value="manual">手动审核</option>
            <option value="first">首次评论审核</option>
          </select>
        </div>
        
        <div class="form-group">
          <label>允许匿名评论</label>
          <select v-model="contentSettings.allowAnonymousComments" class="form-select">
            <option value="true">允许</option>
            <option value="false">不允许</option>
          </select>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="contentSettings.enableCommentLikes" type="checkbox" />
            启用评论点赞
          </label>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="contentSettings.enableArticleLikes" type="checkbox" />
            启用文章点赞
          </label>
        </div>
      </div>

      <!-- 保存按钮 -->
      <div class="settings-actions">
        <button class="btn-save" @click="saveSettings">
          <i class="fas fa-save"></i>
          保存设置
        </button>
        <button class="btn-reset" @click="resetSettings">
          <i class="fas fa-undo"></i>
          重置设置
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import request from '../api/request';

// 响应式数据
const activeTab = ref('basic');
const tabs = [
  { key: 'basic', title: '基本设置', icon: 'fas fa-cog' },
  { key: 'security', title: '安全设置', icon: 'fas fa-shield-alt' },
  { key: 'content', title: '内容设置', icon: 'fas fa-file-alt' }
];

const basicSettings = ref({
  siteName: '',
  siteDescription: '',
  siteKeywords: '',
  siteLogo: '',
  icp: ''
});

const securitySettings = ref({
  minPasswordLength: 8,
  maxLoginAttempts: 5,
  lockoutDuration: 30,
  jwtExpireHours: 24,
  enableCaptcha: false,
  enableTwoFactor: false
});

const contentSettings = ref({
  articlesPerPage: 10,
  summaryLength: 200,
  commentModeration: 'manual',
  allowAnonymousComments: false,
  enableCommentLikes: true,
  enableArticleLikes: true
});

// 获取设置
const fetchSettings = async () => {
  try {
    const response = await request.get('/api/v1/admin/settings', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    const settings = response.data;
    if (settings.basic) Object.assign(basicSettings.value, settings.basic);
    if (settings.security) Object.assign(securitySettings.value, settings.security);
    if (settings.content) Object.assign(contentSettings.value, settings.content);
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '获取设置失败');
  }
};

// 保存设置
const saveSettings = async () => {
  try {
    const settings = {
      basic: basicSettings.value,
      security: securitySettings.value,
      content: contentSettings.value
    };
    await request.put('/api/v1/admin/settings', settings, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    ElMessage.success('设置保存成功');
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '保存设置失败');
  }
};

// 重置设置
const resetSettings = async () => {
  try {
    await request.post('/api/v1/admin/settings/reset', {}, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('设置已重置为默认值');
    fetchSettings();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '重置设置失败');
  }
};

// 处理Logo上传
const handleLogoUpload = (event: any) => {
  const file = event.target.files[0];
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      basicSettings.value.siteLogo = e.target?.result as string;
    };
    reader.readAsDataURL(file);
    
    // 上传到服务器
    uploadLogo(file);
  }
};

// 上传Logo到服务器
const uploadLogo = async (file: File) => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await request.post('/api/v1/upload', formData, {
      headers: { 
        'Content-Type': 'multipart/form-data',
        Authorization: 'Bearer ' + localStorage.getItem('token')
      }
    });
    
    basicSettings.value.siteLogo = response.data.url;
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || 'Logo上传失败');
  }
};

onMounted(() => {
  fetchSettings();
});
</script>

<style scoped>
.admin-settings {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  gap: 20px;
}

.settings-nav {
  width: 250px;
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  height: fit-content;
}

.nav-tab {
  width: 100%;
  padding: 12px 16px;
  border: none;
  background: none;
  text-align: left;
  cursor: pointer;
  border-radius: 8px;
  margin-bottom: 8px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 12px;
  color: #666;
  font-size: 14px;
}

.nav-tab:hover {
  background: #f8f9fa;
  color: #333;
}

.nav-tab.active {
  background: #667eea;
  color: white;
}

.nav-tab i {
  width: 16px;
  text-align: center;
}

.settings-content {
  flex: 1;
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.settings-section h3 {
  margin: 0 0 24px 0;
  color: #333;
  font-size: 1.5rem;
  font-weight: 600;
  padding-bottom: 12px;
  border-bottom: 2px solid #f0f0f0;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
}

.form-input, .form-select, .form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  transition: all 0.3s ease;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-weight: normal;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin: 0;
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
  padding: 30px;
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
  font-size: 32px;
  margin-bottom: 8px;
  color: #ddd;
}

.upload-placeholder p {
  margin: 0;
  font-size: 14px;
}

.image-preview {
  max-width: 100%;
  max-height: 150px;
  border-radius: 6px;
}

.test-email-btn {
  background: #28a745;
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
  margin-top: 16px;
}

.test-email-btn:hover {
  background: #218838;
  transform: translateY(-1px);
}

.settings-actions {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 12px;
}

.btn-save, .btn-reset {
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-save {
  background: #667eea;
  color: white;
}

.btn-save:hover {
  background: #5a6fd8;
  transform: translateY(-1px);
}

.btn-reset {
  background: #f8f9fa;
  color: #666;
  border: 1px solid #ddd;
}

.btn-reset:hover {
  background: #e9ecef;
  color: #333;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-settings {
    flex-direction: column;
  }
  
  .settings-nav {
    width: 100%;
  }
  
  .nav-tab {
    justify-content: center;
  }
  
  .settings-content {
    padding: 20px;
  }
  
  .settings-actions {
    flex-direction: column;
  }
  
  .btn-save, .btn-reset {
    justify-content: center;
  }
}
</style> 