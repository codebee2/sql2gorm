# SQL2GORM

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22.2+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)

**🚀 专业的MySQL表结构转GORM模型代码生成工具**

*支持DSN连接数据库直接读取表结构，自动生成符合Go代码规范的GORM模型*

[快速开始](#-快速开始) • [使用指南](#-使用指南) • [配置说明](#️-配置说明) • [开发指南](#-开发指南)

</div>

---

## ✨ 主要特性

- 🗄️ **DSN直连**：支持直接连接数据库读取表结构，无需手动导出SQL
- 📝 **智能生成**：自动解析表结构，智能生成GORM模型
- 🎨 **代码美化**：自动应用Go代码格式化，生成专业级代码
- ⚙️ **类型映射**：智能映射MySQL类型到Go类型
- 🏷️ **标签管理**：自动添加JSON和GORM标签
- 🔧 **配置灵活**：支持配置文件、命令行参数、环境变量多种配置方式

---

## 🚀 快速开始

### 安装

```bash
# 推荐方式：使用 go install 安装
go install github.com/codebee2/sql2gorm@latest

# 验证安装
sql2gorm --help
```

### 快速体验

```bash
# 使用DSN连接数据库（推荐）
sql2gorm \
  --dsn "user:password@tcp(localhost:3306)/database?charset=utf8mb4" \
  --output models/user.go \
  --table users

# 使用SQL文件
sql2gorm \
  --sql schema.sql \
  --output models/product.go \
  --table products
```

---

## 📖 使用指南

### DSN方式（推荐）

#### 命令行参数
```bash
sql2gorm \
  --dsn "user:password@tcp(localhost:3306)/database?charset=utf8mb4" \
  --output models/user.go \
  --table users \
  --package models
```

#### 配置文件
```yaml
# config-dsn.yaml
dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
output_file: "models/user.go"
table_name: "users"
package_name: "models"
```

```bash
sql2gorm --config config-dsn.yaml
```

#### 环境变量
```bash
export SQL2GORM_DSN="user:password@tcp(localhost:3306)/database?charset=utf8mb4"
export SQL2GORM_OUTPUT_FILE="models/user.go"
export SQL2GORM_TABLE_NAME="users"
export SQL2GORM_PACKAGE_NAME="models"

sql2gorm
```

### SQL文件方式

#### 命令行参数
```bash
sql2gorm \
  --sql schema.sql \
  --output models/product.go \
  --table products \
  --package models
```

#### 配置文件
```yaml
# config.yaml
sql_file: "schema.sql"
output_file: "models/product.go"
table_name: "products"
package_name: "models"
```

```bash
sql2gorm --config config.yaml
```

---

## ⚙️ 配置说明

### 配置参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| `dsn` | string | 否* | 数据库连接字符串 | `user:pass@tcp(localhost:3306)/db` |
| `sql_file` | string | 否* | SQL文件路径 | `schema.sql` |
| `output_file` | string | 是 | 输出Go文件路径 | `models/user.go` |
| `table_name` | string | 是 | 要生成模型的表名 | `users` |
| `package_name` | string | 否 | Go包名 | `models` |

> **注意**: `dsn` 和 `sql_file` 至少需要提供一个，DSN优先级更高

### 环境变量

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| `SQL2GORM_DSN` | 数据库连接字符串 | `user:pass@tcp(localhost:3306)/db` |
| `SQL2GORM_SQL_FILE` | SQL文件路径 | `schema.sql` |
| `SQL2GORM_OUTPUT_FILE` | 输出文件路径 | `models/user.go` |
| `SQL2GORM_TABLE_NAME` | 表名 | `users` |
| `SQL2GORM_PACKAGE_NAME` | 包名 | `models` |

### DSN连接字符串格式

```
username:password@tcp(host:port)/database?param1=value1&param2=value2
```

**常用参数：**
- `charset=utf8mb4` - 字符集
- `parseTime=True` - 时间解析
- `loc=Local` - 时区设置
- `timeout=5s` - 连接超时

**示例：**
```bash
# 本地数据库
root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local

# 远程数据库
user:pass@tcp(db.example.com:3306)/production?charset=utf8mb4&timeout=10s
```

---

## 📊 类型映射规则

### MySQL → Go 类型映射

| MySQL类型 | Go类型 | 说明 |
|-----------|---------|------|
| `int`, `tinyint`, `smallint`, `mediumint`, `bigint` | `int64` | 整数类型 |
| `float`, `double`, `decimal` | `float64` | 浮点数类型 |
| `varchar`, `char`, `text`, `longtext` | `string` | 字符串类型 |
| `date`, `datetime`, `timestamp`, `time` | `string` | 日期时间类型 |
| `blob`, `longblob` | `[]byte` | 二进制类型 |
| `json` | `string` | JSON类型 |

### GORM标签规则

| 标签 | 说明 | 示例 |
|------|------|------|
| `column` | 指定数据库列名 | `gorm:"column:user_name"` |
| `primaryKey` | 标识主键字段 | `gorm:"primaryKey"` |
| `autoIncrement` | 标识自增字段 | `gorm:"autoIncrement"` |
| `default` | 设置字段默认值 | `gorm:"default:0"` |
| `not null` | 标识非空字段 | `gorm:"not null"` |

---

## 💡 使用示例

### 输入SQL示例
```sql
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1启用，0禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 生成的GORM模型
```go
package models

var UserModelIns = new(UserModel)

// UserModel 用户表
type UserModel struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;not null"`
	Username  string `json:"username" gorm:"column:username;not null"`                    // 用户名
	Email     string `json:"email" gorm:"column:email"`                                   // 邮箱
	Status    int64  `json:"status" gorm:"column:status;default:1"`                      // 状态：1启用，0禁用
	CreatedAt string `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at;default:null"`
}

// TableName 表名称
func (own *UserModel) TableName() string {
	return "users"
}
```

---

## 🔧 开发指南

### 项目结构
```
sql2gorm/
├── cmd/           # 命令行处理
├── config/        # 配置管理
├── database/      # 数据库连接
├── generator/     # 代码生成
├── parser/        # SQL解析
├── main.go        # 主程序入口
└── Makefile       # 构建脚本
```

### 构建和安装
```bash
# 方式1：使用 go install 安装（推荐）
go install github.com/codebee2/sql2gorm@latest

# 方式2：从源码构建
git clone https://github.com/codebee2/sql2gorm.git
cd sql2gorm
make deps
make build

# 运行测试
make test
```

### 技术栈
- **sqlparser**: 专业的SQL解析库
- **cobra**: 强大的命令行框架
- **viper**: 灵活的配置管理
- **inflection**: 智能的单词单复数转换
- **go/format**: Go标准代码格式化工具

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 贡献方式
1. 🐛 **报告Bug**：在Issues中报告问题
2. 💡 **功能建议**：提出新功能想法
3. 🔧 **代码贡献**：提交Pull Request
4. 📚 **文档改进**：完善文档和示例
5. 🌟 **Star支持**：给项目点个Star

### 贡献流程
1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

---

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

---

## 🙏 致谢

感谢以下开源项目的支持：
- [sqlparser](https://github.com/xwb1989/sqlparser) - SQL解析库
- [cobra](https://github.com/spf13/cobra) - 命令行框架
- [viper](https://github.com/spf13/viper) - 配置管理
- [inflection](https://github.com/jinzhu/inflection) - 单词转换

---

<div align="center">

**如果这个项目对您有帮助，请给它一个 ⭐ Star！**

Made with ❤️ by [codebee2](https://github.com/codebee2)

</div>