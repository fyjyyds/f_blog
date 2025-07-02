<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <h2 v-if="!sidebarCollapsed">后台管理</h2>
        <h2 v-else>管理</h2>
      </div>
      
      <nav class="sidebar-nav">
        <router-link 
          to="/"
          class="nav-item"
          style="margin-bottom: 10px;"
        >
          <i class="fas fa-home"></i>
          <span v-if="!sidebarCollapsed">返回首页</span>
        </router-link>
        <router-link 
          v-for="item in menuItems" 
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: $route.path.startsWith(item.path) }"
        >
          <i :class="item.icon"></i>
          <span v-if="!sidebarCollapsed">{{ item.title }}</span>
        </router-link>
      </nav>
      
      <div class="sidebar-footer">
        <button class="collapse-btn" @click="toggleSidebar">
          <i :class="sidebarCollapsed ? 'fas fa-expand' : 'fas fa-compress'"></i>
        </button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content" :class="{ expanded: sidebarCollapsed }">
      <!-- 顶部导航栏 -->
      <header class="top-header">
        <div class="header-left">
          <button class="menu-toggle" @click="toggleSidebar">
            <i class="fas fa-bars"></i>
          </button>
          <h1>{{ currentPageTitle }}</h1>
        </div>
        
        <div class="header-right">
          <div class="user-info">
            <span>{{ userInfo.username }}</span>
            <button class="logout-btn" @click="logout">
              <i class="fas fa-sign-out-alt"></i>
              退出
            </button>
          </div>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();

const sidebarCollapsed = ref(false);
const userInfo = ref({
  username: '管理员',
  avatar: ''
});

// 菜单项配置
const menuItems = [
  { path: '/admin/dashboard', title: '仪表盘', icon: 'fas fa-tachometer-alt' },
  { path: '/admin/articles', title: '文章管理', icon: 'fas fa-newspaper' },
  { path: '/admin/users', title: '用户管理', icon: 'fas fa-users' },
  { path: '/admin/categories', title: '分类管理', icon: 'fas fa-folder' },
  { path: '/admin/tags', title: '标签管理', icon: 'fas fa-tags' },
  { path: '/admin/comments', title: '评论管理', icon: 'fas fa-comments' },
  { path: '/admin/banners', title: '轮播图管理', icon: 'fas fa-images' },
  { path: '/admin/settings', title: '系统设置', icon: 'fas fa-cog' }
];

// 计算当前页面标题
const currentPageTitle = computed(() => {
  const currentItem = menuItems.find(item => route.path.startsWith(item.path));
  return currentItem ? currentItem.title : '后台管理';
});

// 切换侧边栏
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value;
};

// 退出登录
const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  ElMessage.success('已退出登录');
  router.push('/login');
};

// 获取用户信息
const getUserInfo = () => {
  const user = localStorage.getItem('user');
  if (user) {
    userInfo.value = JSON.parse(user);
  }
};

onMounted(() => {
  getUserInfo();
});
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: #f5f7fa;
}

.sidebar {
  width: 250px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transition: all 0.3s ease;
  position: fixed;
  height: 100vh;
  z-index: 1000;
}

.sidebar.collapsed {
  width: 60px;
}

.sidebar-header {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-header h2 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.sidebar-nav {
  padding: 20px 0;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border-left-color: rgba(255, 255, 255, 0.5);
}

.nav-item.active {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border-left-color: white;
}

.nav-item i {
  width: 20px;
  margin-right: 12px;
  text-align: center;
}

.nav-item span {
  font-size: 14px;
  font-weight: 500;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.collapse-btn {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.main-content {
  flex: 1;
  margin-left: 250px;
  transition: all 0.3s ease;
}

.main-content.expanded {
  margin-left: 60px;
}

.top-header {
  background: white;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 999;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.menu-toggle {
  background: none;
  border: none;
  font-size: 18px;
  color: #666;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.menu-toggle:hover {
  background: #f5f5f5;
  color: #333;
}

.header-left h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info span {
  font-weight: 500;
  color: #333;
}

.logout-btn {
  background: #f56c6c;
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
}

.logout-btn:hover {
  background: #e74c3c;
}

.content {
  padding: 24px;
  min-height: calc(100vh - 64px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
  }
  
  .sidebar.collapsed {
    transform: translateX(0);
    width: 250px;
  }
  
  .main-content {
    margin-left: 0;
  }
  
  .main-content.expanded {
    margin-left: 0;
  }
}
</style> 