# F_Blog 后端服务

## 环境要求

- Go 1.21+
- MySQL 5.7+

## 安装依赖

```bash
go mod tidy
```

## 环境配置

创建 `.env` 文件：

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=f_blog
DB_CHARSET=utf8mb4
```

## 数据库初始化

1. 创建数据库：
```sql
CREATE DATABASE f_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 运行初始化脚本：
```bash
mysql -u root -p f_blog < scripts/init.sql
```

## 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 测试

### 注册用户
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "12345678"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "testuser",
    "password": "12345678"
  }'
```

## 常见问题

### 1. 登录失败
- 检查数据库连接是否正常
- 确认用户已注册
- 检查密码是否正确

### 2. CORS 错误
- 后端已配置 CORS，支持前端域名访问
- 如果仍有问题，检查前端请求地址是否正确

### 3. 数据库连接失败
- 确认 MySQL 服务正在运行
- 检查 `.env` 文件中的数据库配置
- 确认数据库和表已创建 

## 环境要求

- Go 1.21+
- MySQL 5.7+

## 安装依赖

```bash
go mod tidy
```

## 环境配置

创建 `.env` 文件：

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=f_blog
DB_CHARSET=utf8mb4
```

## 数据库初始化

1. 创建数据库：
```sql
CREATE DATABASE f_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 运行初始化脚本：
```bash
mysql -u root -p f_blog < scripts/init.sql
```

## 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 测试

### 注册用户
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "12345678"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "testuser",
    "password": "12345678"
  }'
```

## 常见问题

### 1. 登录失败
- 检查数据库连接是否正常
- 确认用户已注册
- 检查密码是否正确

### 2. CORS 错误
- 后端已配置 CORS，支持前端域名访问
- 如果仍有问题，检查前端请求地址是否正确

### 3. 数据库连接失败
- 确认 MySQL 服务正在运行
- 检查 `.env` 文件中的数据库配置
- 确认数据库和表已创建 