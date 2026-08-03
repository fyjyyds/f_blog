<template>
  <nav class="navbar">
    <div class="navbar-inner">
      <router-link to="/" class="logo">F.Blog</router-link>

      <div class="nav-right">
        <template v-if="!isLogin">
          <router-link to="/login" class="nav-link">登录</router-link>
          <router-link to="/register" class="nav-btn">注册</router-link>
        </template>
        <template v-else>
          <router-link to="/create-article" class="nav-link">写文章</router-link>
          <div class="user-menu" @click="toggleMenu">
            <img :src="user.avatar || '/default-avatar.png'" class="avatar" alt="" />
            <div class="dropdown" v-if="showMenu" @click.stop>
              <div class="dropdown-item" @click="router.push('/user-center')">个人中心</div>
              <div class="dropdown-item" @click="router.push('/notifications')">通知</div>
              <div v-if="user.role === 'admin'" class="dropdown-item" @click="router.push('/admin')">后台管理</div>
              <div class="dropdown-divider"></div>
              <div class="dropdown-item danger" @click="logout">退出登录</div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import request from '../api/request';

const router = useRouter();
const isLogin = ref(!!localStorage.getItem('token'));
const user = ref<any>({});
const showMenu = ref(false);

const fetchUser = async () => {
  if (!isLogin.value) return;
  try {
    const res = await request.get('/api/v1/user/profile');
    user.value = res.data;
  } catch (e) {
    console.error('获取用户信息失败:', e);
  }
};

const toggleMenu = () => {
  showMenu.value = !showMenu.value;
};

const closeMenu = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.user-menu')) {
    showMenu.value = false;
  }
};

const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
  localStorage.removeItem('role');
  localStorage.removeItem('username');
  localStorage.removeItem('nickname');
  isLogin.value = false;
  user.value = {};
  showMenu.value = false;
  router.push('/');
};

onMounted(() => {
  fetchUser();
  document.addEventListener('click', closeMenu);
});

onUnmounted(() => {
  document.removeEventListener('click', closeMenu);
});
</script>

<style scoped>
.navbar {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 1000;
  height: 56px;
}

.navbar-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 100%;
}

.logo {
  font-size: 20px;
  font-weight: 700;
  color: #1a1a2e;
  text-decoration: none;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-link {
  font-size: 14px;
  color: #374151;
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background 0.2s;
}

.nav-link:hover {
  background: #f3f4f6;
}

.nav-btn {
  font-size: 14px;
  color: #fff;
  background: #5b6abf;
  text-decoration: none;
  padding: 7px 16px;
  border-radius: 6px;
  font-weight: 500;
  transition: background 0.2s;
}

.nav-btn:hover {
  background: #4a59ae;
}

.user-menu {
  position: relative;
  cursor: pointer;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #e5e7eb;
}

.dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 4px 0;
  min-width: 140px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.dropdown-item {
  padding: 8px 16px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition: background 0.15s;
}

.dropdown-item:hover {
  background: #f3f4f6;
}

.dropdown-item.danger {
  color: #ef4444;
}

.dropdown-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 4px 0;
}
</style>
