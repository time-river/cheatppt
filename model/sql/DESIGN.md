
1. 用户表（users）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 用户ID（自增）                                               |
| username       | VARCHAR(255)   | 用户名                                                       |
| email          | VARCHAR(255)   | 邮箱地址                                                     |
| password_hash  | VARCHAR(255)   | 密码哈希值                                                   |
| virtual_currency | INT          | 用户的付费虚拟货币数额                                       |

2. 机器人类型表（robot_types）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 机器人类型ID（自增）                                         |
| provider_name  | VARCHAR(255)   | 提供商名称                                                   |
| display_name   | VARCHAR(255)   | 显示名称                                                     |
| actual_name    | VARCHAR(255)   | 实际名称                                                     |
| access_level   | ENUM           | 访问等级（免费用户、付费用户等级1、付费用户等级2等）         |

3. 聊天上下文表（chat_contexts）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 聊天上下文ID（自增）                                         |
| user_id        | INT            | 用户ID（外键，关联 users 表）                                |
| robot_type_id  | INT            | 机器人类型ID（外键，关联 robot_types 表）                   |

4. 聊天记录表（chat_records）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 聊天记录ID（自增）                                           |
| chat_context_id | INT          | 聊天上下文ID（外键，关联 chat_contexts 表）                 |
| message        | TEXT           | 聊天记录内容                                                 |
| timestamp      | DATETIME       | 聊天记录时间戳                                               |

5. 机器人切换映射表（robot_switch_mapping）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 映射ID（自增）                                               |
| from_robot_type_id | INT          | 源机器人类型ID（外键，关联 robot_types 表）                 |
| to_robot_type_id  | INT          | 目标机器人类型ID（外键，关联 robot_types 表）               |

6. 用户充值记录表（user_recharge_records）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 充值记录ID（自增）                                           |
| user_id        | INT            | 用户ID（外键，关联 users 表）                                |
| amount         | DECIMAL(10, 2) | 充值金额                                                     |
| virtual_currency | INT           | 增加的虚拟货币数额                                           |
| recharge_time  | DATETIME       | 充值时间                                                     |
| payment_method | ENUM           | 付款方式（例如信用卡、PayPal、支付宝等）                     |

7. 每日免费虚拟币表（daily_free_virtual_currency）

| 列名           | 类型           | 说明                                                         |
| -------------- | -------------- | ------------------------------------------------------------ |
| id             | INT PRIMARY KEY | 主键（自增）                                                 |
| user_id        | INT            | 用户ID（外键，关联 users 表）                                |
| date           | DATE           | 免费虚拟币的有效日期                                         |
| free_virtual_currency | INT       | 免费虚拟币额度                                               |

以下是相应的 MySQL 创建表语句：

```sql
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  virtual_currency INT DEFAULT 0,
  INDEX (username)
);

CREATE TABLE robot_types (
  id INT AUTO_INCREMENT PRIMARY KEY,
  provider_name VARCHAR(255) NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  actual_name VARCHAR(255) NOT NULL,
  access_level ENUM('free', 'paid_level_1', 'paid_level_2') NOT NULL,
  INDEX (provider_name)
);

CREATE TABLE chat_contexts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  robot_type_id INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (robot_type_id) REFERENCES robot_types(id),
  INDEX (user_id),
  INDEX (robot_type_id)
);

CREATE TABLE chat_records (
  id INT AUTO_INCREMENT PRIMARY KEY,
  chat_context_id INT NOT NULL,
  message TEXT NOT NULL,
  timestamp DATETIME NOT NULL,
  FOREIGN KEY (chat_context_id) REFERENCES chat_contexts(id),
  INDEX (chat_context_id),
  INDEX (timestamp)
);

CREATE TABLE robot_switch_mapping (
  id INT AUTO_INCREMENT PRIMARY KEY,
  from_robot_type_id INT NOT NULL,
  to_robot_type_id INT NOT NULL,
  FOREIGN KEY (from_robot_type_id) REFERENCES robot_types(id),
  FOREIGN KEY (to_robot_type_id) REFERENCES robot_types(id),
  INDEX (from_robot_type_id),
  INDEX (to_robot_type_id)
);

CREATE TABLE user_recharge_records (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  amount DECIMAL(10, 2) NOT NULL,
  virtual_currency INT NOT NULL,
  recharge_time DATETIME NOT NULL,
  payment_method ENUM('credit_card', 'paypal', 'alipay') NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX (user_id),
  INDEX (recharge_time),
  INDEX (payment_method)
);

CREATE TABLE daily_free_virtual_currency (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  date DATE NOT NULL,
  free_virtual_currency INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX (user_id),
  INDEX (date)
);
```