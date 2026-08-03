<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1>F.Blog</h1>
        <p>登录您的账户</p>
      </div>

      <form class="login-form" @submit.prevent="onLogin">
        <div class="form-group">
          <label>账号</label>
          <input
            v-model="form.username_or_email"
            type="text"
            placeholder="用户名或邮箱"
            required
          />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            required
          />
        </div>

        <div class="form-group">
          <label>验证码</label>
          <div class="captcha-row">
            <input
              v-model="form.captcha"
              type="text"
              class="captcha-input"
              placeholder="请输入验证码"
              required
            />
            <img
              v-if="captchaImg"
              :src="captchaImg"
              class="captcha-img"
              alt="验证码"
              @click="getCaptcha"
              title="点击刷新"
            />
          </div>
        </div>

        <button type="submit" class="login-btn" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>

        <div class="login-footer">
          还没有账号？<router-link to="/register">立即注册</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { login } from '../api/auth';
import { useRouter } from 'vue-router';
import request from '../api/request';

const router = useRouter();
const loading = ref(false);

const form = ref({
  username_or_email: '',
  password: '',
  captcha: ''
});
const captchaId = ref('');
const captchaImg = ref('');

async function getCaptcha() {
  try {
    const res = await request.get('/api/v1/captcha');
    captchaId.value = res.data.captcha_id;
    let img = res.data.captcha_image;
    if (!img.startsWith('data:image/png;base64,')) {
      img = 'data:image/png;base64,' + img;
    }
    captchaImg.value = img;
  } catch (error) {
    console.error('获取验证码失败:', error);
  }
}

onMounted(getCaptcha);

const onLogin = async () => {
  if (!form.value.username_or_email || !form.value.password || !form.value.captcha) {
    ElMessage.error('请填写完整信息');
    return;
  }

  loading.value = true;
  try {
    const res = await login({
      username_or_email: form.value.username_or_email,
      password: form.value.password,
      captcha_id: captchaId.value,
      captcha: form.value.captcha
    });

    if (res.data && res.data.token) {
      localStorage.setItem('token', res.data.token);
      if (res.data.user) {
        localStorage.setItem('user_id', res.data.user.id.toString());
        localStorage.setItem('role', res.data.user.role);
        localStorage.setItem('username', res.data.user.username);
        localStorage.setItem('nickname', res.data.user.nickname || '');
      }
      ElMessage.success('登录成功');
      router.push('/');
    }
  } catch (e: any) {
    const errorMsg = e.response?.data?.error || '登录失败';
    ElMessage.error(errorMsg);
    getCaptcha();
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: #f5f6fa;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 12px;
  padding: 40px 32px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 15px;
  color: #6b7280;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  color: #1a1a2e;
  background: #fff;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #5b6abf;
  box-shadow: 0 0 0 3px rgba(91, 106, 191, 0.1);
}

.form-group input::placeholder {
  color: #9ca3af;
}

.captcha-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.captcha-input {
  flex: 1;
}

.captcha-img {
  height: 40px;
  min-width: 120px;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid #e5e7eb;
  object-fit: contain;
  background: #fff;
}

.captcha-img:hover {
  border-color: #5b6abf;
}

.login-btn {
  width: 100%;
  padding: 11px 0;
  background: #5b6abf;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  margin-top: 4px;
}

.login-btn:hover {
  background: #4a59ae;
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  font-size: 14px;
  color: #6b7280;
  margin-top: 8px;
}

.login-footer a {
  color: #5b6abf;
  text-decoration: none;
  font-weight: 500;
}

.login-footer a:hover {
  text-decoration: underline;
}
</style>
