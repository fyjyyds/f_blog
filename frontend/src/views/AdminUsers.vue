<template>
  <div class="admin-users">
    <!-- 搜索和操作栏 -->
    <div class="action-bar">
      <div class="search-section">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索用户名或邮箱..."
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
          <option value="active">活跃</option>
          <option value="inactive">非活跃</option>
          <option value="banned">已封禁</option>
        </select>
        
        <select v-model="roleFilter" @change="handleFilter" class="filter-select">
          <option value="">全部角色</option>
          <option value="admin">管理员</option>
          <option value="user">普通用户</option>
          <option value="moderator">版主</option>
        </select>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="users-table">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>头像</th>
            <th>用户名</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>状态</th>
            <th>注册时间</th>
            <th>最后登录</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.id">
            <td>{{ user.id }}</td>
            <td>
              <img 
                :src="user.avatar || '/static/avatar/default.jpg'" 
                :alt="user.username"
                class="user-avatar"
              />
            </td>
            <td>{{ user.username }}</td>
            <td>{{ user.email }}</td>
            <td>
              <span class="role-badge" :class="user.role">
                {{ getRoleText(user.role) }}
              </span>
            </td>
            <td>
              <span class="status-badge" :class="user.status">
                {{ getStatusText(user.status) }}
              </span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td>{{ formatDate(user.last_login_at) }}</td>
            <td>
              <div class="action-buttons">
                <button class="btn-edit" @click="editUser(user)">
                  <i class="fas fa-edit"></i>
                </button>
                <button 
                  v-if="user.status !== 'banned'" 
                  class="btn-ban" 
                  @click="banUser(user.id)"
                >
                  <i class="fas fa-ban"></i>
                </button>
                <button 
                  v-else 
                  class="btn-unban" 
                  @click="unbanUser(user.id)"
                >
                  <i class="fas fa-unlock"></i>
                </button>
                <button class="btn-delete" @click="deleteUser(user.id)">
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

    <!-- 编辑用户弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑用户</h3>
          <button class="close-btn" @click="closeEditModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="editForm.username" type="text" class="form-input" />
          </div>
          
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="editForm.email" type="email" class="form-input" />
          </div>
          
          <div class="form-group">
            <label>角色</label>
            <select v-model="editForm.role" class="form-select">
              <option value="user">普通用户</option>
              <option value="moderator">版主</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          
          <div class="form-group">
            <label>状态</label>
            <select v-model="editForm.status" class="form-select">
              <option value="active">活跃</option>
              <option value="inactive">非活跃</option>
              <option value="banned">已封禁</option>
            </select>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeEditModal">取消</button>
          <button class="btn-save" @click="saveUser">保存</button>
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
const users = ref<any[]>([]);
const searchQuery = ref('');
const statusFilter = ref('');
const roleFilter = ref('');
const currentPage = ref(1);
const totalPages = ref(1);
const showEditModal = ref(false);
const editForm = ref<any>({});

// 获取用户列表
const fetchUsers = async () => {
  try {
    const params = {
      page: currentPage.value,
      limit: 10,
      search: searchQuery.value,
      status: statusFilter.value,
      role: roleFilter.value
    };
    
    const response = await request.get('/api/v1/admin/users', {
      params,
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    users.value = response.data.users || [];
    totalPages.value = response.data.totalPages || 1;
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '获取用户列表失败');
  }
};

// 过滤用户
const filteredUsers = computed(() => {
  let filtered = users.value;
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(user => 
      user.username.toLowerCase().includes(query) ||
      user.email.toLowerCase().includes(query)
    );
  }
  
  if (statusFilter.value) {
    filtered = filtered.filter(user => user.status === statusFilter.value);
  }
  
  if (roleFilter.value) {
    filtered = filtered.filter(user => user.role === roleFilter.value);
  }
  
  return filtered;
});

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1;
  fetchUsers();
};

// 过滤处理
const handleFilter = () => {
  currentPage.value = 1;
  fetchUsers();
};

// 分页处理
const changePage = (page: number) => {
  currentPage.value = page;
  fetchUsers();
};

// 编辑用户
const editUser = (user: any) => {
  editForm.value = { ...user };
  showEditModal.value = true;
};

// 保存用户
const saveUser = async () => {
  try {
    await request.put(`/api/v1/admin/users/${editForm.value.id}`, editForm.value, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('用户信息更新成功');
    closeEditModal();
    fetchUsers();
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.error || '更新失败');
  }
};

// 封禁用户
const banUser = async (userId: number) => {
  try {
    await ElMessageBox.confirm('确定要封禁该用户吗？', '确认操作', {
      type: 'warning'
    });
    
    await request.post(`/api/v1/admin/users/${userId}/ban`, {}, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('用户已被封禁');
    fetchUsers();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.error || '操作失败');
    }
  }
};

// 解封用户
const unbanUser = async (userId: number) => {
  try {
    await ElMessageBox.confirm('确定要解封该用户吗？', '确认操作', {
      type: 'warning'
    });
    
    await request.post(`/api/v1/admin/users/${userId}/unban`, {}, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('用户已被解封');
    fetchUsers();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.error || '操作失败');
    }
  }
};

// 删除用户
const deleteUser = async (userId: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？此操作不可恢复！', '确认删除', {
      type: 'error'
    });
    
    await request.delete(`/api/v1/admin/users/${userId}`, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    ElMessage.success('用户已被删除');
    fetchUsers();
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

// 获取角色文本
const getRoleText = (role: string) => {
  const roleMap: { [key: string]: string } = {
    admin: '管理员',
    moderator: '版主',
    user: '用户'
  };
  return roleMap[role] || role;
};

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: { [key: string]: string } = {
    active: '活跃',
    inactive: '非活跃',
    banned: '已封禁'
  };
  return statusMap[status] || status;
};

onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.admin-users {
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

.users-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.role-badge, .status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.role-badge.admin {
  background: #ff6b6b;
  color: white;
}

.role-badge.moderator {
  background: #4ecdc4;
  color: white;
}

.role-badge.user {
  background: #45b7d1;
  color: white;
}

.status-badge.active {
  background: #51cf66;
  color: white;
}

.status-badge.inactive {
  background: #adb5bd;
  color: white;
}

.status-badge.banned {
  background: #ff6b6b;
  color: white;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-edit, .btn-ban, .btn-unban, .btn-delete {
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

.btn-ban {
  color: #ff6b6b;
}

.btn-ban:hover {
  background: rgba(255, 107, 107, 0.1);
}

.btn-unban {
  color: #51cf66;
}

.btn-unban:hover {
  background: rgba(81, 207, 102, 0.1);
}

.btn-delete {
  color: #ff6b6b;
}

.btn-delete:hover {
  background: rgba(255, 107, 107, 0.1);
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

.form-input, .form-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.form-input:focus, .form-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
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
  
  .users-table {
    overflow-x: auto;
  }
  
  table {
    min-width: 800px;
  }
}
</style> 