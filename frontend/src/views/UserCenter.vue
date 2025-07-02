<template>
  <el-button class="back-btn" type="primary" @click="goBack" round>
    <el-icon><ArrowLeft /></el-icon>
    返回
  </el-button>
  <div class="user-center">
    <!-- 个人信息卡片 -->
    <el-card class="profile-card">
      <div class="profile-header">
        <div class="avatar-section">
          <div class="avatar-container" @click="showAvatarDialog = true">
            <el-avatar :src="user.avatar" size="large" class="clickable-avatar" />
            <div class="avatar-overlay">
              <el-icon class="avatar-icon"><Camera /></el-icon>
              <span class="avatar-text">点击更换头像</span>
            </div>
          </div>
          <div class="avatar-tip">点击头像更换</div>
        </div>
        <div class="profile-info">
          <h2>{{ user.nickname || user.username }}</h2>
          <div class="profile-meta">
            <span>邮箱：{{ user.email }}</span>
            <span>性别：{{ genderMap[user.gender] || '-' }}</span>
            <span>生日：{{ user.birthday ? formatDate(user.birthday) : '-' }}</span>
            <span>网站：<a :href="user.website" target="_blank">{{ user.website }}</a></span>
          </div>
          <div class="profile-bio">{{ user.bio }}</div>
        </div>
        <el-button type="primary" @click="showEdit = true">编辑资料</el-button>
        <el-button type="warning" style="margin-left: 12px;" @click="showChangePwd = true">修改密码</el-button>
      </div>
    </el-card>

    <!-- 统计信息 -->
    <div class="stat-bar">
      <div>文章 {{ stats.articles }}</div>
      <div>评论 {{ stats.comments }}</div>
      <div>点赞 {{ stats.likes }}</div>
      <div>关注 {{ stats.followings }}</div>
      <div>粉丝 {{ stats.followers }}</div>
    </div>

    <!-- Tab切换 -->
    <div class="user-tabs-card">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="我的文章" name="articles">
          <div class="user-table-card">
            <el-table :data="articles" style="width: 100%">
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column prop="title" label="标题" />
              <el-table-column prop="created_at" label="创建时间" :formatter="(row: any) => formatDate(row.created_at)" />
            </el-table>
          </div>
        </el-tab-pane>
        <el-tab-pane label="我的评论" name="comments">
          <div class="user-table-card">
            <el-table :data="comments" style="width: 100%">
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column prop="content" label="内容" />
              <el-table-column prop="created_at" label="时间" :formatter="(row: any) => formatDate(row.created_at)" />
            </el-table>
          </div>
        </el-tab-pane>
        <el-tab-pane label="我的点赞" name="likes">
          <div class="user-table-card">
            <el-table :data="likes" style="width: 100%">
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column prop="target_type" label="类型" />
              <el-table-column prop="target_id" label="目标ID" />
              <el-table-column prop="created_at" label="时间" :formatter="(row: any) => formatDate(row.created_at)" />
            </el-table>
          </div>
        </el-tab-pane>
        <el-tab-pane label="我的关注" name="followings">
          <div class="user-table-card">
            <el-table :data="followings" style="width: 100%">
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column prop="following_id" label="用户ID" />
              <el-table-column label="用户名/昵称" :formatter="(row: any) => row.following_nickname || row.following_username" />
              <el-table-column prop="created_at" label="关注时间" :formatter="(row: any) => formatDate(row.created_at)" />
            </el-table>
          </div>
        </el-tab-pane>
        <el-tab-pane label="我的粉丝" name="followers">
          <div class="user-table-card">
            <el-table :data="followers" style="width: 100%">
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column prop="follower_id" label="用户ID" />
              <el-table-column label="用户名/昵称" :formatter="(row: any) => row.follower_nickname || row.follower_username" />
              <el-table-column prop="created_at" label="关注时间" :formatter="(row: any) => formatDate(row.created_at)" />
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 编辑资料弹窗 -->
    <el-dialog v-model="showEdit" title="编辑个人资料">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称"><el-input v-model="editForm.nickname" /></el-form-item>
        <el-form-item label="简介"><el-input v-model="editForm.bio" type="textarea" /></el-form-item>
        <el-form-item label="性别">
          <el-select v-model="editForm.gender" placeholder="请选择">
            <el-option label="男" value="male" />
            <el-option label="女" value="female" />
            <el-option label="保密" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="生日"><el-date-picker v-model="editForm.birthday" type="date" value-format="YYYY-MM-DD" /></el-form-item>
        <el-form-item label="网站"><el-input v-model="editForm.website" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEdit = false">取消</el-button>
        <el-button type="primary" @click="submitEdit">保存</el-button>
      </template>
    </el-dialog>
    <!-- 更换头像弹窗 -->
    <el-dialog v-model="showAvatarDialog" title="更换头像" width="400px">
      <div class="avatar-upload-container">
        <el-avatar :src="user.avatar" size="large" class="preview-avatar" />
        <el-upload
          class="avatar-uploader"
          action="/api/v1/user/avatar"
          :show-file-list="false"
          :headers="uploadHeaders"
          :on-success="onAvatarSuccess"
          :on-error="onAvatarError"
          :before-upload="beforeAvatarUpload"
          name="avatar"
          drag
        >
          <div class="upload-content-only">
            <el-icon class="upload-icon"><Upload /></el-icon>
            <div class="upload-text">
              <span>点击或拖拽文件到此处上传</span>
              <span class="upload-hint">支持 JPG、PNG 格式，文件大小不超过 2MB</span>
            </div>
          </div>
        </el-upload>
      </div>
      <template #footer>
        <el-button @click="showAvatarDialog = false">取消</el-button>
        <el-button type="primary" @click="showAvatarDialog = false">确定</el-button>
      </template>
    </el-dialog>
    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showChangePwd" title="修改密码" width="400px">
      <el-form :model="pwdForm" label-width="90px" :rules="pwdRules" ref="pwdFormRef">
        <el-form-item label="原密码" prop="old_password">
          <el-input v-model="pwdForm.old_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="pwdForm.new_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input v-model="pwdForm.confirm_password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showChangePwd = false">取消</el-button>
        <el-button type="primary" @click="submitChangePwd">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, onUnmounted } from 'vue';
