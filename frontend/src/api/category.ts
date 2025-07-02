import request from './request';

export function getCategories() {
  return request.get('/api/v1/categories');
}

export function getTags() {
  return request.get('/api/v1/tags');
} 

