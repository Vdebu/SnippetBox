# 项目：极简文本与代码片段托管平台
# 适用场景：临时日志存储、代码协作、配置共享等
# 技术栈：Go | HTML/CSS | MySQL

# MySQL 配置
```yaml
mysql:
  version: "8.0"
  host: "localhost"
  port: 3306
  username: "web"
  password: "pass"
  database: "snippetbox"
数据表定义
1. Snippets 表
用于存储用户输入的相关信息。
CREATE TABLE snippets (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATETIME NOT NULL,
  expires DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_snippets_created ON snippets(created);
2. Users 表
用于存储用户的账户信息。
CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL,
  created DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE users 
  ADD CONSTRAINT users_uc_email UNIQUE(email);
3. Sessions 表
用于管理会话信息。
CREATE TABLE sessions (
  token CHAR(43) PRIMARY KEY,
  data BLOB NOT NULL,
  expiry TIMESTAMP(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE INDEX sessions_expiry_idx ON sessions(expiry);
```
# 主页面展示
![home](https://github.com/user-attachments/assets/6a5b80f1-0603-444d-b889-da72128fe487)

# 消息发布页面展示
![snippet](https://github.com/user-attachments/assets/58a3d64b-e6f0-4be5-a302-9ecc7e222c84)

# 账号管理页面展示
![account](https://github.com/user-attachments/assets/eb007d3b-0ca1-4fcb-9082-4c155f7b4969)

# 登入页面展示
![signin](https://github.com/user-attachments/assets/112d44cf-bf6e-42af-9e56-d94d3d4234f0)

# 注册页面展示
![signup](https://github.com/user-attachments/assets/94a3d416-ec64-446d-b8b8-0e7fd6583619)

# 关于页面展示
![about](https://github.com/user-attachments/assets/5305c4e3-fe04-4541-ab81-875bd5ab3329)

