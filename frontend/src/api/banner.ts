import request from './request';

export function getBanners() {
  return request.get('/api/v1/banners');
}
