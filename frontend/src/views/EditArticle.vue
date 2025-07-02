<template>
  <div class="edit-article-container">
    <el-card class="edit-article-card">
      <el-button class="back-btn" @click="goBack" type="primary" :icon="ArrowLeft">返回</el-button>
      <h2>编辑文章</h2>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px" v-loading="loading">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" placeholder="请输入摘要" />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input type="textarea" v-model="form.content" :rows="10" placeholder="请输入内容" />
        </el-form-item>
        <el-form-item label="附件">
          <el-upload
            class="upload-demo"
            action="http://localhost:8080/api/v1/upload"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
            :headers="uploadHeaders"
            :before-upload="beforeUpload"
            name="file"
          >
            <el-button type="primary">上传文件</el-button>
            <span v-if="form.cover" class="file-info">
              <img v-if="isImage(form.cover)" :src="fullUrl(form.cover)" class="cover-preview" />
              <span v-else>已上传：{{ getFileName(form.cover) }}</span>
            </span>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">更新</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { ElMessage } from 'element-plus';
import { useRouter, useRoute } from 'vue-router';
import request from '../api/request';
import { getCategories } from '../api/category';
import { ArrowLeft } from '@element-plus/icons-vue';

const router = useRouter();
const route = useRoute();
const loading = ref(false);
const form = ref({
  title: '',
  summary: '',
  content: '',
  category_id: null,
  cover: '',
  status: ''
});
const formRef = ref();
const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
};

const categories = ref<{ id: number; name: string }[]>([]);

// 上传组件的认证头
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token');
  return token ? { Authorization: 'Bearer ' + token } : {};
});

onMounted(async () => {
  loading.value = true;
  try {
    // 获取分类列表
    const catRes = await getCategories();
    categories.value = catRes.data;
    
    // 获取文章详情
    const articleRes = await request.get(`/api/v1/articles/${route.params.id}`);
    const article = articleRes.data.data;
    
    // 填充表单
    form.value = {
      title: article.title,
      summary: article.summary || '',
      content: article.content,
      category_id: article.category_id,
      cover: article.cover || '',
      status: article.status || ''
    };
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '加载失败');
    router.push('/');
  } finally {
    loading.value = false;
  }
});

function handleUploadSuccess(res: any) {
  form.value.cover = res.url;
  ElMessage.success('上传成功');
}

function beforeUpload(file: File) {
  if (file.size > 10 * 1024 * 1024) {
    ElMessage.error('文件不能超过10MB');
    return false;
  }
  return true;
}

function isImage(url: string) {
  return /\.(jpg|jpeg|png|gif)$/i.test(url);
}
function getFileName(url: string) {
  return url.split('/').pop();
}
function fullUrl(url: string) {
  if (!url) return '';
  if (url.startsWith('http')) return url;
  return 'http://localhost:8080' + url;
}

function getAuthHeaders() {
  if (typeof window !== 'undefined' && window.localStorage) {
    const token = window.localStorage.getItem('token');
    return token ? { Authorization: 'Bearer ' + token } : {};
  }
  return {};
}

const onSubmit = () => {
  if (!form.value.status) {
    form.value.status = 'published';
  }
  formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        await request.put(`/api/v1/articles/${route.params.id}`, form.value);
        ElMessage.success('更新成功');
        router.push(`/article/${route.params.id}`);
      } catch (e: any) {
        ElMessage.error(e.response?.data?.error || '更新失败');
      }
    }
  });
};

function goBack() {
  router.back();
}
</script>

<style scoped>
.edit-article-container {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f6fa;
}
.edit-article-card {
  width: 600px;
  padding: 32px 24px 12px 24px;
}
.back-btn {
  margin-bottom: 12px;
}
.cover-preview {
  width: 80px;
  height: 80px;
  object-fit: cover;
  margin-left: 12px;
  border-radius: 4px;
  vertical-align: middle;
}
.file-info {
  margin-left: 12px;
  color: #666;
}
</style> 
 