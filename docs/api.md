# F_Blog API 接口文档

## 概述

本文档描述了 F_Blog 博客系统的所有 API 接口。系统采用 RESTful API 设计，使用 JSON 格式进行数据交换。

### 基础信息
- **Base URL**: `http://localhost:8080`
- **API Version**: `v1`
- **Content-Type**: `application/json`
- **认证方式**: JWT Bearer Token

### 响应格式
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 错误码说明
- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

---

## 1. 用户认证模块

### 1.1 用户注册

**接口地址**: `POST /api/v1/auth/register`

**请求参数**:
```json
{
  "username": "string",     // 用户名，3-20字符，仅字母数字下划线
  "email": "string",        // 邮箱地址
  "password": "string",     // 密码，8-20字符，需包含字母和数字
  "confirm_password": "string"  // 确认密码
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功，请查收验证邮件",
  "data": {
    "user_id": 1,
    "username": "test_user",
    "email": "test@example.com"
  }
}
```

### 1.2 邮箱验证

**接口地址**: `POST /api/v1/auth/verify-email`

**请求参数**:
```json
{
  "email": "string",    // 邮箱地址
  "code": "string"      // 验证码
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "邮箱验证成功",
  "data": null
}
```

### 1.3 用户登录

**接口地址**: `POST /api/v1/auth/login`

**请求参数**:
```json
{
  "username": "string",  // 用户名或邮箱
  "password": "string"   // 密码
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "test_user",
      "email": "test@example.com",
      "role": "user",
      "avatar": "http://example.com/avatar.jpg"
    }
  }
}
```

### 1.4 用户登出

**接口地址**: `POST /api/v1/auth/logout`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

### 1.5 刷新Token

**接口地址**: `POST /api/v1/auth/refresh-token`

**请求参数**:
```json
{
  "refresh_token": "string"  // 刷新令牌
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "token刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

## 2. 用户管理模块

### 2.1 获取个人资料

**接口地址**: `GET /api/v1/user/profile`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "test_user",
    "email": "test@example.com",
    "nickname": "测试用户",
    "bio": "个人简介",
    "avatar": "http://example.com/avatar.jpg",
    "gender": "male",
    "birthday": "1990-01-01",
    "website": "http://example.com",
    "github": "github_username",
    "twitter": "twitter_username",
    "linkedin": "linkedin_username",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2.2 更新个人资料

**接口地址**: `PUT /api/v1/user/profile`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "nickname": "string",      // 昵称
  "bio": "string",          // 个人简介
  "gender": "string",       // 性别：male/female/other
  "birthday": "string",     // 生日：YYYY-MM-DD
  "website": "string",      // 个人网站
  "github": "string",       // GitHub用户名
  "twitter": "string",      // Twitter用户名
  "linkedin": "string"      // LinkedIn用户名
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "个人资料更新成功",
  "data": null
}
```

### 2.3 修改密码

**接口地址**: `PUT /api/v1/user/password`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "old_password": "string",  // 旧密码
  "new_password": "string"   // 新密码
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

### 2.4 上传头像

**接口地址**: `POST /api/v1/user/avatar`

**请求头**: `Authorization: Bearer <token>`

**请求参数**: `multipart/form-data`
- `file`: 图片文件（JPG/PNG/GIF，最大2MB）

**响应示例**:
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatar_url": "http://example.com/avatars/user_1.jpg"
  }
}
```

### 2.5 删除头像

**接口地址**: `DELETE /api/v1/user/avatar`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "头像删除成功",
  "data": null
}
```

---

## 3. 管理员用户管理

### 3.1 获取用户列表

