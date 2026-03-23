# 云服务器部署指南

## 一、服务器配置要求

### 最低配置
| 项目 | 要求 |
|------|------|
| CPU | 2 核 |
| 内存 | 4 GB |
| 带宽 | 2 Mbps |
| 系统盘 | 40 GB SSD |
| 数据盘 | 50 GB（可选） |

### 推荐配置
| 项目 | 要求 |
|------|------|
| CPU | 4 核 |
| 内存 | 8 GB |
| 带宽 | 5 Mbps |
| 系统盘 | 80 GB SSD |

### 支持的操作系统
- Ubuntu 20.04+ (推荐)
- CentOS 7+
- Debian 11+

---

## 二、购买云服务器流程

### 1. 选择云服务商
- 阿里云 ECS
- 腾讯云 CVM
- 华为云 ECS
- AWS EC2

### 2. 创建实例
1. 登录云服务商控制台
2. 进入 ECS/云服务器 购买页面
3. 选择配置（参考上表）
4. 选择操作系统 Ubuntu 20.04 LTS
5. 设置 root 密码
6. 开放必要端口：22(SSH), 80(HTTP), 443(HTTPS), 8080(后端API)

### 3. 安全组配置
```
入站规则：
- 22端口 - SSH (0.0.0.0/0 或您的IP)
- 80端口 - HTTP (0.0.0.0/0)
- 443端口 - HTTPS (0.0.0.0/0)
- 8080端口 - API (0.0.0.0/0)
```

---

## 三、本地准备工作

### 1. 安装 Git
```bash
# macOS
brew install git

# Ubuntu/Debian
sudo apt update && sudo apt install git -y
```

### 2. 安装 Docker Desktop (macOS/Windows) 或在服务器上安装 Docker

### 3. 克隆项目
```bash
git clone https://github.com/stomyrose/xuexizhushou-.git
cd xuexizhushou-
```

---

## 四、服务器环境配置

### 1. 连接服务器
```bash
ssh root@您的服务器IP
```

### 2. 安装 Docker
```bash
# Ubuntu/Debian
apt update && apt upgrade -y
apt install -y curl wget vim docker.io docker-compose

# 启动 Docker
systemctl start docker
systemctl enable docker

# 添加当前用户到 docker 组
usermod -aG docker $USER
```

### 3. 安装 Docker Compose
```bash
curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
docker-compose --version
```

---

## 五、部署应用

### 1. 上传项目到服务器

**方式一：Git 拉取（推荐）**
```bash
cd /opt
git clone https://github.com/stomyrose/xuexizhushou-.git
cd xuexizhushou-
```

**方式二：SCP 上传**
```bash
# 本地执行
scp -r ./xuexizhushou- root@您的服务器IP:/opt/
```

### 2. 配置环境变量
```bash
cd /opt/xuexizhushou-

# 创建环境变量文件
cat > .env << EOF
# 数据库
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=您的数据库密码
DB_NAME=force_learning

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=您的JWT密钥(至少32位随机字符串)

# 服务
SERVER_PORT=8080
UPLOAD_PATH=/app/uploads

# 支付宝(可选)
ALIPAY_APP_ID=您的支付宝AppID
ALIPAY_PRIVATE_KEY=您的支付宝私钥
ALIPAY_PUBLIC_KEY=支付宝公钥
ALIPAY_NOTIFY_URL=https://您的域名/api/v1/alipay/callback

# 微信支付(可选)
WXPAY_APP_ID=您的微信AppID
WXPAY_MCH_ID=您的微信商户号
WXPAY_API_KEY=您的微信API密钥
WXPAY_NOTIFY_URL=https://您的域名/api/v1/wxpay/callback
EOF
```

### 3. 生成 SSL 证书（可选，使用 HTTPS）
```bash
# 安装 Certbot
apt install -y certbot

# 申请证书(需要有域名)
certbot certonly --standalone -d您的域名 --email您的邮箱 --agree-tos --non-interactive
```

