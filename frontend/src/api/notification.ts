import request from './request';

export function getNotifications() {
  return request.get('/api/v1/user/notifications');
}

export function getNotification(id: number) {
  return request.get(`/api/v1/user/notifications/${id}`);
}

export function markNotificationRead(id: number) {
  return request.post(`/api/v1/user/notifications/${id}/read`);
}

export function markAllNotificationsRead() {
  return request.post('/api/v1/user/notifications/read-all');
}

export function deleteNotification(id: number) {
  return request.delete(`/api/v1/user/notifications/${id}`);
} 