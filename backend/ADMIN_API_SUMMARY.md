# 后端管理员API接口总结

## 📋 已完成的修改

### 1. 新增接口

#### 文章管理
- `GET /api/v1/admin/articles` - 获取所有文章（支持分页、搜索、过滤）
- `POST /api/v1/admin/articles/:id/review` - 审核文章（支持驳回原因）

#### 用户管理
- `GET /api/v1/admin/users` - 获取用户列表（支持分页、搜索、过滤）
- `PUT /api/v1/admin/users/:id` - 更新用户信息
- `POST /api/v1/admin/users/:id/ban` - 封禁用户
- `POST /api/v1/admin/users/:id/unban` - 解封用户
- `DELETE /api/v1/admin/users/:id` - 删除用户

#### 分类管理
- `GET /api/v1/admin/categories` - 获取分类列表
- `POST /api/v1/admin/categories` - 创建分类
- `PUT /api/v1/admin/categories/:id` - 更新分类
- `DELETE /api/v1/admin/categories/:id` - 删除分类

#### 标签管理
- `GET /api/v1/admin/tags` - 获取标签列表
- `POST /api/v1/admin/tags` - 创建标签
- `PUT /api/v1/admin/tags/:id` - 更新标签
- `DELETE /api/v1/admin/tags/:id` - 删除标签

#### 评论管理
- `GET /api/v1/admin/comments` - 获取评论列表（支持分页、状态过滤）
- `POST /api/v1/admin/comments/:id/approve` - 通过评论
- `POST /api/v1/admin/comments/:id/reject` - 驳回评论
- `PUT /api/v1/admin/comments/:id` - 更新评论
- `DELETE /api/v1/admin/comments/:id` - 删除评论

#### 横幅管理
- `GET /api/v1/admin/banners` - 获取横幅列表
- `POST /api/v1/admin/banners` - 创建横幅
- `PUT /api/v1/admin/banners/:id` - 更新横幅
- `DELETE /api/v1/admin/banners/:id` - 删除横幅

#### 统计和活动
- `GET /api/v1/admin/stats` - 获取统计数据
- `GET /api/v1/admin/activities` - 获取最近活动

#### 系统设置
- `GET /api/v1/admin/settings` - 获取系统设置
- `PUT /api/v1/admin/settings` - 更新系统设置
- `POST /api/v1/admin/settings/reset` - 重置设置
- `POST /api/v1/admin/settings/test-email` - 测试邮件

### 2. 新增模型

#### Banner模型 (`backend/internal/model/banner.go`)
```go
type Banner struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"type:varchar(255)" json:"title"`
    Image     string    `gorm:"type:varchar(500)" json:"image"`
    Link      string    `gorm:"type:varchar(500)" json:"link"`
    SortOrder int       `json:"sort_order"`
    Status    string    `gorm:"type:varchar(32)" json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 3. 数据库迁移

在 `backend/internal/database/database.go` 中添加了：
- `&model.Banner{}` - Banner表
- `&model.ArticleTag{}` - 文章标签关联表

### 4. 权限控制

所有管理员接口都使用了以下中间件：
- `middleware.JWTAuth()` - JWT认证
- `middleware.AdminOnly()` - 管理员权限检查

## 🔧 接口特性

### 分页支持
大部分列表接口都支持分页参数：
- `page` - 页码（默认1）
- `page_size` - 每页数量（默认20，最大100）

### 搜索和过滤
- 用户列表：支持按用户名、邮箱、昵称搜索
- 文章列表：支持按标题、内容、摘要搜索，按状态和分类过滤
- 评论列表：支持按状态过滤

### 响应格式
统一的分页响应格式：
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

## 🚀 部署说明

1. 确保数据库连接正常
2. 重启后端服务，数据库会自动迁移新表
3. 使用默认管理员账号登录：`admin/admin123`
4. 访问前端管理界面：`/admin`

## 📝 注意事项

1. **驳回原因**：文章审核驳回时需要提供原因，后端会验证
2. **权限控制**：所有管理员接口都需要admin角色
3. **数据完整性**：删除操作会级联删除相关数据
4. **文件上传**：横幅图片上传使用现有的 `/api/v1/upload` 接口

## 🔮 未来扩展

1. **邮件通知**：文章驳回时自动发送邮件通知作者
2. **操作日志**：记录管理员的所有操作
3. **数据导出**：支持导出用户、文章等数据
4. **批量操作**：支持批量审核、删除等操作
5. **系统监控**：添加系统性能监控接口 