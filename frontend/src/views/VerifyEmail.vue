<template>
  <div class="verify-result">
    <h2 v-if="status==='success'">激活成功！</h2>
    <h2 v-else-if="status==='fail'">激活失败</h2>
    <h2 v-else>正在验证...</h2>
    <p v-if="status==='success'">请返回 <router-link to="/login">登录</router-link></p>
    <p v-else-if="status==='fail'">链接无效或已过期，请重新注册或联系管理员。</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const status = ref('')
const route = useRoute()
const router = useRouter()

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    status.value = 'fail'
    return
  }
  try {
    await axios.get('/api/v1/activate', { params: { token } })
    status.value = 'success'
  } catch (e) {
    status.value = 'fail'
  }
})
</script>

<style scoped>
.verify-result {
  max-width: 400px;
  margin: 100px auto;
  padding: 32px;
  background: #fff;
  border-radius: 16px;
  text-align: center;
  box-shadow: 0 4px 24px rgba(0,0,0,0.08);
}
</style> 