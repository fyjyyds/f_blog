<template>
  <div class="create-article-page">
    <div class="card create-article-card">
      <button class="button back-btn" @click="goBack">返回</button>
      <h2 class="form-title">发布新文章</h2>
      <form @submit.prevent="onSubmit">
        <div class="form-group">
          <label>标题</label>
          <input v-model="form.title" class="form-input" placeholder="请输入标题" />
        </div>
        <div class="form-group">
          <label>摘要</label>
          <input v-model="form.summary" class="form-input" placeholder="请输入摘要" />
        </div>
        <div class="form-group">
          <label>分类</label>
          <select v-model="form.category_id" class="form-input">
            <option disabled value="">请选择分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>标签</label>
          <el-select
            v-model="form.tag_ids"
            multiple
            filterable
            placeholder="请选择标签"
            class="form-input"
            style="width: 100%;"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </div>
        <div class="form-group">
          <label>内容</label>
          <textarea v-model="form.content" class="form-input" rows="8" placeholder="请输入内容"></textarea>
        </div>
        <div class="form-group">
          <label>附件</label>
          <input type="file" class="form-input" @change="onFileChange" />
          <div v-if="form.cover" style="margin-top: 10px;">
            <img v-if="/\.(jpg|jpeg|png|gif)$/i.test(form.cover)" :src="form.cover" class="cover-preview" />
            <a v-else :href="form.cover" target="_blank" class="file-info">下载附件</a>
          </div>
        </div>
        <div class="form-actions">
          <button class="button" type="submit">发布</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';
import { createArticle } from '../api/article';
import { getCategories, getTags } from '../api/category';
import { ArrowLeft } from '@element-plus/icons-vue';
import { ElSelect, ElOption } from 'element-plus';
import 'element-plus/dist/index.css';
import request from '../api/request';

const router = useRouter();
const form = ref({
  title: '',
  summary: '',
  category_id: undefined,
  tag_ids: [],
  content: '',
  cover: ''
});
const formRef = ref();
const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
};

const categories = ref<{ id: number; name: string }[]>([]);
const tags = ref<{ id: number; name: string }[]>([]);

const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token');
  return token ? { Authorization: 'Bearer ' + token } : {};
});

onMounted(async () => {
  const res = await getCategories();
  categories.value = res.data;
  const tagRes = await getTags();
  tags.value = tagRes.data;
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

async function onSubmit() {
  try {
    // 组装参数，类型转换
    const payload = {
      ...form.value,
      category_id: Number(form.value.category_id),
      tag_ids: form.value.tag_ids.map(Number)
    };
    const res = await createArticle(payload);
    const data = res.data;
    if (!data.error && data.id) {
      if (data.status === 'pending') {
        ElMessage.success('文章已提交，等待管理员审核');
      } else {
        ElMessage.success('发布成功');
      }
      router.push('/');
    } else {
      ElMessage.error(data.error || '发布失败');
    }
  } catch (err) {
    ElMessage.error('网络错误或服务器异常');
  }
}

function goBack() {
  router.back();
}

async function onFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (file) {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await request.post('/api/v1/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      if (res.data && res.data.url) {
        form.value.cover = res.data.url;
        ElMessage.success('上传成功');
      } else {
        ElMessage.error('上传失败');
      }
    } catch (err) {
      ElMessage.error('上传失败');
    }
  }
}
</script>

<style scoped>
.create-article-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--background);
}
.create-article-card {
  width: 100%;
  max-width: 600px;
  margin: 40px auto;
  padding: 40px 32px;
  border-radius: var(--border-radius);
  background: var(--card-bg);
  box-shadow: var(--shadow);
  backdrop-filter: var(--card-blur);
  position: relative;
}
.form-title {
  text-align: center;
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 32px;
  color: #fff;
}
.form-group label {
  color: #bfcfff;
  font-weight: 500;
  margin-bottom: 8px;
  display: block;
}
.form-group {
  margin-bottom: 24px;
}
.form-actions {
  text-align: center;
  margin-top: 32px;
}
.back-btn {
  position: absolute;
  left: 32px;
  top: 32px;
  background: transparent;
  color: var(--primary);
  font-size: 16px;
  border: none;
  cursor: pointer;
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