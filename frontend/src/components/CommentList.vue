<template>
  <div class="comment-section card">
    <div class="comment-header">
      <h3 class="comment-title">
        <span class="comment-icon">💬</span>
        {{ comments.length }} 条评论
      </h3>
    </div>
    
    <div class="comment-input-section">
      <div class="input-wrapper">
        <textarea
          v-model="newComment"
          class="comment-textarea"
          placeholder="说点什么吧…"
          rows="3"
        ></textarea>
        <button 
          class="submit-btn"
          @click="submitComment" 
          :disabled="submitting"
        >
          <span class="btn-text">{{ submitting ? '发表中...' : '发表评论' }}</span>
          <div class="btn-glow"></div>
        </button>
      </div>
    </div>
    
    <div v-if="comments.length === 0" class="no-comment">
      <div class="empty-icon">💭</div>
      <p class="empty-text">暂无评论，快来抢沙发吧！</p>
    </div>
    
    <div v-else class="comments-list">
      <CommentItem
        v-for="comment in rootComments"
        :key="comment.id"
        :comment="comment"
        :all-comments="comments"
        :article-id="articleId"
        @refresh="fetchComments"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { ElMessage } from 'element-plus';
import request from '../api/request';
import CommentItem from './CommentItem.vue';

const props = defineProps<{
  articleId: number | string
}>();

const comments = ref<any[]>([]);
const newComment = ref('');
const submitting = ref(false);

const fetchComments = async () => {
  try {
    const res = await request.get(`/api/v1/articles/${props.articleId}/comments`);
    comments.value = res.data.data;
  } catch (error) {
    console.error('获取评论失败:', error);
  }
};

const rootComments = computed(() =>
  comments.value.filter(c => c.parent_id === 0)
);

const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning('评论内容不能为空');
    return;
  }
  
  if (!localStorage.getItem('token')) {
    ElMessage.warning('请先登录');
    return;
  }
  
  submitting.value = true;
  try {
    await request.post(`/api/v1/articles/${props.articleId}/comments`, {
      content: newComment.value,
      parent_id: 0
    });
    ElMessage.success('评论成功');
    newComment.value = '';
    fetchComments();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '评论失败');
  } finally {
    submitting.value = false;
  }
};

onMounted(fetchComments);

defineExpose({ fetchComments });
</script>

<style scoped>
.comment-section {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  padding: 32px;
  margin-top: 32px;
}

.comment-header {
  margin-bottom: 24px;
}

.comment-title {
  font-size: 24px;
  font-weight: 700;
  color: white;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.comment-icon {
  font-size: 28px;
}

.comment-input-section {
  margin-bottom: 32px;
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.comment-textarea {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 16px;
  color: white;
  font-size: 16px;
  font-family: inherit;
  resize: vertical;
  min-height: 100px;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.comment-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
  background: rgba(255, 255, 255, 0.15);
}

.comment-textarea::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.submit-btn {
  align-self: flex-end;
  position: relative;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  min-width: 120px;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.submit-btn:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
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

.submit-btn:not(:disabled):hover::before {
  left: 100%;
}

.btn-text {
  position: relative;
  z-index: 1;
}

.no-comment {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.6);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-text {
  font-size: 16px;
  margin: 0;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

@media (max-width: 768px) {
  .comment-section {
    padding: 20px;
    margin-top: 16px;
  }
  
  .comment-title {
    font-size: 20px;
  }
  
  .comment-icon {
    font-size: 24px;
  }
  
  .submit-btn {
    align-self: stretch;
    padding: 12px 20px;
  }
}
</style>