<template>
  <div class="admin-dashboard">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.title">
        <div class="stat-icon" :style="{ background: stat.color }">
          <i :class="stat.icon"></i>
        </div>
        <div class="stat-content">
          <h3>{{ stat.value }}</h3>
          <p>{{ stat.title }}</p>
          <span class="stat-change" :class="stat.trend">
            <i :class="stat.trend === 'up' ? 'fas fa-arrow-up' : 'fas fa-arrow-down'"></i>
            {{ stat.change }}
          </span>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-header">
          <h3>文章发布趋势</h3>
          <div class="chart-actions">
            <select v-model="articleChartPeriod" @change="updateArticleChart">
              <option value="7">最近7天</option>
              <option value="30">最近30天</option>
              <option value="90">最近90天</option>
            </select>
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="articleChart" width="400" height="200"></canvas>
        </div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>用户注册趋势</h3>
          <div class="chart-actions">
            <select v-model="userChartPeriod" @change="updateUserChart">
              <option value="7">最近7天</option>
              <option value="30">最近30天</option>
              <option value="90">最近90天</option>
            </select>
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="userChart" width="400" height="200"></canvas>
        </div>
      </div>
    </div>

    <!-- 最近活动 -->
    <div class="recent-activities">
      <div class="activities-header">
        <h3>最近活动</h3>
        <button class="refresh-btn" @click="fetchRecentActivities">
          <i class="fas fa-sync-alt"></i>
          刷新
        </button>
      </div>
      
      <div class="activities-list">
        <div 
          v-for="activity in recentActivities" 
          :key="activity.id" 
          class="activity-item"
        >
          <div class="activity-icon" :style="{ background: activity.color }">
            <i :class="activity.icon"></i>
          </div>
          <div class="activity-content">
            <p class="activity-text">{{ activity.text }}</p>
            <span class="activity-time">{{ formatTime(activity.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 快速操作 -->
    <div class="quick-actions">
      <h3>快速操作</h3>
      <div class="actions-grid">
        <button 
          v-for="action in quickActions" 
          :key="action.title"
          class="action-btn"
          @click="action.handler"
        >
          <i :class="action.icon"></i>
          <span>{{ action.title }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import request from '../api/request';

const router = useRouter();

// 统计数据
const stats = ref([
  {
    title: '总文章数',
    value: '0',
    change: '+12%',
    trend: 'up',
    icon: 'fas fa-newspaper',
    color: '#667eea'
  },
  {
    title: '总用户数',
    value: '0',
    change: '+8%',
    trend: 'up',
    icon: 'fas fa-users',
    color: '#764ba2'
  },
  {
    title: '总评论数',
    value: '0',
    change: '+15%',
    trend: 'up',
    icon: 'fas fa-comments',
    color: '#f093fb'
  },
  {
    title: '今日访问',
    value: '0',
    change: '+5%',
    trend: 'up',
    icon: 'fas fa-eye',
    color: '#4facfe'
  }
]);

// 图表相关
const articleChart = ref<HTMLCanvasElement>();
const userChart = ref<HTMLCanvasElement>();
const articleChartPeriod = ref('7');
const userChartPeriod = ref('7');

// 更新热门文章热度
const updatePopularArticles = async () => {
  try {
    await request.post('/api/v1/admin/articles/popular/update');
    ElMessage.success('热门文章热度更新成功');
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '更新失败');
  }
};

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    text: '新用户 "张三" 注册了账号',
    icon: 'fas fa-user-plus',
    color: '#4facfe',
    created_at: new Date(Date.now() - 5 * 60 * 1000)
  },
  {
    id: 2,
    text: '文章 "Vue3 开发指南" 被发布',
    icon: 'fas fa-newspaper',
    color: '#667eea',
    created_at: new Date(Date.now() - 15 * 60 * 1000)
  },
  {
    id: 3,
    text: '用户 "李四" 发表了评论',
    icon: 'fas fa-comment',
    color: '#f093fb',
    created_at: new Date(Date.now() - 30 * 60 * 1000)
  }
]);

