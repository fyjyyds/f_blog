<template>
  <nav class="navbar">
    <div class="navbar-container">
      <div class="logo-section">
        <div class="logo">
          <span class="logo-text">F</span>
          <span class="logo-dot"></span>
          <span class="logo-text">Blog</span>
        </div>
        <div class="logo-subtitle">Future Blog Platform</div>
      </div>
      
      <div class="nav-actions">
        <template v-if="!isLogin">
          <button class="nav-btn login-btn" @click="goLogin">
            <span class="btn-text">登录</span>
            <div class="btn-glow"></div>
          </button>
          <button class="nav-btn register-btn" @click="goRegister">
            <span class="btn-text">注册</span>
            <div class="btn-glow"></div>
          </button>
        </template>
        <template v-else>
          <button class="nav-btn create-btn" @click="goCreateArticle">
            <span class="btn-icon">✏️</span>
            <span class="btn-text">发布文章</span>
            <div class="btn-glow"></div>
          </button>
          <div class="user-section">
            <div class="user-avatar" @click="toggleUserMenu($event)">
              <img :src="user.avatar || '/default-avatar.png'" :alt="user.nickname || user.username" />
              <div class="avatar-ring"></div>
            </div>
            <div class="user-menu" v-if="showUserMenu" @click.stop>
              <div class="menu-item" @click="goUserCenter">
                <span class="menu-icon">👤</span>
                <span>个人中心</span>
              </div>
              <div class="menu-item" @click="goNotifications">
                <span class="menu-icon">🔔</span>
                <span>通知</span>
              </div>
              <div v-if="user.role === 'admin'" class="menu-item" @click="goAdmin">
                <span class="menu-icon">⚙️</span>
                <span>后台管理</span>
              </div>
              <div class="menu-divider"></div>
              <div class="menu-item logout" @click="logout">
                <span class="menu-icon">🚪</span>
                <span>退出登录</span>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ref, computed, onMounted, onUnmounted } from 'vue';
import request from '../api/request';

const router = useRouter();
const username = ref(localStorage.getItem('username') || '');
const isLogin = ref(!!localStorage.getItem('token'));
const user = ref<any>({});
const showUserMenu = ref(false);

const fetchUser = async () => {
  if (!isLogin.value) return;
  try {
    const res = await request.get('/api/v1/user/profile');
    user.value = res.data;
  } catch (error) {
    console.error('获取用户信息失败:', error);
  }
};

const toggleUserMenu = (event: MouseEvent) => {
  event.stopPropagation();
  showUserMenu.value = !showUserMenu.value;
};

const closeUserMenu = () => {
  showUserMenu.value = false;
};

onMounted(() => {
  fetchUser();
  document.addEventListener('click', closeUserMenu);
});

onUnmounted(() => {
  document.removeEventListener('click', closeUserMenu);
});

const goLogin = () => router.push('/login');
const goRegister = () => router.push('/register');
const goCreateArticle = () => router.push('/create-article');
const goUserCenter = () => {
  router.push('/user-center');
  showUserMenu.value = false;
};
const goNotifications = () => {
  router.push('/notifications');
  showUserMenu.value = false;
};
const goAdmin = () => {
  router.push('/admin');
  showUserMenu.value = false;
};
const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
  localStorage.removeItem('role');
  isLogin.value = false;
  showUserMenu.value = false;
  router.push('/login');
};
</script>

<style scoped>
.navbar {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  position: sticky;
  top: 0;
  z-index: 1000;
  padding: 0;
  height: 80px;
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 100%;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.logo {
  display: flex;
  align-items: center;
  gap: 4px;
  position: relative;
}

.logo-text {
  font-size: 28px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 1px;
}

.logo-dot {
  width: 8px;
  height: 8px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.2); opacity: 0.8; }
}

.logo-subtitle {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  margin-top: -4px;
  font-weight: 300;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-btn {
  position: relative;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 12px 20px;
  color: white;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 8px;
  backdrop-filter: blur(10px);
}

.nav-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left 0.5s;
}

.nav-btn:hover::before {
  left: 100%;
}

.nav-btn:hover {
  transform: translateY(-2px);
  border-color: rgba(102, 126, 234, 0.5);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.login-btn {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2), rgba(118, 75, 162, 0.2));
}

.register-btn {
  background: linear-gradient(135deg, rgba(255, 119, 198, 0.2), rgba(120, 219, 255, 0.2));
}

.create-btn {
  background: linear-gradient(135deg, rgba(120, 219, 255, 0.2), rgba(102, 126, 234, 0.2));
}

.btn-icon {
  font-size: 16px;
}

.btn-text {
  position: relative;
  z-index: 1;
}

.user-section {
  position: relative;
}

.user-avatar {
  position: relative;
  cursor: pointer;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  transition: all 0.3s ease;
}

.user-avatar:hover {
  transform: scale(1.1);
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.avatar-ring {
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  border: 2px solid transparent;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2) border-box;
  -webkit-mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: destination-out;
  mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0);
  mask-composite: exclude;
  animation: rotate 3s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.user-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 8px 0;
  min-width: 160px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.menu-item.logout:hover {
  background: rgba(255, 119, 119, 0.2);
  color: #ff7777;
}

.menu-icon {
  font-size: 16px;
}

.menu-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.2);
  margin: 4px 0;
}

@media (max-width: 768px) {
  .navbar-container {
    padding: 0 16px;
  }
  
  .logo-text {
    font-size: 24px;
  }
  
  .logo-subtitle {
    display: none;
  }
  
  .nav-btn {
    padding: 10px 16px;
    font-size: 14px;
  }
  
  .btn-text {
    display: none;
  }
  
  .user-avatar {
    width: 40px;
    height: 40px;
  }
}
</style>
  