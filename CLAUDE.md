# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

f_blog 是一个全栈博客系统，Go 后端 + Vue 3 前端，MySQL 数据库。支持文章发布（含审核流程）、评论、点赞、关注、通知、管理后台等功能。

## Commands

### Backend (Go)
```bash
cd backend
go mod tidy                    # 安装依赖
go run main.go                 # 启动开发服务器 :8080
go build -o f_blog_backend     # 编译二进制
```

### Frontend (Vue 3 + Vite)
```bash
cd frontend
npm install                    # 安装依赖
npm run dev                    # Vite 开发服务器 :5173
npm run build                  # 生产构建
```

### Database
- MySQL 8.0+，字符集 utf8mb4
- 项目根目录 `f_blog.sql` 是完整的 Navicat 数据库转储
- GORM AutoMigrate 在后端启动时自动执行，无需手动迁移
- 手动迁移脚本在 `backend/migrations/`

## Architecture

### 后端结构 (`backend/`)
- **`main.go`** — 入口：加载 .env → 加载配置 → 初始化应用(含默认 admin 账户创建) → Gin 引擎 + CORS → 路由 → 启动
- **`internal/app/init.go`** — 应用初始化：数据库连接、创建默认 admin (`admin`/`admin123`)、JWT 配置、启动 cron 调度器
- **`internal/config/`** — 环境变量配置（数据库、JWT、邮件）
- **`internal/database/`** — GORM 连接 + AutoMigrate，全局 `database.DB` 单例
- **`internal/model/`** — 10 个 GORM 模型（User, Article, Category, Tag, Comment, Like, Follow, Notification, Banner, Setting）
- **`internal/handler/`** — HTTP 处理器（含业务逻辑）+ 路由定义（`routes.go`）
- **`internal/middleware/`** — JWT 认证 + CORS
- **`internal/service/`** — 仅限用户认证、邮件发送、定时任务

### 关键设计模式
- **路由前缀**：所有 API 在 `/api/v1` 下
- **认证**：JWT Bearer token，`middleware.JWTAuth()` 保护需登录路由，`middleware.AdminOnly()` 保护管理路由
- **文章审核**：普通用户创建文章状态为 `pending`，管理员为 `published`；审核通过/拒绝会发通知
- **多态点赞**：Like 模型通过 `target_type`（`article`/`comment`）同时支持文章和评论点赞
- **邮件验证**：注册后用户状态为 `pending`，需点击邮件链接激活
- **验证码**：登录和注册需要 CAPTCHA（`base64Captcha` 生成）
- **通知系统**：评论、点赞、关注、文章审核时在 handler 中内联创建通知
- **设置存储**：`settings` 表用 key-value 存储系统配置（前缀命名如 `security_minPasswordLength`）

### 前端结构 (`frontend/`)
- **Vue 3 Composition API** + TypeScript + Vite 5 + Element Plus
- **无状态管理库**（无 Vuex/Pinia），用 `localStorage` 存 token，组件内 `ref()` 管理状态
- **无前端路由守卫**，权限控制全在后端
- **`src/api/request.ts`** — Axios 实例，自动附加 JWT，拦截 401/403/404/500 错误
- **`src/api/`** — 各模块 API 封装（auth, article, like, banner, category, notification）
- **`src/views/`** — 17 个页面组件（公开页、用户中心、管理后台）
- **Vite 代理**：开发环境 `/api` 和 `/static` 代理到 `http://localhost:8080`
- **UI 风格**：深色毛玻璃风格，紫色/蓝色渐变（`#667eea`/`#764ba2`）

### 环境配置
- 后端通过 `backend/.env` 配置（数据库连接、JWT 密钥、邮件 SMTP）
- CORS 允许的源：`localhost:3000`、`localhost:5173`、`127.0.0.1:5173`

### 静态文件
- 后端 `./static/` 目录直接对外提供服务
- 子目录：`avatars/`、`banners/`、`uploads/`、`fonts/`
