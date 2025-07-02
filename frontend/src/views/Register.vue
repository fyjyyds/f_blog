<template>
  <div class="auth-container">
    <div class="auth-background">
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>
    </div>
    
    <div class="auth-content">
      <div class="auth-card card">
        <div class="auth-header">
          <div class="logo-section">
            <div class="logo-icon">🌟</div>
            <h1 class="auth-title">加入我们</h1>
            <p class="auth-subtitle">创建您的 F.Blog 账户</p>
          </div>
        </div>
        
        <form class="auth-form" @submit.prevent="onRegister">
          <div class="form-group">
            <label class="form-label">用户名</label>
            <div class="input-wrapper">
              <span class="input-icon">👤</span>
              <input 
                v-model="form.username" 
                type="text" 
                class="form-input"
                placeholder="请输入用户名"
                required
              />
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">邮箱</label>
            <div class="input-wrapper">
              <span class="input-icon">📧</span>
              <input 
                v-model="form.email" 
                type="email" 
                class="form-input"
                placeholder="请输入邮箱"
                required
              />
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">密码</label>
            <div class="input-wrapper">
              <span class="input-icon">🔒</span>
              <input 
                v-model="form.password" 
                type="password" 
                class="form-input"
                placeholder="请输入密码"
                required
              />
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">确认密码</label>
            <div class="input-wrapper">
              <span class="input-icon">🔒</span>
              <input 
                v-model="form.confirmPassword" 
                type="password" 
                class="form-input"
                placeholder="请再次输入密码"
                required
              />
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">验证码</label>
            <div class="captcha-wrapper">
              <div class="input-wrapper captcha-input">
                <span class="input-icon">🛡️</span>
                <input 
                  v-model="form.captcha" 
                  type="text" 
                  class="form-input"
                  placeholder="验证码"
                  required
                />
              </div>
              <div class="captcha-image" @click="getCaptcha">
                <img :src="captchaImg" alt="验证码" />
                <div class="captcha-overlay">
                  <span>点击刷新</span>
                </div>
              </div>
            </div>
          </div>
          
          <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
          
          <button type="submit" class="submit-btn">
            <span class="btn-text">注册</span>
            <div class="btn-glow"></div>
          </button>
          
          <div class="auth-footer">
            <span class="footer-text">已有账号？</span>
            <router-link to="/login" class="footer-link">立即登录</router-link>
          </div>
        </form>
      </div>
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
    ElMessage.error('获取验证码失败');
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
  try {
    await register({
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      captcha_id: captchaId.value,
      captcha: form.value.captcha
    });
    ElMessage.success('注册成功，请前往邮箱验证');
    router.push('/login');
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error || '注册失败';
    getCaptcha();
  }
};
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.auth-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
}

.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(255, 119, 198, 0.1), rgba(120, 219, 255, 0.1));
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 80px;
  height: 80px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 120px;
  height: 120px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.shape-3 {
  width: 60px;
  height: 60px;
  top: 80%;
  left: 20%;
  animation-delay: 4s;
}

.shape-4 {
  width: 100px;
  height: 100px;
  top: 10%;
  right: 30%;
  animation-delay: 1s;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.auth-content {
  width: 100%;
  max-width: 450px;
  padding: 20px;
}

.auth-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 24px;
  padding: 40px;
  position: relative;
  overflow: hidden;
}

.auth-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, #ff77c6, transparent);
}

.auth-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.logo-icon {
  font-size: 48px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.auth-title {
  font-size: 32px;
  font-weight: 800;
  color: white;
  margin: 0;
  background: linear-gradient(135deg, #ff77c6 0%, #78daff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.auth-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  font-weight: 400;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.6);
  z-index: 1;
}

.form-input {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 16px 16px 16px 48px;
  color: white;
  font-size: 16px;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.form-input:focus {
  outline: none;
  border-color: #ff77c6;
  box-shadow: 0 0 20px rgba(255, 119, 198, 0.3);
  background: rgba(255, 255, 255, 0.15);
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.captcha-wrapper {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.captcha-input {
  flex: 1;
}

.captcha-image {
  width: 150px;
  height: 44px;
  position: relative;
  cursor: pointer;
  display: inline-block;
}

.captcha-image:hover {
  transform: scale(1.05);
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.captcha-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.captcha-overlay span {
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.captcha-image:hover .captcha-overlay {
  opacity: 1;
}

.submit-btn {
  position: relative;
  background: linear-gradient(135deg, #ff77c6 0%, #78daff 100%);
  border: none;
  border-radius: 12px;
  padding: 16px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  margin-top: 8px;
}

.submit-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.submit-btn:hover::before {
  left: 100%;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(255, 119, 198, 0.4);
}

.submit-btn:active {
  transform: translateY(0);
}

.btn-text {
  position: relative;
  z-index: 1;
}

.auth-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.footer-text {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

.footer-link {
  color: #ff77c6;
  text-decoration: none;
  font-weight: 600;
  margin-left: 4px;
  transition: color 0.3s ease;
}

.footer-link:hover {
  color: #78daff;
}

.error-msg {
  color: #e74c3c;
  margin-top: 8px;
  font-size: 14px;
  text-align: center;
}

@media (max-width: 768px) {
  .auth-content {
    padding: 16px;
  }
  
  .auth-card {
    padding: 24px;
  }
  
  .auth-title {
    font-size: 24px;
  }
  
  .captcha-wrapper {
    flex-direction: column;
    gap: 8px;
  }
  
  .captcha-image {
    width: 100%;
    height: 48px;
  }
}
</style>