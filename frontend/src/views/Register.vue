<template>
  <div class="register-page">
    <div class="register-card">
      <div class="register-header">
        <h1>F.Blog</h1>
        <p>创建您的账户</p>
      </div>

      <form class="register-form" @submit.prevent="onRegister">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.username" type="text" placeholder="3-20位字母数字下划线" required />
        </div>

        <div class="form-group">
          <label>邮箱</label>
          <input v-model="form.email" type="email" placeholder="请输入邮箱" required />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input v-model="form.password" type="password" placeholder="8-20位，包含字母和数字" required />
        </div>

        <div class="form-group">
          <label>确认密码</label>
          <input v-model="form.confirmPassword" type="password" placeholder="请再次输入密码" required />
        </div>

        <div class="form-group">
          <label>验证码</label>
          <div class="captcha-row">
            <input v-model="form.captcha" type="text" class="captcha-input" placeholder="请输入验证码" required />
            <img v-if="captchaImg" :src="captchaImg" class="captcha-img" alt="验证码" @click="getCaptcha" title="点击刷新" />
          </div>
        </div>

        <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>

        <button type="submit" class="register-btn" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>

        <div class="register-footer">
          已有账号？<router-link to="/login">立即登录</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { register } from '../api/auth';
import { useRouter } from 'vue-router';
import request from '../api/request';

const router = useRouter();
const loading = ref(false);

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  captcha: ''
});
const captchaId = ref('');
const captchaImg = ref('');
const errorMsg = ref('');

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

const onRegister = async () => {
  if (!/^[a-zA-Z0-9_]{3,20}$/.test(form.value.username)) {
    errorMsg.value = "用户名格式不正确";
    return;
  }
  if (!/^[\w.-]+@[\w.-]+\.\w+$/.test(form.value.email)) {
    errorMsg.value = "邮箱格式不正确";
    return;
  }
  if (!/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,20}$/.test(form.value.password)) {
    errorMsg.value = "密码需8-20位且包含字母和数字";
    return;
  }
  if (form.value.password !== form.value.confirmPassword) {
    errorMsg.value = "两次密码不一致";
    return;
  }
  if (!form.value.captcha) {
    errorMsg.value = "请输入验证码";
    return;
  }
  errorMsg.value = "";
  loading.value = true;
  try {
    await register({
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      captcha_id: captchaId.value,
      captcha: form.value.captcha
    });
    ElMessage.success('注册成功，请前往邮箱激活账号');
    router.push('/login');
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error || '注册失败';
    getCaptcha();
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: #f5f6fa;
}

.register-card {
  width: 100%;
  max-width: 420px;
  background: #fff;
  border-radius: 12px;
  padding: 40px 32px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 8px;
}

.register-header p {
  font-size: 15px;
  color: #6b7280;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
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

.error-msg {
  color: #ef4444;
  font-size: 13px;
  padding: 8px 12px;
  background: #fef2f2;
  border-radius: 6px;
  border: 1px solid #fecaca;
}

.register-btn {
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

.register-btn:hover {
  background: #4a59ae;
}

.register-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.register-footer {
  text-align: center;
  font-size: 14px;
  color: #6b7280;
  margin-top: 8px;
}

.register-footer a {
  color: #5b6abf;
  text-decoration: none;
  font-weight: 500;
}

.register-footer a:hover {
  text-decoration: underline;
}
</style>
