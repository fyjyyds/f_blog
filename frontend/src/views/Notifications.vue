<template>
  <div class="notifications-page">
    <div class="card notifications-card">
      <div class="header">
        <el-button class="back-btn" size="small" @click="goBack" :icon="ArrowLeft">返回</el-button>
        <h2>我的通知</h2>
        <el-button size="small" @click="handleMarkAllRead" :disabled="!unreadCount">全部标为已读</el-button>
      </div>
      <el-empty v-if="notifications.length === 0" description="暂无通知" />
      <el-list v-else>
        <el-list-item
          v-for="item in notifications"
          :key="item.id"
          :class="{ unread: !item.is_read }"
          @click="openDetail(item)"
        >
          <div class="main">
            <span class="type-tag" :class="item.type">{{ typeText(item.type) }}</span>
            <span class="title">{{ item.title }}</span>
            <span class="content">{{ item.content }}</span>
          </div>
          <div class="meta">
            <span class="date">{{ formatDate(item.created_at) }}</span>
            <el-button
              v-if="!item.is_read"
              size="mini"
              type="primary"
              @click.stop="markRead(item)"
            >标为已读</el-button>
            <el-button
              size="mini"
              type="danger"
              @click.stop="deleteOne(item)"
            >删除</el-button>
          </div>
        </el-list-item>
      </el-list>
    </div>
    <el-dialog v-model="showDetail" title="通知详情" width="400px">
      <div v-if="current">
        <div><b>类型：</b>{{ typeText(current.type) }}</div>
        <div><b>标题：</b>{{ current.title }}</div>
        <div><b>内容：</b>{{ current.content }}</div>
        <div><b>时间：</b>{{ formatDate(current.created_at) }}</div>
        <div v-if="current.data"><b>附加数据：</b>{{ current.data }}</div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { getNotifications, markNotificationRead, markAllNotificationsRead, deleteNotification } from '../api/notification';
import { ElMessage, ElMessageBox } from 'element-plus';
import { ArrowLeft } from '@element-plus/icons-vue';
import { formatDate } from '../utils/date';

const router = useRouter();
const notifications = ref<any[]>([]);
const showDetail = ref(false);
const current = ref<any>(null);

const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length);

function typeText(type: string) {
  switch (type) {
    case 'comment': return '评论';
    case 'like': return '点赞';
    case 'follow': return '关注';
    case 'system': return '系统';
    case 'review': return '审核';
    default: return type;
  }
}

async function fetchList() {
  const res = await getNotifications();
  notifications.value = res.data;
}

function openDetail(item: any) {
  current.value = item;
  showDetail.value = true;
  if (!item.is_read) {
    markRead(item);
  }
}

async function markRead(item: any) {
  await markNotificationRead(item.id);
  item.is_read = true;
  fetchList();
}

async function handleMarkAllRead() {
  await markAllNotificationsRead();
  ElMessage.success('全部已标为已读');
  fetchList();
}

async function deleteOne(item: any) {
  await ElMessageBox.confirm('确定删除该通知吗？', '提示', { type: 'warning' });
  await deleteNotification(item.id);
  ElMessage.success('已删除');
  fetchList();
}

function goBack() {
  router.back();
}

onMounted(fetchList);
</script>

<style scoped>
.notifications-page {
  max-width: 600px;
  margin: 40px auto;
}
.notifications-card {
  padding: 32px;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}
.header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  gap: 12px;
}
.unread {
  background: #f0f6ff;
}
.type-tag {
  display: inline-block;
  min-width: 40px;
  text-align: center;
  border-radius: 4px;
  margin-right: 8px;
  font-size: 12px;
  padding: 2px 6px;
  color: #fff;
}
.type-tag.comment { background: #42b983; }
.type-tag.like { background: #f39c12; }
.type-tag.follow { background: #409eff; }
.type-tag.system { background: #909399; }
.type-tag.review { background: #e74c3c; }
.main { display: flex; align-items: center; gap: 8px; }
.meta { float: right; display: flex; align-items: center; gap: 8px; }
.date { color: #999; font-size: 12px; margin-right: 8px; }
.back-btn {
  margin-right: 12px;
}
</style> 