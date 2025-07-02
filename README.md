# F_Blog - 全栈博客系统

一个基于Go语言开发的全栈博客系统，支持用户管理、文章管理、评论系统、点赞关注等功能。

## 项目结构

```
f_blog/
├── backend/                 # 后端服务
│   ├── main.go             # 应用入口文件
│   ├── go.mod              # Go模块文件
│   ├── README.md           # 后端说明文档
│   ├── internal/           # 内部包
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库相关
│   │   └── model/          # 数据模型
│   └── scripts/            # 脚本目录
│       └── init.sql       # 数据库初始化脚本
├── docs/                   # 文档目录
│   ├── requirements.md     # 需求文档
│   ├── api.md             # API文档
│   └── database.md        # 数据库设计文档
└── README.md              # 项目说明
├── frontend/              #前端说明
    ├── app.vue            #懒得写了
```
如需完成邮件验证，在.env邮件中配置自己的smtp
## 功能特性

### 用户管理
- 用户注册、登录、登出
- JWT身份验证
- 用户资料管理
- 头像上传
- 密码修改

### 文章管理
- 文章发布、编辑、删除
- 富文本编辑器支持
- 文章分类和标签
- 文章状态管理（草稿、发布、私有等）
- 文章搜索和筛选

### 评论系统
- 文章评论
- 评论回复
- 评论审核
- 评论点赞

### 社交功能
- 文章点赞
- 用户关注
- 通知系统
- 个人动态

### 管理功能
- 用户管理
- 文章审核
- 分类标签管理
- 轮播图管理
- 系统统计

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **缓存**: Redis //暂未使用
- **认证**: JWT

### 前端
- **框架**: Vue.js 3 + TypeScript
- **UI组件库**: Element Plus
- **构建工具**: Vite

## 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 2. 后端启动

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 配置数据库
mysql -u root -p f_blog < scripts/init.sql

# 设置环境变量
export DB_PASSWORD=your_password
export JWT_SECRET=your-secret-key

# 运行后端服务
go run main.go
```

### 3. 环境变量配置

创建 `.env` 文件或设置环境变量：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=f_blog
DB_CHARSET=utf8mb4

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_TIME=24h

# 服务器配置
SERVER_PORT=8080
```

## 数据库设计

### 主要表结构

1. **users** - 用户表
2. **articles** - 文章表
3. **categories** - 分类表
4. **tags** - 标签表
5. **comments** - 评论表
6. **likes** - 点赞表
7. **follows** - 关注表
8. **notifications** - 通知表
9. **banners** - 轮播图表
10. **audit_logs** - 审核日志表

详细的数据库设计请参考 [docs/database.md](docs/database.md)

## API文档

完整的API文档请参考 [docs/api.md](docs/api.md)

## 开发指南

### 后端开发

后端代码位于 `backend/` 目录，使用Go语言开发。

- **配置管理**: `backend/internal/config/`
- **数据库操作**: `backend/internal/database/`
- **数据模型**: `backend/internal/model/`

### 添加新的功能

1. 在 `backend/internal/model/` 中添加新的模型
2. 在 `backend/main.go` 中添加新的API路由
3. 实现对应的业务逻辑

### 数据库操作

使用GORM进行数据库操作，数据库连接通过 `backend/internal/database/database.go` 管理。

## 部署

### 后端部署

```bash
cd backend
go build -o f_blog_backend main.go
./f_blog_backend
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或联系项目维护者。 
