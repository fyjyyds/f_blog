import request from './request';

export function listArticles(query = '') {
  return request.get(`/api/v1/articles${query}`);
}

export function getArticle(id: number | string) {
  return request.get(`/api/v1/articles/${id}`);
}

export function createArticle(data: any) {
  return request.post('/api/v1/articles', data);
}

export function updateArticle(id: number | string, data: any) {
  return request.put(`/api/v1/articles/${id}`, data);
}

export function deleteArticle(id: number | string) {
  return request.delete(`/api/v1/articles/${id}`);
}

// 获取热门文章
export function getPopularArticles(limit = 10) {
  return request.get('/api/v1/articles/popular', {
    params: { limit }
  });
}