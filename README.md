# Force Learning - 强制学习系统

一款帮助用户通过定时弹窗方式进行强制性学习的工具，支持跨平台使用（Android/iOS/Web）。

## 核心功能

### 1. 认证模块
- 用户注册（邮箱/手机号）
- 用户登录（JWT Token 认证）
- Token 刷新机制
- 账户状态管理

### 2. 订阅管理
- 多种订阅套餐（月度/季度/年度）
- 直接购买订阅
- 支付宝/微信支付集成
- 订阅状态查询

### 3. 知识库
- 知识文件上传与管理
- 文件分类浏览
- 随机学习内容获取
- 文件下载功能

### 4. 强制学习弹窗
- 可配置的弹窗时间范围
- 自定义学习时长
- 弹窗间隔设置
- 学习进度追踪

### 5. 学习记录
- 本地学习记录存储
- 云端数据同步
- 离线记录批量上传
- 学习时长统计

## 技术架构

### 后端 (Go)
- **框架**: Gin
- **数据库**: PostgreSQL + Redis
- **认证**: JWT
- **支付**: 支付宝/微信支付

### 移动端 (Flutter)
- **状态管理**: Riverpod
- **HTTP 客户端**: Dio
- **本地存储**: Hive + SharedPreferences
- **跨平台**: Android / iOS / Web

## 项目结构

```
workspace/
├── backend/                    # Go 后端服务
│   ├── cmd/server/            # 入口程序
│   ├── internal/
│   │   ├── api/               # HTTP 层
│   │   │   ├── handler/       # 路由处理
│   │   │   ├── middleware/     # 中间件
│   │   │   └── router/        # 路由定义
│   │   ├── service/           # 业务逻辑
│   │   ├── repository/        # 数据访问
│   │   └── model/             # 数据模型
│   └── configs/               # 配置
│
├── force_learning_app/        # Flutter 应用
│   └── lib/
│       ├── core/              # 核心模块
│       │   ├── api/           # API 客户端
│       │   └── storage/       # 本地存储
│       ├── features/          # 功能模块
│       │   ├── auth/          # 认证
│       │   ├── learning/      # 学习
│       │   ├── subscription/  # 订阅
│       │   └── settings/      # 设置
│       └── services/          # 服务层
│
└── docs/
    └── api.md                 # API 文档
```

## 快速部署

### Docker Compose

```bash
# 启动所有服务
docker-compose -f docker-compose.prod.yml up -d

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps
```

### 手动部署

1. **配置环境变量**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=force_learning
export JWT_SECRET=your-secret-key
```

2. **启动后端**
```bash
cd backend
go build -o server ./cmd/server
./server
```

3. **启动 Flutter 应用**
```bash
cd force_learning_app
flutter run
```

## API 文档

详细 API 文档请查看 [docs/api.md](docs/api.md)

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 5432 |
| JWT_SECRET | JWT 密钥 | - |
| ALIPAY_* | 支付宝配置 | - |
| WXPAY_* | 微信支付配置 | - |

## 订阅套餐

| 套餐 | 时长 | 价格 |
|------|------|------|
| 月度订阅 | 30 天 | ¥29.9 |
| 季度订阅 | 90 天 | ¥79.9 |
| 年度订阅 | 365 天 | ¥299.9 |

## 许可证

MIT License
