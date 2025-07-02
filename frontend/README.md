# F_Blog 前端项目

## 技术栈
- Vue 3
- TypeScript
- Element Plus
- Vite
- Axios

## 初始化步骤

1. 安装依赖

```bash
cd frontend
npm install
```

2. 启动开发服务器

```bash
npm run dev
```

3. 目录结构

```
frontend/
  src/
    api/
      auth.ts
    views/
      Register.vue
      Login.vue
    router/
      index.ts
    App.vue
    main.ts
  vite.config.ts
  package.json
  tsconfig.json
```

4. 主要页面
- 用户注册：/register
- 用户登录：/login

5. 代理配置
如后端运行在 8080 端口，Vite 需配置代理转发 API 请求。

```js
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
});
``` 