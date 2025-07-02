import axios from 'axios';

export function getBanners() {
  return axios.get('/api/v1/banners');
}