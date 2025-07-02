import axios from 'axios';
import { ElMessage } from 'element-plus';

// 创建axios实例
const request = axios.create({
  baseURL: 'http://localhost:8080', // 后端服务器地址
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response) {
      // 服务器返回错误状态码
      const { status, data } = error.response;
      
      switch (status) {
        case 401:
          // 未授权，清除token并跳转到登录页
          localStorage.removeItem('token');
          ElMessage.error('登录已过期，请重新登录');
          // 可以在这里跳转到登录页
          break;
        case 403:
          // 只在后台页面弹窗
          if (window.location.pathname.startsWith('/admin')) {
            ElMessage.error('没有权限访问');
          }
          // 其他页面静默处理
          break;
        case 404:
          ElMessage.error('请求的资源不存在');
          break;
        case 500:
          ElMessage.error('服务器内部错误');
          break;
        default:
          ElMessage.error(data?.error || '请求失败');
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      ElMessage.error('网络连接失败，请检查网络');
    } else {
      // 其他错误
      ElMessage.error('请求配置错误');
    }
    return Promise.reject(error);
  }
);

export default request;