import { ElMessage } from 'element-plus';
import { Camera, Upload, ArrowLeft } from '@element-plus/icons-vue';
import { useRouter } from 'vue-router';
import request from '../api/request';
import { formatDate } from '../utils/date';

const router = useRouter();
function goBack() {
  router.back();
}

const user = ref<any>({});
const stats = ref({ articles: 0, comments: 0, likes: 0, followings: 0, followers: 0 });
const articles = ref<any[]>([]);
const comments = ref<any[]>([]);
const likes = ref<any[]>([]);
const followings = ref<any[]>([]);
const followers = ref<any[]>([]);
const activeTab = ref('articles');
const showEdit = ref(false);
const showAvatarDialog = ref(false);
const showChangePwd = ref(false);
const editForm = ref<any>({});
const pwdForm = ref({ old_password: '', new_password: '', confirm_password: '' });
const pwdFormRef = ref();
const pwdRules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码至少8位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== pwdForm.value.new_password) {
          callback(new Error('两次输入密码不一致'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ]
};

const genderMap: Record<string, string> = {
  male: '男',
  female: '女',
  other: '保密'
};

// 动态计算上传headers，确保token是最新的
const uploadHeaders = computed(() => ({
  Authorization: localStorage.getItem('token') ? 'Bearer ' + localStorage.getItem('token') : ''
}));

onMounted(async () => {
  await fetchProfile();
  await fetchStats();
  await fetchTabData('articles');
  window.addEventListener('refresh-user-stats', refreshLikesStatsAndTab);
});

function refreshLikesStatsAndTab() {
  fetchStats();
  if (activeTab.value === 'likes') {
    fetchTabData('likes');
  }
}

onUnmounted(() => {
  window.removeEventListener('refresh-user-stats', refreshLikesStatsAndTab);
});

const fetchProfile = async () => {
  try {
    const res = await request.get('/api/v1/user/profile');
    user.value = res.data;
    editForm.value = { ...user.value };
  } catch (error: any) {
    ElMessage.error('获取用户信息失败');
    console.error('Fetch profile error:', error);
  }
};

const fetchStats = async () => {
  try {
    // 简单统计，实际可后端聚合
    const [a, c, l, f1, f2] = await Promise.all([
      request.get('/api/v1/user/articles'),
      request.get('/api/v1/user/comments'),
      request.get('/api/v1/user/likes?type=article'),
      request.get('/api/v1/user/followings'),
      request.get('/api/v1/user/followers'),
    ]);
    stats.value.articles = a.data.length;
    stats.value.comments = c.data.length;
    stats.value.likes = l.data.length;
    stats.value.followings = f1.data.length;
    stats.value.followers = f2.data.length;
  } catch (error: any) {
    ElMessage.error('获取统计数据失败');
    console.error('Fetch stats error:', error);
  }
};

const fetchTabData = async (tab: string) => {
  try {
    if (tab === 'articles') articles.value = (await request.get('/api/v1/user/articles')).data;
    if (tab === 'comments') comments.value = (await request.get('/api/v1/user/comments')).data;
    if (tab === 'likes') likes.value = (await request.get('/api/v1/user/likes?type=article')).data;
    if (tab === 'followings') followings.value = (await request.get('/api/v1/user/followings')).data;
    if (tab === 'followers') followers.value = (await request.get('/api/v1/user/followers')).data;
  } catch (error: any) {
    ElMessage.error(`获取${tab}数据失败`);
    console.error(`Fetch ${tab} error:`, error);
  }
};

watch(activeTab, (tab) => {
  fetchTabData(tab);
});

const submitEdit = async () => {
  try {
    await request.put('/api/v1/user/profile', editForm.value);
    ElMessage.success('修改成功');
    showEdit.value = false;
    await fetchProfile();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '修改失败');
  }
};