**接口地址**: `GET /api/v1/admin/users`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**查询参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认20）
- `role`: 角色筛选（user/admin）
- `status`: 状态筛选（active/inactive）
- `keyword`: 关键词搜索（用户名/邮箱）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "size": 20,
    "users": [
      {
        "id": 1,
        "username": "test_user",
        "email": "test@example.com",
        "nickname": "测试用户",
        "role": "user",
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 3.2 修改用户角色

**接口地址**: `PUT /api/v1/admin/users/{id}/role`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**:
```json
{
  "role": "string"  // 角色：user/admin
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户角色修改成功",
  "data": null
}
```

### 3.3 启用/禁用用户

**接口地址**: `PUT /api/v1/admin/users/{id}/status`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**:
```json
{
  "status": "string"  // 状态：active/inactive
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "用户状态修改成功",
  "data": null
}
```

---

## 4. 文章管理模块

### 4.1 发布文章

**接口地址**: `POST /api/v1/articles`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "title": "string",           // 文章标题
  "content": "string",         // 文章内容（支持Markdown）
  "summary": "string",         // 文章摘要
  "cover": "string",           // 封面图片URL
  "category_id": 1,            // 分类ID
  "tags": ["tag1", "tag2"],    // 标签数组
  "status": "draft",           // 状态：draft/published/private
  "publish_time": "2024-01-01T00:00:00Z"  // 定时发布时间
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "文章发布成功",
  "data": {
    "id": 1,
    "title": "测试文章",
    "status": "draft"
  }
}
```

### 4.2 编辑文章

**接口地址**: `PUT /api/v1/articles/{id}`

**请求头**: `Authorization: Bearer <token>`

**请求参数**: 同发布文章

**响应示例**:
```json
{
  "code": 200,
  "message": "文章更新成功",
  "data": null
}
```

### 4.3 删除文章

**接口地址**: `DELETE /api/v1/articles/{id}`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "文章删除成功",
  "data": null
}
```

### 4.4 获取文章列表

**接口地址**: `GET /api/v1/articles`

**查询参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认20）
- `sort`: 排序方式（latest/hot/comments）
- `category_id`: 分类ID
- `tag`: 标签名
- `author_id`: 作者ID
- `start_date`: 开始日期
- `end_date`: 结束日期

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "size": 20,
    "articles": [
      {
        "id": 1,
        "title": "测试文章",
        "summary": "文章摘要",
        "cover": "http://example.com/cover.jpg",
        "author": {
          "id": 1,
          "username": "test_user",
          "nickname": "测试用户",
          "avatar": "http://example.com/avatar.jpg"
        },
        "category": {
          "id": 1,
          "name": "技术"
        },
        "tags": ["Go", "Web"],
        "view_count": 100,
        "like_count": 10,
        "comment_count": 5,
        "status": "published",
        "created_at": "2024-01-01T00:00:00Z",
        "published_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 4.5 获取文章详情

**接口地址**: `GET /api/v1/articles/{id}`

**请求头**: `Authorization: Bearer <token>` (可选，未登录仅返回摘要)

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "测试文章",
    "content": "# 文章内容\n\n支持Markdown格式...",
    "summary": "文章摘要",
    "cover": "http://example.com/cover.jpg",
    "author": {
      "id": 1,
      "username": "test_user",
      "nickname": "测试用户",
      "avatar": "http://example.com/avatar.jpg"
    },
    "category": {
      "id": 1,
      "name": "技术"
    },
    "tags": ["Go", "Web"],
    "view_count": 100,
    "like_count": 10,
    "comment_count": 5,
    "status": "published",
    "is_liked": false,
    "is_followed": false,
    "created_at": "2024-01-01T00:00:00Z",
    "published_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4.6 文章搜索

**接口地址**: `GET /api/v1/articles/search`

**查询参数**:
- `q`: 搜索关键词
- `page`: 页码
- `size`: 每页数量
- `category_id`: 分类ID
- `author_id`: 作者ID
- `start_date`: 开始日期
- `end_date`: 结束日期

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "size": 20,
    "articles": [
      {
        "id": 1,
        "title": "测试文章",
        "summary": "文章摘要",
        "highlight": {
          "title": "测试<em>文章</em>",
          "content": "...<em>测试</em>内容..."
        }
      }
    ]
  }
}
```

### 4.7 增加阅读量

**接口地址**: `POST /api/v1/articles/{id}/view`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "view_count": 101
  }
}
```

---

## 5. 分类管理模块

### 5.1 获取分类列表

**接口地址**: `GET /api/v1/categories`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "技术",
      "description": "技术相关文章",
      "parent_id": 0,
      "article_count": 50,
      "children": [
        {
          "id": 2,
          "name": "Go语言",
          "description": "Go语言相关",
          "parent_id": 1,
          "article_count": 20
        }
      ]
    }
  ]
}
```

### 5.2 创建分类（管理员）

**接口地址**: `POST /api/v1/categories`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**:
```json
{
  "name": "string",           // 分类名称
  "description": "string",    // 分类描述
  "parent_id": 0              // 父分类ID，0表示顶级分类
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "分类创建成功",
  "data": {
    "id": 3,
    "name": "新分类"
  }
}
```

### 5.3 编辑分类（管理员）

**接口地址**: `PUT /api/v1/categories/{id}`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**: 同创建分类

**响应示例**:
```json
{
  "code": 200,
  "message": "分类更新成功",
  "data": null
}
```

### 5.4 删除分类（管理员）

**接口地址**: `DELETE /api/v1/categories/{id}`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**响应示例**:
```json
{
  "code": 200,
  "message": "分类删除成功",
  "data": null
}
```

---

## 6. 标签管理模块

### 6.1 获取标签列表

**接口地址**: `GET /api/v1/tags`

**查询参数**:
- `page`: 页码
- `size`: 每页数量
- `keyword`: 关键词搜索

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "size": 20,
    "tags": [
      {
        "id": 1,
        "name": "Go",
        "article_count": 30,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 6.2 创建标签

**接口地址**: `POST /api/v1/tags`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "name": "string"  // 标签名称
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "标签创建成功",
  "data": {
    "id": 2,
    "name": "新标签"
  }
}
```

