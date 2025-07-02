import request from './request';

export function getLikeStatus(targetType: string, targetId: number) {
  return request.get('/api/v1/like/status', {
    params: { target_type: targetType, target_id: targetId }
  });
}

export function like(targetType: string, targetId: number) {
  return request.post('/api/v1/like', { target_type: targetType, target_id: targetId });
}

export function unlike(targetType: string, targetId: number) {
  return request.delete('/api/v1/like', { data: { target_type: targetType, target_id: targetId } });
} 