const onAvatarSuccess = async (response: any) => {
  try {
    ElMessage.success('头像上传成功');
    showAvatarDialog.value = false;
    await fetchProfile();
  } catch (error: any) {
    ElMessage.error('头像更新失败');
    console.error('Avatar update error:', error);
  }
};

const onAvatarError = (error: any) => {
  ElMessage.error('头像上传失败');
  console.error('Avatar upload error:', error);
};

const beforeAvatarUpload = (file: File) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
  if (!isJpgOrPng) {
    ElMessage.error('只能上传jpg/png文件');
    return false;
  }
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    ElMessage.error('图片不能超过2MB');
    return false;
  }
  return true;
};

async function submitChangePwd() {
  // @ts-ignore
  await pwdFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return;
    try {
      const res = await request.post('/api/v1/user/change-password', {
        old_password: pwdForm.value.old_password,
        new_password: pwdForm.value.new_password
      });
      ElMessage.success('密码修改成功，请重新登录');
      showChangePwd.value = false;
      // 清空表单
      pwdForm.value.old_password = '';
      pwdForm.value.new_password = '';
      pwdForm.value.confirm_password = '';
      // 退出登录
      localStorage.removeItem('token');
      setTimeout(() => {
        router.push('/login');
      }, 1000);
    } catch (err: any) {
      ElMessage.error(err?.response?.data?.error || '修改失败');
    }
  });
}
</script>

