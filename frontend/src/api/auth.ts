import request from './request';

export function register(data: { username: string; email: string; password: string; captcha_id: string; captcha: string }) {
  return request.post('/api/v1/auth/register', data);
}

export function login(data: { username_or_email: string; password: string; captcha_id: string; captcha: string }) {
  return request.post('/api/v1/auth/login', data);
} 
