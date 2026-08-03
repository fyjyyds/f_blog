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
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 24px;
  margin-top: 24px;
}

.comment-header {
  margin-bottom: 20px;
}

.comment-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-input-section {
  margin-bottom: 24px;
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-textarea {
  width: 100%;
  background: #fff;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 12px;
  color: #1f2937;
  font-size: 14px;
  font-family: inherit;
  resize: vertical;
  min-height: 80px;
}
.comment-textarea:focus {
  outline: none;
  border-color: #5b6abf;
  box-shadow: 0 0 0 3px rgba(91,106,191,0.1);
}
.comment-textarea::placeholder {
  color: #9ca3af;
}

.submit-btn {
  align-self: flex-end;
  background: #5b6abf;
  border: none;
  border-radius: 6px;
  padding: 9px 20px;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  min-width: 100px;
}
.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.submit-btn:not(:disabled):hover {
  background: #4a59ae;
}

.no-comment {
  text-align: center;
  padding: 40px 20px;
  color: #9ca3af;
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