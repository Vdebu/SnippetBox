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
![home](https://github.com/user-attachments/assets/82684fb7-b53f-4f13-ab74-fccf14ad2e93)
# 消息发布页面展示
![createsnippet](https://github.com/user-attachments/assets/339cce99-6402-48e6-9f51-ab295f097b3a)

# 登入页面展示
![login](https://github.com/user-attachments/assets/37b14493-71d1-4cbe-97a5-e8ecd0e27a86)

# 注册页面展示
![signup](https://github.com/user-attachments/assets/d9f3e572-b612-4e1b-9de1-4f6a93385044)


