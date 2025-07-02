import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import Register from '../views/Register.vue';
import Login from '../views/Login.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/register', component: Register },
  { path: '/login', component: Login },
  { path: '/create-article', component: () => import('../views/CreateArticle.vue') },
  { path: '/article/:id', component: () => import('../views/ArticleDetail.vue') },
  { path: '/edit-article/:id', component: () => import('../views/EditArticle.vue') },
  { path: '/user-center', component: () => import('../views/UserCenter.vue') },
  { path: '/verify-email', component: () => import('../views/VerifyEmail.vue') },
  { path: '/notifications', component: () => import('../views/Notifications.vue') },
  
  // 后台管理路由
  {
    path: '/admin',
    component: () => import('../views/AdminLayout.vue'),
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', component: () => import('../views/AdminDashboard.vue') },
      { path: 'articles', component: () => import('../views/AdminArticle.vue') },
      { path: 'users', component: () => import('../views/AdminUsers.vue') },
      { path: 'categories', component: () => import('../views/AdminCategories.vue') },
      { path: 'tags', component: () => import('../views/AdminTags.vue') },
      { path: 'comments', component: () => import('../views/AdminComments.vue') },
      { path: 'banners', component: () => import('../views/AdminBanners.vue') },
      { path: 'settings', component: () => import('../views/AdminSettings.vue') }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