---

## 7. 管理员文章管理

### 7.1 获取待审核文章列表

**接口地址**: `GET /api/v1/admin/articles`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**查询参数**:
- `page`: 页码
- `size`: 每页数量
- `status`: 状态筛选（pending/approved/rejected）
- `author_id`: 作者ID
- `category_id`: 分类ID

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "size": 20,
    "articles": [
      {
        "id": 1,
        "title": "待审核文章",
        "author": {
          "id": 1,
          "username": "test_user"
        },
        "status": "pending",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 7.2 审核文章

**接口地址**: `PUT /api/v1/admin/articles/{id}/status`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**:
```json
{
  "status": "string",     // 状态：approved/rejected
  "reason": "string"      // 拒绝原因（可选）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "文章审核完成",
  "data": null
}
```

### 7.3 批量操作

**接口地址**: `POST /api/v1/admin/articles/batch-action`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**:
```json
{
  "ids": [1, 2, 3],           // 文章ID数组
  "action": "string",         // 操作：approve/reject/delete
  "reason": "string"          // 拒绝原因（可选）
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "批量操作完成",
  "data": {
    "success_count": 3,
    "failed_count": 0
  }
}
```

---

## 8. 评论系统

### 8.1 获取文章评论

**接口地址**: `GET /api/v1/articles/{id}/comments`

**查询参数**:
- `page`: 页码
- `size`: 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "size": 20,
    "comments": [
      {
        "id": 1,
        "content": "评论内容",
        "user": {
          "id": 1,
          "username": "test_user",
          "nickname": "测试用户",
          "avatar": "http://example.com/avatar.jpg"
        },
        "reply_count": 2,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 8.2 发表评论

**接口地址**: `POST /api/v1/articles/{id}/comments`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "content": "string"  // 评论内容
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "评论发表成功",
  "data": {
    "id": 2,
    "content": "评论内容"
  }
}
```

### 8.3 删除评论

**接口地址**: `DELETE /api/v1/comments/{id}`

**请求头**: `Authorization: Bearer <token>` (作者或管理员权限)

**响应示例**:
```json
{
  "code": 200,
  "message": "评论删除成功",
  "data": null
}
```

### 8.4 回复评论

**接口地址**: `POST /api/v1/comments/{id}/replies`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:
```json
{
  "content": "string"  // 回复内容
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "回复发表成功",
  "data": {
    "id": 3,
    "content": "回复内容"
  }
}
```

### 8.5 获取评论回复

**接口地址**: `GET /api/v1/comments/{id}/replies`

**查询参数**:
- `page`: 页码
- `size`: 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 10,
    "page": 1,
    "size": 20,
    "replies": [
      {
        "id": 3,
        "content": "回复内容",
        "user": {
          "id": 2,
          "username": "reply_user",
          "nickname": "回复用户",
          "avatar": "http://example.com/avatar2.jpg"
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## 9. 轮播图管理

### 9.1 获取轮播图列表

**接口地址**: `GET /api/v1/banners`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "轮播图标题",
      "image": "http://example.com/banner1.jpg",
      "link": "http://example.com/article/1",
      "sort": 1,
      "status": "active"
    }
  ]
}
```

### 9.2 创建轮播图（管理员）

**接口地址**: `POST /api/v1/banners`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**:
```json
{
  "title": "string",      // 标题
  "image": "string",      // 图片URL
  "link": "string",       // 链接地址
  "sort": 1               // 排序
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "轮播图创建成功",
  "data": {
    "id": 2,
    "title": "新轮播图"
  }
}
```

### 9.3 编辑轮播图（管理员）

**接口地址**: `PUT /api/v1/banners/{id}`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**请求参数**: 同创建轮播图

**响应示例**:
```json
{
  "code": 200,
  "message": "轮播图更新成功",
  "data": null
}
```

### 9.4 删除轮播图（管理员）

**接口地址**: `DELETE /api/v1/banners/{id}`

**请求头**: `Authorization: Bearer <token>` (管理员权限)

**响应示例**:
```json
{
  "code": 200,
  "message": "轮播图删除成功",
  "data": null
}
```

---

## 10. 点赞功能

### 10.1 点赞文章

**接口地址**: `POST /api/v1/articles/{id}/like`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "like_count": 11
  }
}
```

### 10.2 取消点赞

**接口地址**: `DELETE /api/v1/articles/{id}/like`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "取消点赞成功",
  "data": {
    "like_count": 10
  }
}
```

### 10.3 获取点赞用户列表

**接口地址**: `GET /api/v1/articles/{id}/likes`

**查询参数**:
- `page`: 页码
- `size`: 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 10,
    "page": 1,
    "size": 20,
    "users": [
      {
        "id": 1,
        "username": "test_user",
        "nickname": "测试用户",
        "avatar": "http://example.com/avatar.jpg",
        "liked_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## 11. 关注功能

### 11.1 关注用户

**接口地址**: `POST /api/v1/users/{id}/follow`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "关注成功",
  "data": null
}
```

### 11.2 取消关注

**接口地址**: `DELETE /api/v1/users/{id}/follow`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "取消关注成功",
  "data": null
}
```

### 11.3 获取粉丝列表

**接口地址**: `GET /api/v1/users/{id}/followers`

**查询参数**:
- `page`: 页码
- `size`: 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "size": 20,
    "followers": [
      {
        "id": 2,
        "username": "follower_user",
        "nickname": "粉丝用户",
        "avatar": "http://example.com/avatar2.jpg",
        "followed_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 11.4 获取关注列表

**接口地址**: `GET /api/v1/users/{id}/following`

**查询参数**:
- `page`: 页码
- `size`: 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 30,
    "page": 1,
    "size": 20,
    "following": [
      {
        "id": 3,
        "username": "following_user",
        "nickname": "关注用户",
        "avatar": "http://example.com/avatar3.jpg",
        "followed_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## 12. 通知系统

### 12.1 获取通知列表

**接口地址**: `GET /api/v1/notifications`

**请求头**: `Authorization: Bearer <token>`

**查询参数**:
- `page`: 页码
- `size`: 每页数量
- `type`: 通知类型（comment/like/follow/system/review）
- `read`: 是否已读（true/false）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "size": 20,
    "notifications": [
      {
        "id": 1,
        "type": "comment",
        "title": "新评论",
        "content": "用户 test_user 评论了你的文章",
        "is_read": false,
        "data": {
          "article_id": 1,
          "comment_id": 1,
          "user_id": 2
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 12.2 标记通知已读

**接口地址**: `PUT /api/v1/notifications/{id}/read`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "标记成功",
  "data": null
}
```

### 12.3 删除通知

**接口地址**: `DELETE /api/v1/notifications/{id}`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

---

## 13. 搜索建议

### 13.1 获取搜索建议

**接口地址**: `GET /api/v1/search/suggestions`

**查询参数**:
- `q`: 搜索关键词

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "suggestions": [
      "Go语言",
      "Golang",
      "Go Web开发"
    ],
    "hot_keywords": [
      "Go",
      "Vue",
      "MySQL",
      "Redis"
    ]
  }
}
```

---

## 14. 文件上传

### 14.1 上传图片

**接口地址**: `POST /api/v1/upload/image`

**请求头**: `Authorization: Bearer <token>`

**请求参数**: `multipart/form-data`
- `file`: 图片文件（JPG/PNG/GIF，最大5MB）

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "http://example.com/uploads/image.jpg",
    "filename": "image.jpg",
    "size": 1024000
  }
}
```

---

## 15. 统计信息

### 15.1 获取首页统计

**接口地址**: `GET /api/v1/stats/home`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_users": 1000,
    "total_articles": 500,
    "total_comments": 2000,
    "today_views": 5000,
    "latest_articles": [
      {
        "id": 1,
        "title": "最新文章",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "hot_articles": [
      {
        "id": 2,
        "title": "热门文章",
        "view_count": 1000
      }
    ]
  }
}
```

---

## 注意事项

1. **认证要求**: 大部分接口需要JWT认证，请在请求头中添加 `Authorization: Bearer <token>`
2. **权限控制**: 管理员接口需要管理员权限
3. **参数验证**: 所有参数都会进行验证，请确保参数格式正确
4. **错误处理**: 请根据返回的错误码进行相应的错误处理
5. **分页**: 支持分页的接口都使用 `page` 和 `size` 参数
6. **时间格式**: 所有时间字段使用 ISO 8601 格式
7. **文件上传**: 文件上传接口使用 `multipart/form-data` 格式
</rewritten_file> 