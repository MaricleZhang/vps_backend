# VPS 后台管理系统

基于 Go 语言的 VPS 后台管理系统，提供用户管理、流量控制、订阅管理、节点管理等功能。

## 功能特性

- ✅ 用户认证与授权 (JWT)
- ✅ 用户订阅管理
- ✅ 流量监控与配额控制
- ✅ VPS 节点管理
- ✅ 账户余额与充值
- ✅ 密码重置
- ✅ 订单管理
- ✅ RESTful API

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT
- **配置**: Viper
- **密码加密**: bcrypt

## 项目结构

```
vps_backend/
├── cmd/server/          # 应用入口
├── internal/
│   ├── api/
│   │   ├── handler/     # HTTP 处理器
│   │   └── router/      # 路由配置
│   ├── service/         # 业务逻辑层
│   ├── model/           # 数据模型
│   ├── middleware/      # 中间件
│   └── util/            # 工具函数
├── pkg/
│   ├── db/              # 数据库连接
│   └── cache/           # Redis 缓存
├── config/              # 配置文件
├── migrations/          # 数据库迁移
├── docker-compose.yml   # Docker 配置
└── Makefile             # 构建脚本
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 启动数据库服务

使用 Docker Compose 启动 PostgreSQL 和 Redis:

```bash
make docker-up
```

### 3. 配置环境

复制配置文件并修改:

```bash
cp .env.example .env
# 编辑 .env 文件，设置数据库连接等配置
```

### 4. 运行应用

```bash
make run
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 认证接口

#### POST /api/auth/login
登录

**请求体:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应:**
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "user"
    }
  }
}
```

#### POST /api/auth/register
注册

#### POST /api/auth/send-reset-code
发送重置密码验证码

#### POST /api/auth/reset-password
重置密码

### 用户接口 (需要认证)

所有以下接口需要在请求头中包含:
```
Authorization: Bearer {token}
```

#### GET /api/user/info
获取用户信息

#### PUT /api/user/info
更新用户信息

#### POST /api/user/change-password
修改密码

### 账户接口

#### GET /api/account/balance
获取账户余额

#### GET /api/account/traffic
获取流量使用情况

#### GET /api/account/stats
获取账户统计信息

#### POST /api/account/recharge
充值

### 订阅接口

#### GET /api/subscriptions
获取用户订阅列表

#### GET /api/subscriptions/plans
获取可用套餐

#### POST /api/subscriptions/purchase
购买订阅

#### POST /api/subscriptions/renew
续费订阅

#### DELETE /api/subscriptions/:id
取消订阅

### 节点接口

#### GET /api/nodes
获取节点列表
- 查询参数: `location`, `protocol`

#### GET /api/nodes/:id
获取节点详情

#### POST /api/nodes/:id/test
测试节点延迟

## 开发指南

### 添加新的数据模型

1. 在 `internal/model/` 中创建模型文件
2. 在 `pkg/db/postgres.go` 的 `AutoMigrate()` 中添加模型

### 添加新的 API 接口

1. 在 `internal/service/` 中实现业务逻辑
2. 在 `internal/api/handler/` 中创建处理器
3. 在 `internal/api/router/router.go` 中注册路由

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make format
```

## 部署

### Docker 部署

```bash
docker build -t vps-backend .
docker run -p 8080:8080 vps-backend
```

### 生产环境配置

1. 修改 `config/config.yaml` 中的配置
2. 设置环境变量或使用配置文件
3. 确保数据库和 Redis 可访问
4. 使用反向代理 (如 Nginx) 进行 SSL 终止

## 数据库架构

系统使用以下核心表:

- `users` - 用户表
- `subscriptions` - 用户订阅
- `subscription_plans` - 订阅套餐
- `nodes` - VPS 节点
- `user_node_access` - 用户节点访问权限
- `traffic_logs` - 流量日志
- `orders` - 订单
- `announcements` - 公告
- `password_resets` - 密码重置

详细的表结构请参考实现计划文档。

## 常见问题

### 数据库连接失败

确保 PostgreSQL 服务正在运行，并检查配置文件中的连接信息。

### Redis 连接失败

Redis 是可选的。如果不使用 Redis，系统仍可正常运行，但某些功能（如缓存）将不可用。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request!