<style scoped>
body {
  background: linear-gradient(135deg, #181c2f 0%, #3a2b5c 100%);
}
.user-center {
  max-width: 1100px;
  margin: 40px auto;
  padding: 32px 0 64px 0;
  background: rgba(30, 32, 60, 0.5);
  border-radius: 32px;
  box-shadow: 0 8px 40px 0 rgba(80, 100, 255, 0.12);
  backdrop-filter: blur(24px);
  position: relative;
}
.profile-card {
  background: rgba(40, 40, 60, 0.7);
  border-radius: 24px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.18);
  border: 1.5px solid rgba(255,255,255,0.08);
  margin-bottom: 32px;
  padding: 32px 40px;
  color: #e3e8ff;
}
.profile-header {
  display: flex;
  align-items: center;
  gap: 40px;
}
.profile-info {
  flex: 1;
}
.profile-meta {
  color: #b3b8e0;
  font-size: 15px;
  margin-bottom: 10px;
  display: flex;
  gap: 24px;
}
.profile-bio {
  color: #fff;
  font-size: 1.1rem;
  margin-bottom: 8px;
  opacity: 0.85;
}
.profile-info h2 {
  font-size: 2.2rem;
  font-weight: 800;
  background: linear-gradient(90deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.stat-bar {
  display: flex;
  gap: 24px;
  margin-bottom: 32px;
  justify-content: center;
}
.stat-bar > div {
  background: linear-gradient(135deg, #667eea99, #764ba299);
  color: #fff;
  border-radius: 16px;
  padding: 18px 32px;
  font-size: 1.1rem;
  font-weight: 600;
  box-shadow: 0 2px 12px #667eea22;
  transition: transform 0.2s, box-shadow 0.2s;
}
.stat-bar > div:hover {
  transform: translateY(-4px) scale(1.06);
  box-shadow: 0 8px 32px #667eea33;
}
.user-tabs-card {
  background: rgba(60, 70, 120, 0.85);
  border-radius: 20px;
  box-shadow: 0 4px 24px 0 rgba(91,127,255,0.18);
  backdrop-filter: blur(18px) saturate(200%);
  -webkit-backdrop-filter: blur(18px) saturate(200%);
  border: 1px solid rgba(255,255,255,0.10);
  margin-bottom: 32px;
  padding: 8px 0 0 0;
}
.el-tabs {
  background: transparent;
  border-radius: 20px;
  --el-tabs-header-height: 60px;
}
.el-tabs__nav {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 32px;
  width: 100%;
  text-align: center;
}
.el-tabs__item {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  background: linear-gradient(135deg, #7faaff44, #b18fff44);
  border-radius: 24px !important;
  color: #e3e8ff;
  font-weight: 800;
  font-size: 1.18rem;
  margin-right: 0;
  padding: 0 44px;
  height: 56px;
  line-height: 56px;
  transition: background 0.3s, color 0.3s, box-shadow 0.3s, transform 0.2s;
  box-shadow: 0 2px 16px #8f5fff33;
  position: relative;
  z-index: 1;
  margin-bottom: 12px;
  border: none !important;
  letter-spacing: 1px;
  outline: none;
  filter: drop-shadow(0 0 12px #8f5fff33);
}
.el-tabs__item.is-active {
  background: linear-gradient(90deg, #7faaff 0%, #b18fff 100%);
  color: #fff !important;
  box-shadow: 0 8px 32px #8f5fff55, 0 0 16px #b18fffcc;
  font-weight: 900;
  z-index: 2;
  transform: scale(1.12);
  filter: drop-shadow(0 0 20px #b18fffcc);
  text-shadow: 0 2px 12px #b18fff88;
}
.el-tabs__item:not(.is-active):hover {
  background: linear-gradient(135deg, #7faaff66, #b18fff66);
  color: #7faaff;
  transform: scale(1.06);
  box-shadow: 0 4px 20px #8f5fff44;
  filter: drop-shadow(0 0 12px #8f5fff44);
}
.el-tabs__active-bar {
  display: none;
}
.el-tabs__content {
  background: none;
  padding: 0 0 24px 0;
}
.el-tab-pane > .el-table {
  margin-top: 0;
  box-shadow: none;
}
.el-tab-pane {
  background: rgba(30, 32, 60, 0.85);
  border-radius: 24px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.12);
  padding: 24px 16px 8px 16px;
  margin: 0 0 24px 0;
  border: 1px solid rgba(255,255,255,0.05);
}
.el-table {
  background: rgba(30, 32, 60, 0.92) !important;
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 4px 24px 0 rgba(91,127,255,0.10);
  margin-top: 0;
  backdrop-filter: blur(16px) saturate(180%);
  -webkit-backdrop-filter: blur(16px) saturate(180%);
  border: 1px solid rgba(255,255,255,0.06);
}
.el-table__header-wrapper th {
  background: rgba(60, 70, 120, 0.92) !important;
  color: #b3cfff !important;
  font-weight: 700;
  font-size: 1rem;
  border-bottom: 1px solid rgba(255,255,255,0.10);
}
.el-table__cell {
  font-size: 1rem;
  color: #e3e8ff;
  background: rgba(30, 32, 60, 0.92) !important;
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.el-table__row {
  background: transparent !important;
  transition: background 0.2s;
}
.el-table__row:hover {
  background: rgba(91,127,255,0.10) !important;
}
.el-dialog {
  background: rgba(40, 40, 60, 0.95) !important;
  border-radius: 20px !important;
  box-shadow: 0 8px 40px #667eea33 !important;
  border: 1.5px solid #667eea55 !important;
  color: #fff;
}
.el-dialog__header {
  color: #667eea;
  font-weight: 700;
}
.el-form-item__label {
  color: #b3b8e0;
}
.avatar-uploader {
  text-align: center;
}
.upload-tip {
  margin-top: 8px;
  color: #666;
  font-size: 14px;
}
.avatar-container {
  position: relative;
  cursor: pointer;
  border-radius: 50%;
  box-shadow: 0 0 0 4px rgba(102,126,234,0.15), 0 8px 32px 0 rgba(80, 100, 255, 0.10);
  transition: box-shadow 0.3s, transform 0.3s;
}
.avatar-container:hover {
  box-shadow: 0 0 0 8px #667eea55, 0 12px 40px 0 #667eea33;
  transform: scale(1.06);
}
.clickable-avatar {
  border: 3px solid #667eea;
  box-shadow: 0 2px 12px #667eea33;
}
.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  color: white;
}
.avatar-container:hover .avatar-overlay {
  opacity: 1;
}
.avatar-icon {
  font-size: 20px;
  margin-bottom: 4px;
}
.avatar-text {
  font-size: 12px;
  text-align: center;
  line-height: 1.2;
}
.avatar-upload-container {
  text-align: center;
}
.preview-avatar {
  display: block;
  margin: 0 auto 16px auto;
  border: 3px solid #409eff;
  box-shadow: 0 2px 8px #0001;
}
.upload-content-only {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100px;
}
.upload-icon {
  font-size: 40px;
  margin-bottom: 8px;
}
.upload-text {
  text-align: center;
}
.upload-hint {
  display: block;
  margin-top: 8px;
  color: #888;
  font-size: 12px;
}
.avatar-tip {
  margin-top: 8px;
  color: #888;
  font-size: 14px;
  text-align: center;
}
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
.back-btn {
  position: fixed;
  top: 32px;
  left: 32px;
  z-index: 2000;
  font-size: 18px;
  color: #fff;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border: none;
  border-radius: 24px;
  box-shadow: 0 4px 24px #667eea44;
  padding: 12px 32px 12px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: box-shadow 0.2s, background 0.2s, transform 0.2s;
  font-weight: 800;
  letter-spacing: 1px;
}
.back-btn:hover {
  background: linear-gradient(90deg, #764ba2, #667eea);
  box-shadow: 0 8px 32px #667eea77, 0 0 16px #764ba2cc;
  color: #fff;
  transform: translateY(-2px) scale(1.06);
}
/* 编辑资料按钮紫色-青色渐变美化 */
.profile-header .el-button[type="primary"] {
  background: linear-gradient(90deg, #a18fff 0%, #43e9fe 100%) !important;
  color: #fff !important;
  border: none !important;
  border-radius: 20px !important;
  box-shadow: 0 2px 16px #43e9fe33;
  font-weight: 800;
  font-size: 1.12rem;
  padding: 13px 36px;
  letter-spacing: 1.2px;
  transition: box-shadow 0.18s, background 0.18s, transform 0.15s, filter 0.15s;
  position: relative;
  overflow: hidden;
  filter: drop-shadow(0 0 10px #43e9fe44);
}
.profile-header .el-button[type="primary"]:hover {
  background: linear-gradient(90deg, #43e9fe 0%, #a18fff 100%) !important;
  box-shadow: 0 6px 24px #43e9fe77, 0 0 12px #a18fff99;
  color: #fff;
  transform: translateY(-2px) scale(1.06);
  filter: brightness(1.08) drop-shadow(0 0 16px #43e9fe99);
}
@media (max-width: 900px) {
  .user-center {
    max-width: 100vw;
    padding: 0 4vw 32px 4vw;
  }
  .profile-header {
    flex-direction: column;
    gap: 20px;
    align-items: flex-start;
  }
  .profile-card {
    padding: 20px 10px;
  }
  .stat-bar {
    flex-wrap: wrap;
    gap: 12px;
  }
  .stat-bar > div {
    padding: 12px 16px;
    font-size: 1rem;
  }
}
/* 强制el-table所有区域为深色毛玻璃，无白色 */
.el-table,
.el-table__body,
.el-table__body-wrapper,
.el-table__row,
.el-table__cell {
  background: rgba(30, 32, 60, 0.92) !important;
}
.el-table th,
.el-table__header-wrapper,
.el-table__header {
  background: rgba(60, 70, 120, 0.92) !important;
}
.el-table tr {
  background: transparent !important;
}
.el-table__row:hover > td {
  background: rgba(91,127,255,0.10) !important;
}
/* 强制Tab栏内容居中，优先级最高 */
:deep(.user-tabs-card .el-tabs__header),
:deep(.user-tabs-card .el-tabs__nav-wrap),
:deep(.user-tabs-card .el-tabs__nav) {
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  width: 100% !important;
  text-align: center !important;
  margin: 0 auto !important;
  float: none !important;
}
:deep(.user-tabs-card .el-tabs__item) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  text-align: center !important;
}
</style> 
  