### 4. 启动服务
```bash
# 使用 Docker Compose 启动所有服务
docker-compose -f docker-compose.prod.yml up -d

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 5. 初始化数据库
```bash
# 后端服务会自动创建表结构，首次启动时会自动迁移
# 等待几秒后检查日志
docker-compose -f docker-compose.prod.yml logs backend
```

---

## 六、Nginx 配置

### 1. 创建 Nginx 配置
```bash
cat > /etc/nginx/sites-available/force-learning << EOF
server {
    listen 80;
    server_name 您的域名;

    client_max_body_size 100M;

    # API 代理到后端
    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Web 静态文件
    location / {
        root /opt/xuexizhushou-/web/dist;
        index index.html;
        try_files \$uri \$uri/ /index.html;
    }

    # 上传文件
    location /uploads/ {
        alias /opt/xuexizhushou-/uploads/;
        expires 30d;
    }
}
EOF
```

### 2. 启用配置
```bash
ln -s /etc/nginx/sites-available/force-learning /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx
```

### 3. 配置 HTTPS（可选）
```bash
cat > /etc/nginx/sites-available/force-learning-ssl << EOF
server {
    listen 443 ssl http2;
    server_name 您的域名;

    ssl_certificate /etc/letsencrypt/live/您的域名/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/您的域名/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 100M;

    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        root /opt/xuexizhushou-/web/dist;
        index index.html;
        try_files \$uri \$uri/ /index.html;
    }

    location /uploads/ {
        alias /opt/xuexizhushou-/uploads/;
        expires 30d;
    }
}

server {
    listen 80;
    server_name 您的域名;
    return 301 https://\$server_name\$request_uri;
}
EOF

ln -s /etc/nginx/sites-available/force-learning-ssl /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx
```

---

## 七、服务管理命令

### 启动服务
```bash
docker-compose -f /opt/xuexizhushou-/docker-compose.prod.yml up -d
```

### 停止服务
```bash
docker-compose -f /opt/xuexizhushou-/docker-compose.prod.yml down
```

### 重启服务
```bash
docker-compose -f /opt/xuexizhushou-/docker-compose.prod.yml restart
```

### 查看日志
```bash
# 所有服务
docker-compose -f /opt/xuexizhushou-/docker-compose.prod.yml logs -f

# 指定服务
docker-compose -f /opt/xuexizhushou-/docker-compose.prod.yml logs -f backend
```

### 更新代码并重启
```bash
cd /opt/xuexizhushou-
git pull origin main
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build
```

---

## 八、备份策略

### 1. 数据库备份
```bash
# 创建备份脚本
cat > /opt/backup-db.sh << EOF
#!/bin/bash
DATE=\$(date +%Y%m%d_%H%M%S)
BACKUP_DIR=/opt/backups
mkdir -p \$BACKUP_DIR

docker exec force_learning_postgres pg_dump -U postgres force_learning > \$BACKUP_DIR/db_\$DATE.sql
find \$BACKUP_DIR -name "db_*.sql" -mtime +7 -delete
EOF

chmod +x /opt/backup-db.sh

# 添加定时任务
crontab -e
# 每天凌晨2点执行备份
0 2 * * * /opt/backup-db.sh
```

### 2. 文件备份
```bash
# 备份上传文件
tar -czf /opt/backups/uploads_$(date +%Y%m%d).tar.gz /opt/xuexizhushou-/uploads/
```

---

## 九、监控与日志

### 1. 查看资源使用
```bash
docker stats
```

### 2. 检查端口占用
```bash
ss -tlnp | grep -E ':(80|443|8080|5432|6379)'
```

### 3. 防火墙配置
```bash
# Ubuntu
ufw allow 22
ufw allow 80
ufw allow 443
ufw allow 8080
ufw enable
```

---

## 十、域名配置

### 1. DNS 解析
在域名服务商控制台添加解析记录：
```
A记录 @ 您的服务器IP
A记录 www 您的服务器IP
```

### 2. 申请免费 SSL 证书
```bash
certbot --nginx -d 您的域名 --email 您的邮箱 --agree-tos --non-interactive
```

---

## 十一、故障排查

### 服务无法启动
```bash
# 查看详细日志
docker-compose -f docker-compose.prod.yml logs

# 检查端口占用
netstat -tlnp | grep 8080
```

### 数据库连接失败
```bash
# 检查 PostgreSQL 容器
docker exec force_learning_postgres psql -U postgres -c "\l"

# 检查连接
docker exec force_learning_backend sh -c "nc -zv postgres 5432"
```

### Nginx 502 错误
```bash
# 检查后端是否运行
curl http://127.0.0.1:8080/health

# 检查 Nginx 日志
tail -f /var/log/nginx/error.log
```

---

## 十二、成本估算（月）

| 项目 | 配置 | 价格(约) |
|------|------|----------|
| 云服务器 | 2核4G 5Mbps | ¥100-200 |
| 域名 | .com | ¥60/年 ≈ ¥5/月 |
| SSL证书 | Let's Encrypt免费 | ¥0 |
| 备份存储 | 10GB | ¥5/月 |
| **合计** | | **¥110-210/月** |