// 快速操作
const quickActions = [
  {
    title: '发布文章',
    icon: 'fas fa-plus',
    handler: () => router.push('/create-article')
  },
  {
    title: '用户管理',
    icon: 'fas fa-users-cog',
    handler: () => router.push('/admin/users')
  },
  {
    title: '文章审核',
    icon: 'fas fa-check-circle',
    handler: () => router.push('/admin/articles')
  },
  {
    title: '更新热门',
    icon: 'fas fa-fire',
    handler: updatePopularArticles
  },
  {
    title: '系统设置',
    icon: 'fas fa-cog',
    handler: () => router.push('/admin/settings')
  }
];

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await request.get('/api/v1/admin/stats', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    
    const data = response.data;
    stats.value[0].value = data.totalArticles || '0';
    stats.value[1].value = data.totalUsers || '0';
    stats.value[2].value = data.totalComments || '0';
    stats.value[3].value = data.todayVisits || '0';
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
};

// 获取最近活动
const fetchRecentActivities = async () => {
  try {
    const response = await request.get('/api/v1/admin/activities', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    });
    recentActivities.value = response.data || recentActivities.value;
  } catch (error) {
    console.error('获取最近活动失败:', error);
  }
};

// 格式化时间
const formatTime = (date: Date) => {
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / (1000 * 60));
  const hours = Math.floor(diff / (1000 * 60 * 60));
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (minutes < 60) {
    return `${minutes} 分钟前`;
  } else if (hours < 24) {
    return `${hours} 小时前`;
  } else {
    return `${days} 天前`;
  }
};

// 更新文章图表
const updateArticleChart = () => {
  // 这里可以调用API获取图表数据
  console.log('更新文章图表，周期:', articleChartPeriod.value);
};

// 更新用户图表
const updateUserChart = () => {
  // 这里可以调用API获取图表数据
  console.log('更新用户图表，周期:', userChartPeriod.value);
};

// 初始化图表
const initCharts = async () => {
  await nextTick();
  
  // 这里可以使用 Chart.js 或其他图表库来绘制图表
  // 示例：简单的Canvas绘制
  if (articleChart.value) {
    const ctx = articleChart.value.getContext('2d');
    if (ctx) {
      // 绘制简单的折线图
      ctx.strokeStyle = '#667eea';
      ctx.lineWidth = 2;
      ctx.beginPath();
      ctx.moveTo(0, 150);
      ctx.lineTo(100, 100);
      ctx.lineTo(200, 120);
      ctx.lineTo(300, 80);
      ctx.lineTo(400, 60);
      ctx.stroke();
    }
  }

  if (userChart.value) {
    const ctx = userChart.value.getContext('2d');
    if (ctx) {
      // 绘制简单的柱状图
      ctx.fillStyle = '#764ba2';
      ctx.fillRect(50, 100, 40, 50);
      ctx.fillRect(120, 80, 40, 70);
      ctx.fillRect(190, 120, 40, 30);
      ctx.fillRect(260, 60, 40, 90);
      ctx.fillRect(330, 140, 40, 10);
    }
  }
};

onMounted(() => {
  fetchStats();
  fetchRecentActivities();
  initCharts();
});
</script>

<style scoped>
.admin-dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
}

.stat-content h3 {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 4px 0;
  color: #333;
}

.stat-content p {
  margin: 0 0 8px 0;
  color: #666;
  font-size: 14px;
}

.stat-change {
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 4px;
}

.stat-change.up {
  color: #67c23a;
}

.stat-change.down {
  color: #f56c6c;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.chart-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.chart-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
}

.chart-actions select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  color: #333;
}

.chart-container {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.recent-activities {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.activities-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.activities-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
}

.refresh-btn {
  background: #667eea;
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.refresh-btn:hover {
  background: #5a6fd8;
}

.activities-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background: #f8f9fa;
  transition: all 0.3s ease;
}

.activity-item:hover {
  background: #e9ecef;
}

.activity-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
}

.activity-content {
  flex: 1;
}

.activity-text {
  margin: 0 0 4px 0;
  color: #333;
  font-size: 14px;
}

.activity-time {
  color: #666;
  font-size: 12px;
}

.quick-actions {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.quick-actions h3 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 1.2rem;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.action-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
  padding: 16px 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  font-weight: 500;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
}

.action-btn i {
  font-size: 18px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .actions-grid {
    grid-template-columns: 1fr;
  }
}
</style> 