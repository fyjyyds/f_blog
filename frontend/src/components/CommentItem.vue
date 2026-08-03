<template>
    <div class="comment-item">
      <div class="comment-header">
        <div class="author-info">
          <div class="author-avatar">
            <span class="avatar-text">{{ comment.display_name.charAt(0) }}</span>
          </div>
          <div class="author-details">
            <span class="author-name">{{ comment.display_name }}</span>
            <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
          </div>
        </div>
        <div class="comment-actions">
          <button v-if="canDelete" class="action-btn delete-btn" @click="onDelete">
            <span class="action-icon">🗑️</span>
            删除
          </button>
          <button v-if="comment.parent_id === 0" class="action-btn reply-btn" @click="showReply = !showReply">
            <span class="action-icon">💬</span>
            回复
          </button>
        </div>
      </div>
      <div class="comment-content">
        <span v-if="comment.reply_to_user" class="reply-info">回复 @{{ comment.reply_to_user }}</span>
        {{ comment.content }}
      </div>
      <div class="comment-footer">
        <button 
          class="like-btn" 
          :class="{ 'liked': commentLiked }" 
          @click="toggleLike"
          :disabled="!isLoggedIn"
        >
          <span class="like-icon">{{ commentLiked ? '❤️' : '🤍' }}</span>
          <span class="like-text">{{ commentLiked ? '已点赞' : '点赞' }}</span>
          <span class="like-count" v-if="commentLikeCount > 0">{{ commentLikeCount }}</span>
        </button>
      </div>
      <div v-if="showReply" class="reply-box">
        <textarea
          v-model="replyContent"
          class="reply-textarea"
          placeholder="回复内容"
          rows="2"
        ></textarea>
        <div class="reply-actions">
          <button class="reply-submit-btn" @click="submitReply" :disabled="replying">
            {{ replying ? '回复中...' : '回复' }}
          </button>
          <button class="reply-cancel-btn" @click="showReply = false">取消</button>
        </div>
      </div>
      <div class="children" v-if="comment.parent_id === 0 && children.length">
        <template v-for="(child, idx) in visibleChildren" :key="child.id">
          <div class="child-comment">
            <div class="comment-header">
              <div class="author-info">
                <div class="author-avatar child-avatar">
                  <span class="avatar-text">{{ child.display_name.charAt(0) }}</span>
                </div>
                <div class="author-details">
                  <span class="author-name">{{ child.display_name }}</span>
                  <span class="comment-date">{{ formatDate(child.created_at) }}</span>
                </div>
              </div>
              <div class="comment-actions">
                <button v-if="child.user_id === userId || role === 'admin'" class="action-btn delete-btn" @click="onDeleteChild(child.id)">
                  <span class="action-icon">🗑️</span>
                  删除
                </button>
                <button class="action-btn reply-btn" @click="replyToChildId = child.id">
                  <span class="action-icon">💬</span>
                  回复
                </button>
              </div>
            </div>
            <div class="comment-content">
              <span v-if="child.reply_to_user" class="reply-info">回复 @{{ child.reply_to_user }}</span>
              {{ child.content }}
            </div>
            <div class="comment-footer">
              <button 
                class="like-btn" 
                :class="{ 'liked': childLiked[child.id] }" 
                @click="toggleChildLike(child.id)"
                :disabled="!isLoggedIn"
              >
                <span class="like-icon">{{ childLiked[child.id] ? '❤️' : '🤍' }}</span>
                <span class="like-text">{{ childLiked[child.id] ? '已点赞' : '点赞' }}</span>
                <span class="like-count" v-if="childLikeCount[child.id] > 0">{{ childLikeCount[child.id] }}</span>
              </button>
            </div>
            <div v-if="replyToChildId === child.id" class="reply-box">
              <div class="reply-to-tip">回复 @{{ child.display_name }}</div>
              <textarea
                v-model="childReplyContent"
                class="reply-textarea"
                placeholder="回复内容"
                rows="2"
              ></textarea>
              <div class="reply-actions">
                <button class="reply-submit-btn" @click="submitChildReply(child)" :disabled="childReplying">
                  {{ childReplying ? '回复中...' : '回复' }}
                </button>
                <button class="reply-cancel-btn" @click="replyToChildId = null">取消</button>
              </div>
            </div>
          </div>
          <div v-if="!expanded && children.length > 1 && idx === 0" class="expand-btn">
            <button class="expand-btn-text" @click="expanded = true">
              展开更多回复 ({{ children.length - 1 }})
            </button>
          </div>
        </template>
        <div v-if="expanded && children.length > 1" class="expand-btn">
          <button class="expand-btn-text" @click="expanded = false">收起</button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue';
  import { ElMessage } from 'element-plus';
  import request from '../api/request';
  import { formatDate } from '../utils/date';
  import { getLikeStatus, like, unlike } from '../api/like';
  
  const props = defineProps<{
    comment: any,
    allComments: any[],
    articleId: number | string
  }>();
  const emit = defineEmits(['refresh']);
  
  const replying = ref(false);
  const replyContent = ref('');
  const showReply = ref(false);
  
  const userId = Number(localStorage.getItem('user_id'));
  const role = localStorage.getItem('role');
  const isLoggedIn = computed(() => !!localStorage.getItem('token'));
  const canDelete = computed(() => props.comment.user_id === userId || role === 'admin');
  
  const expanded = ref(false);
  const children = computed(() =>
    props.allComments.filter(c => c.parent_id === props.comment.id)
  );
  const visibleChildren = computed(() => {
    if (expanded.value || children.value.length === 0) {
      return children.value;
    }
    // 只显示第1条
    return children.value.slice(0, 1);
  });
  
  // 子评论回复相关
  const replyToChildId = ref<number | null>(null);
  const childReplyContent = ref('');
  const childReplying = ref(false);
  
  // 点赞相关
  const commentLiked = ref(false);
  const commentLikeCount = ref(0);
  const childLiked = ref<Record<number, boolean>>({});
  const childLikeCount = ref<Record<number, number>>({});
  
  // 获取评论点赞状态
  const fetchCommentLikeStatus = async () => {
    if (!isLoggedIn.value) return;
    
    try {
      const response = await getLikeStatus('comment', props.comment.id);
      commentLiked.value = response.data.liked;
      commentLikeCount.value = response.data.count;
    } catch (error) {
      console.error('获取评论点赞状态失败:', error);
    }
  };
  
  // 获取子评论点赞状态
  const fetchChildLikeStatus = async (childId: number) => {
    if (!isLoggedIn.value) return;
    
    try {
      const response = await getLikeStatus('comment', childId);
      childLiked.value[childId] = response.data.liked;
      childLikeCount.value[childId] = response.data.count;
    } catch (error) {
      console.error('获取子评论点赞状态失败:', error);
    }
  };
  
  // 切换评论点赞状态
  const toggleLike = async () => {
    if (!isLoggedIn.value) {
      ElMessage.warning('请先登录');
      return;
    }
    
    try {
      if (commentLiked.value) {
        await unlike('comment', props.comment.id);
        commentLikeCount.value--;
      } else {
        await like('comment', props.comment.id);
        commentLikeCount.value++;
      }
      commentLiked.value = !commentLiked.value;
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '操作失败');
    }
  };
  
  // 切换子评论点赞状态
  const toggleChildLike = async (childId: number) => {
    if (!isLoggedIn.value) {
      ElMessage.warning('请先登录');
      return;
    }
    
    try {
      if (childLiked.value[childId]) {
        await unlike('comment', childId);
        childLikeCount.value[childId]--;
      } else {
        await like('comment', childId);
        childLikeCount.value[childId]++;
      }
      childLiked.value[childId] = !childLiked.value[childId];
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '操作失败');
    }
  };
  
  // 初始化点赞状态
  onMounted(() => {
    fetchCommentLikeStatus();
    // 为所有子评论获取点赞状态
    children.value.forEach(child => {
      fetchChildLikeStatus(child.id);
    });
  });
  
  const submitReply = async () => {
    if (!replyContent.value.trim()) {
      ElMessage.warning('回复内容不能为空');
      return;
    }
    
    if (!localStorage.getItem('token')) {
      ElMessage.warning('请先登录');
      return;
    }
    
    replying.value = true;
    try {
      await request.post(`/api/v1/articles/${props.articleId}/comments`, {
        content: replyContent.value,
        parent_id: props.comment.id
      });
      ElMessage.success('回复成功');
      replyContent.value = '';
      showReply.value = false;
      emit('refresh');
    } catch (e: any) {
      ElMessage.error(e.response?.data?.error || '回复失败');
    } finally {
      replying.value = false;
    }
  };
  
  const submitChildReply = async (child: any) => {
    if (!childReplyContent.value.trim()) {
      ElMessage.warning('回复内容不能为空');
      return;
    }
    
    if (!localStorage.getItem('token')) {
      ElMessage.warning('请先登录');
      return;
    }
    
    childReplying.value = true;
    try {
      await request.post(`/api/v1/articles/${props.articleId}/comments`, {
        content: childReplyContent.value,
        parent_id: props.comment.id, // 依然归属于母评论
        reply_to_user: child.display_name // 记录回复给谁
      });
      ElMessage.success('回复成功');
      childReplyContent.value = '';
      replyToChildId.value = null;
      emit('refresh');
    } catch (e: any) {
      ElMessage.error(e.response?.data?.error || '回复失败');
    } finally {
      childReplying.value = false;
    }
  };
  
  const onDelete = async () => {
    try {
      await request.delete(`/api/v1/comments/${props.comment.id}`);
      ElMessage.success('删除成功');
      emit('refresh');
    } catch (e: any) {
      ElMessage.error(e.response?.data?.error || '删除失败');
    }
  };
  
  const onDeleteChild = async (id: number) => {
    try {
      await request.delete(`/api/v1/comments/${id}`);
      ElMessage.success('删除成功');
      emit('refresh');
    } catch (e: any) {
      ElMessage.error(e.response?.data?.error || '删除失败');
    }
  };
  </script>
  
  <style scoped>
  .comment-item {
    background: #fff;
    border: 1px solid #e5e7eb;
    border-radius: 10px;
    padding: 16px;
    margin-bottom: 12px;
    transition: box-shadow 0.15s;
  }
  .comment-item:hover {
    box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  }
  
  .comment-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;
  }
  
  .author-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  
  .author-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: #5b6abf;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    color: #fff;
    font-size: 14px;
  }
  
  .child-avatar {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }
  
  .avatar-text {
    text-transform: uppercase;
  }
  
  .author-details {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  
  .author-name {
    color: #5b6abf;
    font-weight: 600;
    font-size: 14px;
  }
  .comment-date {
    color: #9ca3af;
    font-size: 12px;
  }
  
  .comment-actions {
    display: flex;
    gap: 8px;
  }
  
  .action-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    background: #f3f4f6;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    padding: 5px 10px;
    color: #6b7280;
    font-size: 12px;
    cursor: pointer;
    transition: all 0.15s;
  }
  .action-btn:hover {
    background: #e5e7eb;
  }
  
  .delete-btn:hover {
    background: rgba(255, 107, 107, 0.2);
    border-color: #ff6b6b;
  }
  
  .reply-btn:hover {
    background: rgba(102, 126, 234, 0.2);
    border-color: #667eea;
  }
  
  .action-icon {
    font-size: 12px;
  }
  
  .comment-content {
    color: #1f2937;
    font-size: 14px;
    line-height: 1.6;
    margin-bottom: 10px;
  }
  
  .comment-footer {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
  }
  
  .like-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    padding: 5px 10px;
    color: #6b7280;
    font-size: 12px;
    cursor: pointer;
    transition: all 0.15s;
  }
  .like-btn:hover:not(:disabled) {
    border-color: #f87171;
    color: #f87171;
  }
  .like-btn.liked {
    background: #fef2f2;
    border-color: #f87171;
    color: #ef4444;
  }
  
  .like-btn.liked:hover:not(:disabled) {
    background: rgba(255, 107, 107, 0.3);
  }
  
  .like-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .like-icon {
    font-size: 14px;
  }
  
  .like-text {
    font-weight: 500;
  }
  
  .like-count {
    background: #f3f4f6;
    border-radius: 4px;
    padding: 2px 6px;
    font-size: 11px;
    font-weight: 600;
  }
  
  .reply-info {
    color: #667eea;
    font-size: 13px;
    margin-right: 8px;
    font-weight: 500;
  }
  
  .reply-box {
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 14px;
    margin-top: 10px;
  }
  
  .reply-to-tip {
    color: #667eea;
    font-size: 13px;
    margin-bottom: 8px;
    font-weight: 500;
  }
  
  .reply-textarea {
    width: 100%;
    background: #fff;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    padding: 10px;
    color: #1f2937;
    font-size: 14px;
    font-family: inherit;
    resize: vertical;
    min-height: 60px;
  }
  .reply-textarea:focus {
    outline: none;
    border-color: #5b6abf;
    box-shadow: 0 0 0 3px rgba(91,106,191,0.1);
  }
  .reply-textarea::placeholder {
    color: #9ca3af;
  }
  
  .reply-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
    justify-content: flex-end;
  }
  
  .reply-submit-btn {
    background: #5b6abf;
    border: none;
    border-radius: 6px;
    padding: 7px 14px;
    color: #fff;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
  }
  .reply-submit-btn:hover:not(:disabled) {
    background: #4a59ae;
  }
  .reply-submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .reply-cancel-btn {
    background: #f3f4f6;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    padding: 7px 14px;
    color: #6b7280;
    font-size: 12px;
    cursor: pointer;
  }
  .reply-cancel-btn:hover {
    background: #e5e7eb;
  }
  
  .children {
    margin-left: 20px;
    margin-top: 12px;
    padding-left: 16px;
    border-left: 2px solid #e5e7eb;
  }
  .child-comment {
    background: #f9fafb;
    border: 1px solid #f3f4f6;
    border-radius: 8px;
    padding: 12px;
    margin-bottom: 10px;
  }
  .child-comment:hover {
    background: #f3f4f6;
  }
  
  .expand-btn {
    margin-top: 8px;
  }
  
  .expand-btn-text {
    background: none;
    border: none;
    color: #667eea;
    font-size: 13px;
    cursor: pointer;
    text-decoration: underline;
    transition: color 0.3s ease;
  }
  
  .expand-btn-text:hover {
    color: #764ba2;
  }
  
  @media (max-width: 768px) {
    .comment-item {
      padding: 16px;
    }
    
    .comment-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }
    
    .comment-actions {
      width: 100%;
      justify-content: flex-end;
    }
    
    .children {
      margin-left: 10px;
      padding-left: 10px;
    }
    
    .author-avatar {
      width: 32px;
      height: 32px;
      font-size: 14px;
    }
    
    .child-avatar {
      width: 28px;
      height: 28px;
      font-size: 12px;
    }
  }
  </style>
   